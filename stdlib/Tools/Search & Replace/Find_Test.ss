// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	find: Find
		{
		Find_beep()
			{
			}
		}

	Test_Regex()
		{
		Assert(.find.Regex("", []) is: "")
		Assert(.find.Regex("foo", []) is: "(?i)(?q)foo(?-q)")
		Assert(.find.Regex("foo", [case:]) is: "(?q)foo(?-q)")
		Assert(.find.Regex("foo", [case:, regex:]) is: "foo")
		Assert(.find.Regex("foo", [case:, regex:, word:]) is: "\<foo\>")
		Assert(.find.Regex("foo(", [case:, regex:, word:]) is: "\<foo(")
		Assert(.find.Regex(".foo", [case:, regex:, word:]) is: ".foo\>")
		Assert(.find.Regex(".foo(", [case:, regex:, word:]) is: ".foo(")
		Assert(.find.Regex(`foo\tbar`, [case:, regex:]) is: "foo\tbar") // unescaped
		Assert(.find.Regex('.foo', [case:, regex:, word:]) is: `.foo\>`)
		Assert(.find.Regex('foo(x|y)', [case:, regex:, word:]) is: `\<foo(x|y)\>`)
		}
	Test_Replacement()
		{
		Assert(.find.Replacement('', []) is: '')
		Assert(.find.Replacement('foo', []) is: `\=foo`)
		Assert(.find.Replacement('foo', [regex: false]) is: `\=foo`)
		Assert(.find.Replacement('', [regex:]) is: '')
		Assert(.find.Replacement(`foo\tbar`, [regex:]) is: 'foo\tbar') // unescaped
		}

	text: "function () { abc = 'aabcde' + Abc; return Test_Abc(abc); }"
	Test_DoFind()
		{
		fn = .find.DoFind
		options = [find: 'abc']
		Assert(fn('', 0, options) is: false)
		Assert(fn(.text, 0, options) is: #(14, 3))
		Assert(fn(.text, 15, options) is: #(22, 3))
		Assert(fn(.text, 54, options) is: #(14, 3))

		Assert(fn(.text, 59, options, prev:) is: #(52, 3))
		Assert(fn(.text, 54, options, prev:) is: #(48, 3))
		Assert(fn(.text, 14, options, prev:) is: #(52, 3))

		options = [find: 'abc', word:]
		Assert(fn(.text, 15, options) is: #(31, 3))
		options = [find: 'abc', case:]
		Assert(fn(.text, 15, options) is: #(22, 3))
		options = [find: 'abc', case:, word:]
		Assert(fn(.text, 15, options) is: #(52, 3))

		// test searching by Regular Expression
		options = [find: 'A.+?c']
		Assert(fn(.text, 59, options, prev:) is: false)
		options = [find: 'A.+?c', regex:]
		Assert(fn(.text, 59, options, prev:) is: #(52, 3))
		options = [find: 'A.+?c', regex:, case:]
		Assert(fn(.text, 59, options, prev:) is: #(48, 3))
		options = [find: 'A.+?c', regex:, case:, word:]
		Assert(fn(.text, 59, options, prev:) is: #(31, 3))

		// test searching by Expression
		ast = AstWriteManager(.text)
		options = [find: 'abc', expr:]
		Assert(fn(.text, 0, options) is: false)
		options = [find: "abc", :ast, expr:]
		Assert(fn(.text, 59, options) is: #(14, 3))
		Assert(fn(.text, 15, options) is: #(52, 3))
		options = [find: "a + b", :ast, expr:]
		Assert(fn(.text, 0, options) is: #(20, 14))
		}

	Test_FindAll()
		{
		fn = .find.FindAll
		options = [find: 'abc']
		Assert(fn(.text, options) is: #((14, 3), (22, 3), (31, 3), (48, 3), (52, 3)))

		ast = AstWriteManager(.text)
		options = [find: 'abc', expr:]
		Assert(fn(.text, options) is: false)
		options = [find: "abc", :ast, expr:]
		Assert(fn(.text, options) is: #((14, 3), (52, 3)))
		}

	Test_DoReplace()
		{
		fn = .find.DoReplace
		options = [find: '', replace: 'ABC', case:]
		Assert(fn(.text, .text, 0, 59, options) is: false)

		options = [find: 'efg', replace: 'ABC', case:]
		Assert(fn(.text, .text, 0, 59, options) is: false)

		options = [find: 'Abc', replace: 'ABC', case:]
		Assert(fn(.text, .text, 0, 59, options)
			is: `function () { abc = 'aabcde' + ABC; return Test_ABC(abc); }`)
		Assert(fn(.text, `abc = 'aabcde' + Abc;`, 14, 35, options)
			is: `abc = 'aabcde' + ABC;`)

		ast = AstWriteManager(.text)
		options = [find: 'a+b', replace: '\2 + \1', expr:]
		Assert(fn(.text, .text, 0, 59, options) is: false)
		options = [find: 'a+b', replace: '\2 + \1', :ast, expr:]
		Assert(fn(.text, .text, 0, 59, options)
			like: `function () { abc = Abc + 'aabcde'; return Test_Abc(abc); }`)
		Assert(fn(.text, .text, 0, 25, options) is: false)
		Assert(fn(.text, .text, 25, 59, options) is: false)
		}

	Test_tryIgnoreRegexError()
		{
		fn = Find.Find_tryIgnoreRegexError
		watch = .WatchTable('suneidolog')

		fn({ throw 'regex: error' })
		Assert(.GetWatchTable(watch) isSize: 0)

		fn({ throw 'regexp.cpp: error' })
		Assert(.GetWatchTable(watch) isSize: 0)

		fn({ throw 'error caused by regex' })
		Assert(.GetWatchTable(watch) isSize: 0)

		fn({ throw 'error caused by unforseen reasons' })
		calls = .GetWatchTable(watch)
		Assert(calls isSize: 1)
		Assert(calls[0].sulog_message
			is: 'ERROR: (CAUGHT) Find failed: error caused by unforseen reasons')
		}
	}