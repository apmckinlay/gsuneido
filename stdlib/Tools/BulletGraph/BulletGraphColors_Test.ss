// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		// Purple
		result = BulletGraphColors(0x660033, { it })
		Assert(result.bad is: Object(r: 183, g: 81, b: 132))
		Assert(result.satisfactory is: Object(r: 224, g: 122, b: 173))
		Assert(result.good is: Object(r: 255, g: 153, b: 204))
		Assert(result.value is: Object(r: 102, g: 0, b: 51))

		// Teal/Green
		result = BulletGraphColors(0x32693D, { it })
		Assert(result.bad is: Object(r: 134, g: 189, b: 145))
		Assert(result.satisfactory is: Object(r: 176, g: 231, b: 187))
		Assert(result.good is: Object(r: 207, g: 255, b: 218))
		Assert(result.value is: Object(r: 50, g: 105, b: 61))

		// Black
		result = BulletGraphColors(0x000000, { it })
		Assert(result.bad is: Object(r: 0, g: 0, b: 0))
		Assert(result.satisfactory is: Object(r: 0, g: 0, b: 0))
		Assert(result.good is: Object(r: 0, g: 0, b: 0))
		Assert(result.value is: Object(r: 0, g: 0, b: 0))

		// White
		result = BulletGraphColors(0xFFFFFF, { it })
		Assert(result.bad is: Object(r: 255, g: 255, b:255))
		Assert(result.satisfactory is: Object(r: 255, g: 255, b: 255))
		Assert(result.good is: Object(r: 255, g: 255, b: 255))
		Assert(result.value is: Object(r: 255, g: 255, b: 255))
		}
	}