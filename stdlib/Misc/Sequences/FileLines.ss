// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(file, block)
		{
		File(file)
			{|f|
			return block(Sequence(new .iterator(f)))
			}
		}
	iterator: class
		{
		New(.f)
			{
			}
		Next()
			{
			return (false is line = .f.Readline()) ? this : line
			}
		Dup()
			{
			return new (.Base())(.f)
			}
		Infinite?()
			{
			// return true to prevent instantiation
			// since the point is to avoid reading the whole file into memory
			return true
			}
		}
	}