// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// SuJsWebTest
Test
	{
	Test_one_item()
		{
		s = Stack()
		s.Push("abc")
		Assert(s.Top() is: "abc")
		Assert(s.Pop() is: "abc")
		}
	Test_many_items()
		{
		s = Stack()
		for i in ..50
			s.Push(i)
		Assert(s.Top() is: 49)
		for (i = 49; i >= 0; --i)
			{
			Assert(s.Count() is: (i + 1))
			Assert(s.Pop() is: i)
			}
		}
	Test_members()
		{
		s = Stack()
		items = Object()
		for i in ..10
			items[i] = i
		s.Push(items)
		x = s.Pop()
		Assert(x is: items)
		}
	Test_Top()
		{
		s = Stack()
		for i in ..100
			s.Push(i)
		for i in ..50
			{
			x = Random(99)
			Assert(s.Top(x) is: (99 - x))
			}
		}
	Test_ToString()
		{
		Assert(Display(Stack()) is: 'Stack()')
		Assert(Display(Stack(1, 2, 3)) is: 'Stack(1, 2, 3)')
		Assert('' $ Stack() is: 'Stack()')
		Assert('' $ Stack(1, 2, 3) is: 'Stack(1, 2, 3)')
		}
	}
