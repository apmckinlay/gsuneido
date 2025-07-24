// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		c1 = Channel()
		c2 = Channel()
		Thread()
			{
			sum = 0
			while c1 isnt x = c1.Recv()
				sum += x
			c2.Send(sum)
			}
		for i in ..10
			c1.Send(i)
		c1.Close()
		result = c2.Recv()
		Assert(result is: 45)
		}
	}