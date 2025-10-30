### Ask

``` suneido
(prompt = "", title = "", hwnd = 0, ctrl = 'Field',
    valid = function (value) { return '' })
```

Creates a dialog with title, prompt, a field, where user can enter in data, and OK, Cancel buttons.

The valid function should return "" if the value is valid, or an error message if it is not.

For example:

``` suneido
Ask('Name', 'Login')
```

would display

![](<../../res/ask.png>)