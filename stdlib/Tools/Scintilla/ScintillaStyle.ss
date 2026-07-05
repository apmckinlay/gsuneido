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
		src = .getTextRange(hwnd, start, end)
		.setStyles(hwnd, start, .StyleString(src, query?))
		.foldLevels(hwnd, start, line, src, query?)
		}

	// StyleString returns the style byte for each character of src
	// It takes no hwnd so it can be unit tested and reused
	StyleString(src, query? = false)
		{
		scan = query? ? QueryScanner(src) : Scanner(src)
		styles = ''
		token = ''
		anno = .newAnnoState()
		do
			{
			type = scan.Next2()
			if not query?
				.annotationScan(anno, type, scan)
			styles $= .TokenStyle(type, scan, token, src, .styles).Repeat(scan.Length())
			token = type is '' ? scan.Text() : ''
			} while scan isnt type
		if not query?
			styles = .applyAnnotations(styles, anno.patch)
		return styles
		}

	foldLevels(hwnd, start, line, src, query?)
		{
		level = SendMessage(hwnd, SCI.GETFOLDLEVEL, line, 0) & SC.FOLDLEVELNUMBERMASK
		prev_level = level
		scan = query? ? QueryScanner(src) : Scanner(src)
		do
			{
			type = scan.Next2()
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
			return .identStyle(scan, prev, styles)
		case "ERROR" :
			return .errorStyle(scan, styles)
		default :
			return .defaultStyle(scan, src, styles)
			}
		}

	identStyle(scan, prev, styles)
		{
		if scan.Text() is 'it'
			return styles.KEYWORD
		return prev is '#' ? styles.STRING : styles.DEFAULT
		}

	errorStyle(scan, styles)
		{
		text = scan.Text()
		if text[0] in ("'", '"', '`')
			return styles.STRING
		if text[..2] is "/*"
			return styles.COMMENT
		return styles.DEFAULT
		}

	styles: (
		DEFAULT:	'\x00',
		COMMENT:	'\x01',
		NUMBER:		'\x02',
		STRING:		'\x03',
		KEYWORD:	'\x04',
		OPERATOR:	'\x05',
		WHITESPACE:	'\x06',
		ANNOTATION:	'\x07',
		)

	// --- type annotation highlighting ---------------------------------------
	// foo(x: object, y: boolean|string) :number|other { }


	// A colon is only an annotation when it is in parameter position inside
	// `(...)` (or return position after `)`) is followed by a type name and
	// the list is followed by `{`
	newAnnoState()
		{
		return Object(stack: Object(), pos: 0, patch: Object(), p1: '', curWasDot: false,
			want: false, colon: false, active: false, pp: false, ppPending: false)
		}

	annotationScan(an, type, scan)
		{
		len = scan.Length()
		pos = an.pos
		an.pos += len
		if not String?(type) or type in (#NEWLINE, #WHITESPACE, #COMMENT)
			return
		text = scan.Text()
		memberCall = an.curWasDot // previous token (a name) was a `.member` access
		an.curWasDot = an.p1 is '.'
		an.p1 = text
		r = Object(:pos, :len)
		if .annoWantType(an, type, text, scan, r)
			return
		if .annoAfterParen(an, type, text, r)
			return
		.annoToken(an, type, text, r, memberCall)
		}

	annoWantType(an, type, text, scan, r)
		{
		if an.want is #consume
			{
			if .annoType?(type, scan) or text is '|' or text is '.'
				{ an.active.Add(r); return true }
			an.want = false
			}
		else if an.want is #type
			{
			if .annoType?(type, scan)
				{ an.active.Add(an.colon, r); an.want = #consume; return true }
			an.want = false // colon not followed by a type name
			}
		return false
		}

	// after `)`: a definition only if `{` (optionally `: rettype`) follows
	// returns true when this token is a return annotation colon
	annoAfterParen(an, type, text, r)
		{
		if an.pp is #rettype
			{
			if text is '{'
				.annoFlush(an)
			else
				an.ppPending = false
			an.pp = false
			}
		else if an.pp is #brace
			{
			if text is '{'
				{ .annoFlush(an); an.pp = false }
			else if text is ':' and type is ''
				{ an.pp = #rettype; .annoColon(an, r, an.ppPending); return true }
			else
				{ an.ppPending = false; an.pp = false }
			}
		return false
		}

	annoToken(an, type, text, r, memberCall)
		{
		if type isnt ''
			.annoStep(an, type, text)
		else if text is '('
			an.stack.Add(Object(pending: Object(), slot: #start, def: not memberCall))
		else if text is ')'
			.annoClose(an)
		else if text is ':' and .annoParamColon?(an)
			{ .annoColon(an, r, .annoPending(an)); .annoSetSlot(an, #other) }
		else
			.annoStep(an, type, text)
		}

	annoClose(an)
		{
		top = .annoPop(an)
		if top isnt false and top.def
			{ an.ppPending = top.pending; an.pp = #brace }
		else
			{ an.ppPending = false; an.pp = false }
		}

	annoParamColon?(an)
		{ return .annoSlot(an) is #name and .annoTop(an).def }

	// a type name: an identifier but not a keyword, so value keywords like
	// false / true keep their own colour rather than being read as a type
	annoType?(type, scan)
		{ return type is #IDENTIFIER and not scan.Keyword?() }

	annoColon(an, r, active)
		{ an.want = #type; an.colon = r; an.active = active }

	annoFlush(an)
		{
		if an.ppPending isnt false
			an.patch.Add(@an.ppPending)
		an.ppPending = false
		}

	annoTop(an)
		{ return an.stack.Size() is 0 ? false : an.stack[an.stack.Size() - 1] }

	annoPop(an)
		{
		top = .annoTop(an)
		if top isnt false
			an.stack.Delete(an.stack.Size() - 1)
		return top
		}

	annoSlot(an)
		{ top = .annoTop(an); return top is false ? #none : top.slot }

	annoPending(an)
		{ return .annoTop(an).pending }

	annoSetSlot(an, slot)
		{ top = .annoTop(an); if top isnt false top.slot = slot }

	annoStep(an, type, text)
		{
		top = .annoTop(an)
		if top is false
			return
		if text is ',' and type is ''
			top.slot = #start
		else if top.slot isnt #start
			top.slot = #other
		else
			top.slot = .annoStartSlot(type, text)
		}

	// at a parameter-name position: a name (or its . / @ prefix) stays here,
	// anything else moves on
	annoStartSlot(type, text)
		{
		if type is #IDENTIFIER
			return #name
		if type is '' and (text is '.' or text is '@')
			return #start
		return #other
		}

	applyAnnotations(styles, patch)
		{
		if patch.Size() is 0
			return styles
		patch.Sort!({ |x, y| x.pos < y.pos })
		anno = .styles.ANNOTATION
		parts = Object()
		last = 0
		for r in patch
			{
			if r.pos < last
				continue
			parts.Add(styles[last .. r.pos], anno.Repeat(r.len))
			last = r.pos + r.len
			}
		parts.Add(styles[last ..])
		return parts.Join('')
		}
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
