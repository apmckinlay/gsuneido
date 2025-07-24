// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
#(
	CONVERT: (args: (('Field', #field, allowOther: false), ('Unit', #unit))
		desc: "Usage: CONVERT(Field, Unit)\r\n" $
			"Convert a <Rate> or <Quantity> to the specified unit of measure.", func: 1)
	TONUMBER: (args: (('Field', #field, allowOther: false), ('Unit', #unit))
		desc: "Usage: TONUMBER(Field, Unit)\r\n" $
			"Convert a <Rate> or <Quantity> to the specified unit of measure and " $
			"return the numeric value", func: 1.1)
	UOM: (args: (('Field', #field, allowOther: false))
		desc: "Usage: UOM(Field)\r\n" $
			"Get the unit of measure of a <Rate> or <Quantity>", func: 1.2)
	QUANTITY: (args: (('Value', #number), ('Unit', #unit))
		desc: "Usage: QUANTITY(Value, Unit)\r\n" $
			"Add a Value + Unit of Measure to the formula as a quantity.")
	RATE: (args: (('Value', #number), ('Unit', #unit))
		desc: "Usage: RATE(Value, Unit)\r\n" $
			"Add a Value + Unit of Measure to the formula as a rate.")
	IF: (args: ((false, #if)),
		desc: 'Usage: IF([Field1, Operator1, Value1, Then1],\r\n' $
			'\t[Field2, Operator2, Value2, Then2],\r\n' $
			'\t...,\r\n' $
			'\t[FieldN, OperatorN, ValueN, ThenN],\r\n' $
			'\tElse)\r\n' $
			'Add one or more "Conditions" (Field, Operator, Value) and Return values ' $
			'to the formula. Use the Else field to specify the Return value when none ' $
			'of the conditions are met. ' $
			'The return value for Else will default to empty if not filled in.', func: 2)
	CONTAINS: (args: (('Text', #field1), ('Substring', #text1)),
		desc: "Usage: CONTAINS(Text, Substring)\r\n", func: 3)
	STARTSWITH: (args: (('Text', #field1), ('Substring', #text1)),
		desc: "Usage: STARTSWITH(Text, Substring)\r\n", func: 4)
	ENDSWITH: (args: (('Text', #field1), ('Substring', #text1)),
		desc: "Usage: ENDSWITH(Text, Substring)\r\n", func: 5)
	MATCHES: (args: (('Text', #field1), ('Substring', #text1)),
		desc: "Usage: MATCHES(Text, Substring)\r\n", func: 6)
	AFTERLAST: (args: (('Text', #field1), ('Delimiter', #text1)),
		desc: "Usage: AFTERLAST(Text, Delimiter)\r\n", func: 7),
	BEFOREFIRST: (args: (('Text', #field1), ('Delimiter', #text1)),
		desc: "Usage: BEFOREFIRST(Text, Delimiter)\r\n", func: 8),
	UPPER: (args: (('Text', #field1)),
		desc: "Usage: UPPER(Text)\r\n", func: 8.1),
	DATE: (args: (
		('Date', #date1, format: function(val, record)
			{ Display(Date?(val)
				? val.Format(record.fmt1.Replace('\<yy?y?\>', 'yyyy')) : val) }),
		('Format', #fmt1, format: function(val) { Display(val) }))
		desc: "Usage: DATE(Date, Format)\r\n" $
			"You can use Keyboard Shortcuts, such as t, m, h, y, or r to enter dates " $
			"and use + or - to add or subtract days", func: 9),
	MONTH: (args: (('Field', #field, allowOther: false)),
		desc: "Usage: MONTH(Field)\r\nReturn the month portion of a date", func: 10),
	YEAR: (args: (('Field', #field, allowOther: false)),
		desc: "Usage: YEAR(Field)\r\nReturn the year portion of a date", func: 11)
	WEEKNUMBER: (args: (('Field', #field, allowOther: false)),
		desc: "Usage: WEEKNUMBER(field)\r\nReturn the week number of a date", func: 12)
	DAY: (args: (('Field', #field, allowOther: false)),
		desc: 'Usage: DAY(field)\r\nReturn the day portion of a date', func: 13)
	DAYOFWEEK: (args: (('Field', #field, allowOther: false)),
		desc: 'Usage: DAYOFWEEK(field)\r\nReturn the day of the week of a date', func: 14)
)
