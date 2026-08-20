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
		styles = ""
		token = ""
		anno = .newAnnoState()
		do
			{
			type = scan.Next2()
			if not query?
				.annotationScan(anno, type, scan)
			styles $= .TokenStyle(type, scan, token, src, .styles).Repeat(scan.Length())
			token = type is "" ? scan.Text() : ""
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
			token = type is "" ? scan.Text() : ""
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
		{
		return type is #NEWLINE or type is scan
		}

	chunk: 10_000
	setStyles(hwnd, start, styles)
		{
		SendMessage(hwnd, SCI.STARTSTYLING, start, 0x1f)
		for (i = 0; i < styles.Size(); i += .chunk)
			SendMessageTextIn(hwnd, SCI.SETSTYLINGEX,
				Min(.chunk, styles.Size() - i), styles[i :: .chunk])
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
		case #COMMENT, #NUMBER, #STRING, #WHITESPACE:
			return styles[type]
		case #IDENTIFIER:
			return .identStyle(scan, prev, styles)
		case #ERROR:
			return .errorStyle(scan, styles)
		default:
			return .defaultStyle(scan, src, styles)
			}
		}

	identStyle(scan, prev, styles)
		{
		if scan.Text() is #it
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
		DEFAULT:    '\x00',
		COMMENT:    '\x01',
		NUMBER:     '\x02',
		STRING:     '\x03',
		KEYWORD:    '\x04',
		OPERATOR:   '\x05',
		WHITESPACE: '\x06',
		ANNOTATION: '\x07'
		)

	// --- type annotation highlighting ---------------------------------------
	// foo(x: object, y: boolean|string) :number|other { }
	// A colon is only an annotation when the whole of `(...)` parses as a
	// parameter list, the name before `(` starts its line (or is the
	// function keyword), and a `{` follows. Anything a parameter list can
	// not hold - a literal where a name belongs, a colon with no type after
	// it, a nested call - rules the construct out, so a call taking a block
	// argument, e.g.
	//	QueryApply('table', num: x, update:)
	//		{ ... }
	// is left alone even though it has the same shape as a definition.
	newAnnoState()
		{
		return Object(stack: Object(), pos: 0, patch: Object(), p1: "", nl:,
			defName: false, want: false, colon: false, active: false,
			pp: false, ppPending: false)
		}

	annotationScan(an, type, scan)
		{
		len = scan.Length()
		pos = an.pos
		an.pos += len
		if not String?(type) or type in (#NEWLINE, #WHITESPACE, #COMMENT)
			{
			if type is #NEWLINE
				an.nl = true
			return
			}
		text = scan.Text()
		defName = an.defName // previous token could be a definition name
		prev = an.p1
		an.defName = .annoDefName?(an, type, text, scan)
		an.nl = false
		an.p1 = text
		r = Object(:pos, :len)
		if .annoRetType(an, type, text, scan, r, prev)
			return
		if .annoAfterParen(an, type, text, r)
			return
		.annoToken(an, type, text, scan, r, defName)
		}

	annoDefName?(an, type, text, scan)
		{
		if scan.Keyword?()
			return text is "function"
		return type is #IDENTIFIER and an.nl is true and an.p1 isnt '.'
		}

	// the return type after `)`, which has no parameter list to validate.
	// value keywords count as type names only inside a union, i.e. after `|`,
	// so `:object|false` is all annotation but `Foo(disabled: false)` is not
	annoRetType(an, type, text, scan, r, prev)
		{
		if an.want is #consume
			{
			if .annoType?(type, scan) or text is '|' or text is '.' or
				(prev is '|' and text in ("false", "true"))
				{
				an.active.Add(r)
				return true
				}
			an.want = false
			}
		else if an.want is #type
			{
			if .annoType?(type, scan)
				{
				an.active.Add(an.colon, r)
				an.want = #consume
				return true
				}
			an.want = false // colon not followed by a type name
			}
		return false
		}

	// after `)`: a definition only if `{` (optionally `: rettype`) follows,
	// and the `{` doesn't open a block argument i.e. isn't followed by `|`
	// returns true when this token is a return annotation colon
	annoAfterParen(an, type, text, r)
		{
		if an.pp is #blockcheck
			.annoBlockCheck(an, type, text)
		else if an.pp is #rettype
			if text is '{'
				an.pp = #blockcheck
			else
				{
				an.ppPending = false
				an.pp = false
				}
		else if an.pp is #brace
			if text is '{'
				an.pp = #blockcheck
			else if text is ':' and type is ""
				{
				an.pp = #rettype
				.annoColon(an, r, an.ppPending)
				return true
				}
			else
				{
				an.ppPending = false
				an.pp = false
				}
		return false
		}

	annoBlockCheck(an, type, text)
		{
		if text is '|' and type is ""
			an.ppPending = false
		else
			.annoFlush(an)
		an.pp = false
		}

	annoToken(an, type, text, scan, r, defName)
		{
		if type is "" and text is '('
			{
			.annoNested(an)
			an.stack.Add(Object(pending: Object(), colon: false, def: defName,
				st: defName is true ? #start : #dead))
			}
		else if type is "" and text is ')'
			.annoClose(an)
		else
			.annoParam(an, type, text, scan, r)
		}

	// a parameter list holds no nested call or object literal, so a `(`
	// inside one means the enclosing candidate is not a definition
	annoNested(an)
		{
		top = .annoTop(an)
		if top isnt false
			{
			top.def = false
			top.st = #dead
			}
		}

	annoClose(an)
		{
		top = .annoPop(an)
		// #start, #name, #type, #args, #value are where a list may end
		if top isnt false and top.def and top.st in (#start, #name, #type, #args, #value)
			{
			an.ppPending = top.pending
			an.pp = #brace
			}
		else
			{
			an.ppPending = false
			an.pp = false
			}
		}

	annoParam(an, type, text, scan, r)
		{
		top = .annoTop(an)
		if top is false or top.st is #dead
			return
		top.st = .annoNextSt(top, type, text, scan, r)
		if top.st is #dead
			top.def = false
		}

	// one step of the parameter list grammar, as gSuneido parses it:
	//	'(' ( '@' name | { ['.'] name [':' type {'|' type}] ['=' constant] [','] } ) ')'
	// anything outside it returns #dead, which rules out the definition
	annoNextSt(top, type, text, scan, r)
		{
		name? = .annoType?(type, scan)
		op = type is "" ? text : ""
		switch top.st
			{
		case #start: // the start of a parameter
			if name?
				return #name
			if op is '.'
				return #dot
			return op is '@' ? #at : #dead
		case #dot:
			if name?
				return #name
			return #dead
		case #at: // `@args` has to be the whole list, so `)` must follow
			if name?
				return #args
			return #dead
		case #name:
			if op is ':'
				{
				top.colon = r
				return #colon
				}
			if op is '='
				return #eq
			return op is ',' ? #start : #dead
		case #colon: // the type name after `name:`
			if not name?
				return #dead
			top.pending.Add(top.colon, r)
			return #type
		case #type:
			if op is '|'
				{
				top.pending.Add(r)
				return #union
				}
			if op is ','
				return #start
			return op is '=' ? #eq : #dead
		case #union: // value keywords are type names inside a union
			if not name? and not .annoValueWord?(type, text)
				return #dead
			top.pending.Add(r)
			return #type
		case #eq: // a default value, which has to be a constant
			if type in (#NUMBER, #STRING) or .annoValueWord?(type, text)
				return #value
			return op in ('-', '+', '#') ? #eq : #dead
		case #value:
			return op is ',' ? #start : #dead
		default:
			return #dead
			}
		}

	annoValueWord?(type, text)
		{
		return type is #IDENTIFIER and text in ("false", "true")
		// a type name: an identifier but not a keyword, so value keywords like
		// false / true keep their own colour rather than being read as a type
		}

	annoType?(type, scan)
		{
		return type is #IDENTIFIER and not scan.Keyword?()
		}

	annoColon(an, r, active)
		{
		an.want = #type
		an.colon = r
		an.active = active
		}

	annoFlush(an)
		{
		if an.ppPending isnt false
			an.patch.Add(@an.ppPending)
		an.ppPending = false
		}

	annoTop(an)
		{
		return an.stack.Size() is 0 ? false : an.stack[an.stack.Size() - 1]
		}

	annoPop(an)
		{
		top = .annoTop(an)
		if top isnt false
			an.stack.Delete(an.stack.Size() - 1)
		return top
		}

	applyAnnotations(styles, patch)
		{
		if patch.Size() is 0
			return styles
		patch.Sort!({|x, y| x.pos < y.pos })
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
		parts.Add(styles[last..])
		return parts.Join("")
		}

	defaultStyle(scan, src, styles)
		{
		if scan.Text() is '#' and src[scan.Position()].Alpha?()
			return styles.STRING
		return .operators.Member?(scan.Text()) ? styles.OPERATOR : styles.DEFAULT
		}

	operators: (
		'<':,
		"<=":,
		'>':,
		">=":,
		not:,
		'~':,
		':':,
		'?':,
		"+=":,
		"-=":,
		"$=":,
		"*=":,
		"/=":,
		"%=":,
		"<<=":,
		">>=":,
		"&=":,
		"|=":,
		"^=":,
		'=':,
		"++":,
		"--":,
		'+':,
		'-':,
		'$':,
		'*':,
		'/':,
		'%':,
		"<<":,
		">>":,
		'&':,
		'|':,
		'^':,
		is:,
		isnt:,
		"=~":,
		"!~":
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
