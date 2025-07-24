// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	arr_classes: (
	(code:
		"//Suneido Copyright...
		function (x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x) //line 2
			{
			x = 555555555555555555555555555555555555555555555555555555555555555555555555555 //line 4
			// ===========================================================================================================Should not trigger as too long //Line 5
			}",
		warnings: (warnings: ([name: "stdlib:className:2 - Long line"],
				[name:"stdlib:className:4 - Long line"]),
			rating : 3, desc : "Some lines are too long"),
		fullWarnings: (warnings: ([name: "Line: 2"], [name: "Line: 4"]),
			rating: 3, desc: "Some lines are too long"),
		lineWarnings: ((2),(4))),

	(code: "//Suneido Copyright...
		class
			{
			PubMethod1(y,y,y,y,y,y,y,y,y,y,y,y,y,y,y,y,y,y,y,y,y,y,y,y,y,y,y,y,y,y,y,y,y,y) //Line 4
				{
				x = y*y*y*y*y*y*y*y*y*y*y*y*y*y*y*y*y*y*y*y*y*y*y*y*y*y*y*y*y*y*y*y*y*y*y*y*y*y*y*y*y*y*y*y*y //Line 6
				}
			privMethod1(x,y,z,x,y,z,x,y,z,x,y,z,x,y,z,x,y,z,x,y,z,x,y,z,x,y,z,x,y,z,x,y,z) //Line 8
				{
				x = 5
				x = 6
				// --------------------------------------------------------------------------------Should not trigger as too long
				//Line 12///////////////////////////////////////////////////////////////////////////////////////////////////////// Should not trigger as too long
				//Line 13
				x = 77777777777777777777777777777777777777777777777777777777777777777777771
				if (x < 44444444444444444444444444444444444444444444444444444444444444444444) //Line 15
					return false
				}
			}",
		warnings: (warnings: ([name: "stdlib:className:4 - Long line"],
				[name: "stdlib:className:6 - Long line"],
				[name: "stdlib:className:8 - Long line"],
				[name: "stdlib:className:15 - Long line"],
				[name: "stdlib:className:16 - Long line"]),
			rating: 0, desc: "Some lines are too long"),
		fullWarnings: (warnings: ([name: "Line: 4"],
				[name: "Line: 6"],
				[name: "Line: 8"],
				[name: "Line: 15"],
				[name: "Line: 16"]),
			rating: 0, desc: "Some lines are too long"),
		lineWarnings: ((4), (6), (8), (15), (16))),

	(code: "//Suneido
		class
			{
			//No lines are too long
			FakeMeth1()
				{
				}
			}",
		warnings: (warnings: (), rating: 5, desc: ""),
		fullWarnings: (warnings: (), rating: 5, desc: "All lines are of adequate length"),
		lineWarnings: ())
	)

	Test_Continuous_Line_Size()
		{
		recordData = Record(lib: "stdlib", recordName: "className")
		for test in .arr_classes
			{
			recordData.code = test.code
			warnings = Qc_LineSize(recordData, minimizeOutput?: false)
			lineWarnings = warnings.Extract('lineWarnings')
			Assert(warnings is: test.fullWarnings)
			Assert(lineWarnings isSize: 0)

			warnings = Qc_LineSize(recordData, minimizeOutput?:)
			lineWarnings = warnings.Extract('lineWarnings')
			Assert(warnings is: test.warnings)
			Assert(lineWarnings is: test.lineWarnings)
			}
		}
	}