// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
class
	{
	New(.size = 10, .cmpFn = false)
		{
		if cmpFn is false
			.cmpFn = { |x, y| Cmp(x, y) < 0 }

		.collection = Object()
		}

	Add(x)
		{
		if .size is 0
			return

		if .collection.Size() < .size
			{
			.collection.Add(x)
			.siftUp()
			return
			}
		if .heapCmp(.collection[0], x)
			{
			.collection[0] = x
			.siftDown()
			return
			}
		return
		}

	siftUp()
		{
		i = .collection.Size() - 1
		while i > 0
			{
			parent = ((i - 1) / 2).Floor()
			if not .heapCmp(.collection[i], .collection[parent])
				break
			.collection.Swap(i, parent)
			i = parent
			}
		}
	siftDown()
		{
		i = 0
		n = .collection.Size()
		while true
			{
			best = i
			left = 2 * i + 1
			right = 2 * i + 2
			if left < n and .heapCmp(.collection[left], .collection[best])
				best = left
			if right < n and .heapCmp(.collection[right], .collection[best])
				best = right
			if best is i
				break
			.collection.Swap(i, best)
			i = best
			}
		}
	ToString()
		{
		return 'Heap(' $ .Collection().Join(', ') $ ')'
		}
	Collection()
		{
		return .collection.Copy().Sort!(.cmpFn)
		}
	Count()
		{
		return .collection.Size()
		}

	// heapCmp inverts cmpFn so root is the worst element
	// O(1) peek, O(log n) eviction
	heapCmp(a, b)
		{
		return (.cmpFn)(b, a)
		}
	}
