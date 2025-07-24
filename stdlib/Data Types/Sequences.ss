// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Objects
	{
	HasNamed?()
		{
		return false
		}
	First()
		{
		iter = .Iter()
		if iter isnt x = iter.Next()
			return x
		return // no return value if sequence is empty
		}
	Last()
		{
		iter = .Iter()
		if iter is x = iter.Next()
			return // no return value if sequence is empty
		while iter isnt next = iter.Next()
			x = next
		return x
		}
	Empty?()
		{
		iter = .Iter()
		return iter is iter.Next()
		}
	NotEmpty?()
		{
		iter = .Iter()
		return iter isnt iter.Next()
		}
	Count(value = #(0)) // value doesn't matter, just something you wouldn't use
		{
		n = 0
		if Same?(value, #(0)) // depends on compiler sharing constants
			for unused in this
				++n
		else
			for x in this
				if x is value
					++n
		return n
		}
	Has?(value)
		{
		for x in this
			if x is value
				return true
		return false
		}
	Nth(n)
		{
		Assert(n >= 0)
		for x in this
			if n-- is 0
				return x
		return // no return value if n > size of sequence
		}
	Take(n)
		{
		return Take(this, n)
		}
	Drop(n)
		{
		return Drop(this, n)
		}
	Map(block)
		{
		return Map(this, block)
		}
	Map2(block)
		{
		return Map2(this, block)
		}
	Without(@values)
		{
		args = values.Add(this at: 0)
		return Without(@args)
		}
	Instantiate()
		{
		.Add()
		return this
		}
	}