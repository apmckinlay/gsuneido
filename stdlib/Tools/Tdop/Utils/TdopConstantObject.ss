// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
class
	{
	// OBJECT [HASH, LPAREN, LIST(CONST_MEMBER|CONST_KEYMEMBER), RPAREN]
	// RECORD [HASH, LCURLY, LIST(CONST_MEMBER|CONST_KEYMEMBER), RCURLY]
	// RECORD [HASH, LBRACKET, LIST(CONST_MEMBER|CONST_KEYMEMBER), RBRACKET]
	CallClass(position = -1, length = 0)
		{
		t = _token()
		if t.Token is TDOPTOKEN.LPAREN
			{
			type = TDOPTOKEN.OBJECT
			end = TDOPTOKEN.RPAREN
			}
		else
			{
			type = TDOPTOKEN.RECORD
			end = t.Token is TDOPTOKEN.LCURLY ? TDOPTOKEN.RCURLY : TDOPTOKEN.RBRACKET
			}
		children = TdopAddChild(token: TdopCreateNode(TDOPTOKEN.HASH, :position, :length))
		TdopAddChild(:children, match: t.Token, mustMatch:)
		TdopAddChild(:children, token: .parse(end))
		TdopAddChild(:children, match: end, mustMatch:)
		return TdopCreateNode(type, :children)
		}

	parse(end)
		{
		return TdopCreateList()
			{ |list|
			while _token().Token isnt end and _token() isnt _end
				list.Add(TdopMember(end, Object(TDOPTOKEN.COMMA, TDOPTOKEN.SEMICOLON)))
			}
		}
	}