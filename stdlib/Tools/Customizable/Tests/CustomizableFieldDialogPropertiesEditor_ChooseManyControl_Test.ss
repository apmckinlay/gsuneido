// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_getValidMsg()
		{
		fn = CustomizableFieldDialogPropertiesEditor_ChooseManyControl
		Assert(fn.GetValidMsg(#()) is: "Should have at least two options to choose from")

		list = #(one, two, 'three,four')
		Assert(fn.GetValidMsg(list) is: "Items values can not contain commas")

		list = #(one, two, four)
		Assert(fn.GetValidMsg(list) is: "")
		}
	}