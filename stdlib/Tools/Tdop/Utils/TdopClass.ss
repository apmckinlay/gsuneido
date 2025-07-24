// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
class
	{
	// CLASSDEF [CLASS, COLON, STRING, LCURLY, LIST(CONST_MEMBER), RCURLY]
	CallClass(baseClass = false, leftCurly = false, position = -1, length = 0)
		{
		children = Object()
		TdopAddChild(children, token: TdopCreateNode(TDOPTOKEN.CLASS, :position, :length))
		TdopAddChild(children, match: TDOPTOKEN.COLON)
		if baseClass is false and leftCurly is false
			baseClass = .baseClass()
		if baseClass isnt false
			if not TdopIsGlobal(baseClass.GetDefault(#Value, false))
				throw 'base class must be global defined in library'

		base = baseClass isnt false
			? TdopCreateNode(TDOPTOKEN.STRING, value: baseClass.Value,
				position: baseClass.Position, length: baseClass.Length)
			: TdopCreateNode(TDOPTOKEN.STRING, value: '')
		TdopAddChild(children, token: base)

		if leftCurly is false
			TdopAddChild(children, match: TDOPTOKEN.LCURLY, mustMatch:)
		else
			TdopAddChild(children, token: leftCurly)

		TdopAddChild(children, token: .classBody())
		TdopAddChild(children, match: TDOPTOKEN.RCURLY, mustMatch:)
		return TdopCreateNode(TDOPTOKEN.CLASSDEF, :children)
		}

	baseClass()
		{
		if _token().Match(TDOPTOKEN.LCURLY)
			return false
		t = _token()
		_advance(TDOPTOKEN.IDENTIFIER)
		return t
		}

	classBody()
		{
		return TdopCreateList()
			{ |list|
			while _token() isnt _end and not _token().Match(TDOPTOKEN.RCURLY)
				list.Add(TdopMember(TDOPTOKEN.RCURLY,
					Object(TDOPTOKEN.SEMICOLON, TDOPTOKEN.COMMA) classMember:))
			}
		}
	}
