// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		for ..20
			{
			.ob = Object()
			.out = Object()
			wg = WaitGroup()
			for i in ..4
				wg.Thread(Bind(.thread, i))
			wg.Wait()
			for x in .out
				Assert(x is: .out[0])
			Assert(.ob.x is: .out[0])
			}
		}
	thread(i)
		{
		.out.Add(.ob.GetInit(#x, { i }))
		}
	}