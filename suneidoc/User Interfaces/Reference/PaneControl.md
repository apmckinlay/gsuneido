### PaneControl

``` suneido
(control)
```

Creates a window containing the provided **control**.

For example:

``` suneido
Window(#(Pane
    #(Horz
        #(Field name: 'Value', ystretch: 1, width: 20)
        ystretch: 0
        )
    ))
```