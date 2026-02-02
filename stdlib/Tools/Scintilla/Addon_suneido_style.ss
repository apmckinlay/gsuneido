// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
Addon_base_style
	{
	Name: 'suneido'
	//Styles(itemToColour: (color: "colourToColourAs", bold: bolds the itemToColour))
	Styles: (comment: (color: "comment"),
			number: (color: "number"),
			string: (color: "string"),
			keyword: (color: "keyword", bold:),
			operator: (color: "operator", bold:),
			whitespace: (color: "whitespace"))

	Init()
		{
		super.Init()
		.copyStyles = .DefaultStyles()
		}

	DefaultStyles()
		{
		return Object(
			DEFAULT:	'',
			COMMENT:	'color:' $
				ToCssColor(IDE_ColorScheme.GetColor('comment', 'default')),
			NUMBER:		'color:' $
				ToCssColor(IDE_ColorScheme.GetColor('number', 'default')),
			STRING:		'color:' $
				ToCssColor(IDE_ColorScheme.GetColor('string', 'default')),
			KEYWORD:	'font-weight:bold; color:' $
				ToCssColor(IDE_ColorScheme.GetColor('keyword', 'default')),
			OPERATOR:	'font-weight:bold; color:' $
				ToCssColor(IDE_ColorScheme.GetColor('operator', 'default')),
			WHITESPACE:	'')
		}

	On_Copy()
		{
		s = .GetSelText()
		if s.Size() is 0
			return
		ss = .BuildWithStyles(s, .copyStyles)
		ClipboardWriteHtml(ss, add?:)
		}

	wrap(s, style)
		{
		if style is ''
			return s
		return '<span style="' $ style $ '">' $ XmlEntityEncode(s) $ '</span>'
		}

	BuildWithStyles(s, styles)
		{
		s = s.Detab()
		scan = Scanner(s)
		prev = ''
		ss = `<pre style="line-height:normal;">`
		do
			{
			type = scan.Next2()
			style = ScintillaStyle.TokenStyle(type, scan, prev, s, styles)
			ss $= .wrap(scan.Text(), style)
			}
		while scan isnt type
		ss $= '</pre>'
		}
	}