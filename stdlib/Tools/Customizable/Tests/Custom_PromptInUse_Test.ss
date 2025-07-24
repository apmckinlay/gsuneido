// Copyright (C) 2008 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Setup()
		{
		super.Setup()
		.TearDownIfTablesNotExist('customizable')
		}
	Test_main()
		{
		// speed up tests; test is slow with lots of libraries
		cl = Custom_PromptInUse
			{
			Custom_PromptInUse_libraries()
				{
				return #(stdlib, configlib)
				}
			}
		custFieldPrompt = .TempName() $ " Abc Def" // ensure we have mixed case
		table = .MakeTable('(a) key(a)')
		c = .MakeCustomField(table, 'Text, single line', custFieldPrompt)

		result = cl.PromptInUse(custFieldPrompt, c.field, exclude_custom?: false)
		Assert(result is: '')

		result = cl.PromptInUse(custFieldPrompt, c.field, exclude_custom?: true)
		Assert(result is: '')

		result = cl.PromptInUse(custFieldPrompt, .TempName(), exclude_custom?: true)
		Assert(result is: '')

		result = cl.PromptInUse(custFieldPrompt, .TempName(), exclude_custom?: false)
		Assert(result
			is: 'Prompt: ' $ custFieldPrompt $ ' is already in use by the system.')

		// case variation of original prompt should be considered duplicate
		result = cl.PromptInUse(
			custFieldPrompt.Lower(), .TempName(), exclude_custom?: false)
		Assert(result
			is: 'Prompt: ' $ custFieldPrompt $ ' is already in use by the system.')

		// get an existing field (in stdlib) and try to use that prompt
		if false isnt x = QueryFirst('stdlib
			where name > "Field_a" and group is -1 sort num')
			{
			prompt = Prompt(x.name.AfterFirst('Field_'))
			result = cl.PromptInUse(prompt,	c.field, exclude_custom?: false)
			Assert(result is: "Prompt: " $ prompt $ " is already in use by the system.")
			}
		}
	Test_prompt()
		{
		.test_prompt('Prompt', 'Prompt: "fred" ', "fred")
		.test_prompt('Prompt', 'Control: (Field)\nPrompt: "fred" ', "fred")
		.test_prompt('Prompt', "Prompt:'fred'", "fred")
		.test_prompt('Prompt', 'SelectPrompt: "fred"', false)
		.test_prompt('Prompt', 'Format: "fred"', false)

		.test_prompt('SelectPrompt',
			'Prompt: "fred"\nHeading: "fred2"\nSelectPrompt: "fred3"', "fred3")
		.test_prompt('SelectPrompt', 'Control: (Field)\nSelectPrompt: "fred" ', "fred")
		.test_prompt('SelectPrompt', "SelectPrompt:'fred'", "fred")
		.test_prompt('SelectPrompt', 'Prompt: "fred"', false)
		.test_prompt('SelectPrompt', 'Format: "fred"', false)
		}
	test_prompt(type, text, expected)
		{
		Assert(Custom_PromptInUse.Custom_PromptInUse_prompt(text, type) is: expected)
		}
	}
