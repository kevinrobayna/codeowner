; CodeOwner: @clojure_owner

; This is a single-line comment

;; This is a conventional single-line comment

(comment
  "This is a block-level comment form in Clojure.")

(def x 42)  ; Inline comment

(defn greet [name]
  ;; Another single-line comment
  (str "Hello, " name "!"))

(println (greet "world"))
