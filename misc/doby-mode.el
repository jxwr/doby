;;; doby-mode.el --- Major mode for the Doby programming language

;; Copyright 2013 The Doby Authors. All rights reserved.
;; Use of this source code is dobyverned by a BSD-style
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
;;     - dobydef will not work correctly if multibyte characters are
;;       being used
;;     - Fontification will not handle unicode correctly
;;
;; - Do not use \_< and \_> regexp delimiters directly; use
;;   doby--regexp-enclose-in-symbol
;;
;; - The character `_` must not be a symbol constituent but a
;;   character constituent
;;
;; - Do not use process-lines
;;
;; - Use doby--old-completion-list-style when using a plain list as the
;;   collection for completing-read
;;
;; - Use doby--kill-whole-line instead of kill-whole-line (called
;;   kill-entire-line in XEmacs)
;;
;; - Use doby--position-bytes instead of position-bytes
(defmacro doby--xemacs-p ()
  `(featurep 'xemacs))

(defalias 'doby--kill-whole-line
  (if (fboundp 'kill-whole-line)
      #'kill-whole-line
    #'kill-entire-line))

;; Delete the current line without putting it in the kill-ring.
(defun doby--delete-whole-line (&optional arg)
  ;; Emacs uses both kill-region and kill-new, Xemacs only uses
  ;; kill-region. In both cases we turn them into operations that do
  ;; not modify the kill ring. This solution does depend on the
  ;; implementation of kill-line, but it's the only viable solution
  ;; that does not require to write kill-line from scratch.
  (flet ((kill-region (beg end)
                      (delete-region beg end))
         (kill-new (s) ()))
    (doby--kill-whole-line arg)))

;; declare-function is an empty macro that only byte-compile cares
;; about. Wrap in always false if to satisfy Emacsen without that
;; macro.
(if nil
    (declare-function doby--position-bytes "doby-mode" (point)))
;; XEmacs unfortunately does not offer position-bytes. We can fall
;; back to just using (point), but it will be incorrect as soon as
;; multibyte characters are being used.
(if (fboundp 'position-bytes)
    (defalias 'doby--position-bytes #'position-bytes)
  (defun doby--position-bytes (point) point))

(defun doby--old-completion-list-style (list)
  (mapcar (lambda (x) (cons x nil)) list))

;; GNU Emacs 24 has prog-mode, older GNU Emacs and XEmacs do not, so
;; copy its definition for those.
(if (not (fboundp 'prog-mode))
    (define-derived-mode prog-mode fundamental-mode "Prog"
      "Major mode for editing source code."
      (set (make-local-variable 'require-final-newline) mode-require-final-newline)
      (set (make-local-variable 'parse-sexp-ignore-comments) t)
      (setq bidi-paragraph-direction 'left-to-right)))

(defun doby--regexp-enclose-in-symbol (s)
  ;; XEmacs does not support \_<, GNU Emacs does. In GNU Emacs we make
  ;; extensive use of \_< to support unicode in identifiers. Until we
  ;; come up with a better solution for XEmacs, this solution will
  ;; break fontification in XEmacs for identifiers such as "typeµ".
  ;; XEmacs will consider "type" a keyword, GNU Emacs won't.

  (if (doby--xemacs-p)
      (concat "\\<" s "\\>")
    (concat "\\_<" s "\\_>")))

;; Move up one level of parentheses.
(defun doby-goto-opening-parenthesis (&optional legacy-unused)
  ;; The old implementation of doby-goto-opening-parenthesis had an
  ;; optional argument to speed up the function. It didn't change the
  ;; function's outcome.

  ;; Silently fail if there's no matching opening parenthesis.
  (condition-case nil
      (backward-up-list)
    (scan-error nil)))


(defconst doby-dangling-operators-regexp "[^-]-\\|[^+]\\+\\|[/*&><.=|^]")
(defconst doby-identifier-regexp "[[:word:][:multibyte:]]+")
(defconst doby-label-regexp doby-identifier-regexp)
(defconst doby-type-regexp "[[:word:][:multibyte:]*]+")
(defconst doby-func-regexp (concat (doby--regexp-enclose-in-symbol "func") "\\s *\\(" doby-identifier-regexp "\\)"))
(defconst doby-func-meth-regexp (concat
                               (doby--regexp-enclose-in-symbol "func") "\\s *\\(?:(\\s *"
                               "\\(" doby-identifier-regexp "\\s +\\)?" doby-type-regexp
                               "\\s *)\\s *\\)?\\("
                               doby-identifier-regexp
                               "\\)("))
(defconst doby-builtins
  '("append" "cap"   "close"   "complex" "copy"
    "delete" "imag"  "len"     "make"    "new"
    "panic"  "print" "println" "real"    "recover")
  "All built-in functions in the Doby language. Used for font locking.")

(defconst doby-mode-keywords
  '("break"    "default"     "func"   "interface" "select"
    "case"     "defer"       "doby"     "map"       "struct"
    "chan"     "else"        "goto"   "package"   "switch"
    "const"    "fallthrough" "if"     "range"     "type"
    "continue" "for"         "import" "return"    "var")
  "All keywords in the Doby language.  Used for font locking.")

(defconst doby-constants '("nil" "true" "false" "iota"))
(defconst doby-type-name-regexp (concat "\\(?:[*(]\\)*\\(?:" doby-identifier-regexp "\\.\\)?\\(" doby-identifier-regexp "\\)"))

(defvar doby-dangling-cache)
(defvar doby-dobydoc-history nil)
(defvar doby--coverage-current-file-name)

(defgroup doby nil
  "Major mode for editing Doby code"
  :group 'languages)

(defgroup doby-cover nil
  "Options specific to `cover`"
  :group 'doby)

(defcustom doby-fontify-function-calls t
  "Fontify function and method calls if this is non-nil."
  :type 'boolean
  :group 'doby)

(defcustom doby-mode-hook nil
  "Hook called by `doby-mode'."
  :type 'hook
  :group 'doby)

(defcustom doby-command "doby"
  "The 'doby' command.  Some users have multiple Doby development
trees and invoke the 'doby' tool via a wrapper that sets DOBYROOT and
DOBYPATH based on the current directory.  Such users should
customize this variable to point to the wrapper script."
  :type 'string
  :group 'doby)

(defcustom dobyfmt-command "dobyfmt"
  "The 'dobyfmt' command.  Some users may replace this with 'dobyimports'
from https://github.com/bradfitz/dobyimports."
  :type 'string
  :group 'doby)

(defface doby-coverage-untracked
  '((t (:foreground "#505050")))
  "Coverage color of untracked code."
  :group 'doby-cover)

(defface doby-coverage-0
  '((t (:foreground "#c00000")))
  "Coverage color for uncovered code."
  :group 'doby-cover)
(defface doby-coverage-1
  '((t (:foreground "#808080")))
  "Coverage color for covered code with weight 1."
  :group 'doby-cover)
(defface doby-coverage-2
  '((t (:foreground "#748c83")))
  "Coverage color for covered code with weight 2."
  :group 'doby-cover)
(defface doby-coverage-3
  '((t (:foreground "#689886")))
  "Coverage color for covered code with weight 3."
  :group 'doby-cover)
(defface doby-coverage-4
  '((t (:foreground "#5ca489")))
  "Coverage color for covered code with weight 4."
  :group 'doby-cover)
(defface doby-coverage-5
  '((t (:foreground "#50b08c")))
  "Coverage color for covered code with weight 5."
  :group 'doby-cover)
(defface doby-coverage-6
  '((t (:foreground "#44bc8f")))
  "Coverage color for covered code with weight 6."
  :group 'doby-cover)
(defface doby-coverage-7
  '((t (:foreground "#38c892")))
  "Coverage color for covered code with weight 7."
  :group 'doby-cover)
(defface doby-coverage-8
  '((t (:foreground "#2cd495")))
  "Coverage color for covered code with weight 8.
For mode=set, all covered lines will have this weight."
  :group 'doby-cover)
(defface doby-coverage-9
  '((t (:foreground "#20e098")))
  "Coverage color for covered code with weight 9."
  :group 'doby-cover)
(defface doby-coverage-10
  '((t (:foreground "#14ec9b")))
  "Coverage color for covered code with weight 10."
  :group 'doby-cover)
(defface doby-coverage-covered
  '((t (:foreground "#2cd495")))
  "Coverage color of covered code."
  :group 'doby-cover)

(defvar doby-mode-syntax-table
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
    (modify-syntax-entry ?/ (if (doby--xemacs-p) ". 1456" ". 124b") st)
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
  "Syntax table for Doby mode.")

(defun doby--build-font-lock-keywords ()
  ;; we cannot use 'symbols in regexp-opt because emacs <24 doesn't
  ;; understand that
  (append
   `((,(doby--regexp-enclose-in-symbol (regexp-opt doby-mode-keywords t)) . font-lock-keyword-face)
     (,(doby--regexp-enclose-in-symbol (regexp-opt doby-builtins t)) . font-lock-builtin-face)
     (,(doby--regexp-enclose-in-symbol (regexp-opt doby-constants t)) . font-lock-constant-face)
     (,doby-func-regexp 1 font-lock-function-name-face)) ;; function (not method) name

   (if doby-fontify-function-calls
       `((,(concat "\\(" doby-identifier-regexp "\\)[[:space:]]*(") 1 font-lock-function-name-face) ;; function call/method name
         (,(concat "[^[:word:][:multibyte:]](\\(" doby-identifier-regexp "\\))[[:space:]]*(") 1 font-lock-function-name-face)) ;; bracketed function call
     `((,doby-func-meth-regexp 1 font-lock-function-name-face))) ;; method name

   `(
     (,(concat (doby--regexp-enclose-in-symbol "type") "[[:space:]]*\\([^[:space:]]+\\)") 1 font-lock-type-face) ;; types
     (,(concat (doby--regexp-enclose-in-symbol "type") "[[:space:]]*" doby-identifier-regexp "[[:space:]]*" doby-type-name-regexp) 1 font-lock-type-face) ;; types
     (,(concat "[^[:word:][:multibyte:]]\\[\\([[:digit:]]+\\|\\.\\.\\.\\)?\\]" doby-type-name-regexp) 2 font-lock-type-face) ;; Arrays/slices
     (,(concat "\\(" doby-identifier-regexp "\\)" "{") 1 font-lock-type-face)
     (,(concat (doby--regexp-enclose-in-symbol "map") "\\[[^]]+\\]" doby-type-name-regexp) 1 font-lock-type-face) ;; map value type
     (,(concat (doby--regexp-enclose-in-symbol "map") "\\[" doby-type-name-regexp) 1 font-lock-type-face) ;; map key type
     (,(concat (doby--regexp-enclose-in-symbol "chan") "[[:space:]]*\\(?:<-\\)?" doby-type-name-regexp) 1 font-lock-type-face) ;; channel type
     (,(concat (doby--regexp-enclose-in-symbol "\\(?:new\\|make\\)") "\\(?:[[:space:]]\\|)\\)*(" doby-type-name-regexp) 1 font-lock-type-face) ;; new/make type
     ;; TODO do we actually need this one or isn't it just a function call?
     (,(concat "\\.\\s *(" doby-type-name-regexp) 1 font-lock-type-face) ;; Type conversion
     (,(concat (doby--regexp-enclose-in-symbol "func") "[[:space:]]+(" doby-identifier-regexp "[[:space:]]+" doby-type-name-regexp ")") 1 font-lock-type-face) ;; Method receiver
     (,(concat (doby--regexp-enclose-in-symbol "func") "[[:space:]]+(" doby-type-name-regexp ")") 1 font-lock-type-face) ;; Method receiver without variable name
     ;; Like the original doby-mode this also marks compound literal
     ;; fields. There, it was marked as to fix, but I grew quite
     ;; accustomed to it, so it'll stay for now.
     (,(concat "^[[:space:]]*\\(" doby-label-regexp "\\)[[:space:]]*:\\(\\S.\\|$\\)") 1 font-lock-constant-face) ;; Labels and compound literal fields
     (,(concat (doby--regexp-enclose-in-symbol "\\(goto\\|break\\|continue\\)") "[[:space:]]*\\(" doby-label-regexp "\\)") 2 font-lock-constant-face)))) ;; labels in goto/break/continue

(defvar doby-mode-map
  (let ((m (make-sparse-keymap)))
    (define-key m "}" #'doby-mode-insert-and-indent)
    (define-key m ")" #'doby-mode-insert-and-indent)
    (define-key m "," #'doby-mode-insert-and-indent)
    (define-key m ":" #'doby-mode-insert-and-indent)
    (define-key m "=" #'doby-mode-insert-and-indent)
    (define-key m (kbd "C-c C-a") #'doby-import-add)
    (define-key m (kbd "C-c C-j") #'dobydef-jump)
    (define-key m (kbd "C-x 4 C-c C-j") #'dobydef-jump-other-window)
    (define-key m (kbd "C-c C-d") #'dobydef-describe)
    m)
  "Keymap used by Doby mode to implement electric keys.")

(defun doby-mode-insert-and-indent (key)
  "Invoke the global binding of KEY, then reindent the line."

  (interactive (list (this-command-keys)))
  (call-interactively (lookup-key (current-global-map) key))
  (indent-according-to-mode))

(defmacro doby-paren-level ()
  `(car (syntax-ppss)))

(defmacro doby-in-string-or-comment-p ()
  `(nth 8 (syntax-ppss)))

(defmacro doby-in-string-p ()
  `(nth 3 (syntax-ppss)))

(defmacro doby-in-comment-p ()
  `(nth 4 (syntax-ppss)))

(defmacro doby-goto-beginning-of-string-or-comment ()
  `(goto-char (nth 8 (syntax-ppss))))

(defun doby--backward-irrelevant (&optional stop-at-string)
  "Skips backwards over any characters that are irrelevant for
indentation and related tasks.

It skips over whitespace, comments, cases and labels and, if
STOP-AT-STRING is not true, over strings."

  (let (pos (start-pos (point)))
    (skip-chars-backward "\n\s\t")
    (if (and (save-excursion (beginning-of-line) (doby-in-string-p)) (looking-back "`") (not stop-at-string))
        (backward-char))
    (if (and (doby-in-string-p) (not stop-at-string))
        (doby-goto-beginning-of-string-or-comment))
    (if (looking-back "\\*/")
        (backward-char))
    (if (doby-in-comment-p)
        (doby-goto-beginning-of-string-or-comment))
    (setq pos (point))
    (beginning-of-line)
    (if (or (looking-at (concat "^" doby-label-regexp ":")) (looking-at "^[[:space:]]*\\(case .+\\|default\\):"))
        (end-of-line 0)
      (goto-char pos))
    (if (/= start-pos (point))
        (doby--backward-irrelevant stop-at-string))
    (/= start-pos (point))))

(defun doby--buffer-narrowed-p ()
  "Return non-nil if the current buffer is narrowed."
  (/= (buffer-size)
      (- (point-max)
         (point-min))))

(defun doby-previous-line-has-dangling-op-p ()
  "Returns non-nil if the current line is a continuation line."
  (let* ((cur-line (line-number-at-pos))
         (val (gethash cur-line doby-dangling-cache 'nope)))
    (if (or (doby--buffer-narrowed-p) (equal val 'nope))
        (save-excursion
          (beginning-of-line)
          (doby--backward-irrelevant t)
          (setq val (looking-back doby-dangling-operators-regexp))
          (if (not (doby--buffer-narrowed-p))
              (puthash cur-line val doby-dangling-cache))))
    val))

(defun doby--at-function-definition ()
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
            (setq start-nesting (doby-paren-level))
            (skip-chars-forward "^{")
            (while (> (doby-paren-level) start-nesting)
              (forward-char)
              (skip-chars-forward "^{") 0)
            (if (and (= (doby-paren-level) start-nesting) (= old-point (point)))
                t))))))

(defun doby--indentation-for-opening-parenthesis ()
  "Return the semantic indentation for the current opening parenthesis.

If point is on an opening curly brace and said curly brace
belongs to a function declaration, the indentation of the func
keyword will be returned. Otherwise the indentation of the
current line will be returned."
  (save-excursion
    (if (doby--at-function-definition)
        (progn
          (beginning-of-defun)
          (current-indentation))
      (current-indentation))))

(defun doby-indentation-at-point ()
  (save-excursion
    (let (start-nesting)
      (back-to-indentation)
      (setq start-nesting (doby-paren-level))

      (cond
       ((doby-in-string-p)
        (current-indentation))
       ((looking-at "[])}]")
        (doby-goto-opening-parenthesis)
        (if (doby-previous-line-has-dangling-op-p)
            (- (current-indentation) tab-width)
          (doby--indentation-for-opening-parenthesis)))
       ((progn (doby--backward-irrelevant t) (looking-back doby-dangling-operators-regexp))
        ;; only one nesting for all dangling operators in one operation
        (if (doby-previous-line-has-dangling-op-p)
            (current-indentation)
          (+ (current-indentation) tab-width)))
       ((zerop (doby-paren-level))
        0)
       ((progn (doby-goto-opening-parenthesis) (< (doby-paren-level) start-nesting))
        (if (doby-previous-line-has-dangling-op-p)
            (current-indentation)
          (+ (doby--indentation-for-opening-parenthesis) tab-width)))
       (t
        (current-indentation))))))

(defun doby-mode-indent-line ()
  (interactive)
  (let (indent
        shift-amt
        (pos (- (point-max) (point)))
        (point (point))
        (beg (line-beginning-position)))
    (back-to-indentation)
    (if (doby-in-string-or-comment-p)
        (goto-char point)
      (setq indent (doby-indentation-at-point))
      (if (looking-at (concat doby-label-regexp ":\\([[:space:]]*/.+\\)?$\\|case .+:\\|default:"))
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

(defun doby-beginning-of-defun (&optional count)
  (unless count (setq count 1))
  (let ((first t) failure)
    (dotimes (i (abs count))
      (while (and (not failure)
                  (or first (doby-in-string-or-comment-p)))
        (if (>= count 0)
            (progn
              (doby--backward-irrelevant)
              (if (not (re-search-backward doby-func-meth-regexp nil t))
                  (setq failure t)))
          (if (looking-at doby-func-meth-regexp)
              (forward-char))
          (if (not (re-search-forward doby-func-meth-regexp nil t))
              (setq failure t)))
        (setq first nil)))
    (if (< count 0)
        (beginning-of-line))
    (not failure)))

(defun doby-end-of-defun ()
  (let (orig-level)
    ;; It can happen that we're not placed before a function by emacs
    (if (not (looking-at "func"))
        (doby-beginning-of-defun -1))
    (skip-chars-forward "^{")
    (forward-char)
    (setq orig-level (doby-paren-level))
    (while (>= (doby-paren-level) orig-level)
      (skip-chars-forward "^}")
      (forward-char))))

;;;###autoload
(define-derived-mode doby-mode prog-mode "Doby"
  "Major mode for editing Doby source text.

This mode provides (not just) basic editing capabilities for
working with Doby code. It offers almost complete syntax
highlighting, indentation that is almost identical to dobyfmt and
proper parsing of the buffer content to allow features such as
navigation by function, manipulation of comments or detection of
strings.

In addition to these core features, it offers various features to
help with writing Doby code. You can directly run buffer content
through dobyfmt, read dobydoc documentation from within Emacs, modify
and clean up the list of package imports or interact with the
Playground (uploading and downloading pastes).

The following extra functions are defined:

- `dobyfmt'
- `dobydoc'
- `doby-import-add'
- `doby-remove-unused-imports'
- `doby-goto-imports'
- `doby-play-buffer' and `doby-play-region'
- `doby-download-play'
- `dobydef-describe' and `dobydef-jump'
- `doby-coverage'

If you want to automatically run `dobyfmt' before saving a file,
add the following hook to your emacs configuration:

\(add-hook 'before-save-hook 'dobyfmt-before-save)

If you want to use `dobydef-jump' instead of etags (or similar),
consider binding dobydef-jump to `M-.', which is the default key
for `find-tag':

\(add-hook 'doby-mode-hook (lambda ()
                          (local-set-key (kbd \"M-.\") #'dobydef-jump)))

Please note that dobydef is an external dependency. You can install
it with

doby get code.dobyogle.com/p/rog-doby/exp/cmd/dobydef


If you're looking for even more integration with Doby, namely
on-the-fly syntax checking, auto-completion and snippets, it is
recommended that you look at dobyflymake
\(https://github.com/dougm/dobyflymake), dobycode
\(https://github.com/nsf/dobycode) and yasnippet-doby
\(https://github.com/dominikh/yasnippet-doby)"

  ;; Font lock
  (set (make-local-variable 'font-lock-defaults)
       '(doby--build-font-lock-keywords))

  ;; Indentation
  (set (make-local-variable 'indent-line-function) #'doby-mode-indent-line)

  ;; Comments
  (set (make-local-variable 'comment-start) "// ")
  (set (make-local-variable 'comment-end)   "")
  (set (make-local-variable 'comment-use-syntax) t)
  (set (make-local-variable 'comment-start-skip) "\\(//+\\|/\\*+\\)\\s *")

  (set (make-local-variable 'beginning-of-defun-function) #'doby-beginning-of-defun)
  (set (make-local-variable 'end-of-defun-function) #'doby-end-of-defun)

  (set (make-local-variable 'parse-sexp-lookup-properties) t)
  (if (boundp 'syntax-propertize-function)
      (set (make-local-variable 'syntax-propertize-function) #'doby-propertize-syntax))

  (set (make-local-variable 'doby-dangling-cache) (make-hash-table :test 'eql))
  (add-hook 'before-change-functions (lambda (x y) (setq doby-dangling-cache (make-hash-table :test 'eql))) t t)


  (setq imenu-generic-expression
        '(("type" "^type *\\([^ \t\n\r\f]*\\)" 1)
          ("func" "^func *\\(.*\\) {" 1)))
  (imenu-add-to-menubar "Index")

  ;; Doby style
  (setq indent-tabs-mode t)

  ;; Handle unit test failure output in compilation-mode
  ;;
  ;; Note the final t argument to add-to-list for append, ie put these at the
  ;; *ends* of compilation-error-regexp-alist[-alist]. We want doby-test to be
  ;; handled first, otherwise other elements will match that don't work, and
  ;; those alists are traversed in *reverse* order:
  ;; http://lists.gnu.org/archive/html/bug-gnu-emacs/2001-12/msg00674.html
  (when (and (boundp 'compilation-error-regexp-alist)
             (boundp 'compilation-error-regexp-alist-alist))
    (add-to-list 'compilation-error-regexp-alist 'doby-test t)
    (add-to-list 'compilation-error-regexp-alist-alist
                 '(doby-test . ("^\t+\\([^()\t\n]+\\):\\([0-9]+\\):? .*$" 1 2)) t)))

;;;###autoload
(add-to-list 'auto-mode-alist (cons "\\.doby\\'" 'doby-mode))

(defun doby--apply-rcs-patch (patch-buffer)
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
            (error "invalid rcs patch or internal error in doby--apply-rcs-patch"))
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
                (doby--goto-line (- from line-offset))
                (incf line-offset len)
                (doby--delete-whole-line len)))
             (t
              (error "invalid rcs patch or internal error in doby--apply-rcs-patch")))))))))

(defun dobyfmt ()
  "Formats the current buffer according to the dobyfmt tool."

  (interactive)
  (let ((tmpfile (make-temp-file "dobyfmt" nil ".doby"))
        (patchbuf (get-buffer-create "*Dobyfmt patch*"))
        (errbuf (get-buffer-create "*Dobyfmt Errors*"))
        (coding-system-for-read 'utf-8)
        (coding-system-for-write 'utf-8))

    (with-current-buffer errbuf
      (setq buffer-read-only nil)
      (erase-buffer))
    (with-current-buffer patchbuf
      (erase-buffer))

    (write-region nil nil tmpfile)

    ;; We're using errbuf for the mixed stdout and stderr output. This
    ;; is not an issue because dobyfmt -w does not produce any stdout
    ;; output in case of success.
    (if (zerop (call-process dobyfmt-command nil errbuf nil "-w" tmpfile))
        (if (zerop (call-process-region (point-min) (point-max) "diff" nil patchbuf nil "-n" "-" tmpfile))
            (progn
              (kill-buffer errbuf)
              (message "Buffer is already dobyfmted"))
          (doby--apply-rcs-patch patchbuf)
          (kill-buffer errbuf)
          (message "Applied dobyfmt"))
      (message "Could not apply dobyfmt. Check errors for details")
      (dobyfmt--process-errors (buffer-file-name) tmpfile errbuf))

    (kill-buffer patchbuf)
    (delete-file tmpfile)))


(defun dobyfmt--process-errors (filename tmpfile errbuf)
  ;; Convert the dobyfmt stderr to something understood by the compilation mode.
  (with-current-buffer errbuf
    (goto-char (point-min))
    (insert "dobyfmt errors:\n")
    (while (search-forward-regexp (concat "^\\(" (regexp-quote tmpfile) "\\):") nil t)
      (replace-match (file-name-nondirectory filename) t t nil 1))
    (compilation-mode)
    (display-buffer errbuf)))

;;;###autoload
(defun dobyfmt-before-save ()
  "Add this to .emacs to run dobyfmt on the current buffer when saving:
 (add-hook 'before-save-hook 'dobyfmt-before-save).

Note that this will cause doby-mode to get loaded the first time
you save any file, kind of defeating the point of autoloading."

  (interactive)
  (when (eq major-mode 'doby-mode) (dobyfmt)))

(defun dobydoc--read-query ()
  "Read a dobydoc query from the minibuffer."
  ;; Compute the default query as the symbol under the cursor.
  ;; TODO: This does the wrong thing for e.g. multipart.NewReader (it only grabs
  ;; half) but I see no way to disambiguate that from e.g. foobar.SomeMethod.
  (let* ((bounds (bounds-of-thing-at-point 'symbol))
         (symbol (if bounds
                     (buffer-substring-no-properties (car bounds)
                                                     (cdr bounds)))))
    (completing-read (if symbol
                         (format "dobydoc (default %s): " symbol)
                       "dobydoc: ")
                     (doby--old-completion-list-style (doby-packages)) nil nil nil 'doby-dobydoc-history symbol)))

(defun dobydoc--get-buffer (query)
  "Get an empty buffer for a dobydoc query."
  (let* ((buffer-name (concat "*dobydoc " query "*"))
         (buffer (get-buffer buffer-name)))
    ;; Kill the existing buffer if it already exists.
    (when buffer (kill-buffer buffer))
    (get-buffer-create buffer-name)))

(defun dobydoc--buffer-sentinel (proc event)
  "Sentinel function run when dobydoc command completes."
  (with-current-buffer (process-buffer proc)
    (cond ((string= event "finished\n")  ;; Successful exit.
           (goto-char (point-min))
           (view-mode 1)
           (display-buffer (current-buffer) t))
          ((/= (process-exit-status proc) 0)  ;; Error exit.
           (let ((output (buffer-string)))
             (kill-buffer (current-buffer))
             (message (concat "dobydoc: " output)))))))

;;;###autoload
(defun dobydoc (query)
  "Show doby documentation for a query, much like M-x man."
  (interactive (list (dobydoc--read-query)))
  (unless (string= query "")
    (set-process-sentinel
     (start-process-shell-command "dobydoc" (dobydoc--get-buffer query)
                                  (concat "dobydoc " query))
     'dobydoc--buffer-sentinel)
    nil))

(defun doby-goto-imports ()
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
      (message "No imports or package declaration found. Is this really a Doby file?")
      'fail))))

(defun doby-play-buffer ()
  "Like `doby-play-region', but acts on the entire buffer."
  (interactive)
  (doby-play-region (point-min) (point-max)))

(defun doby-play-region (start end)
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
                       "http://play.dobylang.org/share"
                       (lambda (arg)
                         (cond
                          ((equal :error (car arg))
                           (signal 'doby-play-error (cdr arg)))
                          (t
                           (re-search-forward "\n\n")
                           (kill-new (format "http://play.dobylang.org/p/%s" (buffer-substring (point) (point-max))))
                           (message "http://play.dobylang.org/p/%s" (buffer-substring (point) (point-max)))))))))))

;;;###autoload
(defun doby-download-play (url)
  "Downloads a paste from the playground and inserts it in a Doby
buffer. Tries to look for a URL at point."
  (interactive (list (read-from-minibuffer "Playground URL: " (ffap-url-p (ffap-string-at-point 'url)))))
  (with-current-buffer
      (let ((url-request-method "GET") url-request-data url-request-extra-headers)
        (url-retrieve-synchronously (concat url ".doby")))
    (let ((buffer (generate-new-buffer (concat (car (last (split-string url "/"))) ".doby"))))
      (goto-char (point-min))
      (re-search-forward "\n\n")
      (copy-to-buffer buffer (point) (point-max))
      (kill-buffer)
      (with-current-buffer buffer
        (doby-mode)
        (switch-to-buffer buffer)))))

(defun doby-propertize-syntax (start end)
  (save-excursion
    (goto-char start)
    (while (search-forward "\\" end t)
      (put-text-property (1- (point)) (point) 'syntax-table (if (= (char-after) ?`) '(1) '(9))))))

(defun doby-import-add (arg import)
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
    (replace-regexp-in-string "^[\"']\\|[\"']$" "" (completing-read "Package: " (doby--old-completion-list-style (doby-packages))))))
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
        (case (doby-goto-imports)
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

(defun doby-root-and-paths ()
  (let* ((output (split-string (shell-command-to-string (concat doby-command " env DOBYROOT DOBYPATH"))
                               "\n"))
         (root (car output))
         (paths (split-string (cadr output) ":")))
    (append (list root) paths)))

(defun doby--string-prefix-p (s1 s2 &optional ignore-case)
  "Return non-nil if S1 is a prefix of S2.
If IGNORE-CASE is non-nil, the comparison is case-insensitive."
  (eq t (compare-strings s1 nil nil
                         s2 0 (length s1) ignore-case)))

(defun doby--directory-dirs (dir)
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
                                           (doby--directory-dirs file))
                                     dirs))))))
        dirs)
    '()))


(defun doby-packages ()
  (sort
   (delete-dups
    (mapcan
     (lambda (topdir)
       (let ((pkgdir (concat topdir "/pkg/")))
         (mapcan (lambda (dir)
                   (mapcar (lambda (file)
                             (let ((sub (substring file (length pkgdir) -2)))
                               (unless (or (doby--string-prefix-p "obj/" sub) (doby--string-prefix-p "tool/" sub))
                                 (mapconcat #'identity (cdr (split-string sub "/")) "/"))))
                           (if (file-directory-p dir)
                               (directory-files dir t "\\.a$"))))
                 (if (file-directory-p pkgdir)
                     (doby--directory-dirs pkgdir)))))
     (doby-root-and-paths)))
   #'string<))

(defun doby-unused-imports-lines ()
  ;; FIXME Technically, -o /dev/null fails in quite some cases (on
  ;; Windows, when compiling from within DOBYPATH). Practically,
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
                                   (concat doby-command
                                           (if (string-match "_test\.doby$" buffer-file-truename)
                                               " test -c"
                                             " build -o /dev/null"))) "\n")))))

(defun doby-remove-unused-imports (arg)
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
        (setq lines (doby-unused-imports-lines))
        (dolist (import lines)
          (doby--goto-line import)
          (beginning-of-line)
          (if arg
              (comment-region (line-beginning-position) (line-end-position))
            (doby--delete-whole-line)))
        (message "Removed %d imports" (length lines)))
      (if flymake-state (flymake-mode-on)))))

(defun dobydef--find-file-line-column (specifier other-window)
  "Given a file name in the format of `filename:line:column',
visit FILENAME and doby to line LINE and column COLUMN."
  (if (not (string-match "\\(.+\\):\\([0-9]+\\):\\([0-9]+\\)" specifier))
      ;; We've only been given a directory name
      (funcall (if other-window #'find-file-other-window #'find-file) specifier)
    (let ((filename (match-string 1 specifier))
          (line (string-to-number (match-string 2 specifier)))
          (column (string-to-number (match-string 3 specifier))))
      (with-current-buffer (funcall (if other-window #'find-file-other-window #'find-file) filename)
        (doby--goto-line line)
        (beginning-of-line)
        (forward-char (1- column))
        (if (buffer-modified-p)
            (message "Buffer is modified, file position might not have been correct"))))))

(defun dobydef--call (point)
  "Call dobydef, acquiring definition position and expression
description at POINT."
  (if (doby--xemacs-p)
      (error "dobydef does not reliably work in XEmacs, expect bad results"))
  (if (not (buffer-file-name (doby--coverage-origin-buffer)))
      (error "Cannot use dobydef on a buffer without a file name")
    (let ((outbuf (get-buffer-create "*dobydef*")))
      (with-current-buffer outbuf
        (erase-buffer))
      (call-process-region (point-min)
                           (point-max)
                           "dobydef"
                           nil
                           outbuf
                           nil
                           "-i"
                           "-t"
                           "-f"
                           (file-truename (buffer-file-name (doby--coverage-origin-buffer)))
                           "-o"
                           (number-to-string (doby--position-bytes (point))))
      (with-current-buffer outbuf
        (split-string (buffer-substring-no-properties (point-min) (point-max)) "\n")))))

(defun dobydef-describe (point)
  "Describe the expression at POINT."
  (interactive "d")
  (condition-case nil
      (let ((description (cdr (butlast (dobydef--call point) 1))))
        (if (not description)
            (message "No description found for expression at point")
          (message "%s" (mapconcat #'identity description "\n"))))
    (file-error (message "Could not run dobydef binary"))))

(defun dobydef-jump (point &optional other-window)
  "Jump to the definition of the expression at POINT."
  (interactive "d")
  (condition-case nil
      (let ((file (car (dobydef--call point))))
        (cond
         ((string= "-" file)
          (message "dobydef: expression is not defined anywhere"))
         ((string= "dobydef: no identifier found" file)
          (message "%s" file))
         ((doby--string-prefix-p "dobydef: no declaration found for " file)
          (message "%s" file))
         ((doby--string-prefix-p "error finding import path for " file)
          (message "%s" file))
         (t
          (push-mark)
          (ring-insert find-tag-marker-ring (point-marker))
          (dobydef--find-file-line-column file other-window))))
    (file-error (message "Could not run dobydef binary"))))

(defun dobydef-jump-other-window (point)
  (interactive "d")
  (dobydef-jump point t))

(defun doby--goto-line (line)
  (goto-char (point-min))
  (forward-line (1- line)))

(defun doby--line-column-to-point (line column)
  (save-excursion
    (doby--goto-line line)
    (forward-char (1- column))
    (point)))

(defstruct doby--covered
  start-line start-column end-line end-column covered count)

(defun doby--coverage-file ()
  "Return the coverage file to use, either by reading it from the
current coverage buffer or by prompting for it."
  (if (boundp 'doby--coverage-current-file-name)
      doby--coverage-current-file-name
    (read-file-name "Coverage file: " nil nil t)))

(defun doby--coverage-origin-buffer ()
  "Return the buffer to base the coverage on."
  (or (buffer-base-buffer) (current-buffer)))

(defun doby--coverage-face (count divisor)
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
    (concat "doby-coverage-" (number-to-string n))))

(defun doby--coverage-make-overlay (range divisor)
  "Create a coverage overlay for a RANGE of covered/uncovered
code. Uses DIVISOR to scale absolute counts to a [0,10] scale."
  (let* ((count (doby--covered-count range))
         (face (doby--coverage-face count divisor))
         (ov (make-overlay (doby--line-column-to-point (doby--covered-start-line range)
                                                     (doby--covered-start-column range))
                           (doby--line-column-to-point (doby--covered-end-line range)
                                                     (doby--covered-end-column range)))))

    (overlay-put ov 'face face)
    (overlay-put ov 'help-echo (format "Count: %d" count))))

(defun doby--coverage-clear-overlays ()
  "Remove existing overlays and put a single untracked overlay
over the entire buffer."
  (remove-overlays)
  (overlay-put (make-overlay (point-min) (point-max))
               'face
               'doby-coverage-untracked))

(defun doby--coverage-parse-file (coverage-file file-name)
  "Parse COVERAGE-FILE and extract coverage information and
divisor for FILE-NAME."
  (let (ranges
        (max-count 0))
    (with-temp-buffer
      (insert-file-contents coverage-file)
      (doby--goto-line 2) ;; Skip over mode
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
              (push (make-doby--covered :start-line start-line
                                      :start-column start-column
                                      :end-line end-line
                                      :end-column end-column
                                      :covered (/= count 0)
                                      :count count)
                    ranges)))

          (forward-line)))

      (list ranges (if (> max-count 0) (log max-count) 0)))))

(defun doby-coverage (&optional coverage-file)
  "Open a clone of the current buffer and overlay it with
coverage information gathered via doby test -coverprofile=COVERAGE-FILE.

If COVERAGE-FILE is nil, it will either be inferred from the
current buffer if it's already a coverage buffer, or be prompted
for."
  (interactive)
  (let* ((cur-buffer (current-buffer))
         (origin-buffer (doby--coverage-origin-buffer))
         (dobycov-buffer-name (concat (buffer-name origin-buffer) "<dobycov>"))
         (coverage-file (or coverage-file (doby--coverage-file)))
         (ranges-and-divisor (doby--coverage-parse-file
                              coverage-file
                              (file-name-nondirectory (buffer-file-name origin-buffer))))
         (cov-mtime (nth 5 (file-attributes coverage-file)))
         (cur-mtime (nth 5 (file-attributes (buffer-file-name origin-buffer)))))

    (if (< (float-time cov-mtime) (float-time cur-mtime))
        (message "Coverage file is older than the source file."))

    (with-current-buffer (or (get-buffer dobycov-buffer-name)
                             (make-indirect-buffer origin-buffer dobycov-buffer-name t))
      (set (make-local-variable 'doby--coverage-current-file-name) coverage-file)

      (save-excursion
        (doby--coverage-clear-overlays)
        (dolist (range (car ranges-and-divisor))
          (doby--coverage-make-overlay range (cadr ranges-and-divisor))))

      (if (not (eq cur-buffer (current-buffer)))
          (display-buffer (current-buffer) #'display-buffer-reuse-window)))))

(provide 'doby-mode)
