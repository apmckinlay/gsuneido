## Books

*Books* are Suneido's application *metaphor*.  Suneido's help is a simple example of a book.  The SuneiDoc book contains only HTML pages, but book pages can also contain code that returns a string of HTML, or code that returns a Suneido user interface specification.

![](<../res/bookall.png>)

Application books do not have the toolbar. And in this case the sidebars have been hidden.

![](<../res/mybook.png>)

Note: The Refresh and Goto buttons will only show up when Suneido.User is 'default'

**Tool Bar**

-	Back is similar to Back on a web browser - it takes you to the page you viewed prior to this one. Click on the down arrow to choose from a list of prior pages.
-	Forward is similar to Forward on a web browser - if you have used Back, it takes you back where you've been. Click on the down arrow to choose from a list of pages.
-	Up moves to the page before the current one.
-	Down moves to the page after the current one.
-	Print prints the current page. (Only valid for HTML pages, not program pages.) Right-click to Print Section and Print Preview to print the current page and its *children*.
-	Use the search box to search within a help book.


**Tabs**

The tabs below the tool bar show the main *chapters* of the book. Click on a tab to go to that section of the book. Within nested sub-menus, clicking on the tab you are in will take you "up" one level.

**Bookmarks**

Use the '+' to add a bookmark for the current page. You can also add a bookmark for the current page by simply double clicking on the bookmark area (not on a bookmark).

Use the '-' to remove the bookmark for the current page.

Click on a bookmark to go to that page.

You can drag the bookmarks to change the order, and right click to change their color.

**Tree View**

The optional tree on the left hand side 
shows the *table of contents* of the book.
You can use it to see where you are
and to move around in the book.

### Configuration

Books have a number of optional parts that you can hide or show.
The screenshot above has all its parts showing.

The controllable parts are:
Treeview
: Click on the center arrow on the splitter to hide or show the tree view,
    or drag the splitter to adjust the size.

Bookmarks
: Click on the center arrow on the splitter to hide or show the bookmarks,
    or drag the splitter to adjust the size.
    Click on the '+' button 
    (or double click in the bookmark area)
    to add a bookmark to the current page.
    Click on a bookmark to go to that page.
    Click on the X button to remove the bookmark for the current page.
    

Here is the help book with all the optional parts hidden:

![](<../res/booknone.png>)

Note: If run in a persistent window, like the help is, then your configuration of these parts will be saved and restored. However, if you just use Open a Book from the IDE menu, your settings will **not** be saved.

See also:
[BookControl](<../User Interfaces/Reference/BookControl.md>)