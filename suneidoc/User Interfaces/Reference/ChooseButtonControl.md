### ChooseButtonControl

``` suneido
(text, list, width = false)
```

Similar to a MenuButton but the text shows your choice like an ChooseList.

**text** is the initial value and what is displayed

If **width** is not specified then the button will resize depending on the value chosen.

For example:

``` suneido
Window(#(ChooseButton choose (one two three)))
```

![](<../../res/choosebutton.png>)

Typing the first letter of a value will choose it. If there are several values with the same first letter you can type the letter multiple times to select between them.

Use by 
[ParamsSelectControl](<ParamsSelectControl.md>)