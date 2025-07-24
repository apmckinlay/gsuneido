// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
// see also: Drop
class
	{
	CallClass(iterable, n)
		{
		return Sequence(new .iterator(iterable.Iter(), n))
		}
	iterator: class
		{
		New(.iter, .n)
			{
			}
		Next()
			{
			if --.n < 0 or .iter is value = .iter.Next()
				return this
			return value
			}
		Dup()
			{
			return new (.Base())(.iter.Dup(), .n)
			}
		Infinite?()
			{
			return false
			}
		}
	}
