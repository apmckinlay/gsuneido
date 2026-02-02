// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		list = new DoubleLinkedList()
		Assert(list.ToList() is: #())
		Assert(list.Empty?())

		a = list.Append(#a)
		Assert(list.ToList() is: #(a))
		Assert(not list.Empty?())

		b = list.Append(#b)
		c = list.Append(#c)
		Assert(list.ToList() is: #(a, b, c))

		list.Del(b)
		Assert(list.ToList() is: #(a, c))

		list.Del(a)
		Assert(list.ToList() is: #(c))

		d = list.Insert(#d, c)
		Assert(list.ToList() is: #(c, d))

		x = list.Insert(#x)
		Assert(list.ToList() is: #(x, c, d))
		Assert(not list.Empty?())

		list.Del(x)
		list.Del(d)
		list.Del(c)
		Assert(list.ToList() is: #())
		Assert(list.Empty?())

		a = list.Append(#a)
		b = list.Append(#b)
		c = list.Append(#c)
		d = list.Append(#d)
		Assert(list.ToList() is: #(a, b, c, d))

		list.Extract(b, d)
		newList = DoubleLinkedList(b, d)
		Assert(list.ToList() is: #(a))
		Assert(newList.ToList() is: #(b, c, d))

		Assert(newList.ToList({ |prev, cur| prev is #b and cur is #c }) is: #(b, d))
		Assert(newList.ToList({ |prev, cur/*unused*/| prev is #b }) is: #(b))
		Assert(newList.ToList({ |prev/*unused*/, cur/*unused*/| false }) is: #(b, c, d))
		}
	}