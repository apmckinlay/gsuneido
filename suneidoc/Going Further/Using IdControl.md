## Using IdControl

[IdControl](<../User Interfaces/Reference/IdControl.md>) is used for looking up an id from a master table.  For example, looking up a customer when entering an invoice. IdControl is designed to be used with the "num, abbrev, name" pattern.  If you are not using this pattern, you probably want to use [KeyControl](<../User Interfaces/Reference/KeyControl.md>) instead.

In order to use IdControl correctly, the master table should follow certain conventions.  It should have:

-	A key with a suffix of "_num" that should be a unique number or timestamp. This is an internal identifier that the user normally never sees.
-	A key with a suffix of "_name".  This is the value that will display in the IdControl when a master record is chosen.
-	A unique index with a suffix of "_abbrev".  This allows users to enter an optional shorter version of the name, i.e. a *nickname*.


IdControl also allows prefixes to be entered and will match the prefix entered to a record based on the name and if it can't find a unique match, it will check the prefix against the abbreviation field.

The list button can also be used to choose a valid record from the list. Pressing enter on the field will do the prefix checking and pressing Alt + z will also bring up the list.

Next we will go through a simple example of a customer lookup using an IdControl.  The first thing to do is set up the master  table for customers.  Create the table by running the following from QueryView:

``` suneido
create customers (customer_num, customer_name, customer_abbrev)
    key(customer_num) key(customer_name) index unique(customer_abbrev)
```

The easiest way to get a unique value in the _num key is to use a timestamp rule. Since we only want the rule to fill in the value on the master table, we will rename the customer_num field to customer_num_new when we use it on the master table Access. Here is the rule definition:

**`Rule_customer_num_new`**
``` suneido
function()
    {
    return Timestamp()
    }
```

Run the following from the workspace and enter a few customers:

``` suneido
AccessControl("customers rename customer_num to customer_num_new")
```

Now we are ready to set up an invoice header.  Create a table for this by running the following from QueryView: 

``` suneido
create invoices (ivc_num, customer_num, amount) key(ivc_num)
```

Make sure it has the customer_num field.  This is the field that we will be using the IdControl on. Now all we have to do is set up the datadict entry in your library to specify that the customer_num field will be using IdControl.  The master table will not use the IdControl because we did the renaming. Here is the definition for the datadict entry:

**`Field_customer_num`**
``` suneido
Field_num
    {
    Prompt: "Customer"
    Control: (Id "customers" "customer_num")
    }
```

In the above control specification, the first member, Id, is the name of the control, "customers" is the first argument to IdControl and is the query, "customer_num" is the field from the query whose value will be stored (the _name field will be displayed in the control, by default).  Now test the simple customer lookup by running the following from the WorkSpace:

``` suneido
AccessControl("invoices")
```

That is all you need to get your lookup working using IdControl, although there are many more options that you can use with IdControl. Please refer to the [IdControl](<../User Interfaces/Reference/IdControl.md>) documentation for more information.  If you would like to go further with the invoice example, refer to the [Master-Detail Relationship](<Master-Detail Relationship.md>) section for instructions on linking line items to a header record.