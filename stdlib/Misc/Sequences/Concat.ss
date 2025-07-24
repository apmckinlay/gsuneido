// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(@iterables)
		{
		Assert(iterables.Size() > 0)
		return Sequence(new .iterator(iterables.Map(#Iter)))
		}
	iterator: class
		{
		New(.iters)
			{
			.i = 0
			.iter = .iters[0]
			}
		Next()
			{
			while .iter is x = .iter.Next()
				{
				if ++.i >= .iters.Size()
					return this // eof
				.iter = .iters[.i] // next iter
				}
			return x
			}
		Dup()
			{
			return new (.Base())(.iters.Map(#Dup))
			}
		Infinite?()
			{
			return .iters.Any?(#Infinite?)
			}
		}
	}
