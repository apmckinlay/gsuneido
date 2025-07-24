// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
// NOTE: This is primarily an example of how to interpret coverage results
function (src, cover)
	{
	s = ""
	counts = cover.Any?(Number?)
	for (i = 0; i < src.Size(); ++i)
		{
		j = src.Find('\n', i) + 1
		for (k = i; k < j; k++)
			if cover.Member?(k)
				break
		if counts
			pre = k < j ? cover[k].Pad(6, ' ') $ '  ' : "\t\t" /*= up to 6 digits */
		else
			pre = k < j ? "*\t" : "\t"
		if k < j
			cover.Delete(k)
		s $= pre $ src[i..j]
		i = j
		}
	Assert(cover is: #())
	return s
	}