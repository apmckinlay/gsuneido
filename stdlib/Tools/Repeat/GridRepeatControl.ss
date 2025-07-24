// Copyright (C) 2011 Suneido Software Corp. All rights reserved worldwide.
RepeatControl
	{
	New(table = false, exclude_fields = #(), fields = #(), listField = false,
		columns = 2)
		{
		super(.layout(.fieldPrompt(fields, table, exclude_fields, listField), columns),
			no_minus:, project: .projectFields(columns))
		}

	fieldPrompt(fields, table, exclude_fields, listField)
		{
		return Object('FieldPrompt',
			:fields, :table, :exclude_fields, :listField,
			width: 20)
		}

	fieldPrefix: 'headerfield'
	layout(fieldPrompt, columns)
		{
		layout = Object('Horz')
		for (i = 0; i < columns; i++)
			{
			if i > 0
				layout.Add('Skip')
			fp = fieldPrompt.Copy()
			fp.name = .fieldPrefix $ i
			layout.Add(fp)
			}
		return layout
		}

	projectFields(columns)
		{
		return Seq(columns).Map({ .fieldPrefix $ it })
		}
	}
