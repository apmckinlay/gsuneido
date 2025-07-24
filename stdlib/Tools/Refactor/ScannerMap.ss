// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
function (src, block)
	// src is a string containing source code
	// returns the concatenation of the results of calling block for each token
	{
	out = ''
	prev2 = prev = ''
	scanner = Scanner(src)
	if scanner is token = scanner.Next()
		return ''
	while scanner isnt next = scanner.Next()
		{
		out $= block(prev2, prev, token, next)
		prev2 = prev
		prev = token
		token = next
		}
	out $= block(prev2, prev, token, '')
	return out
	}