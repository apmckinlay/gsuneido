## Using Browse on Access for Master-Detail Relationships

**Category:** User Interface

**Problem**

You have a master table with one or more detail tables that have a foreign key to the master table and you need a way for users to enter the detail records on the master record.

**Ingredients**

AccessControl, BrowseControl

**Recipe**

We will use the example of an invoice. The invoice will have regular line items and also extra charge line items. Here are the tables:

``` suneido
invoices (invoice_num, invoice_date, invoice_customer)  
    key (invoice_num)
```

``` suneido
invoice_lines (line_num, invoice_num, line_desc, line_amount)
    key (line_num)
    index (invoice_num) in invoices
```

``` suneido
invoice_extra (extra_num, invoice_num, extra_desc, extra_amount)
    key (extra_num)
    index (invoice_num) in invoices
```

To handle the keys on the tables, we will define rules that simply return timestamps. This will ensure that the fields values are unique. Also, the key fields will be renamed to prevent the rules from being done on the foreign key fields in the detail tables:

``` suneido
Rule_invoice_num_new
function ()
    {
    return Timestamp()
    }

Rule_line_num_new
function ()
    {
    return Timestamp()
    }

Rule_extra_num_new
function ()
    {
    return Timestamp()
    }
```

It is also important to set up Field_ definitions for the key fields to ensure that the values are always treated as dates (because we used Timestamps).

``` suneido
Field_invoice_num
Field_num
    {
    }

Field_invoice_num_new
Field_invoice_num
    {
    }

Field_line_num_new
Field_num
    {
    }

Field_extra_num_new
Field_num
    {
    }
```

We also need to make a rule to "replicate" the invoice_num. We need 'invoice_num_new2' because you can't have two controls (in this case Browse's) with the same name in a RecordControl because the controls are tracked by name.

``` suneido
Rule_invoice_num_new2
function ()
    {
    return .invoice_num_new
    }
```

Now we should be able to set up a simple Access to enter invoices:

``` suneido
InvoiceAccess
#(Access "invoices rename invoice_num to invoice_num_new"
    (Vert
        invoice_date
        invoice_customer
        (Browse "invoice_lines rename line_num to line_num_new"
            columns: (line_desc, line_amount)
            linkField: "invoice_num",
            name: "invoice_num_new")
        (Browse "invoice_extra rename extra_num to extra_num_new"
            columns: (extra_desc, extra_amount)
            linkField: "invoice_num",
            name: "invoice_num_new2")
        )
    )
```

The 'name' argument to Browse should always be a field (or rule) from the master table query. The 'linkField' argument should be the field from the Browse detail query.

Now you can test the Access by running the following from the WorkSpace:

``` suneido
Window(InvoiceAccess)
```