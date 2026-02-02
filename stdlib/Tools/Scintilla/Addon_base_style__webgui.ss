// Copyright (C) 2025 Axon Development Corporation All rights reserved worldwide.
ScintillaAddon
	{
	Name: 'basic'
	WordChars: "_0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ?!"
	Styles: (comment: (color: "comment"))
	styleTypes: (
		'comment'
		'number'
		'string'
		'keyword'
		'operator'
		'whitespace'
		'variable')

	Query: false
	Init()
		{
		.SetOption('mode', .Name)
		.SetOption('theme', .Name)
		SuRenderBackend().RecordAction(false, 'LoadCssStyles',
			['cm-theme-' $ .Name $ '.css', .buildCss(.Name)])

		defaultBack = .GetSchemeColor('defaultBack')
		.StyleSetBack(0, defaultBack)
		defaultFore = .GetSchemeColor('defaultFore')
		.SetStyleProperty('color', ToCssColor(defaultFore))
		for type in .styleTypes
			{
			if .Styles.Member?(type)
				.SetStyleProperty('--cm-theme-' $ type $ '-color',
					ToCssColor(.GetSchemeColor(.Styles[type].color)))
			}
		}

	buildCss(theme)
		{
		s = ''
		for type in .styleTypes
			{
			if not .Styles.Member?(type)
				continue
			style = .Styles[type]
			css = ''
			if style.Member?(#color)
				css $= '\tcolor: var(--cm-theme-' $ type $ '-color);\r\n'
			if style.GetDefault(#bold, false) is true
				css $= '\tfont-weight: bold;\r\n'
			s $= '.cm-s-' $ theme $ ' .cm-' $ type $ ' {\r\n' $ css $ '}\r\n'
			}
		return s
		}
	}