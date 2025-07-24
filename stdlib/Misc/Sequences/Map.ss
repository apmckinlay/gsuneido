// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(iterable, block)
		{
		if iterable.HasNamed?()
			{
			ProgrammerError("shouldn't use Map with named members")
			return iterable.Copy().Map!(block)
			}
		return Sequence(new .iterator(iterable.Iter(), block))
		}
	iterator: class
		{
		New(.iter, .block)
			{
			}
		Next()
			{
			if .iter is value = .iter.Next()
				return this
			return (.block)(value)
			}
		Dup()
			{
			return new (.Base())(.iter.Dup(), .block)
			}
		Infinite?()
			{
			return .iter.Infinite?()
			}
		}
	}
