;;; doubi-mode-load.el --- automatically extracted autoloads
;;; Commentary:

;; To install doubi-mode, add the following lines to your .emacs file:
;;   (add-to-list 'load-path "PATH CONTAINING doubi-mode-load.el" t)
;;   (require 'doubi-mode-load)
;;
;; After this, doubi-mode will be used for files ending in '.doubi'.
;;
;; To compile doubi-mode from the command line, run the following
;;   emacs -batch -f batch-byte-compile doubi-mode.el
;;
;; See doubi-mode.el for documentation.
;;
;; To update this file, evaluate the following form
;;   (let ((generated-autoload-file buffer-file-name)) (update-file-autoloads "doubi-mode.el"))

;;; Code:


;;;### (autoloads (doubi-download-play doubidoc doubifmt-before-save doubi-mode)
;;;;;;  "doubi-mode" "doubi-mode.el" (20767 50749))
;;; Generated autoloads from doubi-mode.el

(autoload 'doubi-mode "doubi-mode" "\
Major mode for editing Doubi source text.

This mode provides (not just) basic editing capabilities for
working with Doubi code. It offers almost complete syntax
highlighting, indentation that is almost identical to doubifmt,
proper parsing of the buffer content to allow features such as
navigation by function, manipulation of comments or detection of
strings.

Additionally to these core features, it offers various features to
help with writing Doubi code. You can directly run buffer content
through doubifmt, read doubidoc documentation from within Emacs, modify
and clean up the list of package imports or interact with the
Playground (uploading and downloading pastes).

The following extra functions are defined:

- `doubifmt'
- `doubidoc'
- `doubi-import-add'
- `doubi-remove-unused-imports'
- `doubi-doubito-imports'
- `doubi-play-buffer' and `doubi-play-region'
- `doubi-download-play'

If you want to automatically run `doubifmt' before saving a file,
add the following hook to your emacs configuration:

\(add-hook 'before-save-hook 'doubifmt-before-save)

If you're looking for even more integration with Doubi, namely
on-the-fly syntax checking, auto-completion and snippets, it is
recommended to look at doubiflymake
\(https://github.com/dougm/doubiflymake), doubicode
\(https://github.com/nsf/doubicode) and yasnippet-doubi
\(https://github.com/dominikh/yasnippet-doubi)

\(fn)" t nil)

(add-to-list 'auto-mode-alist (cons "\\.d\\'" 'doubi-mode))

(autoload 'doubifmt-before-save "doubi-mode" "\
Add this to .emacs to run doubifmt on the current buffer when saving:
 (add-hook 'before-save-hook 'doubifmt-before-save).

Note that this will cause doubi-mode to get loaded the first time
you save any file, kind of defeating the point of autoloading.

\(fn)" t nil)

(autoload 'doubidoc "doubi-mode" "\
Show doubi documentation for a query, much like M-x man.

\(fn QUERY)" t nil)

(autoload 'doubi-download-play "doubi-mode" "\
Downloads a paste from the playground and inserts it in a Doubi
buffer. Tries to look for a URL at point.

\(fn URL)" t nil)

;;;***

(provide 'doubi-mode-load)
;; Local Variables:
;; version-control: never
;; no-byte-compile: t
;; no-update-autoloads: t
;; coding: utf-8
;; End:
;;; doubi-mode-load.el ends here
