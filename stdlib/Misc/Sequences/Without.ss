// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(@args)
		{
		iterable = args[0]
		remove = args[1..]
		return Sequence(new .iterator(iterable.Iter(), remove))
		}
	iterator: class
		{
		New(.iter, .remove)
			{
			}
		Next()
			{
			while .iter isnt value = .iter.Next()
				if not .remove.Has?(value)
					return value
			return this
			}
		Dup()
			{
			return new (.Base())(.iter.Dup(), .remove)
			}
		Infinite?()
			{
			return .iter.Infinite?()
			}
		}
	}
