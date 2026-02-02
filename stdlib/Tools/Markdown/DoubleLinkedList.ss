// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
class
	{
	New(start = false, end = false)
		{
		.head = Object()
		if start is false
			.head.next = .head
		else
			{
			.head.next = start
			start.prev = .head
			}
		if end is false
			.head.prev = .head
		else
			{
			.head.prev = end
			end.next = .head
			}
		}

	Append(value)
		{
		item = Object(:value, next: .head, prev: .head.prev)
		.head.prev.next = item
		.head.prev = item
		return item
		}

	Insert(value, after = false)
		{
		if after is false
			after = .head
		item = Object(:value, next: after.next, prev: after)
		after.next.prev = item
		after.next = item
		return item
		}

	Del(item)
		{
		item.next.prev = item.prev
		item.prev.next = item.next
		}

	Extract(start, end)
		{
		start.prev.next = end.next
		end.next.prev = start.prev
		}

	Next(cur = false)
		{
		next = cur is false ? .head.next : cur.next
		return Same?(next, .head) ? false : next
		}

	Prev(cur = false)
		{
		prev = cur is false ? .head.prev : cur.prev
		return Same?(prev, .head) ? false : prev
		}

	Empty?()
		{
		return Same?(.head.next, .head)
		}

	ToList(mergeFn = false)
		{
		list = Object()
		for (cur = .Next(); cur isnt false; cur = .Next(cur))
			{
			if mergeFn isnt false and list.NotEmpty?() and mergeFn(list.Last(), cur.value)
				continue
			list.Add(cur.value)
			}
		return list
		}
	}