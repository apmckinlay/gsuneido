// Copyright (C) 2015 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(s)
		{
		return Sequence(new .Iter(s))
		}
	Iter: class
		{
		New(.s)
			{
			.i = 0
			}
		Next()
			{
			// line ending consistent with Readline in file, socket, and pipes
			if .i >= .s.Size()
				return this
			i = .s.Find("\n", .i)
			line = .s[.i .. i].RightTrim('\r')
			.i = i + 1
			return line
			}
		Remainder()
			{
			return .s[.i ..]
			}
		Dup()
			{
			return new (.Base())(.s)
			}
		Infinite?()
			{
			return false
			}
		}
	}
