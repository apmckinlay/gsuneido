// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	magicNumberTests: (
	(code: `class{
		Meth(){
		while (var1 < 44) //line 1 - magic number here
		}
	}`,
		warnings: (warnings: ([name: "stdlib:className:3 - Magic Number"]),
			rating: 4, desc: "Magic numbers found"),
		fullWarnings: (warnings: ([name: "Magic number on line: 3"]),
			rating: 4, desc: "Magic numbers found"),
		lineWarnings: ((35, 2, warning:))),

	(code: `class{
		Meth(){
			x3 += 3 //line 1 - magic number here
			}
		}`,
		warnings: (warnings: ([name: "stdlib:className:3 - Magic Number"]),
			rating: 4, desc: "Magic numbers found"),
		fullWarnings: (warnings: ([name: "Magic number on line: 3"]),
			rating: 4, desc: "Magic numbers found"),
		lineWarnings: ((28, 1, warning:))),

	(code: `class //No magic numbers
		{
		Print("Hello world")
		Print("The following -1 is in the exclusion list")
		if x < -1
			Print("Kool-aid")
		else if x < -2
			Print("Oh yeaaaah")
		x << 16 // bit shifting is not considered magic
		x >> 8 // bit shifting is not considered magic
		x <<= 16 // bit shifting is not considered magic
		x >>= 8 // bit shifting is not considered magic
		}`,
		warnings: (warnings: (), rating: 5, desc: ""),
		fullWarnings: (warnings: (), rating: 5, desc: "No magic numbers were found"),
		lineWarnings: ()),

	(code: 'class{
		privMeth()
		{
		fred = -55
		fred < -44
		fred - 55
		-44 < fred
		fred >= -666 /*= THE DEVIL */
		fred >= -666 /* = THE DEVIL */
		fred >= -666 /*=*/
		fred >= -666 /* = */
		}',
		warnings: (warnings: ([name: "stdlib:className:5 - Magic Number"],
				[name: "stdlib:className:6 - Magic Number"],
				[name: "stdlib:className:7 - Magic Number"],
				[name: "stdlib:className:10 - Magic Number"],
				[name: "stdlib:className:11 - Magic Number"])
			rating: 0,
			desc: "Magic numbers found\n2 occurrences of magic number(s): -44, -666"),
		fullWarnings: (warnings: ([name: "Magic number on line: 5"],
				[name: "Magic number on line: 6"],
				[name: "Magic number on line: 7"],
				[name: "Magic number on line: 10"],
				[name: "Magic number on line: 11"])
			rating: 0,
			desc: "Magic numbers found\n2 occurrences of magic number(s): -44, -666"),
		lineWarnings: ((51, 2, warning:), (64, 2, warning:),
			(71, 2, warning:), (160, 3, warning:), (182, 3, warning:))),

	(code: `class
		{
		Horz()
			{
			abc = Func(a,b, 23, c: 46)//line 1 - magic number here,
			}
		}`,
		warnings: (warnings: ([name: "stdlib:className:5 - Magic Number"]),
			rating: 4, desc: "Magic numbers found"),
		fullWarnings: (warnings: ([name: "Magic number on line: 5"]),
			rating: 4, desc: "Magic numbers found"),
		lineWarnings: ((47, 2, warning:))),

	(code: `class
	{
	MagicMethod (var1, var2, var3)
		{
		while (var1 < 44)
			{
	//      if (var2 < 34)
				do stuff...
			var1++
			}
		if var2 is 10
			Round(2)
			Round(5)
			Round(-55)
		}
	}`,
		warnings: (warnings: ([name: "stdlib:className:5 - Magic Number"],
				[name: "stdlib:className:11 - Magic Number"]),
			rating: 3,
			desc: "Magic numbers found"),
		fullWarnings: (warnings: ([name: "Magic number on line: 5"],
				[name: "Magic number on line: 11"]),
			rating: 3,
			desc: "Magic numbers found"),
		lineWarnings: ((65, 2, warning:), (148, 2, warning:))),

	(code: `class
		{
		privateMagicMethod (x1, x2, x3)
			{
			while (5 < x1)		//Line 5 - Magic number
				x3 += 72		//Line 6 - Magic number
				if x3 is 72		//Line 7 - Magic number
					x3 -= 5		//Line 8 - Magic number
				else //x3 += 5
					x2 *= x3
				x2 /= 3			//Line 11 - Magic number
			for (i = 0; i < 4; i++)//Line 12 - Magic number
				{
				j = i + 5		//Line 14 - Magic number
				k = j + 4
				}
			FuncCall(x: 5)
			Func2Call(55)
			}
		}`,
		warnings: (warnings: ([name: "stdlib:className:5 - Magic Number"],
				[name: "stdlib:className:6 - Magic Number"],
				[name: "stdlib:className:7 - Magic Number"],
				[name: "stdlib:className:8 - Magic Number"],
				[name: "stdlib:className:11 - Magic Number"],
				[name: "stdlib:className:12 - Magic Number"],
				[name: "stdlib:className:14 - Magic Number"],
				[name: "stdlib:className:15 - Magic Number"],
				[name: "stdlib:className:18 - Magic Number"]),
			rating: 0,
			desc: "Magic numbers found\n2 occurrences of magic number(s): 72, 4\n" $
				"3 occurrences of magic number(s): 5"),
		fullWarnings: (warnings: ([name: "Magic number on line: 5"],
				[name: "Magic number on line: 6"],
				[name: "Magic number on line: 7"], [name: "Magic number on line: 8"],
				[name: "Magic number on line: 11"], [name: "Magic number on line: 12"],
				[name: "Magic number on line: 14"], [name: "Magic number on line: 15"],
				[name: "Magic number on line: 18"]),
			rating: 0,
			desc: "Magic numbers found\n2 occurrences of magic number(s): 72, 4\n" $
				"3 occurrences of magic number(s): 5"),
		lineWarnings: ((63, 1, warning:), (107, 2, warning:), (149, 2, warning:),
			(189, 1, warning:), (262, 1, warning:), (311, 1, warning:),
			(363, 1, warning:), (404, 1, warning:), (446, 2, warning:))),

	(code: `class
		{
		privMeth()
			{
			.anotherMethod(-5,6,7,8,9,-2)
			if x < -5 and x > -2 and x isnt 0 and x < -2
				Print("I'm a test")
			if x < -22
				x += -22
			}
		}`,
		warnings: (warnings: ([name: "stdlib:className:5 - Magic Number"],
				[name: "stdlib:className:6 - Magic Number"],
				[name: "stdlib:className:8 - Magic Number"],
				[name: "stdlib:className:9 - Magic Number"]),
			rating: 0,
			desc : "Magic numbers found\n2 occurrences of magic number(s): -5, -22"),
		fullWarnings: (warnings: ([name: "Magic number on line: 5"],
				[name: "Magic number on line: 6"],
				[name: "Magic number on line: 8"],
				[name: "Magic number on line: 9"]),
			rating: 0,
			desc : "Magic numbers found\n2 occurrences of magic number(s): -5, -22"),
		lineWarnings: ((51, 1, warning:), (53, 1, warning:), (55, 1, warning:),
			(57, 1, warning:), (59, 1, warning:), (77, 1, warning:),
			(151, 2, warning:), (165, 2, warning:))),

	(code: 'class
		{
		}',
		warnings: (warnings: (), rating: 5, desc: ""),
		fullWarnings: (warnings: (), rating: 5, desc: "No magic numbers were found"),
		lineWarnings: ())
	)

	Test_ContinuousCC_MagicNumbers_ContinuousChecks()
		{
		recordData = Record(lib: "stdlib", recordName: "className")
		for test in .magicNumberTests
			{
			recordData.code = test.code
			fullWarnings = Qc_MagicNumbers(recordData, minimizeOutput?: false)
			lineWarnings = fullWarnings.Extract('lineWarnings')
			Assert(lineWarnings isSize: 0)
			Assert(fullWarnings is: test.fullWarnings)

			warnings = Qc_MagicNumbers(recordData, minimizeOutput?:)
			lineWarnings = warnings.Extract('lineWarnings')
			Assert(lineWarnings is: test.lineWarnings)
			Assert(warnings is: test.warnings)
			}
		}
	}