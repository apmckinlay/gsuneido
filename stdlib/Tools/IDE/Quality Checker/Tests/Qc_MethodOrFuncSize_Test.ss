// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	arr_classes: (
		(code:
"// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
class
	{
	N: 8 // block size (number of lines)------------------
	CallClass(libs)
		{
		if String?(libs or
				not libs)
			libs = [libs]
		(new this).Detect(libs)
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		}
	Detect(libs)
		{
		.hashes = Object()
		for lib in libs
			.process(lib)
		.output()
		}
	process(lib)
		{
		QueryApply(lib $ ' where name !~ 'Test$'', group: -1)
			{|x|
			.process1(lib, x)
			}
		}
	process1(lib, x)
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
			x = 5
			x = 22
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
	hash(lines, i)
		{
		hasher = Adler32()
		for (j = i, n = 0; n < .N and j < lines.Size(); ++j)
			if lines[j] isnt '' // skip blank lines
				{
				hasher.Update(lines[j])
				++n
				}
		return hasher.Value()
		}
	output()
		{
		dups = .hashes.Values().Filter({ it.Has?(',') })
		dups.Sort!().Each { Print(it.Tr(',', '\t')) }
		Print(dups.Size())
		}
	PubMeth7(){}
	PubMeth8(){}
	PubMeth9(){}
	PubMeth10(){}
	}",
		warnings: (warnings: ([name: "stdlib:className:5 - 61 lines in CallClass"]),
			desc: "Method sizes - Rating affected -> Limit to 40 lines", rating: 2),
		fullWarnings: (warnings: ([name: "61 lines in CallClass"],
				[name: "21 lines in process1"],
				[name: "11 lines in hash"],
				[name: "7 lines in Detect"],
				[name: "7 lines in process"],
				[name: "6 lines in output"],
				[name: "1 lines in PubMeth7"],
				[name: "1 lines in PubMeth8"],
				[name: "1 lines in PubMeth9"],
				[name: "1 lines in PubMeth10"]),
			desc: "Method sizes - Rating affected -> Limit to 40 lines", rating: 2),
		lineWarnings: ((5))),

	(code:
"// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.copyright
	//
			// with comment
// pass in an object for warnings and/or size_warnings to get those results
class
	{
	CallClass(code, name = false, lib = false,
		warnings = false, size_warnings = false)
		{
		if .skip?(code, name)
			return true

		result = .compile(code, name, w = Object())

		if result isnt true
			return result
		if warnings isnt false
			result = .handle_warnings(w, code, name, lib, result, warnings)  +
				.handle_warnings (wwwww)
		if size_warnings isnt false
			.check_sizes(code, size_warnings)
		result = .extra_checks(code, result, warnings)
		x+3
		x+3
		x+3
		x++
		++x
		return result
		}
	compile(code, name, w)
		{
		// won't count commented lines
		x = 5
		type = LibRecordType(code) a
		if name isnt false and (type is 'class' or type is 'function')
			{
			pat = '\<(?q)_' $ name $ '(?-q)\>'
			if code.Size() > uname_pos = code.FindRx(pat)
				/* */
				{
				x
				x+2
				w.Add(uname_pos)
				code = code.Replace(pat, ' ' $ name) // replacement must be same
					length--------
				}
			}
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		try
			code.Compile(w)
		catch (err)
			{
			line = err.Prefix?('syntax error at line ')
				? Number(err.Extract('syntax error at line ([0-9]+)')) - 1
				: 0
			return Object(err, line)
			}
		x = 1
		x = 2
		x = 3
		x = 4
		return true
		}
	skip?(code, name)
		{
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		}
	isWeb?(name)
		{
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		x
		}

	check_sizes(code, size_warnings)
		{
		.check_method_sizes(code, size_warnings)
		.check_line_lengths(code, size_warnings)
		}
	check_method_sizes(code, size_warnings)
		{
		for x in ClassHelp.MethodSizes(code)
			if x.lines > .max_method_size
				size_warnings.Add(x)
		}
	check_line_lengths(code, size_warnings)
		{
		line_number = 0
		for line in code.Lines()
			{
			if line.Detab().RightTrim().Size() > .MaxLineLength and not
				line.Trim().Prefix?('catch')
				size_warnings.Add(line_number)
				{


				}
			}
		}
	PubMeth8(){}
	privMeth9(){}
	privMeth10(){}
	}",
		warnings: (warnings: ([name: "stdlib:className:73 - 61 lines in skip?"],
				[name: "stdlib:className:134 - 61 lines in isWeb?"],
				[name: "stdlib:className:30 - 41 lines in compile"]),
			desc: "Method sizes - Rating affected -> Limit to 40 lines", rating: 0),
		fullWarnings: (warnings: ([name: "61 lines in skip?"],
				[name: "61 lines in isWeb?"],
				[name: "41 lines in compile"],
				[name: "21 lines in CallClass"],
				[name: "12 lines in check_line_lengths"],
				[name: "6 lines in check_method_sizes"],
				[name: "5 lines in check_sizes"],
				[name: "1 lines in PubMeth8"],
				[name: "1 lines in privMeth9"],
				[name: "1 lines in privMeth10"]),
			desc: "Method sizes - Rating affected -> Limit to 40 lines", rating: 0),
		lineWarnings: ((73), (134), (30))),

	(code:
"// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Title: 'Refactor'
	Xstretch: 1
	Ystretch: 1
	New(aRefactor, libview)
		{
		super(.controls(aRefactor))
		library = libview.CurrentTable()
		}
	init()
		{
		if .ref.Init(.data) is false
			.On_Cancel()
		}
	controls(aRefactor)
		{
		.ref = aRefactor
		header = Object('Title', .ref.Name)
		}
	constructed: false
	Edit_Change()
		{
		if .constructed isnt false or true or false or true or false or true or falseeeeee
			.Data.HandleFocus()
		}
	On_OK()
		{
		data = .Data.Get()
		}
	PubMeth6(){ }
	PubMeth7(){ }
	PubMeth8(){ }
	PubMeth9(){ }
	PubMeth10(){ }

	}",
		warnings: (warnings: (), desc: "", rating: 5),
		fullWarnings: (warnings: ([name: "5 lines in New"],
				[name: "5 lines in init"],
				[name: "5 lines in controls"],
				[name: "5 lines in Edit_Change"],
				[name: "4 lines in On_OK"],
				[name: "1 lines in PubMeth6"],
				[name: "1 lines in PubMeth7"],
				[name: "1 lines in PubMeth8"],
				[name: "1 lines in PubMeth9"],
				[name: "1 lines in PubMeth10"]),
			desc: "Method sizes - Rating not affected -> Limit to 40 lines", rating: 5),
		lineWarnings: ()),

	(code: "class
	{
	}",
		warnings: (warnings: (), desc: "", rating: 5),
		fullWarnings: (warnings: (),
			desc: "No methods found to check method sizes", rating: 5),
		lineWarnings: ()),

	(code: "function ()
	{
	return 1
	}",
		warnings: (warnings: (), desc: "", rating: 5),
		fullWarnings: (warnings: ([name: "4 lines in function"]),
			desc: "Method sizes - Rating not affected -> Limit to 40 lines", rating: 5),
		lineWarnings: ()),

	(code: "//hello
// test comment
function ()
	{
	x
	x
	x
	x
	x
	x
	x
	x
	x
	x
	x
	x
	x
	x
	x
	x
	x
	x
	x
	x
	x
	x
	x
	x
	x
	x
	x
	x
	x
	x
	x
	x
	x
	x
	x
	x
	x
	x
	x
	x
	x
	x
	x
	x
	x
	return 1
	}",
		warnings: (warnings: ([name: "stdlib:className:3 - 49 lines in function"]),
			desc: "Method sizes - Rating affected -> Limit to 40 lines", rating: 0),
		fullWarnings: (warnings: ([name: "49 lines in function"]),
			desc: "Method sizes - Rating affected -> Limit to 40 lines", rating: 0),
		lineWarnings: ((3)))
	)

	Test_MethodOrFuncSize()
		{
		recordData = Record(recordName: "className", lib: "stdlib")
		for test in .arr_classes
			{
			recordData.code = test.code
			warnings = Qc_MethodOrFuncSize(recordData, minimizeOutput?:)
			lineWarnings = warnings.Extract('lineWarnings')
			Assert(warnings is: test.warnings)
			Assert(lineWarnings is: test.lineWarnings)

			fullWarnings = Qc_MethodOrFuncSize(recordData, lineWarnings,
				minimizeOutput?: false)
			lineWarnings = fullWarnings.Extract('lineWarnings')
			Assert(fullWarnings is: test.fullWarnings)
			Assert(lineWarnings isSize: 0)
			}
		}
	}