;;; doubi-mode.el --- Major mode for the Doubi programming language

;; Copyright 2013 The Doubi Authors. All rights reserved.
;; Use of this source code is doubiverned by a BSD-style
;; license that can be found in the LICENSE file.

(require 'cl)
(require 'etags)
(require 'ffap)
(require 'ring)
(require 'url)

;; XEmacs compatibility guidelines
;; - Minimum required version of XEmacs: 21.5.32
;;   - Feature that cannot be backported: POSIX character classes in
;;     regular expressions
;;   - Functions that could be backported but won't because 21.5.32
;;     covers them: plenty.
;;   - Features that are still partly broken:
;;     - doubidef will not work correctly if multibyte characters are
;;       being used
;;     - Fontification will not handle unicode correctly
;;
;; - Do not use \_< and \_> regexp delimiters directly; use
;;   doubi--regexp-enclose-in-symbol
;;
;; - The character `_` must not be a symbol constituent but a
;;   character constituent
;;
;; - Do not use process-lines
;;
;; - Use doubi--old-completion-list-style when using a plain list as the
;;   collection for completing-read
;;
;; - Use doubi--kill-whole-line instead of kill-whole-line (called
;;   kill-entire-line in XEmacs)
;;
;; - Use doubi--position-bytes instead of position-bytes
(defmacro doubi--xemacs-p ()
  `(featurep 'xemacs))

(defalias 'doubi--kill-whole-line
  (if (fboundp 'kill-whole-line)
      #'kill-whole-line
    #'kill-entire-line))

;; Delete the current line without putting it in the kill-ring.
(defun doubi--delete-whole-line (&optional arg)
  ;; Emacs uses both kill-region and kill-new, Xemacs only uses
  ;; kill-region. In both cases we turn them into operations that do
  ;; not modify the kill ring. This solution does depend on the
  ;; implementation of kill-line, but it's the only viable solution
  ;; that does not require to write kill-line from scratch.
  (flet ((kill-region (beg end)
                      (delete-region beg end))
         (kill-new (s) ()))
    (doubi--kill-whole-line arg)))

;; declare-function is an empty macro that only byte-compile cares
;; about. Wrap in always false if to satisfy Emacsen without that
;; macro.
(if nil
    (declare-function doubi--position-bytes "doubi-mode" (point)))
;; XEmacs unfortunately does not offer position-bytes. We can fall
;; back to just using (point), but it will be incorrect as soon as
;; multibyte characters are being used.
(if (fboundp 'position-bytes)
    (defalias 'doubi--position-bytes #'position-bytes)
  (defun doubi--position-bytes (point) point))

(defun doubi--old-completion-list-style (list)
  (mapcar (lambda (x) (cons x nil)) list))

;; GNU Emacs 24 has prog-mode, older GNU Emacs and XEmacs do not, so
;; copy its definition for those.
(if (not (fboundp 'prog-mode))
    (define-derived-mode prog-mode fundamental-mode "Prog"
      "Major mode for editing source code."
      (set (make-local-variable 'require-final-newline) mode-require-final-newline)
      (set (make-local-variable 'parse-sexp-ignore-comments) t)
      (setq bidi-paragraph-direction 'left-to-right)))

(defun doubi--regexp-enclose-in-symbol (s)
  ;; XEmacs does not support \_<, GNU Emacs does. In GNU Emacs we make
  ;; extensive use of \_< to support unicode in identifiers. Until we
  ;; come up with a better solution for XEmacs, this solution will
  ;; break fontification in XEmacs for identifiers such as "typeµ".
  ;; XEmacs will consider "type" a keyword, GNU Emacs won't.

  (if (doubi--xemacs-p)
      (concat "\\<" s "\\>")
    (concat "\\_<" s "\\_>")))

;; Move up one level of parentheses.
(defun doubi-goto-opening-parenthesis (&optional legacy-unused)
  ;; The old implementation of doubi-goto-opening-parenthesis had an
  ;; optional argument to speed up the function. It didn't change the
  ;; function's outcome.

  ;; Silently fail if there's no matching opening parenthesis.
  (condition-case nil
      (backward-up-list)
    (scan-error nil)))


(defconst doubi-dangling-operators-regexp "[^-]-\\|[^+]\\+\\|[/*&><.=|^]")
(defconst doubi-identifier-regexp "[[:word:][:multibyte:]]+")
(defconst doubi-label-regexp doubi-identifier-regexp)
(defconst doubi-type-regexp "[[:word:][:multibyte:]*]+")
(defconst doubi-func-regexp (concat (doubi--regexp-enclose-in-symbol "func") "\\s *\\(" doubi-identifier-regexp "\\)"))
(defconst doubi-func-meth-regexp (concat
                               (doubi--regexp-enclose-in-symbol "func") "\\s *\\(?:(\\s *"
                               "\\(" doubi-identifier-regexp "\\s +\\)?" doubi-type-regexp
                               "\\s *)\\s *\\)?\\("
                               doubi-identifier-regexp
                               "\\)("))
(defconst doubi-builtins
  '("append" "cap"   "close"   "complex" "copy"
    "delete" "imag"  "len"     "make"    "new"
    "panic"  "print" "println" "real"    "recover")
  "All built-in functions in the Doubi language. Used for font locking.")

(defconst doubi-mode-keywords
  '("break"    "default"     "func"   "interface" "select"
    "case"     "defer"       "doubi"     "map"       "struct"
    "chan"     "else"        "goto"   "package"   "switch"
    "const"    "fallthrough" "if"     "range"     "type"
    "continue" "for"         "import" "return"    "var")
  "All keywords in the Doubi language.  Used for font locking.")

(defconst doubi-constants '("nil" "true" "false" "iota"))
(defconst doubi-type-name-regexp (concat "\\(?:[*(]\\)*\\(?:" doubi-identifier-regexp "\\.\\)?\\(" doubi-identifier-regexp "\\)"))

(defvar doubi-dangling-cache)
(defvar doubi-doubidoc-history nil)
(defvar doubi--coverage-current-file-name)

(defgroup doubi nil
  "Major mode for editing Doubi code"
  :group 'languages)

(defgroup doubi-cover nil
  "Options specific to `cover`"
  :group 'doubi)

(defcustom doubi-fontify-function-calls t
  "Fontify function and method calls if this is non-nil."
  :type 'boolean
  :group 'doubi)

(defcustom doubi-mode-hook nil
  "Hook called by `doubi-mode'."
  :type 'hook
  :group 'doubi)

(defcustom doubi-command "doubi"
  "The 'doubi' command.  Some users have multiple Doubi development
trees and invoke the 'doubi' tool via a wrapper that sets DOUBIROOT and
DOUBIPATH based on the current directory.  Such users should
customize this variable to point to the wrapper script."
  :type 'string
  :group 'doubi)

(defcustom doubifmt-command "doubifmt"
  "The 'doubifmt' command.  Some users may replace this with 'doubiimports'
from https://github.com/bradfitz/doubiimports."
  :type 'string
  :group 'doubi)

(defface doubi-coverage-untracked
  '((t (:foreground "#505050")))
  "Coverage color of untracked code."
  :group 'doubi-cover)

(defface doubi-coverage-0
  '((t (:foreground "#c00000")))
  "Coverage color for uncovered code."
  :group 'doubi-cover)
(defface doubi-coverage-1
  '((t (:foreground "#808080")))
  "Coverage color for covered code with weight 1."
  :group 'doubi-cover)
(defface doubi-coverage-2
  '((t (:foreground "#748c83")))
  "Coverage color for covered code with weight 2."
  :group 'doubi-cover)
(defface doubi-coverage-3
  '((t (:foreground "#689886")))
  "Coverage color for covered code with weight 3."
  :group 'doubi-cover)
(defface doubi-coverage-4
  '((t (:foreground "#5ca489")))
  "Coverage color for covered code with weight 4."
  :group 'doubi-cover)
(defface doubi-coverage-5
  '((t (:foreground "#50b08c")))
  "Coverage color for covered code with weight 5."
  :group 'doubi-cover)
(defface doubi-coverage-6
  '((t (:foreground "#44bc8f")))
  "Coverage color for covered code with weight 6."
  :group 'doubi-cover)
(defface doubi-coverage-7
  '((t (:foreground "#38c892")))
  "Coverage color for covered code with weight 7."
  :group 'doubi-cover)
(defface doubi-coverage-8
  '((t (:foreground "#2cd495")))
  "Coverage color for covered code with weight 8.
For mode=set, all covered lines will have this weight."
  :group 'doubi-cover)
(defface doubi-coverage-9
  '((t (:foreground "#20e098")))
  "Coverage color for covered code with weight 9."
  :group 'doubi-cover)
(defface doubi-coverage-10
  '((t (:foreground "#14ec9b")))
  "Coverage color for covered code with weight 10."
  :group 'doubi-cover)
(defface doubi-coverage-covered
  '((t (:foreground "#2cd495")))
  "Coverage color of covered code."
  :group 'doubi-cover)

(defvar doubi-mode-syntax-table
  (let ((st (make-syntax-table)))
    (modify-syntax-entry ?+  "." st)
    (modify-syntax-entry ?-  "." st)
    (modify-syntax-entry ?%  "." st)
    (modify-syntax-entry ?&  "." st)
    (modify-syntax-entry ?|  "." st)
    (modify-syntax-entry ?^  "." st)
    (modify-syntax-entry ?!  "." st)
    (modify-syntax-entry ?=  "." st)
    (modify-syntax-entry ?<  "." st)
    (modify-syntax-entry ?>  "." st)
    (modify-syntax-entry ?/ (if (doubi--xemacs-p) ". 1456" ". 124b") st)
    (modify-syntax-entry ?*  ". 23" st)
    (modify-syntax-entry ?\n "> b" st)
    (modify-syntax-entry ?\" "\"" st)
    (modify-syntax-entry ?\' "\"" st)
    (modify-syntax-entry ?`  "\"" st)
    (modify-syntax-entry ?\\ "\\" st)
    ;; It would be nicer to have _ as a symbol constituent, but that
    ;; would trip up XEmacs, which does not support the \_< anchor
    (modify-syntax-entry ?_  "w" st)

    st)
  "Syntax table for Doubi mode.")

(defun doubi--build-font-lock-keywords ()
  ;; we cannot use 'symbols in regexp-opt because emacs <24 doesn't
  ;; understand that
  (append
   `((,(doubi--regexp-enclose-in-symbol (regexp-opt doubi-mode-keywords t)) . font-lock-keyword-face)
     (,(doubi--regexp-enclose-in-symbol (regexp-opt doubi-builtins t)) . font-lock-builtin-face)
     (,(doubi--regexp-enclose-in-symbol (regexp-opt doubi-constants t)) . font-lock-constant-face)
     (,doubi-func-regexp 1 font-lock-function-name-face)) ;; function (not method) name

   (if doubi-fontify-function-calls
       `((,(concat "\\(" doubi-identifier-regexp "\\)[[:space:]]*(") 1 font-lock-function-name-face) ;; function call/method name
         (,(concat "[^[:word:][:multibyte:]](\\(" doubi-identifier-regexp "\\))[[:space:]]*(") 1 font-lock-function-name-face)) ;; bracketed function call
     `((,doubi-func-meth-regexp 1 font-lock-function-name-face))) ;; method name

   `(
     (,(concat (doubi--regexp-enclose-in-symbol "type") "[[:space:]]*\\([^[:space:]]+\\)") 1 font-lock-type-face) ;; types
     (,(concat (doubi--regexp-enclose-in-symbol "type") "[[:space:]]*" doubi-identifier-regexp "[[:space:]]*" doubi-type-name-regexp) 1 font-lock-type-face) ;; types
     (,(concat "[^[:word:][:multibyte:]]\\[\\([[:digit:]]+\\|\\.\\.\\.\\)?\\]" doubi-type-name-regexp) 2 font-lock-type-face) ;; Arrays/slices
     (,(concat "\\(" doubi-identifier-regexp "\\)" "{") 1 font-lock-type-face)
     (,(concat (doubi--regexp-enclose-in-symbol "map") "\\[[^]]+\\]" doubi-type-name-regexp) 1 font-lock-type-face) ;; map value type
     (,(concat (doubi--regexp-enclose-in-symbol "map") "\\[" doubi-type-name-regexp) 1 font-lock-type-face) ;; map key type
     (,(concat (doubi--regexp-enclose-in-symbol "chan") "[[:space:]]*\\(?:<-\\)?" doubi-type-name-regexp) 1 font-lock-type-face) ;; channel type
     (,(concat (doubi--regexp-enclose-in-symbol "\\(?:new\\|make\\)") "\\(?:[[:space:]]\\|)\\)*(" doubi-type-name-regexp) 1 font-lock-type-face) ;; new/make type
     ;; TODO do we actually need this one or isn't it just a function call?
     (,(concat "\\.\\s *(" doubi-type-name-regexp) 1 font-lock-type-face) ;; Type conversion
     (,(concat (doubi--regexp-enclose-in-symbol "func") "[[:space:]]+(" doubi-identifier-regexp "[[:space:]]+" doubi-type-name-regexp ")") 1 font-lock-type-face) ;; Method receiver
     (,(concat (doubi--regexp-enclose-in-symbol "func") "[[:space:]]+(" doubi-type-name-regexp ")") 1 font-lock-type-face) ;; Method receiver without variable name
     ;; Like the original doubi-mode this also marks compound literal
     ;; fields. There, it was marked as to fix, but I grew quite
     ;; accustomed to it, so it'll stay for now.
     (,(concat "^[[:space:]]*\\(" doubi-label-regexp "\\)[[:space:]]*:\\(\\S.\\|$\\)") 1 font-lock-constant-face) ;; Labels and compound literal fields
     (,(concat (doubi--regexp-enclose-in-symbol "\\(goto\\|break\\|continue\\)") "[[:space:]]*\\(" doubi-label-regexp "\\)") 2 font-lock-constant-face)))) ;; labels in goto/break/continue

(defvar doubi-mode-map
  (let ((m (make-sparse-keymap)))
    (define-key m "}" #'doubi-mode-insert-and-indent)
    (define-key m ")" #'doubi-mode-insert-and-indent)
    (define-key m "," #'doubi-mode-insert-and-indent)
    (define-key m ":" #'doubi-mode-insert-and-indent)
    (define-key m "=" #'doubi-mode-insert-and-indent)
    (define-key m (kbd "C-c C-a") #'doubi-import-add)
    (define-key m (kbd "C-c C-j") #'doubidef-jump)
    (define-key m (kbd "C-x 4 C-c C-j") #'doubidef-jump-other-window)
    (define-key m (kbd "C-c C-d") #'doubidef-describe)
    m)
  "Keymap used by Doubi mode to implement electric keys.")

(defun doubi-mode-insert-and-indent (key)
  "Invoke the global binding of KEY, then reindent the line."

  (interactive (list (this-command-keys)))
  (call-interactively (lookup-key (current-global-map) key))
  (indent-according-to-mode))

(defmacro doubi-paren-level ()
  `(car (syntax-ppss)))

(defmacro doubi-in-string-or-comment-p ()
  `(nth 8 (syntax-ppss)))

(defmacro doubi-in-string-p ()
  `(nth 3 (syntax-ppss)))

(defmacro doubi-in-comment-p ()
  `(nth 4 (syntax-ppss)))

(defmacro doubi-goto-beginning-of-string-or-comment ()
  `(goto-char (nth 8 (syntax-ppss))))

(defun doubi--backward-irrelevant (&optional stop-at-string)
  "Skips backwards over any characters that are irrelevant for
indentation and related tasks.

It skips over whitespace, comments, cases and labels and, if
STOP-AT-STRING is not true, over strings."

  (let (pos (start-pos (point)))
    (skip-chars-backward "\n\s\t")
    (if (and (save-excursion (beginning-of-line) (doubi-in-string-p)) (looking-back "`") (not stop-at-string))
        (backward-char))
    (if (and (doubi-in-string-p) (not stop-at-string))
        (doubi-goto-beginning-of-string-or-comment))
    (if (looking-back "\\*/")
        (backward-char))
    (if (doubi-in-comment-p)
        (doubi-goto-beginning-of-string-or-comment))
    (setq pos (point))
    (beginning-of-line)
    (if (or (looking-at (concat "^" doubi-label-regexp ":")) (looking-at "^[[:space:]]*\\(case .+\\|default\\):"))
        (end-of-line 0)
      (goto-char pos))
    (if (/= start-pos (point))
        (doubi--backward-irrelevant stop-at-string))
    (/= start-pos (point))))

(defun doubi--buffer-narrowed-p ()
  "Return non-nil if the current buffer is narrowed."
  (/= (buffer-size)
      (- (point-max)
         (point-min))))

(defun doubi-previous-line-has-dangling-op-p ()
  "Returns non-nil if the current line is a continuation line."
  (let* ((cur-line (line-number-at-pos))
         (val (gethash cur-line doubi-dangling-cache 'nope)))
    (if (or (doubi--buffer-narrowed-p) (equal val 'nope))
        (save-excursion
          (beginning-of-line)
          (doubi--backward-irrelevant t)
          (setq val (looking-back doubi-dangling-operators-regexp))
          (if (not (doubi--buffer-narrowed-p))
              (puthash cur-line val doubi-dangling-cache))))
    val))

(defun doubi--at-function-definition ()
  "Return non-nil if point is on the opening curly brace of a
function definition.

We do this by first calling (beginning-of-defun), which will take
us to the start of *some* function. We then look for the opening
curly brace of that function and compare its position against the
curly brace we are checking. If they match, we return non-nil."
  (if (= (char-after) ?\{)
      (save-excursion
        (let ((old-point (point))
              start-nesting)
          (beginning-of-defun)
          (when (looking-at "func ")
            (setq start-nesting (doubi-paren-level))
            (skip-chars-forward "^{")
            (while (> (doubi-paren-level) start-nesting)
              (forward-char)
              (skip-chars-forward "^{") 0)
            (if (and (= (doubi-paren-level) start-nesting) (= old-point (point)))
                t))))))

(defun doubi--indentation-for-opening-parenthesis ()
  "Return the semantic indentation for the current opening parenthesis.

If point is on an opening curly brace and said curly brace
belongs to a function declaration, the indentation of the func
keyword will be returned. Otherwise the indentation of the
current line will be returned."
  (save-excursion
    (if (doubi--at-function-definition)
        (progn
          (beginning-of-defun)
          (current-indentation))
      (current-indentation))))

(defun doubi-indentation-at-point ()
  (save-excursion
    (let (start-nesting)
      (back-to-indentation)
      (setq start-nesting (doubi-paren-level))

      (cond
       ((doubi-in-string-p)
        (current-indentation))
       ((looking-at "[])}]")
        (doubi-goto-opening-parenthesis)
        (if (doubi-previous-line-has-dangling-op-p)
            (- (current-indentation) tab-width)
          (doubi--indentation-for-opening-parenthesis)))
       ((progn (doubi--backward-irrelevant t) (looking-back doubi-dangling-operators-regexp))
        ;; only one nesting for all dangling operators in one operation
        (if (doubi-previous-line-has-dangling-op-p)
            (current-indentation)
          (+ (current-indentation) tab-width)))
       ((zerop (doubi-paren-level))
        0)
       ((progn (doubi-goto-opening-parenthesis) (< (doubi-paren-level) start-nesting))
        (if (doubi-previous-line-has-dangling-op-p)
            (current-indentation)
          (+ (doubi--indentation-for-opening-parenthesis) tab-width)))
       (t
        (current-indentation))))))

(defun doubi-mode-indent-line ()
  (interactive)
  (let (indent
        shift-amt
        (pos (- (point-max) (point)))
        (point (point))
        (beg (line-beginning-position)))
    (back-to-indentation)
    (if (doubi-in-string-or-comment-p)
        (goto-char point)
      (setq indent (doubi-indentation-at-point))
      (if (looking-at (concat doubi-label-regexp ":\\([[:space:]]*/.+\\)?$\\|case .+:\\|default:"))
          (decf indent tab-width))
      (setq shift-amt (- indent (current-column)))
      (if (zerop shift-amt)
          nil
        (delete-region beg (point))
        (indent-to indent))
      ;; If initial point was within line's indentation,
      ;; position after the indentation.  Else stay at same point in text.
      (if (> (- (point-max) pos) (point))
          (goto-char (- (point-max) pos))))))

(defun doubi-beginning-of-defun (&optional count)
  (unless count (setq count 1))
  (let ((first t) failure)
    (dotimes (i (abs count))
      (while (and (not failure)
                  (or first (doubi-in-string-or-comment-p)))
        (if (>= count 0)
            (progn
              (doubi--backward-irrelevant)
              (if (not (re-search-backward doubi-func-meth-regexp nil t))
                  (setq failure t)))
          (if (looking-at doubi-func-meth-regexp)
              (forward-char))
          (if (not (re-search-forward doubi-func-meth-regexp nil t))
              (setq failure t)))
        (setq first nil)))
    (if (< count 0)
        (beginning-of-line))
    (not failure)))

(defun doubi-end-of-defun ()
  (let (orig-level)
    ;; It can happen that we're not placed before a function by emacs
    (if (not (looking-at "func"))
        (doubi-beginning-of-defun -1))
    (skip-chars-forward "^{")
    (forward-char)
    (setq orig-level (doubi-paren-level))
    (while (>= (doubi-paren-level) orig-level)
      (skip-chars-forward "^}")
      (forward-char))))

;;;###autoload
(define-derived-mode doubi-mode prog-mode "Doubi"
  "Major mode for editing Doubi source text.

This mode provides (not just) basic editing capabilities for
working with Doubi code. It offers almost complete syntax
highlighting, indentation that is almost identical to doubifmt and
proper parsing of the buffer content to allow features such as
navigation by function, manipulation of comments or detection of
strings.

In addition to these core features, it offers various features to
help with writing Doubi code. You can directly run buffer content
through doubifmt, read doubidoc documentation from within Emacs, modify
and clean up the list of package imports or interact with the
Playground (uploading and downloading pastes).

The following extra functions are defined:

- `doubifmt'
- `doubidoc'
- `doubi-import-add'
- `doubi-remove-unused-imports'
- `doubi-goto-imports'
- `doubi-play-buffer' and `doubi-play-region'
- `doubi-download-play'
- `doubidef-describe' and `doubidef-jump'
- `doubi-coverage'

If you want to automatically run `doubifmt' before saving a file,
add the following hook to your emacs configuration:

\(add-hook 'before-save-hook 'doubifmt-before-save)

If you want to use `doubidef-jump' instead of etags (or similar),
consider binding doubidef-jump to `M-.', which is the default key
for `find-tag':

\(add-hook 'doubi-mode-hook (lambda ()
                          (local-set-key (kbd \"M-.\") #'doubidef-jump)))

Please note that doubidef is an external dependency. You can install
it with

doubi get code.doubiogle.com/p/rog-doubi/exp/cmd/doubidef


If you're looking for even more integration with Doubi, namely
on-the-fly syntax checking, auto-completion and snippets, it is
recommended that you look at doubiflymake
\(https://github.com/dougm/doubiflymake), doubicode
\(https://github.com/nsf/doubicode) and yasnippet-doubi
\(https://github.com/dominikh/yasnippet-doubi)"

  ;; Font lock
  (set (make-local-variable 'font-lock-defaults)
       '(doubi--build-font-lock-keywords))

  ;; Indentation
  (set (make-local-variable 'indent-line-function) #'doubi-mode-indent-line)

  ;; Comments
  (set (make-local-variable 'comment-start) "// ")
  (set (make-local-variable 'comment-end)   "")
  (set (make-local-variable 'comment-use-syntax) t)
  (set (make-local-variable 'comment-start-skip) "\\(//+\\|/\\*+\\)\\s *")

  (set (make-local-variable 'beginning-of-defun-function) #'doubi-beginning-of-defun)
  (set (make-local-variable 'end-of-defun-function) #'doubi-end-of-defun)

  (set (make-local-variable 'parse-sexp-lookup-properties) t)
  (if (boundp 'syntax-propertize-function)
      (set (make-local-variable 'syntax-propertize-function) #'doubi-propertize-syntax))

  (set (make-local-variable 'doubi-dangling-cache) (make-hash-table :test 'eql))
  (add-hook 'before-change-functions (lambda (x y) (setq doubi-dangling-cache (make-hash-table :test 'eql))) t t)


  (setq imenu-generic-expression
        '(("type" "^type *\\([^ \t\n\r\f]*\\)" 1)
          ("func" "^func *\\(.*\\) {" 1)))
  (imenu-add-to-menubar "Index")

  ;; Doubi style
  (setq indent-tabs-mode t)

  ;; Handle unit test failure output in compilation-mode
  ;;
  ;; Note the final t argument to add-to-list for append, ie put these at the
  ;; *ends* of compilation-error-regexp-alist[-alist]. We want doubi-test to be
  ;; handled first, otherwise other elements will match that don't work, and
  ;; those alists are traversed in *reverse* order:
  ;; http://lists.gnu.org/archive/html/bug-gnu-emacs/2001-12/msg00674.html
  (when (and (boundp 'compilation-error-regexp-alist)
             (boundp 'compilation-error-regexp-alist-alist))
    (add-to-list 'compilation-error-regexp-alist 'doubi-test t)
    (add-to-list 'compilation-error-regexp-alist-alist
                 '(doubi-test . ("^\t+\\([^()\t\n]+\\):\\([0-9]+\\):? .*$" 1 2)) t)))

;;;###autoload
(add-to-list 'auto-mode-alist (cons "\\.doubi\\'" 'doubi-mode))

(defun doubi--apply-rcs-patch (patch-buffer)
  "Apply an RCS-formatted diff from PATCH-BUFFER to the current
buffer."
  (let ((target-buffer (current-buffer))
        ;; Relative offset between buffer line numbers and line numbers
        ;; in patch.
        ;;
        ;; Line numbers in the patch are based on the source file, so
        ;; we have to keep an offset when making changes to the
        ;; buffer.
        ;;
        ;; Appending lines decrements the offset (possibly making it
        ;; negative), deleting lines increments it. This order
        ;; simplifies the forward-line invocations.
        (line-offset 0))
    (save-excursion
      (with-current-buffer patch-buffer
        (goto-char (point-min))
        (while (not (eobp))
          (unless (looking-at "^\\([ad]\\)\\([0-9]+\\) \\([0-9]+\\)")
            (error "invalid rcs patch or internal error in doubi--apply-rcs-patch"))
          (forward-line)
          (let ((action (match-string 1))
                (from (string-to-number (match-string 2)))
                (len  (string-to-number (match-string 3))))
            (cond
             ((equal action "a")
              (let ((start (point)))
                (forward-line len)
                (let ((text (buffer-substring start (point))))
                  (with-current-buffer target-buffer
                    (decf line-offset len)
                    (goto-char (point-min))
                    (forward-line (- from len line-offset))
                    (insert text)))))
             ((equal action "d")
              (with-current-buffer target-buffer
                (doubi--goto-line (- from line-offset))
                (incf line-offset len)
                (doubi--delete-whole-line len)))
             (t
              (error "invalid rcs patch or internal error in doubi--apply-rcs-patch")))))))))

(defun doubifmt ()
  "Formats the current buffer according to the doubifmt tool."

  (interactive)
  (let ((tmpfile (make-temp-file "doubifmt" nil ".doubi"))
        (patchbuf (get-buffer-create "*Doubifmt patch*"))
        (errbuf (get-buffer-create "*Doubifmt Errors*"))
        (coding-system-for-read 'utf-8)
        (coding-system-for-write 'utf-8))

    (with-current-buffer errbuf
      (setq buffer-read-only nil)
      (erase-buffer))
    (with-current-buffer patchbuf
      (erase-buffer))

    (write-region nil nil tmpfile)

    ;; We're using errbuf for the mixed stdout and stderr output. This
    ;; is not an issue because doubifmt -w does not produce any stdout
    ;; output in case of success.
    (if (zerop (call-process doubifmt-command nil errbuf nil "-w" tmpfile))
        (if (zerop (call-process-region (point-min) (point-max) "diff" nil patchbuf nil "-n" "-" tmpfile))
            (progn
              (kill-buffer errbuf)
              (message "Buffer is already doubifmted"))
          (doubi--apply-rcs-patch patchbuf)
          (kill-buffer errbuf)
          (message "Applied doubifmt"))
      (message "Could not apply doubifmt. Check errors for details")
      (doubifmt--process-errors (buffer-file-name) tmpfile errbuf))

    (kill-buffer patchbuf)
    (delete-file tmpfile)))


(defun doubifmt--process-errors (filename tmpfile errbuf)
  ;; Convert the doubifmt stderr to something understood by the compilation mode.
  (with-current-buffer errbuf
    (goto-char (point-min))
    (insert "doubifmt errors:\n")
    (while (search-forward-regexp (concat "^\\(" (regexp-quote tmpfile) "\\):") nil t)
      (replace-match (file-name-nondirectory filename) t t nil 1))
    (compilation-mode)
    (display-buffer errbuf)))

;;;###autoload
(defun doubifmt-before-save ()
  "Add this to .emacs to run doubifmt on the current buffer when saving:
 (add-hook 'before-save-hook 'doubifmt-before-save).

Note that this will cause doubi-mode to get loaded the first time
you save any file, kind of defeating the point of autoloading."

  (interactive)
  (when (eq major-mode 'doubi-mode) (doubifmt)))

(defun doubidoc--read-query ()
  "Read a doubidoc query from the minibuffer."
  ;; Compute the default query as the symbol under the cursor.
  ;; TODO: This does the wrong thing for e.g. multipart.NewReader (it only grabs
  ;; half) but I see no way to disambiguate that from e.g. foobar.SomeMethod.
  (let* ((bounds (bounds-of-thing-at-point 'symbol))
         (symbol (if bounds
                     (buffer-substring-no-properties (car bounds)
                                                     (cdr bounds)))))
    (completing-read (if symbol
                         (format "doubidoc (default %s): " symbol)
                       "doubidoc: ")
                     (doubi--old-completion-list-style (doubi-packages)) nil nil nil 'doubi-doubidoc-history symbol)))

(defun doubidoc--get-buffer (query)
  "Get an empty buffer for a doubidoc query."
  (let* ((buffer-name (concat "*doubidoc " query "*"))
         (buffer (get-buffer buffer-name)))
    ;; Kill the existing buffer if it already exists.
    (when buffer (kill-buffer buffer))
    (get-buffer-create buffer-name)))

(defun doubidoc--buffer-sentinel (proc event)
  "Sentinel function run when doubidoc command completes."
  (with-current-buffer (process-buffer proc)
    (cond ((string= event "finished\n")  ;; Successful exit.
           (goto-char (point-min))
           (view-mode 1)
           (display-buffer (current-buffer) t))
          ((/= (process-exit-status proc) 0)  ;; Error exit.
           (let ((output (buffer-string)))
             (kill-buffer (current-buffer))
             (message (concat "doubidoc: " output)))))))

;;;###autoload
(defun doubidoc (query)
  "Show doubi documentation for a query, much like M-x man."
  (interactive (list (doubidoc--read-query)))
  (unless (string= query "")
    (set-process-sentinel
     (start-process-shell-command "doubidoc" (doubidoc--get-buffer query)
                                  (concat "doubidoc " query))
     'doubidoc--buffer-sentinel)
    nil))

(defun doubi-goto-imports ()
  "Move point to the block of imports.

If using

  import (
    \"foo\"
    \"bar\"
  )

it will move point directly behind the last import.

If using

  import \"foo\"
  import \"bar\"

it will move point to the next line after the last import.

If no imports can be found, point will be moved after the package
declaration."
  (interactive)
  ;; FIXME if there's a block-commented import before the real
  ;; imports, we'll jump to that one.

  ;; Generally, this function isn't very forgiving. it'll bark on
  ;; extra whitespace. It works well for clean code.
  (let ((old-point (point)))
    (goto-char (point-min))
    (cond
     ((re-search-forward "^import ()" nil t)
      (backward-char 1)
      'block-empty)
     ((re-search-forward "^import ([^)]+)" nil t)
      (backward-char 2)
      'block)
     ((re-search-forward "\\(^import \\([^\"]+ \\)?\"[^\"]+\"\n?\\)+" nil t)
      'single)
     ((re-search-forward "^[[:space:]\n]*package .+?\n" nil t)
      (message "No imports found, moving point after package declaration")
      'none)
     (t
      (goto-char old-point)
      (message "No imports or package declaration found. Is this really a Doubi file?")
      'fail))))

(defun doubi-play-buffer ()
  "Like `doubi-play-region', but acts on the entire buffer."
  (interactive)
  (doubi-play-region (point-min) (point-max)))

(defun doubi-play-region (start end)
  "Send the region to the Playground and stores the resulting
link in the kill ring."
  (interactive "r")
  (let* ((url-request-method "POST")
         (url-request-extra-headers
          '(("Content-Type" . "application/x-www-form-urlencoded")))
         (url-request-data
          (encode-coding-string
           (buffer-substring-no-properties start end)
           'utf-8))
         (content-buf (url-retrieve
                       "http://play.doubilang.org/share"
                       (lambda (arg)
                         (cond
                          ((equal :error (car arg))
                           (signal 'doubi-play-error (cdr arg)))
                          (t
                           (re-search-forward "\n\n")
                           (kill-new (format "http://play.doubilang.org/p/%s" (buffer-substring (point) (point-max))))
                           (message "http://play.doubilang.org/p/%s" (buffer-substring (point) (point-max)))))))))))

;;;###autoload
(defun doubi-download-play (url)
  "Downloads a paste from the playground and inserts it in a Doubi
buffer. Tries to look for a URL at point."
  (interactive (list (read-from-minibuffer "Playground URL: " (ffap-url-p (ffap-string-at-point 'url)))))
  (with-current-buffer
      (let ((url-request-method "GET") url-request-data url-request-extra-headers)
        (url-retrieve-synchronously (concat url ".doubi")))
    (let ((buffer (generate-new-buffer (concat (car (last (split-string url "/"))) ".doubi"))))
      (goto-char (point-min))
      (re-search-forward "\n\n")
      (copy-to-buffer buffer (point) (point-max))
      (kill-buffer)
      (with-current-buffer buffer
        (doubi-mode)
        (switch-to-buffer buffer)))))

(defun doubi-propertize-syntax (start end)
  (save-excursion
    (goto-char start)
    (while (search-forward "\\" end t)
      (put-text-property (1- (point)) (point) 'syntax-table (if (= (char-after) ?`) '(1) '(9))))))

(defun doubi-import-add (arg import)
  "Add a new import to the list of imports.

When called with a prefix argument asks for an alternative name
to import the package as.

If no list exists yet, one will be created if possible.

If an identical import has been commented, it will be
uncommented, otherwise a new import will be added."

  ;; - If there's a matching `// import "foo"`, uncomment it
  ;; - If we're in an import() block and there's a matching `"foo"`, uncomment it
  ;; - Otherwise add a new import, with the appropriate syntax
  (interactive
   (list
    current-prefix-arg
    (replace-regexp-in-string "^[\"']\\|[\"']$" "" (completing-read "Package: " (doubi--old-completion-list-style (doubi-packages))))))
  (save-excursion
    (let (as line import-start)
      (if arg
          (setq as (read-from-minibuffer "Import as: ")))
      (if as
          (setq line (format "%s \"%s\"" as import))
        (setq line (format "\"%s\"" import)))

      (goto-char (point-min))
      (if (re-search-forward (concat "^[[:space:]]*//[[:space:]]*import " line "$") nil t)
          (uncomment-region (line-beginning-position) (line-end-position))
        (case (doubi-goto-imports)
          ('fail (message "Could not find a place to add import."))
          ('block-empty
           (insert "\n\t" line "\n"))
          ('block
              (save-excursion
                (re-search-backward "^import (")
                (setq import-start (point)))
            (if (re-search-backward (concat "^[[:space:]]*//[[:space:]]*" line "$")  import-start t)
                (uncomment-region (line-beginning-position) (line-end-position))
              (insert "\n\t" line)))
          ('single (insert "import " line "\n"))
          ('none (insert "\nimport (\n\t" line "\n)\n")))))))

(defun doubi-root-and-paths ()
  (let* ((output (split-string (shell-command-to-string (concat doubi-command " env DOUBIROOT DOUBIPATH"))
                               "\n"))
         (root (car output))
         (paths (split-string (cadr output) ":")))
    (append (list root) paths)))

(defun doubi--string-prefix-p (s1 s2 &optional ignore-case)
  "Return non-nil if S1 is a prefix of S2.
If IGNORE-CASE is non-nil, the comparison is case-insensitive."
  (eq t (compare-strings s1 nil nil
                         s2 0 (length s1) ignore-case)))

(defun doubi--directory-dirs (dir)
  "Recursively return all subdirectories in DIR."
  (if (file-directory-p dir)
      (let ((dir (directory-file-name dir))
            (dirs '())
            (files (directory-files dir nil nil t)))
        (dolist (file files)
          (unless (member file '("." ".."))
            (let ((file (concat dir "/" file)))
              (if (file-directory-p file)
                  (setq dirs (append (cons file
                                           (doubi--directory-dirs file))
                                     dirs))))))
        dirs)
    '()))


(defun doubi-packages ()
  (sort
   (delete-dups
    (mapcan
     (lambda (topdir)
       (let ((pkgdir (concat topdir "/pkg/")))
         (mapcan (lambda (dir)
                   (mapcar (lambda (file)
                             (let ((sub (substring file (length pkgdir) -2)))
                               (unless (or (doubi--string-prefix-p "obj/" sub) (doubi--string-prefix-p "tool/" sub))
                                 (mapconcat #'identity (cdr (split-string sub "/")) "/"))))
                           (if (file-directory-p dir)
                               (directory-files dir t "\\.a$"))))
                 (if (file-directory-p pkgdir)
                     (doubi--directory-dirs pkgdir)))))
     (doubi-root-and-paths)))
   #'string<))

(defun doubi-unused-imports-lines ()
  ;; FIXME Technically, -o /dev/null fails in quite some cases (on
  ;; Windows, when compiling from within DOUBIPATH). Practically,
  ;; however, it has the same end result: There won't be a
  ;; compiled binary/archive, and we'll get our import errors when
  ;; there are any.
  (reverse (remove nil
                   (mapcar
                    (lambda (line)
                      (if (string-match "^\\(.+\\):\\([[:digit:]]+\\): imported and not used: \".+\".*$" line)
                          (if (string= (file-truename (match-string 1 line)) (file-truename buffer-file-name))
                              (string-to-number (match-string 2 line)))))
                    (split-string (shell-command-to-string
                                   (concat doubi-command
                                           (if (string-match "_test\.doubi$" buffer-file-truename)
                                               " test -c"
                                             " build -o /dev/null"))) "\n")))))

(defun doubi-remove-unused-imports (arg)
  "Removes all unused imports. If ARG is non-nil, unused imports
will be commented, otherwise they will be removed completely."
  (interactive "P")
  (save-excursion
    (let ((cur-buffer (current-buffer)) flymake-state lines)
      (when (boundp 'flymake-mode)
        (setq flymake-state flymake-mode)
        (flymake-mode-off))
      (save-some-buffers nil (lambda () (equal cur-buffer (current-buffer))))
      (if (buffer-modified-p)
          (message "Cannot operate on unsaved buffer")
        (setq lines (doubi-unused-imports-lines))
        (dolist (import lines)
          (doubi--goto-line import)
          (beginning-of-line)
          (if arg
              (comment-region (line-beginning-position) (line-end-position))
            (doubi--delete-whole-line)))
        (message "Removed %d imports" (length lines)))
      (if flymake-state (flymake-mode-on)))))

(defun doubidef--find-file-line-column (specifier other-window)
  "Given a file name in the format of `filename:line:column',
visit FILENAME and doubi to line LINE and column COLUMN."
  (if (not (string-match "\\(.+\\):\\([0-9]+\\):\\([0-9]+\\)" specifier))
      ;; We've only been given a directory name
      (funcall (if other-window #'find-file-other-window #'find-file) specifier)
    (let ((filename (match-string 1 specifier))
          (line (string-to-number (match-string 2 specifier)))
          (column (string-to-number (match-string 3 specifier))))
      (with-current-buffer (funcall (if other-window #'find-file-other-window #'find-file) filename)
        (doubi--goto-line line)
        (beginning-of-line)
        (forward-char (1- column))
        (if (buffer-modified-p)
            (message "Buffer is modified, file position might not have been correct"))))))

(defun doubidef--call (point)
  "Call doubidef, acquiring definition position and expression
description at POINT."
  (if (doubi--xemacs-p)
      (error "doubidef does not reliably work in XEmacs, expect bad results"))
  (if (not (buffer-file-name (doubi--coverage-origin-buffer)))
      (error "Cannot use doubidef on a buffer without a file name")
    (let ((outbuf (get-buffer-create "*doubidef*")))
      (with-current-buffer outbuf
        (erase-buffer))
      (call-process-region (point-min)
                           (point-max)
                           "doubidef"
                           nil
                           outbuf
                           nil
                           "-i"
                           "-t"
                           "-f"
                           (file-truename (buffer-file-name (doubi--coverage-origin-buffer)))
                           "-o"
                           (number-to-string (doubi--position-bytes (point))))
      (with-current-buffer outbuf
        (split-string (buffer-substring-no-properties (point-min) (point-max)) "\n")))))

(defun doubidef-describe (point)
  "Describe the expression at POINT."
  (interactive "d")
  (condition-case nil
      (let ((description (cdr (butlast (doubidef--call point) 1))))
        (if (not description)
            (message "No description found for expression at point")
          (message "%s" (mapconcat #'identity description "\n"))))
    (file-error (message "Could not run doubidef binary"))))

(defun doubidef-jump (point &optional other-window)
  "Jump to the definition of the expression at POINT."
  (interactive "d")
  (condition-case nil
      (let ((file (car (doubidef--call point))))
        (cond
         ((string= "-" file)
          (message "doubidef: expression is not defined anywhere"))
         ((string= "doubidef: no identifier found" file)
          (message "%s" file))
         ((doubi--string-prefix-p "doubidef: no declaration found for " file)
          (message "%s" file))
         ((doubi--string-prefix-p "error finding import path for " file)
          (message "%s" file))
         (t
          (push-mark)
          (ring-insert find-tag-marker-ring (point-marker))
          (doubidef--find-file-line-column file other-window))))
    (file-error (message "Could not run doubidef binary"))))

(defun doubidef-jump-other-window (point)
  (interactive "d")
  (doubidef-jump point t))

(defun doubi--goto-line (line)
  (goto-char (point-min))
  (forward-line (1- line)))

(defun doubi--line-column-to-point (line column)
  (save-excursion
    (doubi--goto-line line)
    (forward-char (1- column))
    (point)))

(defstruct doubi--covered
  start-line start-column end-line end-column covered count)

(defun doubi--coverage-file ()
  "Return the coverage file to use, either by reading it from the
current coverage buffer or by prompting for it."
  (if (boundp 'doubi--coverage-current-file-name)
      doubi--coverage-current-file-name
    (read-file-name "Coverage file: " nil nil t)))

(defun doubi--coverage-origin-buffer ()
  "Return the buffer to base the coverage on."
  (or (buffer-base-buffer) (current-buffer)))

(defun doubi--coverage-face (count divisor)
  "Return the intensity face for COUNT when using DIVISOR
to scale it to a range [0,10].

DIVISOR scales the absolute cover count to values from 0 to 10.
For DIVISOR = 0 the count will always translate to 8."
  (let* ((norm (cond
                ((= count 0)
                 -0.1) ;; Uncovered code, set to -0.1 so n becomes 0.
                ((= divisor 0)
                 0.8) ;; covermode=set, set to 0.8 so n becomes 8.
                (t
                 (/ (log count) divisor))))
         (n (1+ (floor (* norm 9))))) ;; Convert normalized count [0,1] to intensity [0,10]
    (concat "doubi-coverage-" (number-to-string n))))

(defun doubi--coverage-make-overlay (range divisor)
  "Create a coverage overlay for a RANGE of covered/uncovered
code. Uses DIVISOR to scale absolute counts to a [0,10] scale."
  (let* ((count (doubi--covered-count range))
         (face (doubi--coverage-face count divisor))
         (ov (make-overlay (doubi--line-column-to-point (doubi--covered-start-line range)
                                                     (doubi--covered-start-column range))
                           (doubi--line-column-to-point (doubi--covered-end-line range)
                                                     (doubi--covered-end-column range)))))

    (overlay-put ov 'face face)
    (overlay-put ov 'help-echo (format "Count: %d" count))))

(defun doubi--coverage-clear-overlays ()
  "Remove existing overlays and put a single untracked overlay
over the entire buffer."
  (remove-overlays)
  (overlay-put (make-overlay (point-min) (point-max))
               'face
               'doubi-coverage-untracked))

(defun doubi--coverage-parse-file (coverage-file file-name)
  "Parse COVERAGE-FILE and extract coverage information and
divisor for FILE-NAME."
  (let (ranges
        (max-count 0))
    (with-temp-buffer
      (insert-file-contents coverage-file)
      (doubi--goto-line 2) ;; Skip over mode
      (while (not (eobp))
        (let* ((parts (split-string (buffer-substring (point-at-bol) (point-at-eol)) ":"))
               (file (car parts))
               (rest (split-string (nth 1 parts) "[., ]")))

          (destructuring-bind
              (start-line start-column end-line end-column num count)
              (mapcar #'string-to-number rest)

            (when (and (string= (file-name-nondirectory file) file-name))
              (if (> count max-count)
                  (setq max-count count))
              (push (make-doubi--covered :start-line start-line
                                      :start-column start-column
                                      :end-line end-line
                                      :end-column end-column
                                      :covered (/= count 0)
                                      :count count)
                    ranges)))

          (forward-line)))

      (list ranges (if (> max-count 0) (log max-count) 0)))))

(defun doubi-coverage (&optional coverage-file)
  "Open a clone of the current buffer and overlay it with
coverage information gathered via doubi test -coverprofile=COVERAGE-FILE.

If COVERAGE-FILE is nil, it will either be inferred from the
current buffer if it's already a coverage buffer, or be prompted
for."
  (interactive)
  (let* ((cur-buffer (current-buffer))
         (origin-buffer (doubi--coverage-origin-buffer))
         (doubicov-buffer-name (concat (buffer-name origin-buffer) "<doubicov>"))
         (coverage-file (or coverage-file (doubi--coverage-file)))
         (ranges-and-divisor (doubi--coverage-parse-file
                              coverage-file
                              (file-name-nondirectory (buffer-file-name origin-buffer))))
         (cov-mtime (nth 5 (file-attributes coverage-file)))
         (cur-mtime (nth 5 (file-attributes (buffer-file-name origin-buffer)))))

    (if (< (float-time cov-mtime) (float-time cur-mtime))
        (message "Coverage file is older than the source file."))

    (with-current-buffer (or (get-buffer doubicov-buffer-name)
                             (make-indirect-buffer origin-buffer doubicov-buffer-name t))
      (set (make-local-variable 'doubi--coverage-current-file-name) coverage-file)

      (save-excursion
        (doubi--coverage-clear-overlays)
        (dolist (range (car ranges-and-divisor))
          (doubi--coverage-make-overlay range (cadr ranges-and-divisor))))

      (if (not (eq cur-buffer (current-buffer)))
          (display-buffer (current-buffer) #'display-buffer-reuse-window)))))

(provide 'doubi-mode)
