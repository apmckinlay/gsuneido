// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
class
	{
	New(text, .wantWhitespace = false, .wantComments = false)
		{
		.setup_white()
		.scanner = Scanner(text)
		.cur = []
		.next = []
		.prev = .prev2 = []
		.Next()
		}
	setup_white()
		{
		.white = Object(WHITESPACE:, NEWLINE:, COMMENT:).Set_default(false)
		if .wantComments
			.white.Delete(#COMMENT)
		}
	Next()
		{
		.prev2 = .prev
		if not .white[.cur.type]
			.prev = .cur
		.cur = .next
		.next = Object(token: .next(),
			position: .scanner.Position(),
			keyword?: .scanner.Keyword?(),
			type: .scanner.Type())
		return .cur.token is .scanner ? this : .cur.token
		}
	next()
		{
		do
			token = .scanner.Next()
			while token isnt .scanner and .skip?()
		return token
		}
	skip?()
		{
		return .wantWhitespace ? false : .white[.scanner.Type()]
		}
	Token()
		{ return .cur.token }
	Type()
		{ return .cur.type }
	Keyword?()
		{ return .cur.keyword? }
	Position()
		{ return .cur.position }
	Ahead()
		{ return .next.token is .scanner ? '' : .next.token }
	AheadPos()
		{ return .next.token is .scanner ? '' : .next.position }
	AheadType()
		{ return .next.token is .scanner ? '' : .next.type }
	Prev()
		{ return .prev.token }
	Prev2()
		{ return .prev2.token }

	Getter_(member)
		{
		switch member.BeforeFirst('_')
			{
		case #Prev:
			return .prev[member.AfterFirst('_').Lower()]
		case #Prev2:
			return .prev2[member.AfterFirst('_').Lower()]
		default:
			return
			}
		}
	}
