### AddressFormat

``` suneido
(data, prefix = "", suffix = "", justify = 'left')
```

Takes an object containing address information and prints it in a proper format.  Note that members have to have appropriate names (address1, address2, city, state_prov, zip_postal) in order for AddressFormat to work properly.

**justify** can be either 'left' (the default) or 'center'.

For example:

``` suneido
Params.On_Preview(Object('Address' 
    data: #(test_address1: '123 Main Street'
        test_address2: ''
        test_city: 'Saskatoon'
        test_state_prov: 'SK'
        test_zip_postal: 'S7N 1L6') 
    prefix: 'test_'))
```

will print:

``` suneido
123 Main Street
Saskatoon SK  S7N 1L6
```

See also: [AddressControl](<../../User Interfaces/Reference/AddressControl.md>)