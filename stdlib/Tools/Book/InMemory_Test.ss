// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		cl = InMemory
			{
			InMemory_serverEval?() { return false }
			}
		prev = Suneido.GetDefault(#InMemory, Object())
		.AddTeardown({ Suneido.InMemory = prev })
		Suneido.Delete(#InMemory)

		url = cl.Add(s ="hello world")
		Assert(cl.Get(url) is: s)
		Assert(cl.Get(url $ '.ext') is: s)

		cl.Remove(url)
		Assert(Suneido.InMemory isSize: 0)
		}
	}