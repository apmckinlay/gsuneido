// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	NumParamTests: (
		(code: "// Copyright (C) 2017 Suneido Software Corp. All rights reserved
//	 worldwide
function (x1, x2, x5, x7, x8, x9, x10, x11, x4 = Object(), x3 = [], x6 = false)
	{
	x = 5
	y = 55
	if (x > y)
		{
		return 55
		}
	return (y > x)
	}",
		warnings: (warnings: ([name: "stdlib:className:3 - 8 parameters in function"]),
			rating: 0, desc: "Number of parameters - Rating affected -> Limit to 5"),
		fullWarnings: (warnings: ([name: "8 parameters in function"]),
			rating: 0, desc: "Number of parameters - Rating affected -> Limit to 5"),
		lineWarnings: ((3))),

		(code: "function (x1, x2, x3)
	{
	x = 5
	y = 55
	if (x > y)
		{
		return 55
		}
	return (y > x)
	}",
		warnings: (warnings: (), rating: 5, desc: ""),
		fullWarnings: (warnings: ([name: '3 parameters in function']), rating: 5,
			desc: 'Number of parameters - Rating not affected -> Limit to 5'),
		lineWarnings: ()),

		(code: "class {
	N: 8 // block size (number of lines)
	CallClass(libs, libarry, lib3, lib4, lib5, lib6, lib7, lib8,
		default1 = false, default2 = true, default3 = 'bunny', default4 = 0)
		{
		if String?(libs)
			libs = [libs]
		(new this).Detect(libs)
		}
	Detect(libs, libs2, libs3, libs4, x = true, y = false, z = false, q = 5)
		{
		.hashes = Object()
		for lib in libs
			.process(lib)
		.output()
		}
	process(lib, bil, lib3, lib4, lib5,)
		{
		QueryApply(lib $ ' where name !~ 'Test$'', group: -1)
			{|x|
			.process1(lib, x)
			}
		}
	process1(lib, x, var3, var4)
		{
		last = -999
		lines = x.text.Lines().Map!(#Trim)
		for (i = 0; i < lines.Size() - .N; ++i)
			{
			if not .hashes.Member?(hash)
				.hashes[hash] = name
			else if i - last >= .N
				{
				last = i
				.hashes[hash] $= ', ' $ name
				}
			}
		}
	PubMeth5(){ }
	PubMeth6(){ }
	PubMeth7(){ }
	PubMeth8(){ }
	PubMeth9(){ }
	PubMeth10(){ }
	}",
		warnings: (warnings: ([name: "stdlib:className:3 - 8 parameters in CallClass"]),
			rating: 2, desc: "Number of parameters - Rating affected -> Limit to 5"),
		fullWarnings: (warnings: ([name: "8 parameters in CallClass"],
				[name: "5 parameters in process"],
				[name: "4 parameters in Detect"],
				[name: "4 parameters in process1"],
				[name: "0 parameters in PubMeth5"],
				[name: "0 parameters in PubMeth6"],
				[name: "0 parameters in PubMeth7"],
				[name: "0 parameters in PubMeth8"],
				[name: "0 parameters in PubMeth9"],
				[name: "0 parameters in PubMeth10"]),
			rating: 2, desc: "Number of parameters - Rating affected -> Limit to 5"),
		lineWarnings: ((3))),

		(code: "class {
	N: 8 // block size (number of lines)
	CallClass(libs, libarry, lib3, lib4, lib5, lib6,
		default1 = false, default2 = true, default3 = 'bunny', default4 = 0)
		{
		if String?(libs)
			libs = [libs]
		(new this).Detect(libs)
		}
	Detect(libs, libs2, libs3, libs4, x = true, y = false, z = false, q = 5)
		{
		.hashes = Object()
		for lib in libs
			.process(lib)
		.output()
		}
	process(lib, bil, lib3, lib4, lib5,)
		{
		QueryApply(lib $ ' where name !~ 'Test$'', group: -1)
			{|x|
			.process1(lib, x)
			}
		}
	process1(lib, x, var3, var4)
		{
		last = -999
		lines = x.text.Lines().Map!(#Trim)
		for (i = 0; i < lines.Size() - .N; ++i)
			{
			if not .hashes.Member?(hash)
				.hashes[hash] = name
			else if i - last >= .N
				{
				last = i
				.hashes[hash] $= ', ' $ name
				}
			}
		}
	PubMeth5(){ }
	PubMeth6(){ }
	PubMeth7(){ }
	PubMeth8(){ }
	PubMeth9(){ }
	PubMeth10(){ }
	}",
		warnings: (warnings: ([name: "stdlib:className:3 - 6 parameters in CallClass"]),
			rating: 5, desc: "Number of parameters - Rating not affected -> Limit to 5"),
		fullWarnings: (warnings: ([name: "6 parameters in CallClass"],
			[name: "5 parameters in process"],
			[name: "4 parameters in Detect"],
			[name: "4 parameters in process1"],
			[name: "0 parameters in PubMeth5"],
			[name: "0 parameters in PubMeth6"],
			[name: "0 parameters in PubMeth7"],
			[name: "0 parameters in PubMeth8"],
			[name: "0 parameters in PubMeth9"],
			[name: "0 parameters in PubMeth10"]),
			rating: 5, desc: "Number of parameters - Rating not affected -> Limit to 5"),
		lineWarnings: ((3))),

	(code: "class
	{
	func1(a, b, c, d, e, f, g)
		{
		// test comment
		// test comment
		// test comment
		// test comment
		// test comment
		// test comment
		}
	func3(a, b, c, d, e, f) { }
	func2(a, b, c, d, e, f, g, h, i, j) { }
	}",
		warnings: (warnings: ([name: "stdlib:className:13 - 10 parameters in func2"],
				[name: "stdlib:className:3 - 7 parameters in func1"],
				[name: "stdlib:className:12 - 6 parameters in func3"]),
			rating: 0, desc: "Number of parameters - Rating affected -> Limit to 5"),
		fullWarnings: (warnings: ([name: "10 parameters in func2"],
				[name: "7 parameters in func1"],
				[name: "6 parameters in func3"]),
			rating: 0, desc: "Number of parameters - Rating affected -> Limit to 5"),
		lineWarnings: ((13), (3), (12))),

	(code: "function () { return 0 }",
		warnings: (warnings: (), rating: 5, desc: ""),
		fullWarnings: (warnings: ([name: "0 parameters in function"]),
			rating: 5, desc: "Number of parameters - Rating not affected -> Limit to 5"),
		lineWarnings: ()),

	(code: "class { }",
		warnings: (warnings: (), rating: 5, desc: ""),
		fullWarnings: (warnings: (), rating: 5,
			desc: "Number of parameters - Rating not affected -> Limit to 5"),
		lineWarnings: ())
	)


	Test_RetrieveParamsList()
		{
		recordData = Record(lib: 'stdlib', recordName: 'className')
		for test in .NumParamTests
			{
			recordData.code = test.code
			fullWarnings = Qc_NumParams(recordData, minimizeOutput?: false)
			lineWarnings = fullWarnings.Extract('lineWarnings')
			Assert(fullWarnings is: test.fullWarnings)
			Assert(lineWarnings isSize: 0)

			warnings = Qc_NumParams(recordData, minimizeOutput?:)
			lineWarnings = warnings.Extract('lineWarnings')
			Assert(warnings is: test.warnings)
			Assert(lineWarnings is: test.lineWarnings)
			}
		}
	}