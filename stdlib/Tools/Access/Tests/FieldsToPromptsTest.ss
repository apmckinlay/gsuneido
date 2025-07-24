// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		.MakeLibraryRecord(
			[name: 'Field_FieldsToPrompts1', text: 'Field_string { Prompt: one }'],
			[name: 'Field_FieldsToPrompts2', text: 'Field_string { Prompt: two }'],
			[name: 'Field_FieldsToPrompts3', text: 'Field_string { Prompt: three }'])
		fields = #(FieldsToPrompts1, FieldsToPrompts2, FieldsToPrompts3)
		prompts = FieldsToPrompts(fields, map = Object())
		Assert(prompts is: #('one', 'three', 'two'))
		Assert(map
			is: #(one: FieldsToPrompts1, two: FieldsToPrompts2, three: FieldsToPrompts3))
		}
	}