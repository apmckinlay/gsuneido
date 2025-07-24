// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(code, tokens, quickCheck)
		{
		for regex in quickCheck
			if code !~ regex
				return #()
		lines = Object()
		.ForEachMatch(code, tokens)
			{|from, to|
			from = code.LineFromPosition(from)
			to = code.LineFromPosition(to)
			lines.Add(Seq(from, to + 1))
			}
		return lines
		}

	// calls the block for each match, passing position range
	// complicated by being single pass,
	// but potentially have multiple potential matches in progress
	// e.g. searching for "ab" in "aab"
	ForEachMatch(code, tokens, block)
		{
		start = tokens[0]
		ntokens = tokens.Size()
		len = Object() // lengths of partial matches in progress
		pos = Object() // starting position of partial matches
		for (scan = FindByTokenScan.Iterator(code); scan isnt tok = scan.Next(); )
			{
			for (i = len.Size() - 1; i >= 0; --i)
				{
				if tok isnt tokens[len[i]]
					{
					len.Delete(i) // match failed
					pos.Delete(i)
					}
				else if ++len[i] is ntokens
					{
					block(pos[i], scan.Position())
					len.Delete(all:) // clear partial matches in progress
					pos.Delete(all:) // since we don't want overlapping matches
					}
				}
			if tok is start // start of new possible match
				{
				scanpos = scan.Position()
				startpos = scanpos - tok.Size()
				if ntokens is 1
					block(startpos, scanpos)
				else
					{
					len.Add(1)
					pos.Add(startpos)
					}
				}
			}
		}
	}