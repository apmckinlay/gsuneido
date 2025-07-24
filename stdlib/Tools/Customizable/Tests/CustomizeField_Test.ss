// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Setup()
		{
		.tranTable = .MakeTable('(test_field_a, test_field_b, test_field_c) key()')
		.MakeLibraryRecord(
			[name: 'Field_test_field',
				text: 'Field_number { Prompt: "Test Field" }'],
			[name: 'Field_test_field_string',
				text: 'Field_string { Prompt: "Test Field String" }'],
			[name: 'Field_test_field_a',
				text: 'Field_number { Prompt: "Test Field A" }'],
			[name: 'Field_test_field_b',
				text: 'Field_number { Prompt: "Test Field B" }'],
			[name: 'Field_test_field_c',
				text: 'Field_number { Prompt: "Test Field C" }'])
		.selectFields = SelectFields(QueryColumns(.tranTable), joins: false)
		}

	code1: 'FormulaAdd(Object(type: "NUMBER", value: .test_field_a),' $
		'FormulaNeg(Object(type: "NUMBER", value: .test_field_b)))'
	code2: 'FormulaCat(Object(type: "NUMBER", value: 12),' $
		'Object(type: "STRING", value: " Dollars\r\n Amount"))'
	code3: 'FormulaRate(FormulaIf({FormulaGt(Object(type: "NUMBER", value: .test_field_a),' $
		'Object(type: "NUMBER", value: .test_field_b))},' $
		'{Object(type: "NUMBER", value: .test_field_c)},{Object(type: "NUMBER", value: 1)}),' $
		'Object(type: "STRING", value: "unit"))'
	code(code, field)
		{
		return 'function()
	{
		return FormulaReturn({' $ code $ '}, "' $ field $ '")
	}'
		}
	Test_TranslateFormula_and_ValidateCode()
		{
		.testTranslate('Test Field A + -Test Field B', 'test_field',
			expected: .code(.code1, 'test_field'), fields: 'test_field_a,test_field_b')

		.testTranslate('12 $ " Dollars\r\n Amount"', 'test_field_string',
			expected: .code(.code2, 'test_field_string'), fields: '')

		.testTranslate('RATE(IF(Test Field A > Test Field B, Test Field C, 1), "unit")',
			'test_field_string',
			expected: .code(.code3, 'test_field_string'),
			fields: 'test_field_a,test_field_b,test_field_c')

		.testTranslate('Test Field A + Test Field B + * + Test Field C', 'test_field',
			expected: false)

		.testTranslate('Test Field A + Test Field E', 'test_field',
			expected: false)

		.testTranslate('rec.test_field_c', 'test_field',
			expected: #(err: "Formula: Invalid FormulaSymbol: ."))

		.testTranslate('Query1(Test Field A)', 'test_field',
			expected: #(err: 'Formula: Invalid Function Query1'))

		// test validate binary op
		.testTranslate('1 + "a"', 'test_field',
			expected: #(err:
				'Formula: Operation not supported: "<Number> + <String>"'))

		// test validate unary op
		.testTranslate('not Test Field A', 'test_field',
			expected: #(err: 'Formula: Operation not supported: "not <Number>"'))

		// test validate call
		.testTranslate('CONVERT()', 'test_field',
			expected: #(err: 'Formula: CONVERT missing arguments'))

		.testTranslate('CONVERT(1 $ " unit", "unit")', 'test_field',
			expected: #(err:
				'Formula: CONVERT Field must be a <Quantity> or <Rate>'))

		.testTranslate('RATE(IF(true, 1, "2"), "unit")', 'test_field',
			expected: #(err: 'Formula: RATE Value must be a <Number>'))

		// test validate reserved words
		.testTranslate('is', 'test_field',
			expected: #(err: 'Formula: Cannot use the reserved keyword is'))
		.testTranslate('and', 'test_field',
			expected: #(err: 'Formula: Cannot use the reserved keyword and'))

		.testTranslate('true', 'test_field',
			expected: #(err: 'Formula: Test Field cannot assign <Boolean> to <Number>'))
		.testTranslate('RATE(2, "unit")', 'test_field',
			expected: #(err: 'Formula: Test Field cannot assign <Uom_rate> to <Number>'))
		}

	testTranslate(formula, field, expected, fields = false)
		{
		result = CustomizeField.TranslateFormula(.selectFields, formula, field)
		Assert(result.formulaCode is: expected)
		if fields isnt false
			Assert(result.fields is: fields)
		Assert(CustomizeField.ValidateCode(result.formulaCode)
			is: expected is false
				? "Invalid or missing operators in formula."
				: Object?(expected)
					? expected.err
					: '')
		}

	Test_SetFormulas()
		{
		cl = CustomizeField
			{
			CustomizeField_formulasNotAllowed?() { return false }
			}
		Assert(cl.HasCustomFieldFormula?(false) is: false)

		formula = 'Test Field A + Test Field B + Test Field C'
		.MakeCustomizeField(.tranTable, 'test_field', formula)
		Customizable.ResetCustomizedCache(.tranTable)

		rec = [test_field_a: 10, test_field_b: 20, test_field_c: 30, test_field: 999]
		Assert(rec.test_field is: 999)
		Assert(cl.HasCustomFieldFormula?(.tranTable))
		cl.SetFormulas(.tranTable, rec, false)
		Assert(rec.test_field is: 999)
		rec.test_field_a = 50
		Assert(rec.test_field is: 100)

		rec = [test_field_a: 10, test_field_b: 20, test_field_c: 30, test_field: 999,
			test_field_protect: 'yes, protected']
		Assert(rec.test_field is: 999)
		cl.SetFormulas(.tranTable, rec, 'test_field_protect')
		rec.test_field_a = 50
		Assert(rec.test_field is: 999)

		rec = [test_field_a: 10, test_field_b: 20, test_field_c: 30, test_field: 999,
			test_field_protect: #(test_field:)]
		Assert(rec.test_field is: 999)
		cl.SetFormulas(.tranTable, rec, 'test_field_protect')
		rec.test_field_a = 50
		Assert(rec.test_field is: 999)

		rec = [test_field_a: 10, test_field_b: 20, test_field_c: 30, test_field: 999,
			test_field_protect: #('allbut', test_field:)]
		Assert(rec.test_field is: 999)
		cl.SetFormulas(.tranTable, rec, 'test_field_protect')
		rec.test_field_a = 50
		Assert(rec.test_field is: 100)

		Assert(rec.test_field__protect)
		rec.test_field_a = ''
		Assert(rec.test_field__protect)
		rec.test_field_b = rec.test_field_c = ''
		Assert(rec.test_field__protect is: false)
		}

	Test_encoding()
		{
		cl = CustomizeField
			{
			CustomizeField_formulasNotAllowed?() { return false }
			}
		Assert(cl.HasCustomFieldFormula?(false) is: false)

		formula = 'Test Field A + Test Field B + Test Field C'
		.MakeCustomizeField(.tranTable, 'test_field_string', formula)
		Customizable.ResetCustomizedCache(.tranTable)

		rec = [test_field_a: 10, test_field_b: 20, test_field_c: 30,
			test_field_string: 'test1']
		Assert(cl.HasCustomFieldFormula?(.tranTable))
		cl.SetFormulas(.tranTable, rec, false)
		Assert(rec.test_field_string is: 'test1')
		rec.test_field_a = 50
		Assert(rec.test_field_string is: '100')
		}

	Test_check_custom_fields()
		{
		acc = CustomizeField
			{
			CustomizeField_customized_protected?(field,
				data /*unused*/, recordCtrl /*unused*/, protectField /*unused*/)
				{
				if field is 'custom_test_readonly'
					return true
				return false
				}

			CheckCustomFields(data, customFields)
				{
				ctrl = Mock()
				ctrl.When.Get().Return(data)
				super.CheckCustomFields(customFields, ctrl, protectField: false)
				}
			}
		check = acc.CheckCustomFields
		Assert(check([], false) is: "")
		Assert(check([], #(custom_test: #(readonly:))) is: "")
		Assert(check([], #(custom_test: #(mandatory:))) is: "Required: custom_test")
		Assert(check([], #(custom_test_readonly: #(mandatory:))) is: "")
		Assert(check([custom_test: 'test'], #(custom_test: #(mandatory:))) is: "")

		valid = check([], #(custom_test: #(mandatory:), custom_test2: #(mandatory:)))
		Assert(valid has: "custom_test")
		Assert(valid has: "custom_test2")

		valid = check([custom_test2: "test2"],
			#(custom_test: #(mandatory:), custom_test2: #(mandatory:)
				custom_test_readonly: #(mandatory:)))
		Assert(valid has: "custom_test")
		Assert(valid hasnt: "custom_test2")
		Assert(valid hasnt: "custom_test_readonly")
		}
	}
