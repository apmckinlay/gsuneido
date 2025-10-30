### InfoControl

``` suneido
(def = 0, mandatory = false)
```

A composite control consisting of a ComboBox to choose the type of information,
and a second field to enter the information.
The initial type can be specified with def as a number from 0 to 7.

![](<../../res/infocontrol.gif>)

The allowed types and the corresponding controls are:

|  |  |  | 
| :---- | :---- | :---- |
| Work: | Phone | work phone number | 
| Fax: | Phone | fax phone number | 
| Home: | Phone | home phone number | 
| Cell: | Phone | cell phone number | 
| Pager: | Phone | pager phone number | 
| Email: | MailLink | email address | 
| Website: | HttpLink | website address | 
| Other: | Field | anything | 


The *value* of the control is a string of the form "type: value",
for example: "Work: 249-5050".  
If the right hand field is empty, the value is "".
If an InfoControl is set to "", the **def** type will be selected.

If the **mandatory** parameter is true, then must enter a value in the field