// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_getValidMsg()
		{
		fn = CustomizableFieldDialogPropertiesEditor_ChooseListControl
		Assert(fn.GetValidMsg(#()) is: "Should have at least two options to choose from")

		list = Object()
		for (i = 0; i <= 301; i++)
			list.Add(String(i))

		Assert(fn.GetValidMsg(list) is: 'Cannot have more than 300 items\n\n' $
			'Use the "Text, from custom table" type if you need more items\n\n' $
			'Please contact Axon for assistance')

		list.Delete(0, 1)
		Assert(fn.GetValidMsg(list) is: "")
		}
	}