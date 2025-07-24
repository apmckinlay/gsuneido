// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
// For efficiency, the block can combine map and filter
// To exclude an element, the block can return Nothing()
class
	{
	CallClass(s, block = function (scan) { scan.Text() })
		{
		return Sequence(new .iterator(s, block))
		}
	iterator: class
		{
		New(.s, .block)
			{
			.scan = Scanner(s)
			}
		Next()
			{
			while .scan isnt .scan.Next2()
				// if block has no return value we keep looping
				// need the assignment so it doesn't propagate the lack of return value
				try return unused = (.block)(.scan)
			return this
			}
		Dup()
			{
			return new (.Base())(.s, .block)
			}
		Infinite?()
			{
			return false
			}
		}
	}