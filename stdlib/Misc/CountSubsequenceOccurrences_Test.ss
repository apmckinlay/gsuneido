// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		fn = CountSubsequenceOccurrences

		Assert(fn("", "a") is: 0)
		Assert(fn("abc", "") is: 1)
		Assert(fn("abc", "d") is: 0)
		Assert(fn("abc", "abc") is: 1)
		Assert(fn("abcdef", "ace") is: 1)
		Assert(fn("abcabc", "abc") is: 4)
		Assert(fn("aaaa", "aa") is: 6)
		Assert(fn("aaaa", "aaa") is: 4)
		Assert(fn("ababab", "ab") is: 6)
		Assert(fn("xyz", "yx") is: 0)

		Assert(fn(#(), #(a)) is: 0)
		Assert(fn(#(a, b, c), #()) is: 1)
		Assert(fn(#(a, b, c), #(b)) is: 1)
		Assert(fn(#(a, b, a, b, a, b), #(a, b)) is: 6)
		}
	}