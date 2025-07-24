// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
class
	{
	ResetStyles()
		{
		for window in Suneido.GetDefault(#Persistent, #()).GetDefault(#Windows, #())
			{
			try
				window.ResetStyle()
			catch (err)
				Print(err)
			}
		}

	DefaultStyle: #(
		name: 			"default",
		comment:		0x007f00,
		number:			0x007f7f,
		string:			0x7f007f,
		keyword:		0x7f0000,
		braceGoodFore: 	0xffffff,
		braceGoodBack: 	0x80e080,
		braceBadFore:	0x000000,
		braceBadBack: 	0xa0a0ff,
		defaultFore:	0x000000,
		defaultBack:	0xffffff,
		cursorLine:		0xfef2e8,
		occurrence:		0x000044,
		lineNumberFore: 0x000000,
		lineNumberBack:	0xd5d5d5,
		foldMargin:		0xe9e9e9,
		longLineMargin:	0x000000,
		selectedBack:	0xd0d0d0,
		operator: 		0x7f0000,
		whitespace: 	0xffffff,
		warning:		0x154f91,
		error:			0x3c14dc,

		// html specific styles
		tag:			0x800000,
		unknownTag:		0x0000ff,
		attr:			0x808000,
		unknownAttr: 	0x0000ff,
		insideTag:		0x800080,
		unquotedVal:	0xff00ff,
		sgml:			0x800000,
		entity:			0xee9933,

		// changed lines
		delete:			0xc828ff,
		add:			0x00c800,
		modify:			0xff0000)

	GetTheme(option = '')
		{
		if option is ''
			option = .GetCurrent()

		themes = Plugins().Contributions('ColorSchemes', 'theme')
		chosen = themes.FindIf({ it.name is option })
		return chosen is false
			? .DefaultStyle
			: themes[chosen].Copy().MergeNew(.DefaultStyle)
		}

	GetCurrent()
		{ return IDESettings.Get(#ide_color_scheme, #default) }

	GetColor(c, option = '') // static
		{
		theme = .GetTheme(option)
		if false is color = theme.GetDefault(c, false)
			color = .DefaultStyle.GetDefault(c, false)
		return color
		}

	IsDark?()
		{
		themeBackColor = .GetColor(#defaultBack)
		return RGBColors.GetContrast(themeBackColor) is CLR.WHITE
		}
	}