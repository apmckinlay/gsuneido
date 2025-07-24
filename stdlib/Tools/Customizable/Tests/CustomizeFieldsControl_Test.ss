// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_getFormulaExcludeFields()
		{
		fn = CustomizeFieldsControl.CustomizeFieldsControl_getFormulaExcludeFields

		sfOb = #(cols: #(), excludeFields: #())
		Assert(fn(sfOb) is: #(), msg: 'empty')

		sfOb = #(cols: #(), excludeFields: #(excludeTest))
		Assert(fn(sfOb) is: #(excludeTest), msg: 'only exclude')

		sfOb = #(cols: #(text), excludeFields: #(excludeTest))
		Assert(fn(sfOb) is: #(excludeTest), msg: 'with text')

		sfOb = #(cols: #(text, num), excludeFields: #(excludeTest))
		Assert(fn(sfOb) is: #(excludeTest), msg: 'no num')

		fieldKey = .TempName() $ '_num'
		.MakeLibraryRecord([name: 'Field_' $ fieldKey,
			text: `Field_num { Control: (Key) }`])

		fieldId = .TempName() $ '_num_renamed'
		.MakeLibraryRecord([name: 'Field_' $ fieldId,
			text: `Field_num { Control: (Id) }`])

		fieldId2 = .TempName() $ '_non_formula_field'
		.MakeLibraryRecord([name: 'Field_' $ fieldId2,
			text: `Field_string { NoFormulas: true }`])

		fieldCustom = 'custom_' $ .TempName()
		.MakeLibraryRecord([name: 'Field_' $ fieldCustom,
			text: `Field_num { Control: (Id) }`])

		fieldCustom2 = 'custom_' $ .TempName()
		.MakeLibraryRecord([name: 'Field_' $ fieldCustom2, // choosedatetime control
			text: `Field_custfield_date_modified { Prompt: "test" }`])

		sfOb = Object(cols: Object('text', 'num', 'custom_nodatadictstring',
			fieldKey, fieldId, fieldId2, fieldCustom, fieldCustom2),
			excludeFields: #(excludeTest))
		Assert(fn(sfOb).Sort!(),
			is: Object(fieldKey, fieldId, fieldId2, fieldCustom, 'excludeTest').Sort!(),
			msg: 'all test')
		}

	Test_Set_options()
		{
		m = CustomizeFieldsControl.Set_options

		m(control = [#Field, mandatory:])
		Assert(control.mandatory is: false)

		m(control = [#Info, mandatory:])
		Assert(control.mandatory is: false)
		Assert(control.allowOnlyType)

		m(control = [#ChooseMany, mandatory:])
		Assert(control.mandatory is: false)
		Assert(control.saveNone is: false)

		m(control = [#ChooseManyControl, mandatory:])
		Assert(control.mandatory is: false)
		Assert(control.saveNone is: false)

		m(control = [#RadioButtons, mandatory:])
		Assert(control.mandatory is: false)

		m(control = [#RadioButtons, mandatory:, noInitalValue:])
		Assert(control.mandatory)

		m(control = [#RadioButtons, mandatory:, noInitalValue: false])
		Assert(control.mandatory is: false)

		.testUOMOptions(m, #UOM)
		.testUOMOptions(m, #UOMControl)
		}

	testUOMOptions(m, uom)
		{
		m(control = [uom, mandatory:])
		Assert(control.mandatory is: false)

		m(control = [uom, [mandatory:], mandatory:])
		Assert(control.mandatory is: false)
		Assert(control[1].mandatory is: false)

		m(control = [uom, [mandatory:], [mandatory:], mandatory:])
		Assert(control.mandatory is: false)
		Assert(control[1].mandatory is: false)
		Assert(control[2].mandatory is: false)
		}
	}
