## Adding Context Sensitive Help

**Category:** User Interface

**Problem**

You want to add context sensitive help to a book.

**Ingredients**

[BookControl](<../User Interfaces/Reference/BookControl.md>)

**Recipe**

<u>Step 1: Create a Help Book</u>

First create your help book. Give the book the same name as your main book with "Help" added on the end.  For example, if your main book was named "Accounting", your help book would be named "AccountingHelp".

[BookControl](<../User Interfaces/Reference/BookControl.md>) handles the rest of the work with the code in the On_Help method, triggered by the F1 key.

<u>Step 2: Make it Context Sensitive</u>

When you press F1 for help, the help book will be opened, or if it is already open, activated (brought to the top). If the help book contains a page with the same name as the current book page, that page will be displayed.

Note: The headings in your help book must match the headings in your main book or the context sensitive help will not work. For example, if you have a page in your main book called "Processing", the corresponding help record in your help book must be called "Processing".

<u>Adding Help Buttons with Context Sensitivity to a Book Page</u>

Add an event to the button.

For example:

``` suneido
<INPUT Type="Button" value="Help" onClick="openHelp()">
```

In the openHelp function, the OpenBook function will be called.  To call the OpenBook function in suneido.

The JavaScript code will look something like this:

``` suneido
function openHelp()
    {
    document.location="suneido:/eval?OpenBook('HelpBook', 'PageName')";
    }
```

A page containing only a context sensitive help button would look something like this:

``` suneido
<html>
<body>
<script language="JavaScript">
function openHelp()
    {
    document.location="suneido:/eval?OpenBook('HelpBook', 'PageName')";
    }
</script>
<INPUT Type="Button" value="Help" onClick="openHelp()">
</body>
</html>
```