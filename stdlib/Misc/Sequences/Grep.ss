// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(iterable, regex, block = false)
		{
		return Sequence(new .iterator(iterable.Iter(), regex, block))
		}
	iterator: class
		{
		New(.iter, .regex, .block)
			{
			}
		i: -1
		Next()
			{
			do
				{
				if .iter is value = .iter.Next()
					return this
				++.i
				}
				while value !~ .regex
			return .block is false ? value : (.block)(.i, value)
			}
		Dup()
			{
			return new (.Base())(.iter.Dup(), .regex, .block)
			}
		Infinite?()
			{
			return .iter.Infinite?()
			}
		}
	}
