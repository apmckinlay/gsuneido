// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
class
	{
	New(src, .Symbols)
		{
		.scan = Scanner(src)
		.End = TdopSymbol('<end>')
		.Advance()
		.Advance()
		}

	token: false
	newline: false
	position: false
	length: false
	aheadToken: false
	aheadNewline: false
	aheadPosition: false
	aheadLength: false
	Token()
		{
		return .token
		}
	Position()
		{
		return .position
		}
	Length()
		{
		return .length
		}
	IsNewline()
		{
		return .newline
		}
	Ahead()
		{
		return .aheadToken
		}

	Advance()
		{
		.token = .aheadToken
		.newline = .aheadNewline
		.position = .aheadPosition
		.length = .aheadLength
		.advance()
		}

	skipNewline: false
	advance()
		{
		tok = .skip()
		switch tok
			{
		case #IDENTIFIER:
			id = .scan.Value()
			if .Symbols.Member?(id) and id not in ('STRING', 'NUMBER', 'IDENTIFIER')
				token = .Symbols[id].Copy()
			else
				token = (.Symbols[tok])(id)
		case #NUMBER, #STRING:
			token = (.Symbols[tok])(.scan.Text())
		case .scan:
			token = .End
		default: // operator
			id = .scan.Text()
			if not .Symbols.Member?(id)
				throw 'Unexpected Symbol ' $ id
			token = .Symbols[id].Copy()
			}
		.aheadLength = token.Length = .scan.Text().Size()
		.aheadPosition = token.Position = .scan.Position() - token.Length + 1
		.aheadToken = token
		.aheadNewline = .skipNewline
		}

	skip()
		{
		.skipNewline = false
		do
			{
			tok = .scan.Next2()
			if tok is #NEWLINE
				.skipNewline = true
			}
			while tok in (#WHITESPACE, #NEWLINE, #COMMENT)
		return tok
		}
	}