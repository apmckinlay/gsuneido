// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_promptFromInfo()
		{
		fn = Prompt.Prompt_promptFromInfo
		Assert(fn([prefix: '', suffix: '', baseField: '',
			promptMethod: 'Prompt_Test.PromptMethod']), is: '')
		Assert(fn([prefix: '', suffix: '', baseField: 'Test Prompt',
			promptMethod: 'Prompt_Test.PromptMethod']), is: 'Test Prompt')
		Assert(fn([prefix: 'Prefix', suffix: '', baseField: 'Test Prompt',
			promptMethod: 'Prompt_Test.PromptMethod']), is: 'Prefix Test Prompt')
		Assert(fn([prefix: '', suffix: 'Suffix', baseField: 'Test Prompt',
			promptMethod: 'Prompt_Test.PromptMethod']), is: 'Test Prompt Suffix')
		Assert(fn([prefix: 'Prefix', suffix: 'Suffix', baseField: 'Test Prompt',
			promptMethod: 'Prompt_Test.PromptMethod']), is: 'Prefix Test Prompt Suffix')
		Assert(fn([prefix: 'Prefix', suffix: 'Suffix', baseField: '',
			promptMethod: 'Prompt_Test.PromptMethod']), is: 'Prefix Suffix')
		Assert(fn([prefix: 'Prefix', suffix: '', baseField: '',
			promptMethod: 'Prompt_Test.PromptMethod']), is: 'Prefix')
		Assert(fn([prefix: '', suffix: 'Suffix', baseField: '',
			promptMethod: 'Prompt_Test.PromptMethod']), is: 'Suffix')
		}

	PromptMethod(field)
		{
		return field
		}
	}