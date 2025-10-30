### CustomizableControl

``` suneido
(table, name = false)
```

Allow end users to add fields to a table and design a layout for the fields.

`name` defaults to `table`

User defined layouts are saved in a table called `customizable` under the specified `name`.

For example, if we create a table with a timestamp key field:

``` suneido
Database("create mytable (test_timestamp) key(test_timestamp)")
```

Then we can use Customizable with it like this:

``` suneido
Window(#(Access mytable (Customizable mytable)))
```

![](<../../res/customizable.png>)

The layout is a simple form of wysiwyg that is converted to a [FormControl](<FormControl.md>) layout. Valid field names are recognized and converted to prompts and controls, anything else becomes static text. Fields that start in the same column will be placed in the same Form group.

The standard field types are:

-	Checkmark (
	[CheckBoxControl](<CheckBoxControl.md>), 
	[CheckMarkFormat](<../../Reports/Reference/CheckMarkFormat.md>))
-	Number, no decimals (
	[NumberControl](<NumberControl.md>), 
	[NumberFormat](<../../Reports/Reference/NumberFormat.md>))
-	Number, 1 decimal (
	[NumberControl](<NumberControl.md>), 
	[NumberFormat](<../../Reports/Reference/NumberFormat.md>))
-	Number, 2 decimals (
	[NumberControl](<NumberControl.md>), 
	[NumberFormat](<../../Reports/Reference/NumberFormat.md>))
-	Number, 3 decimals (
	[NumberControl](<NumberControl.md>), 
	[NumberFormat](<../../Reports/Reference/NumberFormat.md>))
-	Number, 4 decimals (
	[NumberControl](<NumberControl.md>), 
	[NumberFormat](<../../Reports/Reference/NumberFormat.md>))
-	Text, single line (
	[FieldControl](<FieldControl.md>), 
	[TextFormat](<../../Reports/Reference/TextFormat.md>))
-	Text, multi line (
	[EditorControl](<EditorControl.md>), 
	[WrapFormat](<../../Reports/Reference/WrapFormat.md>))
-	Short Date (
	[ChooseDateControl](<ChooseDateControl.md>), 
	[ShortDateFormat](<../../Reports/Reference/ShortDateFormat.md>))
-	Long Date (
	[ChooseDateControl](<ChooseDateControl.md>), 
	[LongDateFormat](<../../Reports/Reference/LongDateFormat.md>))


Field types are defined with plugins so it is easy to add additional ones. The standard ones are defined in Plugin_FieldTypes.

A typical use of Customizable is to add a "Custom" tab on an Access that contains a CustomizableControl. This allows users to add their own fields to your application's tables. They can then use these custom fields in Select or Reporter.

**Note:** CustomizableControl should be regarded as "beta". It is usable, but it needs more error checking/handling and lacks any way to rename or delete fields.