// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
// some simple code formatting
// e.g. converts 'for(i=0;i<10;++i)' to 'for (i = 0; i < 10; ++i)'
// does NOT do anything with indenting or line breaks
// removes trailing whitespace and semicolons on lines
// ensures a space after keywords, semicolons, and commas
// ensures a space before/after binary operators
class
	{
	map: ('==': 'is', '!=': 'isnt', '!': 'not', '&&': 'and', '||': 'or')
	spaceAfter: (break:, catch:, class:, continue:, for:, forever:,
		if:, return:, switch:, try:, while:, not:, ';':, ',':)
	binaryOps: ('and':, '&':, '&=':, '|':, '|=':, '^':, '^=',
		'$':, '$=':, '/':, '/=':, '=':, '>':, '>=':, 'is':, 'isnt':,
		'<<':, '<<=':, '<':, '<=':, '=~':, '-=':, '%':, '%=':,
		'*':, '*=':, '!~':, 'or':, '+=':, '?':, '>>':, '>>=':)
		// '+', '-', ':' omitted since not always binary operators
	CallClass(text)
		{
		text = text.Replace('[ \t]+$', '')
		s = ""
		state = Object(parensNest: 0, insideBlockArgs: false)
		scan = ScannerWithContext(text, wantWhitespace:)
		formats = Object(.formatBlock, .removeSemicolon, .addSpaces, .nop)
		while scan isnt token = scan.Next()
			{
			token = .map.GetDefault(token, token)
			if token is '('
				++state.parensNest
			else if token is ')'
				--state.parensNest

			for f in formats
				if false isnt result = f(token, s, :scan, :state)
					{
					s = result
					break
					}
			}

		// strip trailing whitespace on lines
		s = s.Replace('[ \t]+$', '')

		// convert {...} to put curly braces on separate lines
		s = s.Replace('^(\t+)({(\|.*\|)?)[ \t]*(\S.*?)[ \t]*}$', '\1\2\r\n\1\4\r\n\1}')

		// ensure single trailing newline
		s = s.RightTrim('\r\n') $ '\r\n'

		return s
		}

	formatBlock(token, s, scan, state)
		{
		if token is '|' and scan.Prev() is '{'
			{
			state.insideBlockArgs = true
			while scan.Ahead().White?()
				scan.Next()
			return s.RightTrim(' \t') $ token
			}
		if token is '|' and state.insideBlockArgs
			{
			state.insideBlockArgs = false
			return s.RightTrim(' \t') $ token
			}
		return false
		}

	removeSemicolon(token, s, scan, state)
		{
		if .semicolonAtEnd?(token, s, scan, state.parensNest)
			return s // omit ';'
		if .semicolonAtEndOfBlock?(token, s, scan)
			return s.RightTrim(' \t;') $ ' }'
		if .semicolonBeforeComment?(token, s, scan)
			// omit semicolon at end of line before comment
			return s.BeforeLast(';') $ s.AfterLast(';') $ token
		return false
		}

	semicolonBeforeComment?(token, s, scan)
		{
		return token.Prefix?('//') and scan.Prev() is ';' and
			s.AfterLast('\n').Has1of?('^ \t;')
		}

	semicolonAtEnd?(token, s, scan, parensNest)
		{
		return token is ';' and scan.Ahead().Has?('\n') and
			not s.AfterLast('\n').White?() and parensNest is 0
		}

	semicolonAtEndOfBlock?(token, s, scan)
		{
		return token is '}' and scan.Prev() is ';' and s !~ '\n\s*\Z'
		}

	missingSpaceAfterKeyword(token, scan)
		{
		return not scan.Ahead().White?() and .spaceAfter.Member?(token) and
			scan.Ahead() isnt '' // end
		}

	addSpaces(token, s, scan)
		{
		if .missingSpaceAfterKeyword(token, scan)
			return s $ token $ ' '
		if .binaryOps.Member?(token)
			return .ensureSpacesOnBinaryOps(token, s, scan)
		return false
		}

	ensureSpacesOnBinaryOps(token, s, scan)
		{
		if not s[-1].White?()
			s $= ' '
		s $= token
		if not scan.Ahead().White?() and scan.Ahead() isnt ':'
			s $= ' '
		return s
		}

	nop(token, s)
		{
		return s $ token
		}
	}
