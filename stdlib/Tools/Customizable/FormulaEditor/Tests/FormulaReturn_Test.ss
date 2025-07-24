// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		.MakeLibraryRecord(
			[name: 'Field_test_field_number_1',
				text: 'Field_number { Prompt: "Test Field 1", Control_mask: false }'],
			[name: 'Field_test_field_number_2',
				text: 'Field_number { Prompt: "Test Field 2", Control_mask: "##.#" }'],
			[name: 'Field_test_field_string',
				text: 'Field_string { Prompt: "Test Field 3" }'],
			[name: 'Field_test_field_date',
				text: 'Field_date { Prompt: "Test Field 4" }'],
			[name: 'Field_test_field_boolean',
				text: 'Field_boolean { Prompt: "Test Field 5" }'],
			[name: 'Field_test_field_datetime',
				text: 'Field_date_time { Prompt: "Test Field 6" }'],
			[name: 'Field_test_multiline_string',
				text: 'Field_string { Prompt: "Test Field 7" ' $
					'Control: (ScintillaRichWordAddonsControl) }'])

		.check("string", "test_field_string", "string")
		.check(123, "test_field_string", "123")
		.check(true, "test_field_string", "true")
		.check('', "test_field_string", "")
		.check(#20180101, "test_field_string", #20180101.ShortDate())
		.check(#20180101.010101, "test_field_string", (#20180101.010101).ShortDateTime())

		bigStr = 'A'.Repeat(555)
		.checkThrow(bigStr, 'test_field_string',
			'There is a problem calculating Test Field 3\r\n' $
			'Value exceeds max size for single line text field')
		.check(bigStr, 'test_multiline_string', bigStr)

		.check(123, "test_field_number_1", 123)
		.check('', "test_field_number_1", 0)
		.checkThrow(1/0, "test_field_number_1",
			'There is a problem calculating Test Field 1\r\n' $
				'Invalid <Number> value: inf')
		.check(1/0, "test_field_number_1", '#', isFormat:)
		.check(-1/0, "test_field_number_1", '#', isFormat:)
		.checkThrow("invalid", "test_field_number_1",
			'There is a problem calculating Test Field 1\r\n' $
				'Invalid <Number> value: "invalid"')
		.check(11.91, "test_field_number_2", 11.9)
		.checkThrow(123, "test_field_number_2",
			'There is a problem calculating Test Field 2\r\n' $
				'Invalid <Number> value: 123. Maximum digits before decimal is 2')
		.checkThrow(-12, "test_field_number_2",
			'There is a problem calculating Test Field 2\r\n' $
				'The format for this field does not support displaying -12')

		.check(#20180101, 'test_field_date', #20180101)
		.check(#20180101.235959, 'test_field_date', #20180101)
		.check(#20180101.010101, 'test_field_date', #20180101)
		.check('', 'test_field_date', '')
		.checkThrow('str', 'test_field_date',
			'There is a problem calculating Test Field 4\r\n' $
				'Invalid <Date> value: "str"')

		.check(true, 'test_field_boolean', true)
		.check('', 'test_field_boolean', false)
		.check(false, 'test_field_boolean', false)
		.checkThrow('str', 'test_field_boolean',
			'There is a problem calculating Test Field 5\r\n' $
				'Invalid <Boolean> value: "str"')

		.check(#20180101, 'test_field_datetime', #20180101)
		.check(#20180101.235959, 'test_field_datetime', #20180101.235959)
		.check(#20180101.010101, 'test_field_datetime', #20180101.010101)
		.check('', 'test_field_datetime', '')
		.checkThrow('str', 'test_field_datetime',
			'There is a problem calculating Test Field 6\r\n' $
				'Invalid <Date> value: "str"')
		}

	check(value, field, expected, isFormat = false)
		{
		Assert(FormulaReturn({ Object(:value) }, field, :isFormat) is: expected)
		}

	checkThrow(value, field, expected, isFormat = false)
		{
		fr = FormulaReturn
			{
			FormulaReturn_handleError(msg, field /*unused*/, isFormat/*unused*/)
				{
				throw msg
				}
			}
		Assert({ fr({ Object(:value) }, field, :isFormat) } throws: expected)
		}
	}
