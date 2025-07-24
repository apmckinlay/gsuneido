// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
// Could be done by Assocs/Sort/Map
// but more efficient to call block with two arguments
class
	{
	CallClass(iterable, block)
		{
		return Sequence(new .iterator(iterable.Iter(), block))
		}
	iterator: class
		{
		i: -1
		New(.iter, .block)
			{
			}
		Next()
			{
			++.i
			if .iter is value = .iter.Next()
				return this
			return (.block)(.i, value)
			}
		Dup()
			{
			return new (.Base())(.iter.Dup(), .block)
			}
		Infinite?()
			{
			return false
			}
		}
	}
