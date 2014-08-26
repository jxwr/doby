;;; doby-mode-load.el --- automatically extracted autoloads
;;; Commentary:

;; To install doby-mode, add the following lines to your .emacs file:
;;   (add-to-list 'load-path "PATH CONTAINING doby-mode-load.el" t)
;;   (require 'doby-mode-load)
;;
;; After this, doby-mode will be used for files ending in '.doby'.
;;
;; To compile doby-mode from the command line, run the following
;;   emacs -batch -f batch-byte-compile doby-mode.el
;;
;; See doby-mode.el for documentation.
;;
;; To update this file, evaluate the following form
;;   (let ((generated-autoload-file buffer-file-name)) (update-file-autoloads "doby-mode.el"))

;;; Code:


;;;### (autoloads (doby-download-play dobydoc dobyfmt-before-save doby-mode)
;;;;;;  "doby-mode" "doby-mode.el" (20767 50749))
;;; Generated autoloads from doby-mode.el

(autoload 'doby-mode "doby-mode" "\
Major mode for editing Doby source text.

This mode provides (not just) basic editing capabilities for
working with Doby code. It offers almost complete syntax
highlighting, indentation that is almost identical to dobyfmt,
proper parsing of the buffer content to allow features such as
navigation by function, manipulation of comments or detection of
strings.

Additionally to these core features, it offers various features to
help with writing Doby code. You can directly run buffer content
through dobyfmt, read dobydoc documentation from within Emacs, modify
and clean up the list of package imports or interact with the
Playground (uploading and downloading pastes).

The following extra functions are defined:

- `dobyfmt'
- `dobydoc'
- `doby-import-add'
- `doby-remove-unused-imports'
- `doby-dobyto-imports'
- `doby-play-buffer' and `doby-play-region'
- `doby-download-play'

If you want to automatically run `dobyfmt' before saving a file,
add the following hook to your emacs configuration:

\(add-hook 'before-save-hook 'dobyfmt-before-save)

If you're looking for even more integration with Doby, namely
on-the-fly syntax checking, auto-completion and snippets, it is
recommended to look at dobyflymake
\(https://github.com/dougm/dobyflymake), dobycode
\(https://github.com/nsf/dobycode) and yasnippet-doby
\(https://github.com/dominikh/yasnippet-doby)

\(fn)" t nil)

(add-to-list 'auto-mode-alist (cons "\\.d\\'" 'doby-mode))

(autoload 'dobyfmt-before-save "doby-mode" "\
Add this to .emacs to run dobyfmt on the current buffer when saving:
 (add-hook 'before-save-hook 'dobyfmt-before-save).

Note that this will cause doby-mode to get loaded the first time
you save any file, kind of defeating the point of autoloading.

\(fn)" t nil)

(autoload 'dobydoc "doby-mode" "\
Show doby documentation for a query, much like M-x man.

\(fn QUERY)" t nil)

(autoload 'doby-download-play "doby-mode" "\
Downloads a paste from the playground and inserts it in a Doby
buffer. Tries to look for a URL at point.

\(fn URL)" t nil)

;;;***

(provide 'doby-mode-load)
;; Local Variables:
;; version-control: never
;; no-byte-compile: t
;; no-update-autoloads: t
;; coding: utf-8
;; End:
;;; doby-mode-load.el ends here
