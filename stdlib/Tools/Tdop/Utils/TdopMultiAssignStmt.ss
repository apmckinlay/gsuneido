// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
function (first)
	{
	children = Object()
	elems = TdopCreateList()
		{ |list|
		elemChildren = Object()
		TdopAddChild(elemChildren, token: first)
		TdopAddChild(elemChildren, match: TDOPTOKEN.COMMA)
		list.Add(TdopCreateNode(TDOPTOKEN.EXPR_ELEM, children: elemChildren))

		while _token() isnt _end and not _token().Match(TDOPTOKEN.EQ)
			{
			elemChildren = Object()
			TdopAddChild(elemChildren, match: TDOPTOKEN.IDENTIFIER, mustMatch:)
			TdopAddChild(elemChildren, match: TDOPTOKEN.COMMA)
			list.Add(TdopCreateNode(TDOPTOKEN.EXPR_ELEM, children: elemChildren))
			}
		}

	TdopAddChild(children, token:  elems)
	TdopAddChild(children, match: TDOPTOKEN.EQ, mustMatch:)

	call = TdopStmtExpr()
	Assert(call.Match(TDOPTOKEN.CALL))
	TdopAddChild(children, token: call)
	return TdopCreateNode(TDOPTOKEN.MULTIASSIGNSTMT, :children)
	}