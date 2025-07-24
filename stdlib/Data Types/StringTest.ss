// Copyright (C) 2011 Suneido Software Corp. All rights reserved worldwide.
// Tests for builtin methods, see also StringsTest
// SuJsWebTest
Test
	{
	Test_negative_indexes()
		{
		Assert('abcd'[0] is: 'a')
		Assert('abcd'[3] is: 'd')
		Assert('abcd'[4] is: '')
		Assert('abcd'[-1] is: 'd')
		Assert('abcd'[-4] is: 'a')
		Assert('abcd'[-5] is: '')
		}
	Test_ranges()
		{
		Assert('abcd'[1 .. 3] is: 'bc')
		Assert('abcd'[1 .. 9] is: 'bcd')
		Assert('abcd'[1 ..] is: 'bcd')
		Assert('abcd'[6 .. 9] is: '')
		Assert('abcd'[2 .. 1] is: '')
		Assert('abcd'[-3 .. -1] is: 'bc')
		Assert('abcd'[1 .. -1] is: 'bc')
		Assert('abcd'[.. -2] is: 'ab')
		Assert('abcd'[-2 ..] is: 'cd')

		Assert('abcd'[1 :: 2] is: 'bc')
		Assert('abcd'[:: 2] is: 'ab')
		Assert('abcd'[-2 :: 1] is: 'c')
		Assert('abcd'[1 :: -1] is: '')
		Assert('abcd'[1 :: 9] is: 'bcd')
		Assert('abcd'[1 ::] is: 'bcd')
		Assert('abcd'[9 :: 1] is: '')
		}
	Test_Asc()
		{
		Assert("Alert".Asc() is: 65)
		}
	Test_Compile()
		{
		// SuJsWebTest Excluded
		Assert("123".Compile() is: 123)
		Assert("(1 2)".Compile() is: #(1, 2))
		Assert(Function?("function () { }".Compile()))
		Assert(Class?("class { F() { } }".Compile()))
		}
	Test_MapN()
		{
		Assert(''.MapN(1, { it.Capitalize() }) is: '')
		Assert(''.MapN(3, { it.Capitalize() }) is: '')
		Assert('helloworld'.MapN(1, { it.Capitalize() }) is: 'HELLOWORLD')
		Assert('helloworld'.MapN(2, { it.Capitalize() }) is: 'HeLlOwOrLd')
		Assert('helloworld'.MapN(1, {|unused| '' }) is: '')
		Assert('helloworld'.MapN(3, {|unused| 'x' }) is: 'xxxx')
		Assert('helloworld'.MapN(3, { String(it.Size()) }) is: '3331')
		Assert('helloworld'.MapN(1, { it }) is: 'helloworld')
		Assert('helloworld'.MapN(3, { it }) is: 'helloworld')
		Assert('helloworld'.MapN(999, { it }) is: 'helloworld')
		s = 'helloworld'.Repeat(1000)
		Assert(s.MapN(1, { it }) is: s)
		Assert(s.MapN(3, { it }) is: s)
		}
	Test_StringCall()
		{
		Assert("Size"("hello") is: 5) // equivalent to "hello".Size()
		Assert({ ""() } throws: "string call requires 'this' argument")
		}
	Test_Replace_Single_Arg()
		{
		Assert('hello world'.Replace('o') is: 'hell wrld')
		}
	Test_regex()
		{
		a = 'aZ'
		Assert(a !~ '[a-Z]') // backwards empty range
		}
	Test_regex_case_insensitive()
		{
		a = 'a' // use variable so not evaluated at compile time

		matchesRegardlessOfIgnoreCase = function (str, pat)
			{
			Assert(str =~ pat)
			Assert(str =~ "(?i)" $ pat)
			}
		matchesRegardlessOfIgnoreCase(a, '[[:lower:]]')
		matchesRegardlessOfIgnoreCase(a, '[a-z]')

		matchesOnlyWithIgnoreCase = function (str, pat)
			{
			Assert(str !~ pat)
			Assert(str =~ "(?i)" $ pat)
			}
		matchesOnlyWithIgnoreCase(a, '[A-Z]')
		matchesOnlyWithIgnoreCase(a, '[[:upper:]]')

		// check for bug if ignore case converts range
		Assert(a =~ '(?i)[5-M]')
		Assert(a =~ '(?i)[M-}]')
		}
	Test_Match()
		{
		Assert("hello world".Match('o')[0][0] is: 4)
		Assert("hello world".Match('o', prev:)[0][0] is: 7)
		Assert("hello world".Match('o', 5)[0][0] is: 7)
		Assert("hello world".Match('o', 6, prev:)[0][0] is: 4)
		Assert("hello world".Match('o', 8) is: false)
		Assert("hello world".Match('o', 3, prev:) is: false)

		Assert("Hello World".Match("lo ") is: Object(Object(3,3)))
		Assert("Hello World".Match("t") is: false)
		}
	Test_Find_methods()
		{
		Assert("hello world".Find('o') is: 4)
		Assert("hello world".Find('o', 4) is: 4)
		Assert("hello world".Find('o', 5) is: 7)
		Assert("hello world".FindLast('o') is: 7)
		Assert("hello world".FindLast('o', 7) is: 7)
		Assert("hello world".FindLast('o', 5) is: 4)

		s = "this is a test"
		fail = s.Size()
		Assert(s.Find("") is: 0)
		Assert(s.Find("is") is: 2)
		Assert(s.FindLast("is") is: 5)
		Assert(s.Find("xyz") is: fail)
		Assert(s.FindLast("xyz") is: false)
		Assert(s.Find1of("si") is: 2)
		Assert(s.FindLast1of("si") is: 12)
		Assert(s.Find1of("xy") is: fail)
		Assert(s.FindLast1of("xy") is: false)
		Assert(s.Find1of("^this ") is: 8)
		Assert(s.FindLast1of("^tse") is: 9)
		Assert(s.FindLast1of("^ aeihst") is: false)
		Assert(s.Find("", 2) is: 2)
		Assert("hellohello".Find("lo") is: 3)
		Assert("hellohello".Find("he", 2) is: 5)

		for ..10 // in case of intermittent errors
			{
			// first
			fail = 4
			pos = [-99, 0, 2, 99]
			find = [1, 1, fail, fail]
			for i in pos.Members()
				{
				Assert('abcd'.Find('b', pos[i])					is: find[i])
				Assert('abcd'.Find('b', pos: pos[i])			is: find[i])
				Assert('abcd'.Find('x', pos[i])					is: fail)
				Assert('abcd'.Find('x', pos: pos[i])			is: fail)
				Assert('abcd'.Find1of('b', pos[i])				is: find[i])
				Assert('abcd'.Find1of('b', pos: pos[i])			is: find[i])
				Assert('abcd'.Find1of('x', pos[i]) 				is: fail)
				Assert('abcd'.Find1of('x', pos: pos[i]) 		is: fail)
				Assert('abcd'.Find1of('^acd', pos[i]) 			is: find[i])
				Assert('abcd'.Find1of('^acd', pos: pos[i])		is: find[i])
				Assert('abcd'.Find1of('^abcd', pos[i])			is: fail)
				Assert('abcd'.Find1of('^abcd', pos: pos[i])		is: fail)
				}
			// last
			fail = false
			pos = [-99, 0, 1, 99]
			find = [fail, fail, 1, 1]
			for i in pos.Members()
				{
				Assert('abcd'.FindLast('b', pos[i])					is: find[i])
				Assert('abcd'.FindLast('b', pos: pos[i])			is: find[i])
				Assert('abcd'.FindLast('x', pos[i])					is: fail)
				Assert('abcd'.FindLast('x', pos: pos[i])			is: fail)
				Assert('abcd'.FindLast1of('b', pos[i])				is: find[i])
				Assert('abcd'.FindLast1of('b', pos: pos[i])			is: find[i])
				Assert('abcd'.FindLast1of('x', pos[i]) 				is: fail)
				Assert('abcd'.FindLast1of('x', pos: pos[i]) 		is: fail)
				Assert('abcd'.FindLast1of('^acd', pos[i]) 		is: find[i])
				Assert('abcd'.FindLast1of('^acd', pos: pos[i])	is: find[i])
				Assert('abcd'.FindLast1of('^abcd', pos[i])		is: fail)
				Assert('abcd'.FindLast1of('^abcd', pos: pos[i])	is: fail)
				}
			}

		Assert("foobar".FindLast("") is: 6)
		Assert("foobar".FindLast("", 4) is: 4)
		}

	Test_Number?()
		{
		test = function (s, expected)
			{
			Assert(s.Number?() is: expected msg: s)
			if expected
				{
				Assert(("+" $ s).Number?() msg: "+" $ s)
				Assert(("-" $ s).Number?() msg: "-" $ s)
				Assert(("z" $ s).Number?() is: false msg: "z" $ s)
				Assert((s $ "z").Number?() is: false msg: s $ "z")
				Assert((" " $ s).Number?() is: false msg: "blank+ " $ s)
				Assert((s $ " ").Number?() is: false msg: s $ " +blank")
				}
			}
		test("0", true)
		test("6", true)
		test("007", true)
		test("123", true)
		test("123.", true)
		test(".123", true)
		test("123.465", true)
		test("1e6", true)
		test("1.5e6", true)
		test("1.5e-6", true)
		test("1.5e+6", true)
		test("1.5e-23", true)
		test("123e", false)
		test("123e-", false)
		test("123e+", false)

		test("", false)
		test(".", false)
		test("+", false)
		test("-", false)
		test("-.", false)
		test("+-.", false)
		test("1.2.3", false)
		test("e", false)
		test("e5", false)

		test("1.0", true)

		test("0x", false)
		test("0x123.456", false)
		test("0xZ12", false)

		test("0x123", true)
		test("123_456", true)
		test("1_2_3", true)
		test("1.2_3", true)
		test("0x1_2_3", true)
		}

	Test_Prefix?()
		{
		Assert("".Prefix?(""))
		Assert("hello".Prefix?(""))
		Assert("hello".Prefix?("h"))
		Assert("hello".Prefix?("hel"))
		Assert("hello".Prefix?("hello"))
		Assert(not "hello".Prefix?("x"))
		Assert(not "hello".Prefix?("hex"))
		Assert(not "hello".Prefix?("hellox"))
		Assert("hello".Prefix?("lo", 3))
		Assert("hello world".Prefix?('hell', 0))
		Assert("hello world".Prefix?('or', 7))
		Assert("hello world".Prefix?('hell', -99))
		Assert("hello world".Prefix?('or', -4))
		}
	Test_Suffix?()
		{
		Assert("".Suffix?(""))
		Assert("hello".Suffix?(""))
		Assert("hello".Suffix?("o"))
		Assert("hello".Suffix?("lo"))
		Assert("hello".Suffix?("hello"))
		Assert(not "hello".Suffix?("x"))
		Assert(not "hello".Suffix?("xlo"))
		Assert(not "hello".Suffix?("xhello"))
		}
	Test_Has?()
		{
		Assert("hello" has: "")
		Assert("hello" has: "lo")
		Assert("hello" has: "lo")
		Assert("hello" has: "ell")
		Assert("hello" has: "hello")
		Assert("hello" hasnt: "loX")
		Assert("hello" hasnt: "Xhe")
		Assert("hello" hasnt: "helloX")
		Assert("hello" hasnt: "Xhello")
		Assert("hello" hasnt: "XhelloX")
		Assert("hello" hasnt: "xxx")
		Assert("" has: "")
		}
	Test_Split()
		{
		Assert("one,two,,four".Split(',') is: #('one', 'two', '', 'four'))
		}
	Test_Detab()
		{
		Assert("	test3".Detab() is: "    test3")
		}
	Test_Entab()
		{
		Assert("hello".Entab() is: "hello")
		Assert("  hello".Entab() is: "  hello")
		Assert("    hello".Entab() is: "\thello")
		Assert("      hello".Entab() is: "\t  hello")
		Assert("\thello".Entab() is: "\thello")
		Assert("  \thello".Entab() is: "\thello")
		Assert("\t\thello".Entab() is: "\t\thello")
		Assert("hello  \t".Entab() is: "hello")
		Assert("    test3".Entab() is: "	test3")
		}
	Test_Size()
		{
		Assert("testing" isSize: 7, msg: "test string Size1")
		Assert("\ntesting" isSize: 8, msg: "test string Size2")
		Assert("" isSize: 0, msg: "test string Size3")
		}
	Test_String?()
		{
		Assert(String?(false) is: false)
		Assert(String?(34) is: false)
		Assert(String?("\ntest"))
		}
	Test_Substr()
		{
		Assert("hello world"[-5 ..] is: "world", msg: "test string Substr1")
		Assert("hello world"[3 :: 2] is: "lo", msg: "test string Substr2")
		Assert("hello world"[.. 5] is: "hello", msg: "test string Substr3")
		}
	Test_Tr()
		{
		s = "Hello World"
		Assert(s is: s.Tr(""))
		Assert(s is: s.Tr("a"))
		Assert(s is: s.Tr("\x00"))
		Assert(s is: s.Tr("a" is: "A"))
		Assert(s is: s.Tr("o", "o"))
		Assert("HelloWorld" is: s.Tr(" "))
		Assert("Hello world" is: s.Tr("W", "w"))
		Assert("HELLO WORLD" is: s.Tr("a-z", "A-Z"))
		Assert("H. W." is: s.Tr("a-z", "."))

		s = "hello\x00world"
		Assert(s is: s.Tr(""))
		Assert(s is: s.Tr("a"))
		Assert(s is: s.Tr("a", "A"))
		Assert(s is: s.Tr("o", "o"))
		Assert("helloworld" is: s.Tr("\x00"))

		Assert('hello\xff'.Tr('\x7f-\xff') is: 'hello')
		Assert('hello'.Tr('^\x20-\xff') is: 'hello')
		Assert('hello\x7f'.Tr('\x70-\x7f') is: 'hello')
		}
	Test_Replace()
		{
		Assert("Hello World".Replace("Hello", "GoodBye") is: "GoodBye World")
		Assert("Hello 12345".Replace("[0-9]", "#") is: "Hello #####")
		Assert("abc".Replace('.', { |s| s.Upper() }, 1) is: "Abc")
		Assert("abc".Replace('.', { |s| s.Upper() }) is: "ABC")
		Assert("hello".Replace('he', "\\") is: "\\llo")
		}
	Test_Repeat()
		{
		Assert("test".Repeat(5) is: "testtesttesttesttest", msg: "test string repeat1")
		Assert("test".Repeat(1) is: "test", msg: "test string repeat2")
		Assert("test".Repeat(0) is: "", msg: "test string repeat3")
		}
	Test_Extract()
		{
		Assert("hello world".Extract(".....$") is: "world", msg: "test string Extract1")
		Assert("hello world".Extract("w(..)ld") is: "or", msg: "test string Extract2")
		}
	Test_Eval()
		{
		// SuJsWebTest Excluded
		Assert("true".Eval())
		Assert("false".Eval() is: false)
		Assert("(1 is 1)".Eval())
		Assert("(1 is 2)".Eval() is: false)
		}
	Test_regex_greedy()
		{
		Assert("xxx".Extract("x.?x") is: "xxx")
		Assert("xxx".Extract("x.??x") is: "xx")
		Assert("ab<cd>ef<gh>ij".Extract("<.+>") is: "<cd>ef<gh>")
		Assert("ab<cd>ef<gh>ij".Extract("<.+?>") is: "<cd>")
		Assert("ab<cd>ef<gh>ij".Extract("<.*>") is: "<cd>ef<gh>")
		Assert("ab<cd>ef<gh>ij".Extract("<.*?>") is: "<cd>")
		}
	Test_regex_literal()
		{
		Assert("abc".Replace("..", "&&") is: "ababc")
		Assert("abc".Replace("..", "\=&&") is: "&&c")
		Assert("a..d".Replace("(?q)..", "\=&&") is: "a&&d")
		Assert("too T.o t.on".Replace("(?i)(?q)t.o", "***") is: "too *** ***n")
		Assert("too T.o t.on".Replace("(?i)\<(?q)t.o(?-q)\>", "***") is: "too *** t.on")
		}
	Test_NthLine()
		{
		Assert("".NthLine(0) is: "")
		Assert("".NthLine(5) is: "")
		s = "one\ntwo\r\nthree\nfour"
		Assert(Seq(5).Map({ s.NthLine(it) }) is: #(one, two, three, four, ''))
		}
	Test_string_coerce()
		{
		s = ''; s = '' // prevent executing at compile time
		Assert(s $ true is: 'true')
		Assert(123 $ s is: '123')

		cant = function (block)
			{ Assert(block throws: "can't convert") }
		cant({ class{}.Method?(123) })
		cant({ class{}.Method?(false) })
		cant({ 123 =~ s })
		cant({ s =~ true })
		}
	}