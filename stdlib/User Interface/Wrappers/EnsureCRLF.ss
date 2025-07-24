// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
function (origin)
	{
	s = ''
	size = origin.Size()
	cur = 0
	while cur < origin.Size()
		{
		i = origin.Find1of('\r\n', cur)
		s $= origin[cur..i]
		if i is size
			break
		if origin[i] is '\r'
			i++
		if i isnt size and origin[i] is '\n'
			i++
		s $= '\r\n'
		cur = i
		}
	return s
	}