// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
class
	{
	Index: 'index_wiki'

	CallClass()
		{
		builder = Ftsearch.Create()
		QueryApply('wiki')
			{ |x|
			builder.Add(x.num, .HandleWikiWords(x.name), .HandleWikiWords(x.text))
			}
		idx = builder.Pack()
		PutFile(.Index, idx)
		return idx
		}

	wikiWord: `\<[A-Z][a-z0-9]+[A-Z][A-Za-z0-9]*\>`
	HandleWikiWords(s)
		{
		// process words that start with a capital, and
		return s.Replace(.wikiWord)
			{|s|
			// this is similar to WikiFormatTitle
			// except we don't split numbers since Ftsearch does that
			t = s.Replace('([^A-Z])([A-Z])', `\1 \2`).
				Replace("([A-Z])([A-Z][^A-Z])", `\1 \2`)
			// always include the original word as well as the parts
			// so we can search on it
			t is s ? s : s $ ' ' $ t
			}
		}
	}
