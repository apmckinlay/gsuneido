// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
function ()
	{
	if Suneido.GetDefault('wiki_regenerate', false)
		{
		Suneido[IndexWiki.Index] = Ftsearch.Load(IndexWiki())
		Suneido.wiki_regenerate = false
		}
	return true
	}
