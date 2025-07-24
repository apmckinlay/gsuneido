// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
// Normally accessed by: number.Of(block)
class
	{
	CallClass(n, block)
		{
		return Sequence(new .iterator(n, block))
		}
	iterator: class
		{
		New(.n, .block)
			{
			}
		Next()
			{
			return --.n >= 0 ? (.block)() : this
			}
		Dup()
			{
			return new (.Base())(.n, .block)
			}
		Infinite?()
			{
			return false
			}
		}
	}