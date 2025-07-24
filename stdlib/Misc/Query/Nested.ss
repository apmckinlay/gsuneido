// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
class
	{
	Count(tree)
		{
		n = 0
		for x in tree
			n += Object?(x) ? 1 + .Count(x) : 1
		return n
		}
	Visit(tree, block)
		{
		.visit(0, tree, block)
		}
	visit(i, tree, block)
		{
		j = 0
		for x in tree
			{
			block(x, :i, :j, ob: tree)
			i++
			j++
			if j > 100 /*= catch loops */
				throw 'too many'
			if Object?(x)
				i = .visit(i, x, block)
			}
		return i
		}
	Random(tree)
		{
		pick = Random(.Count(tree))
		.Visit(tree)
			{|x/*unused*/, i, j, ob|
			if i is pick
				return [ob, j]
			}
		}
	FlatStr(tree)
		{
		if not Object?(tree)
			return String(tree)
		s = ""
		for x in tree
			s $= (Object?(x) ? '(' $ .FlatStr(x) $ ')' : x) $ ' '
		return s.RightTrim()
		}
	}