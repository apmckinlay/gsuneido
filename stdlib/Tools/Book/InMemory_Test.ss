// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		prev = Suneido.GetDefault(#InMemory, Object())
		.AddTeardown({ Suneido.InMemory = prev })
		Suneido.Delete(#InMemory)

		url = InMemory.Add(s ="hello world")
		Assert(InMemory.Get(url) is: s)
		Assert(InMemory.Get(url $ '.ext') is: s)

		InMemory.Remove(url)
		Assert(Suneido.InMemory isSize: 0)
		}
	}