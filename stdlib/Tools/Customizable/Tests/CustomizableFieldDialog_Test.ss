// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_OK()
		{
		cfd = Mock(CustomizableFieldDialog)

		cfd.When.notChanged?([anyArgs:]).CallThrough()
		cfd.When.valid?([anyArgs:]).CallThrough()
		cfd.When.findType([anyArgs:]).CallThrough()

		cfd['CustomizableFieldDialog_types'] = CustomFieldTypes()
		cfd['CustomizableFieldDialog_colnme'] = ""
		cfd['CustomizableFieldDialog_originalData'] = #()
		cfd['Window'] = Object(Hwnd: 0)

		cfd.When.AlertError([anyArgs:]).Throw('Test tried to call Alert')
		cfd.When.On_Cancel().Return(0)
		cfd.When.beep().Return(0)
		cfd.When.promptInUse?([anyArgs:]).Return(false)
		cfd.When.promptValid?([anyArgs:]).Return(true)
		cfd.When.hasConversionFunction?([anyArgs:]).Return(false)

		cfd.When.FindControl('ctllbl').Return(FakeObject(Get: "Text, single line"))
		cfd.When.FindControl('colpro').Return(FakeObject(Get: "Entered Prompt"))

		validPeditorNoChange = FakeObject(Valid?: true, Get: #())
		validPeditorWithChange = FakeObject(Valid?: true, Get: #(width: 50))
		invalidPeditor = FakeObject(Valid?: false, Get: #())
		validTrueNoChange = FakeObject(Valid:, Dirty?: false,
			Get: Object(colpro: "Entered Prompt"))
		validTrueWithChange = FakeObject(Valid:, Dirty?: true,
			Get: Object(colpro: "Original Prompt"))
		invalidWithChange = FakeObject(Valid: false, Dirty?: true,
			Get: Object(colpro: "Original Prompt"))

		// Test OK Clicked when everything is valid, and there was no change
		cfd['Data'] = validTrueNoChange
		cfd.When.FindControl('peditor').Return(validPeditorNoChange)
		Assert(cfd.Eval(CustomizableFieldDialog.OK) is: false)
		// On_Cancel should ONLY get called when everthing valid, and no data changed
		cfd.Verify.Times(1).On_Cancel()

		// Test OK Clicked when everything is valid, but there was a data change
		cfd['Data'] = validTrueWithChange
		result = #(colpro: "Entered Prompt", options: #(), colnme: "",
			ctllbl: "Text, single line", fldbse: "Field_string_custom")
		Assert(cfd.Eval(CustomizableFieldDialog.OK) is: result)
		cfd.Verify.Times(1).On_Cancel()

		// Test OK Clicked when everything is valid, but peditor changed
		cfd.When.FindControl('peditor').Return(validPeditorWithChange)
		cfd['Data'] = validTrueNoChange
		result = #(colpro: "Entered Prompt", options: #(width: 50),
			colnme: "", ctllbl: "Text, single line", fldbse: "Field_string_custom")
		Assert(cfd.Eval(CustomizableFieldDialog.OK) is: result)
		cfd.Verify.Times(1).On_Cancel()

		// Test OK Clicked, but peditor not valid
		cfd.When.FindControl('peditor').Return(invalidPeditor)
		Assert(cfd.Eval(CustomizableFieldDialog.OK) is: false)
		cfd.Verify.Times(1).On_Cancel()

		// Test OK Clicked, but Data not valid
		cfd.When.FindControl('peditor').Return(validPeditorWithChange)
		cfd['Data'] = invalidWithChange
		Assert(cfd.Eval(CustomizableFieldDialog.OK) is: false)
		cfd.Verify.Times(1).On_Cancel()

		// Data changed, and is valid, but prompt already in use
		cfd.When.promptInUse?([anyArgs:]).Return(true)
		cfd.When.FindControl('peditor').Return(validPeditorWithChange)
		cfd['Data'] = validTrueWithChange
		Assert(cfd.Eval(CustomizableFieldDialog.OK) is: false)
		cfd.Verify.Times(1).On_Cancel()

		// Data changed, and is valid, but prompt itself is invalid
		cfd.When.promptInUse?([anyArgs:]).Return(false)
		cfd.When.promptValid?([anyArgs:]).Return(false)
		cfd.When.FindControl('peditor').Return(validPeditorWithChange)
		cfd['Data'] = validTrueWithChange
		Assert(cfd.Eval(CustomizableFieldDialog.OK) is: false)
		cfd.Verify.Times(1).On_Cancel()
		}
	}