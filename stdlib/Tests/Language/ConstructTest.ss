// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// requires ConstructTestClass
Test
	{
	Test_one()
		{
		tests = Object(
			[ConstructTestClass,						"okay"],
			['ConstructTestClass', 						"okay"],
			[Object(ConstructTestClass), 				"okay"],
			[Object('ConstructTestClass'), 				"okay"],
			[Object(ConstructTestClass, 'arg'), 		"okayarg"],
			[Object('ConstructTestClass', 'arg'), 		"okayarg"],
			['ConstructTest', 'Class', 		   			"okay", 	multiArgs:],
			[Object('ConstructTest'), 'Class', 			"okay", 	multiArgs:],
			)
		for t in tests
			if t.GetDefault('multiArgs', false) is true
				Assert(Construct(t[0], t[1]).F() is: t[2])
			else
				Assert(Construct(t[0]).F() is: t[1])
		}
	Test_two()
		{
		Assert(Construct("", "ConstructTestClass").F() is: "okay")
		Assert(Construct(#(""), "ConstructTestClass").F() is: "okay")

		// extra named args should be ignored as usual
		Assert(Construct(#("", a: 1, b: 2), "ConstructTestClass").F() is: "okay")

		Assert({ Construct(#("", 123, 456), "ConstructTestClass").F() }
			throws: "too many arguments")
		}
	}