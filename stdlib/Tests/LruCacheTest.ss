// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		getfn = function (i) { return i }
		.testSampleData(getfn)
		}
	data: (11, 10, 8, 13, 2, 6, 18, 3, 4, 8, 5, 6, 15, 19, 10, 3,
		17, 9, 0, 14, 6, 10, 17, 19, 10, 10, 18, 8, 12, 19, 5, 11,
		19, 9, 1, 1, 18, 14, 12, 15, 3, 14, 0, 5, 7, 9, 7, 16, 18, 18)

	Test_GetN()
		{
		getfn = function (i, j) { Assert(i is: -j); return i }
		lru = LruCache(getfn, 20)
		for n in .data
			Assert(lru.Get(n, -n) is: n)
		Assert(lru.GetMissRate() is: .4)
		}

	testSampleData(getfn)
		{
		lru = LruCache(getfn, 20)
		for n in .data
			Assert(lru.Get(n) is: n)
		Assert(lru.GetMissRate() is: .4)
		return lru
		}
	}
