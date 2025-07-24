// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
ScintillaAddon
	{
	WordChars: "_0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ?!"
	Styles: (comment: (color: "comment"))
	styleMap: (
		1: 'comment'
		2: 'number'
		3: 'string'
		4: 'keyword'
		5: 'operator'
		6: 'whitespace')

	Init()
		{
		.SetWordChars(.WordChars)
		.SetLexer(SCLEX.CONTAINER)

		// base font color (fore/back)
		defaultBack = .GetSchemeColor('defaultBack')
		defaultFore = .GetSchemeColor('defaultFore')

		.DefineStyle(0, defaultFore, back: defaultBack)
		.DefineStyle(SC.STYLE_DEFAULT, defaultFore, back: defaultBack)
		.DefineStyle(SC.STYLE_LINENUMBER, .GetSchemeColor('lineNumberFore'),
			back: .GetSchemeColor('lineNumberBack'))
		.SetCaretFore(defaultFore)
		.SetCaretWidth(ScaleWithDpiFactor(1.4/*=width*/))
		.SetElementColour(SC.ELEMENT_SELECTION_BACK,
			.GetSchemeColor('selectedBack'))
		.SetElementColour(SC.ELEMENT_SELECTION_INACTIVE_BACK,
			.GetSchemeColor('selectedBack'))

		for i in .styleMap.Members()
			{
			style = .styleMap[i]
			if .Styles.Member?(style)
				{
				.DefineStyle(i, .GetSchemeColor(.Styles[style].color),
					back: defaultBack, bold: .Styles[style].GetDefault("bold", false))
				}
			else
				.DefineStyle(i, defaultFore, back: defaultBack)
			}
		}
	Query: false
	Style(from, to)
		{
		ScintillaStyle(.Hwnd, from, to, .Query)
		}
	}