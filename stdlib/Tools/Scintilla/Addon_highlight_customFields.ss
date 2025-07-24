// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
ScintillaAddon
	{
	customizable: false
	field_dicts: false
	text: false
	matches: false

	New(@args)
		{
		super(@args)
		.customizable = args[1].GetDefault('customizable', false)
		.field_dicts = args[1].GetDefault('field_dicts', false)
		.useMarkdown = args[1].GetDefault('useMarkdown', false)
		if .field_dicts isnt false
			.field_dicts = .field_dicts.Copy().Map!(SelectFields.GetFieldPrompt)
		}

	styleLevel: 40
	Init()
		{
		.indic_customField = .IndicatorIdx(level: .styleLevel)
		}

	Styling()
		{
		return [[level: .styleLevel, indicator: [INDIC.ROUNDBOX, fore: CLR.Highlight]]]
		}

	UpdateUI()
		{
		.reset()
		if .customizable isnt false
			.markOccurences(SelectFields(.customizable.CustomFields()).Prompts())
		else if .field_dicts isnt false
			.markOccurences(.field_dicts)
		}

	markOccurences(fields)
		{
		for field in fields
			{
			if .useMarkdown isnt false
				field = '<' $ field $ '>'
			.setIndicator(field, .text, .matches, .setIndic, .indic_customField)
			}
		}

	reset()
		{
		.ClearIndicator(.indic_customField)
		.matches = Object().Set_default(Object())
		.text = .Get()
		}

	setIndicator(prompt, text, matches, setIndic, indic)
		{
		text.ForEachMatch("(?q)" $ prompt)
			{ |m|
			m = m[0]
			cur = matches[prompt]
			cur.Add(m)
			setIndic(indic, m[0], m[1])
			}
		}

	setIndic(indic, pos, len)
		{
		.SetIndicator(indic, pos, len)
		}
	}
