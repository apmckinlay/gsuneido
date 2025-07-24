// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		name = "WaitGroup_Test"
		Suneido.WaitGroup_Test = 0
		wg = WaitGroup()
		mu = Mutex()
		mu.Do()
			{
			for ..2
				{
				wg.Add()
				Thread(:name)
					{
					mu.Do({ })
					++Suneido.WaitGroup_Test
					wg.Done()
					}
				}
			for ..2
				wg.Thread(:name)
					{
					mu.Do({ })
					++Suneido.WaitGroup_Test
					}
			Assert(Thread.List().CountIf({ it.Suffix?("WaitGroup_Test") }) is: 4,
				msg: "before")
			}
		wg.Wait(1)
		Assert(Suneido.WaitGroup_Test is: 4, msg: "after")
		}
	Teardown()
		{
		Suneido.Delete(#WaitGroup_Test)
		}
	}