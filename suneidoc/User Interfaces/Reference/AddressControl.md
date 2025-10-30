### AddressContol

``` suneido
(prefix = '', suffix = '', title = '')
```

Combines the controls for the
'address1', 'address2', 'city', 'state_prov', and 'zip_postal' fields.
The supplied prefix and suffix are used to derive the field names.

If a title is specified, the fields are contained in a groupbox
with the title as the label.

For example:

``` suneido
Window(#('Address'))
```

Would produce:
![](<../../res/AddressControl.gif>)