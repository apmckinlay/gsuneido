### PairControl

``` suneido
( left_control, right_control )
```

Displays two controls side by side.  Calculates Descent for baseline allignment in FormControls.  Also for Forms, Left is set so that grouped pairs are alligned vertically on the line between the left and right controls, like:

``` suneido
  one   two
three   four
 five   six
```

Pair is commonly used to display a field's prompt and control.  
For example, Window.Construct automatically looks up field names in the data dictionary 
and creates a Pair with the field's prompt and control.