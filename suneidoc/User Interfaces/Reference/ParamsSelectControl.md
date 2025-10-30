### ParamsSelectControl

``` suneido
(field)
```

Builds a row with a prompt, a ComboBox for specifying the operator, and some type of 
control for entering a value. This is used for specifying parameters on a report.  

The field parameter can be a field name or a control specification.
If field is a field name, the input control will be retrieved by Datadict(). If field is a control,
then that control is used as the input control.