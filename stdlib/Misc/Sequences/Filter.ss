// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(iterable, block)
		{
		return Sequence(new .iterator(iterable.Iter(), block))
		}
	iterator: class
		{
		New(.iter, .block)
			{
			}
		Next()
			{
			while .iter isnt value = .iter.Next()
				if ((.block)(value))
					return value
			return this
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
