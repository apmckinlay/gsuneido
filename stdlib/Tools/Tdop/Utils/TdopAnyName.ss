// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
function(token)
	{
	if token.Token is TDOPTOKEN.IDENTIFIER or
		token.Token is TDOPTOKEN.STRING
		return true
	if token.Member?(#Value) and String?(token.Value) and token.Value[0].Alpha?()
		return true
	return false
	}