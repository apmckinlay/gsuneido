### Control Rules

Rules can be used to send messages to controls (within a RecordControl) 
by using a special naming convention.  
For example, to have a rule determine whether a 'name' field is protected or not:

``` suneido
Rule_name__protect
function ()
    { return .id is ""; }
```

**Note:** There must be **two** underscores between the field name and the method name.  
The method name is capitalized before being called.

In general, a rule called:

``` suneido
   "Rule_" $ fieldname $ "__" $ method
```

will translate to:

``` suneido
   control[Method](value)
```

Note: For this to work, the control must support the method.
Currently, FieldControl supports "valid" and "protect" methods.