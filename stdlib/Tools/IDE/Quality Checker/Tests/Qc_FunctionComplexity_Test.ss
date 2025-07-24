// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	tests: ((code: "// Copyright (C) 2017 Suneido Software Corp.
		//All rights reserved worldwide.
function (x1, x2, x3, x4, x5, x6)
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
		fullWarnings: (warnings: ([name: "Function McCabe complexity: 2"]), rating: 5,
			desc: "McCabe function complexity - Rating not affected -> Limit to 8"),
		lineWarnings: ()),

	(code: "// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
class
	{
	N: 8 // block size (number of lines)
	CallClass(libs, libarry, lib3, lib4, lib5, lib6, lib7, lib8)
		{
		if String?(libs)
			libs = [libs]
		(new this).Detect(libs)
		}
	Detect(libs, libs2, libs3, libs4)
		{
		.hashes = Object()
		for lib in libs
			{
			.process(lib)
			if (libs is libs and libs or libs and libs or libs or libs and libs and
				libs and lib)
				{
				test = 4
				}
			}
		.output()
		}
	process(lib, bil, lib3, lib4, lib5, lib6)
		{
		QueryApply(lib $ ' where name !~ 'Test$'', group: -1)
			{|x|
			if (bil and lib3 or lib4 and x and lib5 or lib6 and x or bil and
				lib3 or lib4 and lib5 or lib6 and bil and lib5 and lib6)
			.process1(lib, x)
			}
		}
	process1(lib, x, var3, var4)
		{
		last = -999
		lines = x.text.Lines().Map!(#Trim)
		for (i = 0; i < lines.Size() - .N; ++i)
			{
			// don't start block with blank or '}' or ')' line
			if lines[i] is '' or lines[i] is '}' or lines[i] is ')'
				continue
			name = lib $ ':' $ x.name $ ':' $ (i + 1)
			hash = .hash(lines, i)
			if not .hashes.Member?(hash)
				.hashes[hash] = name
			else if i - last >= .N
				{
				last = i
				.hashes[hash] $= ', ' $ name
				}
			}
		}
	// hash .N non-blank lines
	hash(lines, i, var3, var4, x5, y6, y7)
		{
		hasher = Adler32()
		for (j = i, n = 0; n < .N and j < lines.Size(); ++j)
			if lines[j] isnt '' // skip blank lines
				{
				hasher.Update(lines[j])
				++n
				if (var3 and var4 or var4 and x5 and y6 or var3 and x5 or y7 and
				var3 or var4 and y6 or var3 and var3 and var3 and x5 or y7)
					y6
				}
			try
				{
				x5 = 5
				}
			catch (e)
				{
				x5 = 5
				}
		return hasher.Value()
		}
	output(l,m,n,o,p)
		{
		dups = .hashes.Values().Filter({ it.Has?(',') })
		dups.Sort!().Each { Print(it.Tr(',', '\t')) }
		Print(dups.Size())
		}
	PubMeth7(){ }
	privMeth8(){ }
	PubMeth9(){ }
	privMeth10(){ }
	}",
		fullWarnings: (warnings: (
				[name: "McCabe complexity is 21 for hash"],
				[name: "McCabe complexity is 16 for process"],
				[name: "McCabe complexity is 11 for Detect"],
				[name: "McCabe complexity is 7 for process1"],
				[name: "McCabe complexity is 2 for CallClass"],
				[name: "McCabe complexity is 1 for output"],
				[name: "McCabe complexity is 1 for PubMeth7"],
				[name: "McCabe complexity is 1 for privMeth8"],
				[name: "McCabe complexity is 1 for PubMeth9"],
				[name: "McCabe complexity is 1 for privMeth10"]),
			rating: 0,
			desc: "McCabe function complexity - Rating affected -> Limit to 8"),
		warnings: (warnings: (
				[name: "stdlib:className:55 - McCabe complexity is 21 for hash"],
				[name: "stdlib:className:25 - McCabe complexity is 16 for process"],
				[name: "stdlib:className:11 - McCabe complexity is 11 for Detect"]),
			rating: 0,
			desc: "McCabe function complexity - Rating affected -> Limit to 8"),
		lineWarnings: ((55), (25), (11))),

	(code: `function (x,y,z)
		{
		if (x is 5 and y is 5 and z is 5 and x is 10 or x is 8 or y is 99 or z is 10
			and true)
			x = x + y
		z = z + x
		return x + y + z
		}`,
		fullWarnings: (warnings: ([name: "Function McCabe complexity: 9"]),
			rating: 0,
			desc: "McCabe function complexity - Rating affected -> Limit to 8")
		warnings: (
			warnings: ([name: "stdlib:className:1 - Function McCabe complexity: 9"]),
			rating: 0,
			desc: "McCabe function complexity - Rating affected -> Limit to 8"),
		lineWarnings: ((1))),

	(code: "// Copyright (C) 2013 Suneido Software Corp.
			All rights reserved worldwide.
	class
		{
		PubMeth1(x)
			{
			switch (x)
				{
			case 1: DoSomething()
			case 1: DoSomething()
			case 1: DoSomething()
			case 1: DoSomething()
			case 1: DoSomething()
			case 1: DoSomething()
			case 1: DoSomething()
			case 1: DoSomething()
			case 1: DoSomething()
			case 1: DoSomething()
			case 1: DoSomething()
			case 1: DoSomething()
				}
			}
		PubMeth2(x,y,z)
			{
			if x or y or z or x or y or z or x or y
			}
		PubMeth3(){ }
		PubMeth4(){ }
		PubMeth5(){ }
		PubMeth6(){ }
		PubMeth7(){ }
		PubMeth8(){ }
		PubMeth9(){ }
		PubMeth10(){ }
		}",
		fullWarnings: (warnings: (
				[name: "McCabe complexity is 13 for PubMeth1"],
				[name: "McCabe complexity is 9 for PubMeth2"],
				[name: "McCabe complexity is 1 for PubMeth3"],
				[name: "McCabe complexity is 1 for PubMeth4"],
				[name: "McCabe complexity is 1 for PubMeth5"],
				[name: "McCabe complexity is 1 for PubMeth6"],
				[name: "McCabe complexity is 1 for PubMeth7"],
				[name: "McCabe complexity is 1 for PubMeth8"],
				[name: "McCabe complexity is 1 for PubMeth9"],
				[name: "McCabe complexity is 1 for PubMeth10"]
				),
			rating: 4,
			desc: "McCabe function complexity - Rating affected -> Limit to 8")
		warnings: (warnings: (
				[name: "stdlib:className:5 - McCabe complexity is 13 for PubMeth1"],
				[name: "stdlib:className:23 - McCabe complexity is 9 for PubMeth2"]),
			rating: 4,
			desc: "McCabe function complexity - Rating affected -> Limit to 8"),
		lineWarnings: ((5), (23))),

	(code: "class { }",
		fullWarnings: (warnings: (), rating: 5,
			desc: "No methods found to run McCabe complexity on"),
		warnings: (warnings: (), rating: 5, desc: ""),
		lineWarnings: ()),

	(code: "function () { return 0 }",
		fullWarnings: (warnings: ([name: "Function McCabe complexity: 1"]), rating: 5,
			desc: "McCabe function complexity - Rating not affected -> Limit to 8"),
		warnings: (warnings: (), rating: 5, desc: ""),
		lineWarnings: ())
	)

	Test_McCabeComplexity()
		{
		recordData = Record(recordName: "className", lib: "stdlib")
		for test in .tests
			{
			recordData.code = test.code
			fullWarnings = Qc_FunctionComplexity(recordData, minimizeOutput?: false)
			lineWarnings = fullWarnings.Extract('lineWarnings')
			Assert(fullWarnings is: test.fullWarnings)
			Assert(lineWarnings isSize: 0)

			recordData.code = test.code
			warnings = Qc_FunctionComplexity(recordData, minimizeOutput?:)
			lineWarnings = warnings.Extract('lineWarnings')
			Assert(warnings is: test.warnings)
			Assert(lineWarnings is: test.lineWarnings)
			}
		}
	}