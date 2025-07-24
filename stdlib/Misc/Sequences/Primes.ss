// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
// an infinite sequence of the prime numbers
// e.g. Primes().Nth(10) or Primes().Take(10)
class
	{
	CallClass()
		{
		return Sequence(new .iterator)
		}
	iterator: class
		{
		inprev: false
		New()
			{
			.prev = Object()
			.i = 1
			}
		Next()
			{
			forever
				if .prime?(++.i)
					return .i
			}
		prime?(n)
			{
			limit = n.Sqrt()
			for p in .prev
				if p > limit
					break
				else if n % p is 0 // divisible by previous prime
					return false
			.prev.Add(n)
			return true
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