// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(hwnd, start, end, query? = false)
		{
		// back up to start of previous line
		line = SendMessage(hwnd, SCI.LINEFROMPOSITION, start, 0)
		if start > 0 and line > 0
			--line
		start = SendMessage(hwnd, SCI.POSITIONFROMLINE, line, 0)

		// back up to outside string or comment (based on style)
		while start > 0 and line > 0
			{
			style = SendMessage(hwnd, SCI.GETSTYLEAT, start, 0)
			if style isnt 1 /*= comment */ and style isnt 3 /*= string */
				break
			start = SendMessage(hwnd, SCI.POSITIONFROMLINE, --line, 0)
			}
		.Style(hwnd, start, line, end, query?)
		}

	Style(hwnd, start, line, end, query? = false)
		{
		level = SendMessage(hwnd, SCI.GETFOLDLEVEL, line, 0) & SC.FOLDLEVELNUMBERMASK
		prev_level = level

		src = .getTextRange(hwnd, start, end)
		scan = query? ?  QueryScanner(src) :  Scanner(src)
		styles = ''
		token = ''
		do
			{
			type = scan.Next2()
			styles $= .TokenStyle(type, scan, token, src, .styles).Repeat(scan.Length())

			token = type is '' ? scan.Text() : ''
			if token is '{'
				++level
			else if token is '}'
				--level
			else if .levelEnd?(type, scan)
				{
				pos = start + scan.Position()
				result = .setFoldLevel(line, hwnd, pos, level, prev_level)
				line = result.line
				prev_level = result.prev_level
				}
			} while scan isnt type
		.setStyles(hwnd, start, styles)
		}

	levelEnd?(type, scan)
		{ return type is #NEWLINE or type is scan }

	chunk: 10000
	setStyles(hwnd, start, styles)
		{
		SendMessage(hwnd, SCI.STARTSTYLING, start, 0x1f)
		for (i = 0; i < styles.Size(); i += .chunk)
			SendMessageTextIn(hwnd, SCI.SETSTYLINGEX,
				Min(.chunk, styles.Size() - i), styles[i::.chunk])
		}

	getTextRange(hwnd, start, end)
		{
		s = ""
		for (i = start; i < end; i += .chunk)
			s $= SendMessageTextRange(hwnd, SCI.GETTEXTRANGE, i, Min(end, i + .chunk))
		return s
		}

	TokenStyle(type, scan, prev, src, styles)
		{
		if scan.Keyword?()
			return styles.KEYWORD
		switch type
			{
		case "COMMENT", "NUMBER", "STRING", "WHITESPACE" :
			return styles[type]
		case "IDENTIFIER" :
			return scan.Text() is 'it'
				? styles.KEYWORD
				: prev is '#'
					? styles.STRING
					: styles.DEFAULT
		case "ERROR":
			text = scan.Text()
			if text[0] in ("'", '"', '`')
				return styles.STRING
			if text[..2] is "/*"
				return styles.COMMENT
			return styles.DEFAULT
		default :
			return .defaultStyle(scan, src, styles)
			}
		}

	styles: (
		DEFAULT:	'\x00',
		COMMENT:	'\x01',
		NUMBER:		'\x02',
		STRING:		'\x03',
		KEYWORD:	'\x04',
		OPERATOR:	'\x05',
		WHITESPACE:	'\x06',
		)
	defaultStyle(scan, src, styles)
		{
		if scan.Text() is '#' and src[scan.Position()].Alpha?()
			return styles.STRING
		return .operators.Member?(scan.Text()) ? styles.OPERATOR : styles.DEFAULT
		}

	operators: (
		'<':,
		'<=':,
		'>':,
		'>=':,
		'not':,
		'~':,
		':':,
		'?':,
		'+=':,
		'-=':,
		'$=':,
		'*=':,
		'/=':,
		'%=':,
		'<<=':,
		'>>=':,
		'&=':,
		'|=':,
		'^=':,
		'=':,
		'++':,
		'--':,
		'+':,
		'-':,
		'$':,
		'*':,
		'/':,
		'%':,
		'<<':,
		'>>':,
		'&':,
		'|':,
		'^':,
		'is':,
		'isnt':,
		'=~':,
		'!~':,
		)
	setFoldLevel(line, hwnd, pos, level, prev_level)
		{
		i = line
		line = SendMessage(hwnd, SCI.LINEFROMPOSITION, pos, 0)
		flags = level > prev_level ? SC.FOLDLEVELHEADERFLAG : 0
		for (; i <= line; ++i)
			{
			SendMessage(hwnd, SCI.SETFOLDLEVEL, i, prev_level | flags)
			flags = 0 // header only on first
			prev_level = level
			}
		return Object(:line, :prev_level)
		}
	}
