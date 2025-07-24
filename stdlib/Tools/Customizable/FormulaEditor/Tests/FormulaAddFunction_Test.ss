// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_convertArgs()
		{
		fn = FormulaAddFunction.FormulaAddFunction_convertArgs
		unitChooser = FormulaAddFunction.FormulaAddFunction_unitChooser
		textDisplay = FormulaAddFunction.FormulaAddFunction_textDisplay
		uomList = Uom_Conversions.Members()
		prompts = #(a, b, c)

		ob = Object()
		fn(Object(), ob, prompts)
		Assert(ob is: Object())

		ob = Object()
		args = #((Field, field1), (Unit, unit2), (Number, number),
			(Checkbox, checkbox), (Text, text1), (Date, date1), (false, default))
		expected = Object(
			Object('Pair', Object('Static', 'Field'),
				Object('ChooseList', prompts, name: 'field1')),
			Object('Pair', Object('Static', 'Unit'),
				Object(unitChooser, uomList, allowOther:, name: 'unit2')),
			Object('Pair', Object('Static', 'Number'),
				 Object('Number', mask: false, width: 20, name: 'number')),
			Object('Pair', Object('Static', 'Checkbox'),
				Object('CheckBox', name: 'checkbox')),
			Object('Pair', Object('Static', 'Text'),
				Object(textDisplay, width: 12, name: 'text1')),
			Object('Pair', Object('Static', 'Date'),
				Object('ChooseDate', name: 'date1')),
			Object('Field', name: 'default')
			)
		fn(args, ob, prompts)
		Assert(ob is: expected)
		}

	Test_extractArgs()
		{
		fn = FormulaAddFunction.FormulaAddFunction_extractArgs
		dateFormat = GetContributions('FormulaReserved').DATE.args[0].format
		fmtFormat = GetContributions('FormulaReserved').DATE.args[1].format
		args = Object(#(Field, field1, default: '""'), #(Unit, unit2),
			#(Number, number, default: 111), #(Checkbox, checkbox, default: true),
			Object(#Date, #date1, format: dateFormat),
			Object(#Format, #fmt1, format: fmtFormat), #(false, default))
		values = [fmt1: 'MM dd yyyy']
		list = Object()
		fn(args, values, list)
		Assert(list is: Object('""', '', 111, true, '""', '"MM dd yyyy"', ''))

		date = #20180120
		values = [field1: "field1", unit2: "km", number: 0, checkbox: false,
			date1: date, fmt1: 'yyyy MM dd']
		list = Object()
		fn(args, values, list)
		Assert(list is: Object("field1", 'km', 0, false, '"2018 01 20"', '"yyyy MM dd"',
			''))
		}
	}
