// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
// see also: Take
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
			do
				if .iter is value = .iter.Next()
					return this
				while .n-- > 0
			return value
			}
		Dup()
			{
			return new (.Base())(.iter.Dup(), .n)
			}
		Infinite?()
			{
			return .iter.Infinite?()
			}
		}
	}
