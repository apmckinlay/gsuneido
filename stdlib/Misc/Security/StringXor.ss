// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
function (str, key)
	{
	keysize = key.Size()
	Assert(keysize > 0)
	if keysize is 1
		{
		k = key[0].Asc()
		return str.Map({ (it.Asc() ^ k).Chr() })
		}
	else
		{
		ko = Object()
		for c in key
			ko.Add(c.Asc())
		i = 0
		return str.Map({ (it.Asc() ^ ko[i++ % keysize]).Chr() })
		// NOTE: side effects (i++) but ok with string.Map
		}
	}