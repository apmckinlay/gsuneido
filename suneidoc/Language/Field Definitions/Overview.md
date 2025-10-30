### Overview

Field definitions are used by the User Interface and Reports to provide additional information about fields, i.e. a *data dictionary*.

A field definition is a "class". It should almost always inherit from one of the "Base" classes:

-	Field_string
-	Field_number
-	Field_date
-	Field_boolean


If you use a "standard" field name such as email, you'll notice that it "automatically" has a prompt and a specific control. This is because stdlib contains a number of definitions for common fields (look in the Datadict folder). You can use these directly, by naming your field appropriately (e.g. email) or by inheriting from them, which also allows you to override options. For example:

**`Field_customer_email`**
``` suneido
Field_email
    {
    Prompt: "Customer Email"
    }
```

A field definition can include several members. All are optional and will inherit from the base if not defined.
Prompt
: A string to display as a "label" in front of the field

SelectPrompt
: Defaults to Prompt. If you have duplicate Prompt's on an 
[AccessControl](<../../User Interfaces/Reference/AccessControl/AccessControl.md>) (e.g. two "City" fields) you can specify different prompts to use in Select so the user can tell them apart (e.g. "Origin City" and "Destination City").

Heading
: Defaults to Prompt. A string to display as a "column heading" on reports

Control
: the user interface control to use

Format
: The report format used to print the value of the field on reports and to display the value in a 
[BrowseControl](<../../User Interfaces/Reference/BrowseControl/BrowseControl.md>)

The name of the record that contains the field definition should start with "Field_", followed by the name of the field the definition is for.

For example, the name of the record that contains the definition for a field called *firstname* would be:

``` suneido
Field_firstname
```

If this field were to be used to store first names as strings, the field definition could inherit from Field_string.  The prompt should be set.  In this case, it would most likely be set to "First Name".  The code would look like this:

``` suneido
Field_string
    {
    Prompt: "First Name"
    }
```

If the field was not to be modified, the control could be overridden to be readonly.  The code would then look like this:

``` suneido
Field_string
    {
    Prompt: "First Name"
    Control: (Field readonly:)
    }
```

If when the field printed on a report, the font was too small, the format could be overridden in the field definition.  The code would then look like this:

``` suneido
Field_string
    {
    Prompt: "First Name"
    Control: (Field readonly:)
    Format: (Text font: #(name: "Arial" size: 16 weight: 400))
    }
```

To change the column heading that appears when the field in on a [BrowseControl](<../../User Interfaces/Reference/BrowseControl/BrowseControl.md>) or the prompt that appears in the select choose list on as [AccessControl](<../../User Interfaces/Reference/AccessControl/AccessControl.md>), the Heading and SelectPrompt could be added.  The code would then look like this:

``` suneido
Field_string
    {
    Prompt: "First Name"
    SelectPrompt: "Customer First Name"
    Heading: "First\nName"
    Control: (Field readonly:)
    Format: (Text font: (name: "Arial", size: 16, weight: 400))
    }
```

As above, a report Heading can be made multi-line by including newlines (\n).