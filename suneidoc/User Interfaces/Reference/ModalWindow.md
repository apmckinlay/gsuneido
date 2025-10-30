### ModalWindow

``` suneido
(control, title = false, border = 5, closeButton? = true,
        onDestroy = false, keep_size = true, useDefaultSize = false)
```

A [Window](<Window.md>) that disables its parent while it's open. This is similar to [Dialog](<Dialog.md>) except that ModalWindow does not run a nested message loop, so it is non-blocking. This means it can't "return" a result like a Dialog.

It automatically centers itself on its parent, and it automatically saves and restores its size (if stretchable).