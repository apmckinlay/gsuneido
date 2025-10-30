## Create a Book

*Books* are part of Suneido's application framework.  Suneido's help is a simple example of a book.  The Suneido User's Manual book contains only HTML pages, but book pages can also contain code that returns a string of HTML, or code that returns a Suneido user interface specification.

You can look inside the Suneido User's Manual book by using Edit a Book on the WorkSpace IDE menu and choosing the `suneidoc` book.  You will get a BookEdit window with a tree pane on the left, a Scintilla control for entering page contents at the top right, and a preview pane at the bottom right.

Close this BookEdit and we'll create a new book.  Choose Edit a Book again, pick "New Book...", and enter `mybook` as the book name.  **Note:** If you loaded in the supplied book you can just select it from the list.

Books should start with a *cover* page.  Choose New Item and enter "Cover" for the name.  Then enter your cover page, for example:

``` suneido
<h1>My Book</h1>
```

**Note:** You only need to give the contents of the body i.e. you don't need \<html>, \<head>, or \<body> tags. (Unless you want to specify the \<head> contents.)

Next, most books have a *Contents* page:

``` suneido
<h1>Contents</h1>
<p><a href="Edit My Contacts">Edit My Contacts</a></p>
<p><a href="Print My Contacts">Print My Contacts</a></p>
```

Now we need to add the actual program pages.  **Note:** The names of the pages that you enter in the tree control must match the href's.  The Edit My Contacts page should contain:

``` suneido
My_ContactsAccess
```

and the Print My Contacts page should contain:

``` suneido
My_ContactsReport
```

One step remains.  By default, book pages are displayed in alphabetical order.  If you collapse and then expand the mybook folder in Book Edit you'll see the default order.  Use Set Order from the Tools menu to set the order of the pages, for example Cover - 1, Contents - 2, Edit My Contents - 3, Print My Contents - 4.  Collapse and expand to check the order.

![](<../res/mybookedit.png>)

Now you can open the book as a user would.  Choose Open a Book from the IDE menu, pick `mybook`, and voila, your book application.  Books display the top level of pages as tabs at the top.  If you expanded on this application you might want to create a Contacts page and put the Edit and Print underneath that.

![](<../res/mybook.png>)

Note: The Refresh and Goto buttons will only show up when Suneido.User is 'default'

**Note:** For more complex HTML you may want to use some other tool for creating HTML.  You can either paste in the HTML, or use Insert File on the Edit menu.