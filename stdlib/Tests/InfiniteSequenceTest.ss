// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		ob = Object(InfiniteSequenceStart, false)
		for (i = 0; i < 1000; ++i)
			{
			n = 1 + Random(ob.Size() - 1)
			before = ob[n - 1]
			after = ob[n]
			Assert(after is false or before < after)
			mid = InfiniteSequence(before, after)
			Assert(mid > before and (after is false or mid < after))
			ob.Add(mid at: n)
			displayed = Display(mid)
			Assert(displayed.Eval() is: mid)
			}
		prev = ''
		for n in ob
			{
			Assert(n > prev or n is false)
			prev = n
			}
		}

	Test_insertAllAtEnd()
		{
		ob = Object(InfiniteSequenceStart, false)
		for (i = 0; i < 1000; ++i)
			{
			n = ob.Size() - 1
			before = ob[n - 1]
			after = ob[n]
			Assert(after is false or before < after)
			mid = InfiniteSequence(before, after)
			Assert(mid > before and (after is false or mid < after))
			ob.Add(mid at: n)
			displayed = Display(mid)
			Assert(displayed.Eval() is: mid)
			}
		}

	Test_two()
		{
		m = InfiniteSequence(b = IntToStr(1) $ '\xff', a = IntToStr(2))
		Assert(b < m and m < a)
		}
	}