// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	cases: (
		((), ()),
		((a, (), b, ()), (a, b)),
		((a, (b), b, ()), (b, a)),
		((a, (b, b, b), b, ()), (b, a)),
		((a, (b), b, (), c, (a)), (b, a, c)),
		((a, (b), b, (), c, (b)), (b, a, c)),
		((2, (5), 0, (5, 4), 1, (3, 4), 3, (2), 4, (), 5, ()), (4, 5, 2, 0, 3, 1))
		)
	Test_main()
		{
		fn = TopologicalSort
		for c in .cases
			Assert(fn(.convert(c[0])).Map({ it.name }) is: c[1])

		Assert({ fn(.convert(#(a, (), a, ()))) },
			throws: 'Find duplicate name: a')
		Assert({ fn(.convert(#(a, (b, c)))) },
			throws: 'Find unknown deps:
b in a
c in a')
		Assert({ fn(.convert(#(a, (b), b, (c), c, (a)))) },
			throws: 'Find circles:
a -> b -> c -> a')
		Assert({ fn(.convert(#(a, (b), b, (c, d), c, (a), d, (b)))) },
			throws: 'Find circles:
a -> b -> c -> a
b -> d -> b')
		Assert({ fn(.convert(#(a, (b, d), b, (c, e), c, (a), d, (a)))) },
			throws: 'Find circles:
a -> b -> c -> a
a -> d -> a
Find unknown deps:
e in b')
		}
	convert(ob)
		{
		converted = Object()
		for (i = 0; i < ob.Size(); i += 2)
			converted.Add(Object(name: ob[i], deps: ob[i+1]))
		return converted
		}
	}