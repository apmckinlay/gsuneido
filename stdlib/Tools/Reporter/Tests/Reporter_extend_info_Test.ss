// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_CallClass()
		{
		infoFields = Object()
		Assert(Reporter_extend_info('Email', infoFields) is: '')

		infoFields = Object('Work: 111-111-1111', 'Email: test1@test.com',
			'', '', '', 'Email: test2@test.com')
		Assert(Reporter_extend_info('Work', infoFields) is: '111-111-1111')
		Assert(Reporter_extend_info('Email', infoFields)
			is: 'test1@test.com, test2@test.com')
		Assert(Reporter_extend_info('Fax', infoFields) is: '')

		infoFields = Object('Email: test1@test.com, test2@test.com',
			'Email: test3@test.com')
		Assert(Reporter_extend_info('Email', infoFields)
			is: 'test1@test.com, test2@test.com, test3@test.com')
		}

	Test_Extend()
		{
		infoFields = Object()
		fields = Object()
		Assert(Reporter_extend_info.Extend(infoFields, fields) is: "")

		infoFields = Object(
			Object(field: 'reporter_info_testing_Work', prefix: 'testing', type: 'Work',
				prompt: 'Test Work', deps: Object('testing_info1', 'testing_info2'))
			Object(field: 'reporter_info_testing1_Work', prefix: 'testing1', type: 'Work',
				prompt: 'Test 1 Work', deps: Object('testing1_info3', 'testing1_info4'))
			Object(field: 'reporter_info_testing_Email', prefix: 'testing', type: 'Email'
				prompt: 'Test Email', deps: Object('testing_info1', 'testing_info2'))
			Object(field: 'reporter_info_testing1_Email', prefix: 'testing1',
				type: 'Email', prompt: 'Test 1 Email',
				deps: Object('testing1_info3', 'testing1_info4')))

		fields = Object()
		Assert(Reporter_extend_info.Extend(infoFields, fields) is: "")

		fields = Object('reporter_info_testing_Email')
		Assert(Reporter_extend_info.Extend(infoFields, fields) is:
			'\nextend reporter_info_testing_Email = Reporter_extend_info("Email", ' $
			'Object(testing_info1, testing_info2))')

		fields = Object('reporter_info_testing_Email', 'reporter_info_testing1_Work')
		Assert(Reporter_extend_info.Extend(infoFields, fields) is:
			'\nextend ' $
			'reporter_info_testing1_Work = Reporter_extend_info("Work", ' $
			'Object(testing1_info3, testing1_info4)),\n' $
			'reporter_info_testing_Email = Reporter_extend_info("Email", ' $
			'Object(testing_info1, testing_info2))'
			)
		}

	Test_AddFields()
		{
		sf = SelectFields(#())
		Assert(Reporter_extend_info.AddFields(sf) is: Object())
		Assert(sf.Fields is: Object())

		.test("name", "testing", "testing")
		.test("num", "testing1", "testing1")
		.test("num", "Name", "")
		.test("num", "Date/Time Created", "")
		.test("num", "Fred Date/Time Created", "Fred")
		}

	fields:
		#(
			#('testing_info1', 'info1')
			#('testing_info2', 'info2')
			#('testing_info3', 'info3')
			#('testing_info4', 'info4')
			#('testing_info5', 'info5')
		)

	test(suffix, prompt, expectedPrompt)
		{
		super.Teardown()
		.MakeLibraryRecord([name: "Field_testing_" $ suffix,
			text: `Field_string { Prompt: "` $ prompt $ `"}`])
		sf = SelectFields(#())
		.fields.Each()
			{
			sf.AddField(@it)
			}
		infoFields = Reporter_extend_info.AddFields(sf)
		expectedFields = Object()
		expectedPrompts = Object()
		InfoTypes.Each()
			{
			expectedFields.Add("reporter_info_testing_" $ it[..-1])
			expectedPrompts.Add(Opt(expectedPrompt, ' ') $ it[..-1])
			}
		Assert(infoFields.Map({ it.field }) is: expectedFields)
		Assert(infoFields.Map({ it.prompt }) is: expectedPrompts)
		Assert(infoFields[0].deps equalsSet: .fields.Map({ it[0] }))
		Assert(sf.Fields.Values()
			equalsSet: .fields.Map({ it[0] }).Append(expectedFields))
		}
	}