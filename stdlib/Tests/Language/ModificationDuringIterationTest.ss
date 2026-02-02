// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		omdi = 'object modified during iteration'
		Assert({
			ob = Object(12, 34, a: 56, b: 78)
			for unused in ob.Members()
				ob.Add(99)
			}, throws: omdi, msg: 'Add')
		Assert({
			ob = Object(12, 34, a: 56, b: 78)
			for unused in ob.Members()
				ob.x = 99
			}, throws: omdi, msg: 'set')
		Assert({
			ob = Object(12, 34, a: 56, b: 78)
			for unused in ob.Members()
				ob.Delete(0)
			}, throws: omdi, msg: 'vec Delete')
		Assert({
			ob = Object(12, 34, a: 56, b: 78)
			for unused in ob.Members()
				ob.Delete(#a)
			}, throws: omdi, msg: 'map Delete')
		}
	}