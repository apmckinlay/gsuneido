// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Setup()
		{
		ServerSuneido.DeleteMember('SuneidoLog_MessageCount')
		.ddNoPrompt = .MakeDatadict()
		.ddEmptyPrompt = .MakeDatadict(Prompt: '')
		.ddWithPrompt = .MakeDatadict(Prompt: 'DD Prompt')

		.ddEmptySelectPrompt = .MakeDatadict(SelectPrompt: '')
		.ddWithSelectPrompt = .MakeDatadict(
			Prompt: 'DD Prompt',
			SelectPrompt: 'DD SelectPrompt')
		.ddWithOnlySelectPrompt = .MakeDatadict(SelectPrompt: 'DD Only SelectPrompt')

		.ddAllEmpty = .MakeDatadict(Prompt: '', SelectPrompt: '', Heading: '')
		.ddWithAll = .MakeDatadict(
			Prompt: 'DD Prompt',
			SelectPrompt: 'DD SelectPrompt',
			Heading: 'DD Heading')

		.ddOnlyHeading1 = .MakeDatadict(
			Prompt: '',
			SelectPrompt: '',
			Heading: 'DD Only Heading1')
		.ddOnlyHeading2 = .MakeDatadict(Heading: 'DD Only Heading2')
		.ddExcludeSelect = .MakeDatadict(
			Prompt: 'DD ExcludeSelect Prompt',
			SelectPrompt: 'DD ExcludeSelect SelectPrompt',
			ExcludeSelect:)
		.ddInternal = .MakeDatadict(Prompt: 'DD Internal', baseClass: 'Field_internal')
		.ddInternalCustom = .MakeDatadict(
			fieldName: 'custom_999999',
			baseClass: 'Field_string_custom',
			Prompt: 'DD Custom Internal',
			Internal:)
		.ddCustom = .MakeDatadict(
			fieldName: 'custom_999996',
			baseClass: 'Field_string_custom',
			Prompt: 'DD Prompt')
		.ddCustomDuplicate = .MakeDatadict(
			fieldName: 'custom_999997',
			baseClass: 'Field_string_custom',
			Prompt: 'DD Prompt')
		.ddCustomSelect = .MakeDatadict(
			fieldName: 'custom_999998',
			baseClass: 'Field_string_custom',
			Prompt: 'DD Prompt',
			SelectPrompt: 'DD Prompt ~ CustomTable')
		.ddDifferentPrompt = .MakeDatadict(
			Prompt: 'DD Prompt2',
			SelectPrompt: 'DD SelectPrompt2')
		}

	Test_main()
		{
		Assert(Datadict('hours') is: Field_hours)
		.testSummarizeField('total_hours', Field_hours, 'Total')
		.testSummarizeField("max_hours", Field_hours, 'Max')
		Assert(Datadict("hours_1") is: Field_hours)
		Assert(Datadict("dsjfdjsdkljsd") is: Field_string)
		.testSummarizeField("max_dsjfdjsdkljsd", Field_string, 'Max')
		.testSummarizeField("total_dsjfdjsdkljsd", Field_number, 'Total')
		.testSummarizeField("average_dsjfdjsdkljsd", Field_number, 'Average')
		Assert(Datadict("desc") is: Field_desc)
		.testSummarizeField("max_desc", Field_desc, 'Max')

		d = Datadict("total_desc")
		Assert(d base: Field_number)
		Assert(d.Prompt is: 'Total Description')
		d = Datadict("average_desc")
		Assert(d base: Field_number)
		Assert(d.Prompt is: 'Average Description')
		}

	testSummarizeField(field, expectedDDClass, func)
		{
		d = Datadict(field)
		Assert(d base: expectedDDClass)
		for m in #(Prompt, Heading, SelectPrompt)
			if d.Member?(m)
				Assert(d[m] is: func $ ' ' $ expectedDDClass[m])
		}

	Test_PropertyInjection()
		{
		fieldName = .TempName()
		.MakeLibraryRecord([name: "Field_" $ fieldName,
			text: `class
				{
				Control: (Field)
				Format: (Text)
				}`])
		dd = Datadict(fieldName)
		Assert(dd.Control isSize: 1)
		Assert(dd.Format isSize: 1)

		fieldName = .TempName()
		.MakeLibraryRecord([name: "Field_" $ fieldName,
			text: `class
				{
				Control: (Field)
				Control_readonly: true
				Control_width() { return 10 }
				Format: (Text)
				Format_justify: "right"
				SelectControl: (Field)
				}`])

		dd = Datadict(fieldName)
		Assert(dd.Control isSize: 3)
		Assert(dd.Control.readonly)
		Assert(dd.Control.width is: 10)
		Assert(dd.Format isSize: 2)
		Assert(dd.Format.justify is: 'right')

		fieldName2 = .TempName()
		.MakeLibraryRecord([name: "Field_" $ fieldName2,
			text: `Field_` $ fieldName $
				`{
				Control_width() { return 12 }
				Format_justify: "left"
				SelectControl_width: 14
				SelectControl_weight() { return 600 }
				}`])

		dd = Datadict(fieldName2)
		Assert(dd.Control isSize: 3)
		Assert(dd.Control.readonly)
		Assert(dd.Control.width is: 12)
		Assert(dd.Format isSize: 2)
		Assert(dd.Format.justify is: 'left')
		Assert(dd.SelectControl.width is: 14)
		Assert(dd.SelectControl.weight is: 600)

		fieldName = .TempName()
		.MakeLibraryRecord([name: "Field_" $ fieldName,
			text: `class
				{
				Control: (Field)
				Control_readonly: true
				Control_width() { return 10 }
				Format_justify: "right"
				}`])

		// exception could be uninitialized member: "Format" (c) or
		// "member not found: Format" (j)
		Assert({ Datadict(fieldName) } throws: 'Format')
		}

	Test_ddValues()
		{
		fieldName1 = .TempName()
		.MakeLibraryRecord([name: "Field_" $ fieldName1,
			text: `class
				{
				Heading: 'Test'
				Prompt: 'Test 1'
				Control: (Field)
				Control_readonly: true
				Control_width() { return 10 }
				Format: (Text)
				Format_justify: 'right'
				}`])

		fieldName2 = .TempName()
		.MakeLibraryRecord([name: "Field_" $ fieldName2,
			text: `Field_` $ fieldName1 $
				`{
				SelectPrompt: 'Test'
				Prompt: 'Test 2'
				Control_width() { return 12 }
				Format_justify: 'left'
				}`])

		ddVals = Datadict(fieldName2,
			#(Format_justify, Control_readonly, Control_width, SelectPrompt, Prompt,
				Heading, NA))
		Assert(ddVals.Format_justify is: 'left')
		Assert(ddVals.Control_readonly)
		Assert((ddVals.Control_width)() is: 12)
		Assert(ddVals.SelectPrompt is: 'Test')
		Assert(ddVals.Prompt is: 'Test 2')
		Assert(ddVals.Heading is: 'Test')
		Assert(ddVals hasntMember: 'NA')

		dd = Datadict(Display(Timestamp()), #(Prompt))
		// Despite being passed members to get, the code should simply return
		// Field_string as this field does not actually exist
		Assert(dd is: Field_string)

		err = 'Cannot get Control or Format. ' $
				'Get whole datadict to ensure any injects are included'
		Assert( { Datadict(fieldName2, #(Control)) } throws: err)
		Assert( { Datadict(fieldName2, #(Control, Format)) } throws: err)
		Assert( { Datadict(fieldName2, #(Format)) } throws: err)
		Assert( { Datadict(fieldName2, #(Prompt, Format, Heading)) } throws: err)
		}
	Test_lower_with_inject()
		{
		fieldName = .TempName()
		.MakeLibraryRecord([name: "Field_" $ fieldName,
			text: `class
				{
				Prompt: 'Test Field'
				SelectPrompt: 'Select Test Field'
				Control: (Field)
				Control_readonly: true
				}`])

		d = Datadict(fieldName $ '_lower!')
		Assert(d base: Global('Field_' $ fieldName))
		Assert(d.Prompt is: 'Test Field*')
		Assert(d.SelectPrompt is: 'Select Test Field*')
		}

	Test_Prompt()
		{
		sulog = .WatchTable('suneidolog')
		cl = Datadict { Datadict_programmerError(msg) { SuneidoLog(msg) } }
		Assert(cl.Prompt(.ddNoPrompt) is: .ddNoPrompt)
		Assert(.GetWatchTable(sulog).Last().sulog_message
			is: "no Prompt for: " $ .ddNoPrompt)
		Assert(.GetWatchTable(sulog) isSize: 1)

		Assert(cl.Prompt(.ddEmptyPrompt) is: '')
		Assert(.GetWatchTable(sulog) isSize: 1)

		Assert(cl.Prompt(.ddWithPrompt) is: "DD Prompt")
		Assert(.GetWatchTable(sulog) isSize: 1)

		Assert(cl.Prompt(.ddInternal) is: "DD Internal")
		Assert(.GetWatchTable(sulog) isSize: 2)
		Assert(.GetWatchTable(sulog).Last().sulog_message
			is: .ddInternal $ ' should have been excluded due to tag: Internal')

		Assert(cl.Prompt(.ddInternalCustom) is: "DD Custom Internal")
		Assert(.GetWatchTable(sulog) isSize: 3)
		Assert(.GetWatchTable(sulog).Last().sulog_message is: .ddInternalCustom $
			' should have been excluded due to tag: Internal')

		// should not log if custom field was deleted in last day
		ddInternalDeletedMsg = Customizable.DeletedCustomFieldMessage(.ddInternalCustom)
		SuneidoLog(ddInternalDeletedMsg)

		Assert(cl.Prompt(.ddInternalCustom) is: "DD Custom Internal")
		Assert(.GetWatchTable(sulog) isSize: 4)

		QueryDo('update suneidolog where sulog_message is ' $
			Display(ddInternalDeletedMsg) $ ' set sulog_timestamp = #20200101')
		Assert(cl.Prompt(.ddInternalCustom) is: "DD Custom Internal")
		Assert(.GetWatchTable(sulog) isSize: 5)
		}

	// SelectPrompt Priority:
	// SelectPrompt > Prompt > Heading
	Test_SelectPrompt()
		{
		sulog = .WatchTable('suneidolog')
		cl = Datadict { Datadict_programmerError(msg) { SuneidoLog(msg) } }

		// Datadicts with just Prompts
		Assert(cl.SelectPrompt(.ddNoPrompt) is: .ddNoPrompt)
		Assert(.GetWatchTable(sulog).Last().sulog_message
			is: "no SelectPrompt for: " $ .ddNoPrompt)

		Assert(cl.SelectPrompt(.ddEmptyPrompt) is: .ddEmptyPrompt)
		Assert(.GetWatchTable(sulog).Last().sulog_message
			is: "no SelectPrompt for: " $ .ddEmptyPrompt)
		Assert(.GetWatchTable(sulog) isSize: 2)

		Assert(cl.SelectPrompt(.ddWithPrompt) is: "DD Prompt")
		Assert(.GetWatchTable(sulog) isSize: 2)

		// Datadicts with SelectPrompts. SelectPrompt "" will use "" but report error.
		Assert(cl.SelectPrompt(.ddAllEmpty) is: "")
		Assert(.GetWatchTable(sulog).Last().sulog_message
			is: "no SelectPrompt for: " $ .ddAllEmpty)
		Assert(.GetWatchTable(sulog) isSize: 3)

		Assert(cl.SelectPrompt(.ddWithAll) is: "DD SelectPrompt")
		Assert(.GetWatchTable(sulog) isSize: 3)

		// Datadicts with only Heading
		// uses empty SelectPrompt rather than the heading
		Assert(cl.SelectPrompt(.ddOnlyHeading1) is: "")
		Assert(.GetWatchTable(sulog).Last().sulog_message
			is: "no SelectPrompt for: " $ .ddOnlyHeading1)
		Assert(cl.SelectPrompt(.ddOnlyHeading2) is: "DD Only Heading2")

		Assert(cl.SelectPrompt(.ddInternal) is: "DD Internal")
		Assert(.GetWatchTable(sulog).Last().sulog_message
			is: .ddInternal $ ' should have been excluded due to tag: Internal')
		Assert(.GetWatchTable(sulog) isSize: 5)

		Assert(cl.SelectPrompt(.ddExcludeSelect) is:
			"DD ExcludeSelect SelectPrompt")
		Assert(.GetWatchTable(sulog).Last().sulog_message
			is: .ddExcludeSelect $ ' should have been excluded due to tag: ExcludeSelect')
		Assert(.GetWatchTable(sulog) isSize: 6)
		}

	// Heading Priority:
	// Heading > Prompt
	Test_Heading()
		{
		sulog = .WatchTable('suneidolog')
		cl = Datadict { Datadict_programmerError(msg) { SuneidoLog(msg) } }

		// Datadicts with just Prompts
		Assert(cl.Heading(.ddNoPrompt) is: .ddNoPrompt)
		Assert(.GetWatchTable(sulog).Last().sulog_message
			is: "no Heading for: " $ .ddNoPrompt)
		Assert(.GetWatchTable(sulog) isSize: 1)

		Assert(cl.Heading('Fred') is: 'Fred')
		Assert(cl.Heading('4Fred') is: '4Fred')
		Assert(.GetWatchTable(sulog) isSize: 1) // Capital letter words treated as valid

		// Will use a prompt or heading explicitly set to "", but will log error
		Assert(cl.Heading(.ddEmptyPrompt) is: "")
		Assert(.GetWatchTable(sulog).Last().sulog_message
			is: "no Heading for: " $ .ddEmptyPrompt)
		Assert(cl.Heading(.ddAllEmpty) is: "")
		Assert(.GetWatchTable(sulog).Last().sulog_message
			is: "no Heading for: " $ .ddAllEmpty)
		Assert(.GetWatchTable(sulog) isSize: 3)

		Assert(cl.Heading(.ddOnlyHeading1) is: "DD Only Heading1")
		Assert(cl.Heading(.ddOnlyHeading2) is: "DD Only Heading2")
		Assert(.GetWatchTable(sulog) isSize: 3)

		Assert(cl.Heading(.ddInternal) is: "DD Internal")
		Assert(.GetWatchTable(sulog).Last().sulog_message
			is: .ddInternal $ ' should have been excluded due to tag: Internal')
		}

	// PromptOrHeading Priority:
	// Prompt > Heading > SelectPrompt
	Test_PromptOrHeading()
		{
		sulog = .WatchTable('suneidolog')
		cl = Datadict { Datadict_programmerError(msg) { SuneidoLog(msg) } }

		// Datadicts with just Prompts
		Assert(cl.PromptOrHeading(.ddNoPrompt) is: .ddNoPrompt)
		Assert(.GetWatchTable(sulog).Last().sulog_message
			is: "no PromptOrHeading for: " $ .ddNoPrompt)

		// Will NOT use a prompt or heading explicitly set to "", but log error
		Assert(cl.PromptOrHeading(.ddEmptyPrompt) is: .ddEmptyPrompt)
		Assert(.GetWatchTable(sulog).Last().sulog_message
			is: "no PromptOrHeading for: " $ .ddEmptyPrompt)
		Assert(cl.PromptOrHeading(.ddAllEmpty) is: "")
		Assert(.GetWatchTable(sulog).Last().sulog_message
			is: "no PromptOrHeading for: " $ .ddAllEmpty)
		Assert(.GetWatchTable(sulog) isSize: 3)

		Assert(cl.PromptOrHeading(.ddOnlyHeading1) is: "DD Only Heading1")
		Assert(cl.PromptOrHeading(.ddOnlyHeading2) is: "DD Only Heading2")
		Assert(.GetWatchTable(sulog) isSize: 3)

		Assert(cl.PromptOrHeading(.ddWithOnlySelectPrompt) is:
			"DD Only SelectPrompt")

		Assert(.GetWatchTable(sulog) isSize: 3)
		Assert(cl.PromptOrHeading(.ddInternal) is: "DD Internal")
		Assert(.GetWatchTable(sulog).Last().sulog_message
			is: .ddInternal $ ' should have been excluded due to tag: Internal')
		}

	Test_GetFieldPrompt()
		{
		sulog = .WatchTable('suneidolog')
		cl = Datadict { Datadict_programmerError(msg) { SuneidoLog(msg) } }

		Assert(cl.GetFieldPrompt(.ddNoPrompt) is: .ddNoPrompt)
		Assert(.GetWatchTable(sulog).Last().sulog_message
			is: "no SelectPrompt for: " $ .ddNoPrompt)
		Assert(.GetWatchTable(sulog) isSize: 1)

		Assert(cl.GetFieldPrompt(.ddWithSelectPrompt) is: "DD SelectPrompt")
		Assert(.GetWatchTable(sulog) isSize: 1)

		Assert(cl.GetFieldPrompt(.ddInternal) is: "DD Internal")
		Assert(.GetWatchTable(sulog).Last().sulog_message
			is: .ddInternal $ ' should have been excluded due to tag: Internal')
		Assert(.GetWatchTable(sulog) isSize: 2)

		Assert(cl.GetFieldPrompt(.ddExcludeSelect)
			is: "DD ExcludeSelect SelectPrompt")
		Assert(.GetWatchTable(sulog) isSize: 2)

		Assert(cl.GetFieldPrompt(.ddExcludeSelect, excludeTags: #(ExcludeSelect))
			is: "DD ExcludeSelect SelectPrompt")
		Assert(.GetWatchTable(sulog).Last().sulog_message
			is: .ddExcludeSelect $ ' should have been excluded due to tag: ExcludeSelect')
		}

	Test_GetPromptMap()
		{
		f = Datadict.GetPromptMap
		mock = Mock(Datadict)
		mock.When.GetFieldPrompt([anyArgs:]).CallThrough()
		mock.When.SelectPrompt([anyArgs:]).CallThrough()

		map = mock.Eval(f, #())
		Assert(map isSize: 0)
		mock.Verify.Never().suneidologOnce([anyArgs:])

		columns = Object(.ddWithPrompt)
		map = mock.Eval(f, columns)
		Assert(map isSize: 1)
		Assert(map[.ddWithPrompt] is: 'DD Prompt')
		mock.Verify.Never().suneidologOnce([anyArgs:])

		// one "standard" field, one custom field. Duplication can be fixed by
		// switching the custom field to use SelectPrompt
		columns = Object(.ddWithPrompt,	.ddCustomSelect)
		map = mock.Eval(f, columns)
		Assert(map isSize: 2)
		Assert(map[.ddWithPrompt] is: 'DD Prompt')
		Assert(map[.ddCustomSelect] is: 'DD Prompt ~ CustomTable')
		mock.Verify.Never().suneidologOnce([anyArgs:])

		// Same scenario as above. order of object passed in shouldnt affect map
		columns = Object(.ddCustomSelect, .ddWithPrompt)
		map = mock.Eval(f, columns)
		Assert(map[.ddWithPrompt] is: 'DD Prompt')
		Assert(map[.ddCustomSelect] is: 'DD Prompt ~ CustomTable')
		mock.Verify.Never().suneidologOnce([anyArgs:])

		// same scenario as above, except two custom fields.
		columns = Object(.ddCustomSelect, .ddCustom)
		map = mock.Eval(f, columns)
		Assert(map[.ddCustom] is: 'DD Prompt')
		Assert(map[.ddCustomSelect] is: 'DD Prompt ~ CustomTable')
		mock.Verify.Never().suneidologOnce([anyArgs:])

		// order of object passed in shouldnt affect map
		columns = Object(.ddCustom, .ddCustomSelect)
		map = mock.Eval(f, columns)
		Assert(map[.ddCustom] is: 'DD Prompt')
		Assert(map[.ddCustomSelect] is: 'DD Prompt ~ CustomTable')
		mock.Verify.Never().suneidologOnce([anyArgs:])

		// two duplicate custom fields. Neither custom field has a SelectPrompt.
		// Should not log this scenario until customer data has been fixed
		columns = Object(.ddCustom, .ddCustomDuplicate)
		map = mock.Eval(f, columns)
		Assert(map[.ddCustom] is: 'DD Prompt')
		Assert(map[.ddCustomDuplicate] is: 'DD Prompt')
		mock.Verify.Never().suneidologOnce([anyArgs:])

		// two standard fields with same prompt. This scenario should log
		columns = Object(.ddWithAll, .ddWithSelectPrompt)
		map = mock.Eval(f, columns)
		Assert(map[.ddWithAll] is: 'DD SelectPrompt')
		Assert(map[.ddWithSelectPrompt] is: 'DD SelectPrompt')
		mock.Verify.Times(1).suneidologOnce('ERROR: Duplicate Prompt: ' $
			.ddWithSelectPrompt $ ' & ' $ .ddWithAll)

		// one standard field with same prompt as custom field. Custom field does not
		// have SelectPrompt so the duplication can't be fixed. This scenario should log
		columns = Object(.ddCustom, .ddWithPrompt)
		map = mock.Eval(f, columns)
		Assert(map[.ddCustom] is: 'DD Prompt')
		Assert(map[.ddWithPrompt] is: 'DD Prompt')
		mock.Verify.Times(1).suneidologOnce('ERROR: Duplicate Prompt: ' $
			.ddWithPrompt $ ' & ' $ .ddCustom)

		// same as above, opposite order
		columns = Object(.ddWithPrompt, .ddCustom)
		map = mock.Eval(f, columns)
		Assert(map[.ddCustom] is: 'DD Prompt')
		Assert(map[.ddWithPrompt] is: 'DD Prompt')
		mock.Verify.Times(1).suneidologOnce('ERROR: Duplicate Prompt: ' $
			.ddCustom $ ' & ' $ .ddWithPrompt)

		mock = Mock(Datadict)
		mock.When.GetFieldPrompt([anyArgs:]).CallThrough()
		mock.When.SelectPrompt([anyArgs:]).CallThrough()
		columns = Object(.ddWithSelectPrompt, .ddDifferentPrompt, .ddCustomSelect,
			.ddCustom)
		map = mock.Eval(f, columns)
		Assert(map isSize: 4)
		Assert(map[.ddCustom] is: 'DD Prompt')
		Assert(map[.ddDifferentPrompt] is: 'DD SelectPrompt2')
		Assert(map[.ddCustomSelect] is: 'DD Prompt ~ CustomTable')
		Assert(map[.ddWithSelectPrompt] is: 'DD SelectPrompt')
		mock.Verify.Never().suneidologOnce([anyArgs:])
		}

	Teardown()
		{
		QueryDo('delete suneidolog where sulog_timestamp is #20200101')
		super.Teardown()
		}
	}
