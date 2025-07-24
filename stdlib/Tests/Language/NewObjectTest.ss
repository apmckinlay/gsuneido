// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		v1 = 1
		v2 = 2
		v3 = 3
		v4 = 4
		Assert([] is: #{})

		Assert([1] is: #{1})
		Assert([v1] is: #{1})

		Assert([1, 2] is: #{1, 2})
		Assert([v1, v2] is: #{1, 2})
		Assert([1, v2] is: #{1, 2})
		Assert([v1, 2] is: #{1, 2})

		Assert([a: 1] is: #{a: 1})
		Assert([a: v1] is: #{a: 1})

		Assert([a: 1, b: 2] is: #{a: 1, b: 2})
		Assert([a: v1, b: 2] is: #{a: 1, b: 2})
		Assert([a: 1, b: v2] is: #{a: 1, b: 2})
		Assert([a: v1, b: v2] is: #{a: 1, b: 2})

		Assert([1, a: 2] is: #{1, a: 2})
		Assert([v1, a: v2] is: #{1, a: 2})
		Assert([1, a: v2] is: #{1, a: 2})
		Assert([v1, a: 2] is: #{1, a: 2})

		Assert([1, 2, a: 3, b: 4] is: #{1, 2, a: 3, b: 4})
		Assert([v1, v2, a: 3, b: 4] is: #{1, 2, a: 3, b: 4})
		Assert([1, 2, a: v3, b: v4] is: #{1, 2, a: 3, b: 4})
		Assert([v1, 2, a: 3, b: v4] is: #{1, 2, a: 3, b: 4})
		Assert([1, 2, a: 3, b: 4] is: #{1, 2, a: 3, b: 4})
		Assert([1, v2, a: v3, b: 4] is: #{1, 2, a: 3, b: 4})
		Assert([v1, v2, a: v3, b: v4] is: #{1, 2, a: 3, b: 4})
		}
	}