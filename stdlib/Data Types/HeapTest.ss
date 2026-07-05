// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_empty()
		{
		h = Heap()
		Assert(h.Count() is: 0)
		Assert(h.Collection() is: #())
		Assert(h.ToString() is: 'Heap()')
		}

	Test_min_keeps_smallest()
		{
		h = Heap(size: 3)
		for x in #(5, 3, 8, 1, 9, 2)
			h.Add(x)
		Assert(h.Collection() is: #(1, 2, 3))
		Assert(h.Count() is: 3)
		Assert(h.ToString() is: 'Heap(1, 2, 3)')
		}

	Test_max_keeps_largest()
		{
		maxFirst = {|x, y| Cmp(y, x) < 0 }
		h = Heap(size: 3, cmpFn: maxFirst)
		for x in #(5, 3, 8, 1, 9, 2)
			h.Add(x)
		Assert(h.Collection() is: #(9, 8, 5))
		Assert(h.ToString() is: 'Heap(9, 8, 5)')
		}

	Test_under_capacity()
		{
		h = Heap()
		for x in #(3, 1, 2)
			h.Add(x)
		Assert(h.Collection() is: #(1, 2, 3))
		Assert(h.Count() is: 3)
		}

	Test_capacity_cap()
		{
		h = Heap(size: 10)
		for unused in ..20
			h.Add(Random(1000))
		Assert(h.Count() is: 10)
		}

	Test_duplicates_kept()
		{
		h = Heap(size: 3)
		for x in #(1, 1, 1, 1)
			h.Add(x)
		Assert(h.Collection() is: #(1, 1, 1))
		}

	Test_string_comparator()
		{
		h = Heap(size: 2)
		for s in #('banana', 'apple', 'cherry', 'dapple')
			h.Add(s)
		Assert(h.Collection() is: #('apple', 'banana'))
		}

	Test_eviction_boundary()
		{
		h = Heap(size: 2)
		h.Add(1)
		h.Add(2)
		h.Add(3)
		Assert(h.Collection() is: #(1, 2))
		h.Add(0)
		Assert(h.Collection() is: #(0, 1))
		}

	Test_more()
		{
		// -- size zero
		h = Heap(size: 0)
		h.Add(5)
		h.Add(3)
		Assert(h.Count() is: 0)
		Assert(h.Collection() is: #())
		Assert(h.ToString() is: 'Heap()')

		// -- size one
		h = Heap(size: 1)
		h.Add(5)
		h.Add(3)
		Assert(h.Collection() is: #(3))
		h.Add(7)
		Assert(h.Collection() is: #(3))

		// -- collection is independent copy
		h = Heap(size: 3)
		for x in #(5, 3, 8)
			h.Add(x)
		c = h.Collection()
		c.Add(99)
		Assert(h.Collection() is: #(3, 5, 8))

		// -- descending input
		// each new element is better than the previous root forces SiftUp to swap all the way
		h = Heap(size: 5)
		for x in #(50, 40, 30, 20, 10, 5, 2)
			h.Add(x)
		Assert(h.Collection() is: #(2, 5, 10, 20, 30))
		Assert(h.Count() is: 5)

		// -- ascending input
		// each new element is worse than all existing  never replaces root once full
		h = Heap(size: 3)
		for x in #(1, 2, 3, 4, 5, 6)
			h.Add(x)
		Assert(h.Collection() is: #(1, 2, 3))

		// -- interleaved
		h = Heap(size: 4)
		for x in #(5, 1, 8, 2, 9, 3, 7, 4)
			h.Add(x)
		Assert(h.Collection() is: #(1, 2, 3, 4))

		// -- max heap interleaved
		h = Heap(size: 4, cmpFn: Gt)
		for x in #(5, 1, 8, 2, 9, 3, 7, 4)
			h.Add(x)
		Assert(h.Collection() is: #(9, 8, 7, 5))

		// -- large random
		h = Heap(size: 50)
		for ..500
			h.Add(Random(100000))
		Assert(h.Count() is: 50)
		Assert(h.Collection().Sorted?())

		// -- max heap large random
		h = Heap(size: 30, cmpFn: Gt)
		for ..300
			h.Add(Random(100000))
		Assert(h.Count() is: 30)
		Assert(h.Collection().Reverse!().Sorted?())
		}
	}
