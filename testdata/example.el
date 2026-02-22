; CodeOwner: @elisp_owner

; This is a single-line comment

;; This is a conventional single-line comment

;;; Section comment in Emacs Lisp.

(defvar x 42)  ; Inline comment

(defun greet (name)
  ;; Another single-line comment
  (format "Hello, %s!" name))

(message (greet "world"))
