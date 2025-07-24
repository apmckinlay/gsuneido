// Copyright (C) 2015 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_SplitCoord_and_valid?()
		{
		split = CoordControl.SplitCoord
		valid = CoordControl.CoordControl_valid?

		Assert(split("") is: ob = #(x: '', y: ''))
		Assert(valid(ob.x, ob.y))
		Assert(valid(ob.x, ob.y, mandatory:) is: false)

		Assert(split('bob,') is: ob = #(x: 'bob', y: ''))
		Assert(valid(ob.x, ob.y) is: false)
		Assert(valid(ob.x, ob.y, mandatory:) is: false)

		Assert(split('bob,bill') is: ob = #(x: 'bob', y: 'bill'))
		Assert(valid(ob.x, ob.y) is: false)
		Assert(valid(ob.x, ob.y, mandatory:) is: false)

		Assert(split('bob,2') is: ob = #(x: 'bob', y: 2))
		Assert(valid(ob.x, ob.y) is: false)
		Assert(valid(ob.x, ob.y, mandatory:) is: false)

		Assert(split('2,bill') is: ob = #(x: 2, y: 'bill'))
		Assert(valid(ob.x, ob.y) is: false)
		Assert(valid(ob.x, ob.y, mandatory:) is: false)

		Assert(split('2,') is: ob = #(x: 2, y: ''))
		Assert(valid(ob.x, ob.y) is: false)
		Assert(valid(ob.x, ob.y, mandatory:) is: false)

		Assert(split(',4') is: ob = #(x: '', y: 4))
		Assert(valid(ob.x, ob.y) is: false)
		Assert(valid(ob.x, ob.y, mandatory:) is: false)

		Assert(split('2,3') is: ob = #(x: 2, y: 3))
		Assert(valid(ob.x, ob.y))
		Assert(valid(ob.x, ob.y, mandatory:))
		}
	}