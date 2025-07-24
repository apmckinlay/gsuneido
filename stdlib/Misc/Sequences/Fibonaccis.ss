// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
// an infinite sequence of the fibonacci numbers
// e.g. Fibonaccis().Nth(10) or Fibonaccis().Take(10)
class
	{
	CallClass()
		{
		return Sequence(new .iterator)
		}
	iterator: class
		{
		i: 0
		prev1: 0
		prev2: 1
		Next()
			{
			switch .i++
				{
			case 0:
				return 0
			case 1:
				return 1
			default:
				n = .prev1 + .prev2
				.prev1 = .prev2
				.prev2 = n
				return n
				}
			}
		Dup()
			{
			return new (.Base())()
			}
		Infinite?()
			{
			return true
			}
		}
	}