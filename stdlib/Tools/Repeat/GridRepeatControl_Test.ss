// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_layout()
		{
		fp = GridRepeatControl.GridRepeatControl_fieldPrompt('f', 't', 'ef', 'lf')
		f = { GridRepeatControl.GridRepeatControl_layout(fp, it) }
		g = { Object('FieldPrompt', table: 't', exclude_fields: 'ef',
			fields: 'f', listField: 'lf', width: 20, name: 'headerfield' $ it) }
		Assert(f(0) is: #(Horz))
		Assert(f(1) is: ['Horz', g(0)])
		Assert(f(3) is: ['Horz', g(0), 'Skip', g(1), 'Skip', g(2)])
		}

	Test_projectFields()
		{
		f = GridRepeatControl.GridRepeatControl_projectFields
		Assert(f(0) is: #())
		Assert(f(1) is: #(headerfield0))
		Assert(f(3) is: #(headerfield0, headerfield1, headerfield2))
		}
	}
