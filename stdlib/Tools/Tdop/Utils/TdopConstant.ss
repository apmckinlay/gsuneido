// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass()
		{
		switch (_token().Token)
			{
		case TDOPTOKEN.SUB, TDOPTOKEN.ADD:
			return .addSub()
		case TDOPTOKEN.NUMBER:
			t = _token()
			_advance()
			return t
		case TDOPTOKEN.STRING:
			return .strings()
		case TDOPTOKEN.HASH:
			return .hash()
		case TDOPTOKEN.LPAREN, TDOPTOKEN.LCURLY, TDOPTOKEN.LBRACKET:
			return TdopConstantObject()
		case TDOPTOKEN.IDENTIFIER:
			return .identifier()
		default:
			throw 'unexpected token ' $ Display(_token())
			}
		}

	addSub()
		{
		sign = _token()
		_advance()
		if not _token().Match(TDOPTOKEN.NUMBER)
			throw 'expected ' $ TDOPTOKEN.NUMBER $ ' but got ' $ Display(_token())
		return sign.Nud()
		}

	strings()
		{
		res = _token()
		_advance(TDOPTOKEN.STRING)
		while _token().Match(TDOPTOKEN.CAT) and _ahead().Match(TDOPTOKEN.STRING)
			{
			t = _token()
			_advance()
			res = t.Led(res)
			}
		return res
		}

	hash()
		{
		t = _token()
		_advance()
		return t.Nud()
		}

	identifier()
		{
		t = _token()
		_advance()
		if t.Match(TDOPTOKEN.TRUE) or t.Match(TDOPTOKEN.FALSE) or
			t.Match(TDOPTOKEN.FUNCTION) or t.Match(TDOPTOKEN.CLASS) or
			t.Match(TDOPTOKEN.DLL) or t.Match(TDOPTOKEN.CALLBACK) or
			t.Match(TDOPTOKEN.STRUCT)
			return t.Nud()
		else if _token().Match(TDOPTOKEN.LCURLY)
			return TdopClass(t)
		else
			return TdopCreateNode(TDOPTOKEN.STRING, value: t.Value,
				position: t.Position, length: t.Length)
		}
	}