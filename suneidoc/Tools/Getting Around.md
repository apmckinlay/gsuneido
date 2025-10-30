## Getting Around

Suneido includes a number of functions to help you move around both in your own code and in the standard library (stdlib). You have likely discovered some of these, but some of them may not be obvious. Also, these functions are in various places in the IDE and so any documentation is scattered. Although it may not seem that important to be familiar with these functions, knowing how to quickly get to where you want to be can make your programming much more efficient and pleasant.

### Go To Definition

Probably the most important function is Go To Definition from any Scintilla source code editor, such as Library View. This function is available from the right-click context menu, or by pressing F12. If you have some text selected it will try to locate that text. Otherwise, it will select the word at the cursor and try to locate that. The word selection is equivalent to double clicking at that point.

Go To Definition does not just look for an exact match to the word; it looks for the following variations:

-	name
-	name $ "Control"
-	name $ "Format"
-	"Rule_" $ name.UnCapitalize()
-	"Field_" $ name.UnCapitalize()
-	"Trigger_" $ name


It also handles:

-	super and super.Method
-	.private_method
-	.Public_method (If the method is inherited it will go to the method in the parent class)
-	Global.Name


All the libraries in the database are searched. Once it has searched for these variations:

-	If it can't find the word, nothing happens.
-	If it finds a single match, it will go to that record in Library View.
-	If several matches are found, a list will pop up so you can choose which one to go to.


It gets a Library View as follows:

-	If you don't have a Library View open it will open one.
-	If the Go To Definition was done from a Library View, that Library View will be used.
-	Otherwise, the Library View that was opened first will be used.


If the name you go to is a previous definition reference like _Name, then Go To Definition is smart enough to know that it has to look in previous libraries. This means if your Go To takes you to a record that derives from or calls the previous definition, you can easily work your way through the previous definitions.

### WorkSpace Find

The WorkSpace has a Find tab (on the right) that can be used to search your libraries in a variety of ways. Double-click on a result to go to it.

Note: If you are writing other tools that output references to source, consider using the same format as Find in Library (i.e. library:name:line) so that Go To Definition will work with it.

### Go To Line

To move around within library records Library View has a Go To Line command on the Edit menu that is also available via CTRL+G. This command is not used a lot since you don't usually know line numbers (unless you get something like an error message that contains a line number). It can be useful if someone points you to a certain line. You can also use Show/Hide Line Numbers on the right-click context menu.

### Show References

Sometimes when you're looking at a library record, you want to know where that record is used. You could use WorkSpace Find, but that means switching to the WorkSpace and typing in the name of the record. The Show References (on the Library View Tool menu or 'R' on the tool bar) function provides a quick way to see where the current record is used. It will display a list of all the locations and allow you to go to them. (Tip: Double-click a line to go to it.) The list is divided into two parts by a blank line. The first part shows all definitions of the name you're currently on. Normally there will only be a single definition, but it is possible to overload names in different libraries, in which case there would be several entries. The second part shows all the records that refer to the current name.

Similar to the variations used by Go To Definition, Show References strips off standard prefixes (Field_, Rule_, Trigger_) and suffixes (Control, Format) before searching. Occasionally, this may result in some extra undesired matches but this is usually not a big problem.

### Library View Outline
Click on items in the Library View outline pane (on the right) to go to them. 
### Debug

The Debug window also provides a Go To Definition button. This button will take you to the location selected in the call stack.

Note: You can also use Go To Definition from the right-click context menu of the upper source code pane.

### Test Runner

From the Test Runner (Gui) you can easily jump to the definitions of the tests using F12 or Go To Definition on the right-click context menu of the list.

### Version Control

Version Control allows Go To Definition from several places. You can right-click on an item in the local changes list and choose Go To Definition from the context menu. Or you can use the Go To Definition button when viewing the differences between two versions. This will take you to the line of the currently selected difference.

### Book

When Suneido.User is "default" (which it is by default), program pages in BookControl will have a Go To button which will take you to the definition of the page you're on.

Note: This only works if what is in the book is simple (e.g. "Name" or "Name()"). If you put more complicated code in the book, Go To may not be able to figure out where to go.

### Tips

When you are running code from the WorkSpace, Go To Definition is a quick way to get to the where things are defined, either in your own code, or in the standard library (stdlib).

You may even want to leave a list of names on your WorkSpace to serve as quick links to library records.

Don't forget you can use Go To Definition from QueryView (e.g. on field or table names), BookEdit, or SchemaView.