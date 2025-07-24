// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		Assert(CustomPromptValid("Test Prompt") is: "")
		Assert(CustomPromptValid("Test Prompt2") is: "")
		Assert(CustomPromptValid("Test Prompt?") is: "")
		Assert(CustomPromptValid("Test/Prompt") is: "")
		Assert(CustomPromptValid("TestPrompt#") is: "")
		Assert(CustomPromptValid(" Test Prompt2 ") is: "")
		Assert(CustomPromptValid("", true) is: "")
		Assert(CustomPromptValid(" ", true)
			is: "Custom Prompt must contain more than just blank spaces.")
		Assert(CustomPromptValid("") is: "Custom Prompt can not be empty.")
		Assert(CustomPromptValid("    ")
			is: "Custom Prompt must contain more than just blank spaces.")
		Assert(CustomPromptValid("1234")
			is: "Custom Prompt must contain at least one alpha character.")
		Assert(CustomPromptValid(" 12 34 ")
			is: "Custom Prompt must contain at least one alpha character.")
		Assert(CustomPromptValid("Test*")
			is: "Custom Prompt can not contain special characters.")
		Assert(CustomPromptValid("Test Prompt:")
			is: "Custom Prompt can not contain special characters.")
		Assert(CustomPromptValid("Test Prompt:", description: 'Role')
			is: "Role can not contain special characters.")
		Assert(CustomPromptValid("a".Repeat(50)) is: "")
		Assert(CustomPromptValid("a".Repeat(51))
			is: 'Custom Prompt must be less than 50 characters. \r\n\r\n' $
			'Please consider using "Tooltips" or ' $
			'typing static text in the layouts directly.')
		}
	}