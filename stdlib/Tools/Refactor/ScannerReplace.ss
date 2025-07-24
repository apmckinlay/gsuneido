// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
function (text, from, to)
	{
	if not Object?(from)
		from = [from]
	if not Object?(to)
		to = [to]
	Assert(from.Size() is to.Size(),
		"ScannerReplace from and to must be the same size")
	return ScannerMap(text)
		{ |prev2/*unused*/, prev/*unused*/, token, next/*unused*/|
		if false isnt i = from.Find(token)
			token = to[i]
		token
		}
	}