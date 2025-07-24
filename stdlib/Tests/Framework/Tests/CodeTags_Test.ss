// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_ExtractTags()
		{
		test = function (code, tags)
			{ Assert(CodeTags.ExtractTags(code) is: tags) }
		test("", #())
		test("// copyright
			function ()
				{ }", #())
		test("// TAGS: client windows", #(client, windows))
		test("// copyright\n" $
			"// TAGS:  client  windows\n" $
			"123", #(client, windows))
		}
	Test_MatchTags()
		{
		test = function (codetags, expected)
			{ Assert(CodeTags.MatchTags(codetags, _systags) is: expected) }
		_systags = #()
		test(#(), true)
		test(#(foo), false)
		test(#(foo, bar), false)
		_systags = #(windows)
		test(#(), true)
		test(#(windows), true)
		test(#(windows), true)
		test(#(foo), false)
		test(#(foo) false)
		test(#("!foo"), true)
		test(#("!foo !bar"), true)
		}
	}
