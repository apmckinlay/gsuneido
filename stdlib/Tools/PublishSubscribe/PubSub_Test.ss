// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		flag = Object(count: 0)
		unsub = PubSub.Subscribe(#test, { flag.count++ })
		PubSub.Publish(#test)
		Assert(flag.count is: 1)
		PubSub.Publish(#test)
		Assert(flag.count is: 2)
		unsub.Unsubscribe()
		}
	}