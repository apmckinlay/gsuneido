// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_equality()
		{
		c = class { F(){ } }
		instance1 = c()
		Assert(instance1.F is: instance1.F)
		instance2 = c()
		Assert(instance1.F isnt: instance2.F)
		}
	}