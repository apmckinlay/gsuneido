// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_line_continuation()
		{
		.test("for p in Provinces\n{ }")
		.test("for p in F()\n{ }")
		.test("x \n ? y : z")
		}

	Test_if_comment_bug()
		{
		.test("if x < y\n z")
		.test("if x < y\n\n z")
		.test("if x < y\n /**/ z")
		.test("if x < y\n /**/\n z")
		.test("if x < y\n //...\n z")
		}

	test(s)
		{
		("function (x, y, z) {\n" $ s $ "\n}").Compile()
		}

	Test_warnings()
		{
		test = function (s, expected)
			{
			s = "function () { " $ s $ " }"
			s.Compile(warnings = Object())
			Assert(warnings is: expected)
			}
		test("", #())
		test("xyz", #("ERROR: used but not initialized: xyz @14"))
		test("_xyz", #())
		test("_xyz = 123", #())
		test("xyz = 0", #("WARNING: initialized but not used: xyz @14"))
		test("Not__Defined", #("ERROR: can't find: Not__Defined @14"))
		}

	Test_cant_find()
		{
		.MakeLibraryRecord(#(name: CompileTestFunc,
			text: 'function () { _CompileTestFunc() }'))
		Assert({ Global('CompileTestFunc') } throws: "can't find")
		}

	Test_invalid_reference()
		{
		Assert({ .test("_Control") } throws: 'invalid reference to _Control')

		.MakeLibraryRecord(#(name: CompileTestFunc2,
			text: 'function () { _Not_defined() }'))
		Assert({ Global('CompileTestFunc2') }
			throws: "invalid reference to _Not_defined")
		}

	Test_trailing_block_bug()
		{
		Seq(3).Map({|x| x + 1 })
		Seq(3).Map() {|x| x + 1 }
		Seq(3).Map() /* - */ {|x| x + 1 }
		Seq(3).Map()
			{|x| x + 1 }
		Seq(3).Map() // comment
			{|x| x + 1 }
		'Seq(3).Map()
			// comment

			// comment
			{|x| x + 1 }'.Eval() // using Eval to handle old exe
		}
	}