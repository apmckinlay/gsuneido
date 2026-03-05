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
				Thread(Bind(.f1, mu, wg), :name)
				}
			for ..2
				wg.Thread(Bind(.f2, mu), :name)
			Assert(Thread.List().CountIf({ it.Suffix?("WaitGroup_Test") }) is: 4,
				msg: "before")
			}
		wg.Wait(1)
		Assert(Suneido.WaitGroup_Test is: 4, msg: "after")
		}
	f1(mu, wg)
		{
		mu.Do({ })
		++Suneido.WaitGroup_Test
		wg.Done()
		}
	f2(mu)
		{
		mu.Do({ })
		++Suneido.WaitGroup_Test
		}
	Teardown()
		{
		Suneido.Delete(#WaitGroup_Test)
		}
	}