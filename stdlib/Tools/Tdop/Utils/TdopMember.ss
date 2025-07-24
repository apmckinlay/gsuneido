// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
class
	{
	// CONST_KEYMEMBER [STRING|NUMBER, COLON, constant, COMMA|SEMICOLON]
	// CONST_MEMBER [constant, COMMA|SEMICOLON]
	CallClass(end, seps, classMember = false)
		{
		children = Object()
		sign = .checkAddSub()
		ahead = _ahead()
		memberName = .checkMemberName(ahead, sign, classMember)
		if memberName isnt false
			{
			TdopAddChild(children, token: memberName)
			TdopAddChild(children, match: TDOPTOKEN.COLON,
				mustMatch: not .isMethod(classMember, ahead))
			sign = false
			}
		if classMember is true and memberName is false
			throw "class members must be named"

		memberValue = .checkMemberValue(ahead, sign, classMember, end,
			memberName isnt false)
		TdopAddChild(children, token: memberValue)
		TdopAddChild(children, match: seps)
		type = memberName is false ? TDOPTOKEN.CONST_MEMBER : TDOPTOKEN.CONST_KEYMEMBER
		return TdopCreateNode(type, :children)
		}

	checkAddSub()
		{
		sign = false
		if _token().Match(TDOPTOKEN.SUB) or _token().Match(TDOPTOKEN.ADD)
			{
			sign = _token()
			_advance()
			if not _token().Match(TDOPTOKEN.NUMBER)
				throw 'expected ' $ TDOPTOKEN.NUMBER $ ' but got ' $ Display(_token())
			}
		return sign
		}

	checkMemberName(ahead, sign, classMember)
		{
		memberName = false
		if ahead.Match(TDOPTOKEN.COLON) or .isMethod(classMember, ahead)
			{
			memberName = _token()
			if not (TdopAnyName(memberName) or memberName.Match(TDOPTOKEN.NUMBER))
				throw 'unexpected member name ' $ Display(memberName)
			if memberName.Match(TDOPTOKEN.IDENTIFIER)
				memberName = TdopCreateNode(TDOPTOKEN.STRING, value: memberName.Value,
					position: memberName.Position, length: memberName.Length)
			if sign isnt false
				memberName = TdopCreateNode(TDOPTOKEN.UNARYOP,
					children: Object(TdopCreateNode(sign), memberName))
			_advance()
			}
		return memberName
		}

	checkMemberValue(ahead, sign, classMember, end, allowDefault)
		{
		if .isMethod(classMember, ahead)
			memberValue = TdopFunction()
		else if not _token().Match(end) and not _token().Match(TDOPTOKEN.COMMA)
			{
			memberValue = TdopConstant()
			if sign isnt false
				memberValue = TdopCreateNode(TDOPTOKEN.UNARYOP,
					children: Object(TdopCreateNode(sign), memberValue))
			}
		else if allowDefault is true
			memberValue = TdopCreateNode(TDOPTOKEN.TRUE)
		else
			throw 'unexpected member ' $ Display(_token())
		return memberValue
		}

	isMethod(classMember, ahead)
		{
		return classMember is true and ahead.Match(TDOPTOKEN.LPAREN)
		}
	}
