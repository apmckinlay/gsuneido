// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_deps()
		{
		r = Record(a: 123)
		Assert(r.a is: 123)
		r.SetDeps(#a, 'b')
		Assert(r.a is: 123)
		r.b = 456
		Assert(r.a is: 123)
		}
	Test_invalidate()
		{
		r = Record(a: 123)
		Assert(r.a is: 123)
		r.Invalidate('a')
		Assert(r.a is: 123)
		}
	}