### Use CallClass for Dialogs

Rather than put the [Dialog](<../../User Interfaces/Reference/Dialog.md>) call in the application code, possibly duplicated multiple times, put the Dialog code in the [CallClass](<../../Language/Classes/CallClass.md>) of the [Controller](<../../User Interfaces/Reference/Controller.md>).

For example:

``` suneido
Controller
    {
    CallClass(...)
        {
        return Dialog(0, Object(this, ...))
        }
    New(...)
        {
        ...
        }
    ...
```

This is especially useful when the Dialog call is more complex, e.g. using styles or sizes.

**Note**: This approach can also be used for [Window](<../../User Interfaces/Reference/Window.md>) as well as Dialog