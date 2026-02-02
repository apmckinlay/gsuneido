// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// Tests for methods in Strings, see also StringTest
// SuJsWebTest
Test
	{
	Test_As()
		{
		Assert("String".As("This is text") is: "This is text")
		}
	Test_Lines()
		{
		Assert("one\ntwo\r\nthree\n\n".Lines() is: #(one, two, three, ''))
		}
	Test_LineFromPosition()
		{
		Assert("".LineFromPosition(0) is: 0)
		Assert("xxx".LineFromPosition(1) is: 0)
		Assert("one\ntwo".LineFromPosition(6) is: 1)
		Assert("one\ntwo".LineFromPosition(99) is: 1)
		}
	Test_RemoveBlankLines()
		{
		Assert("".RemoveBlankLines() is: "")
		Assert("abc".RemoveBlankLines() is: "abc")
		Assert("abc\ndef".RemoveBlankLines() is: "abc\ndef")
		Assert("abc\n\ndef".RemoveBlankLines() is: "abc\ndef")
		Assert("abc\n   \ndef".RemoveBlankLines() is: "abc\ndef")
		Assert("abc\n   \n\n\n\ndef".RemoveBlankLines() is: "abc\ndef")
		Assert("abc\r\n\r\ndef".RemoveBlankLines() is: "abc\r\ndef")
		Assert("\r\n\r\n".RemoveBlankLines() is: "")
		}
	Test_RemovePrefix()
		{
		Assert("".RemovePrefix("") is: "")
		Assert("".RemovePrefix("foo") is: "")
		Assert("foo".RemovePrefix("") is: "foo")
		Assert("foo".RemovePrefix("foo") is: "")
		Assert("foobar".RemovePrefix("foo") is: "bar")
		Assert("bar".RemovePrefix("foo") is: "bar")
		Assert("barfoo".RemovePrefix("foo") is: "barfoo")
		}
	Test_RemoveSuffix()
		{
		Assert("".RemoveSuffix("") is: "")
		Assert("".RemoveSuffix("foo") is: "")
		Assert("foo".RemoveSuffix("") is: "foo")
		Assert("foo".RemoveSuffix("foo") is: "")
		Assert("barfoo".RemoveSuffix("foo") is: "bar")
		Assert("bar".RemoveSuffix("foo") is: "bar")
		Assert("foobar".RemoveSuffix("foo") is: "foobar")
		}
	Test_Trim()
		{
		Assert(" \t\n".Trim() is: "")
		Assert("hello".Trim() is: "hello")
		Assert(" hello".Trim() is: "hello")
		Assert("hello  ".Trim() is: "hello")
		Assert("	hello   ".Trim() is: "hello")
		}
	Test_RightTrim()
		{
		Assert(" \t\n".RightTrim() is: "")
		Assert("hello".RightTrim() is: "hello")
		Assert(" hello".RightTrim() is: " hello")
		Assert("hello  ".RightTrim() is: "hello")
		Assert("	hello   ".RightTrim() is: "	hello")
		}
	Test_LeftTrim()
		{
		Assert(" \t\n".LeftTrim() is: "")
		Assert(" hello".LeftTrim() is: "hello")
		Assert("hello  ".LeftTrim() is: "hello  ")
		Assert("	hello   ".LeftTrim() is: "hello   ")
		}
	Test_SplitCSV()
		{
		line = '"one",two,3, " f,o,u,r","""five""","56"'

		// split with implicit fields
		record = line.SplitCSV()
		Assert(record[0] is: 'one')
		Assert(record[1] is: 'two')
		Assert(record[2] is: 3) // number
		Assert(record[3] is: ' f,o,u,r')
		Assert(record[4] is: '"five"')
		Assert(record[5] is: '56') // string

		// split with explicit fields
		record = line.SplitCSV( #(five, four, three, two, one))
		Assert(record['five'] is: 'one')
		Assert(record['four'] is: 'two')
		Assert(record['three'] is: 3)
		Assert(record['two'] is: ' f,o,u,r')
		Assert(record['one'] is: '"five"')

		// test string vals option (doesn't convert strings to numbers)
		line = 'val1,2,3,"4","five"'
		record = line.SplitCSV(string_vals:)
		Assert(record[0] is: 'val1')
		Assert(record[1] is: '2')
		Assert(record[2] is: '3')
		Assert(record[3] is: '4')
		Assert(record[4] is: 'five')
		}
	Test_SplitFixedLength()
		{
		// simply calls FixedLength.Split - see FixedLength_Test
		}
	Test_ReplaceSubstr()
		{
		s = "0123456789"
		Assert(s.ReplaceSubstr(0, 0, 'abc') is: 'abc0123456789')
		Assert(s.ReplaceSubstr(0, 3, 'abc') is: 'abc3456789')
		Assert(s.ReplaceSubstr(4, 0, 'abc') is: '0123abc456789')
		Assert(s.ReplaceSubstr(4, 3, 'abc') is: '0123abc789')
		}
	Test_Unescape()
		{
		Assert("".Unescape() is: "")
		Assert("hello".Unescape() is: "hello")
		Assert(`hello\n\tworld`.Unescape() is: "hello\n\tworld")
		Assert(`\t`.Unescape() is: "\t")
		Assert(`\n`.Unescape() is: "\n")
		Assert(`\*`.Unescape() is: "\*")
		}
	Test_SplitOnFirst()
		{
		Assert(''.SplitOnFirst() is: #('',''))
		Assert('hello'.SplitOnFirst() is: #('hello',''))
		Assert('hello world'.SplitOnFirst() is: #('hello','world'))
		Assert('hello there world'.SplitOnFirst() is: #('hello','there world'))
		Assert('hello/there/world'.SplitOnFirst('/') is: #('hello','there/world'))
		Assert('hello//there//world'.SplitOnFirst('//') is: #('hello', 'there//world'))
		}
	Test_SplitOnLast()
		{
		Assert(''.SplitOnLast() is: #('',''))
		Assert('hello'.SplitOnLast() is: #('hello',''))
		Assert('hello world'.SplitOnLast() is: #('hello','world'))
		Assert('hello there world'.SplitOnLast() is: #('hello there','world'))
		Assert('hello/there/world'.SplitOnLast('/') is: #('hello/there','world'))
		Assert('hello//there//world'.SplitOnLast('//') is: #('hello//there', 'world'))
		}
	Test_BeforeFirst()
		{
		Assert(''.BeforeFirst(' ') is: '')
		Assert('hello'.BeforeFirst(' ') is: 'hello')
		Assert('hello world'.BeforeFirst(' ') is: 'hello')
		Assert('hello there world'.BeforeFirst(' ') is: 'hello')
		Assert('hello/there/world'.BeforeFirst('/') is: 'hello')
		Assert('hello//there//world'.BeforeFirst('//') is: 'hello')
		}
	Test_AfterFirst()
		{
		Assert(''.AfterFirst(' ') is: '')
		Assert('hello'.AfterFirst(' ') is: '')
		Assert('hello world'.AfterFirst(' ') is: 'world')
		Assert('hello there world'.AfterFirst(' ') is: 'there world')
		Assert('hello/there/world'.AfterFirst('/') is: 'there/world')
		Assert('hello//there//world'.AfterFirst('//') is: 'there//world')
		}
	Test_Numeric?()
		{
		Assert(not "".Numeric?())
		Assert("1".Numeric?())
		Assert("1234567890".Numeric?())
		Assert(not "123xyz".Numeric?())
		}
	Test_Alpha?()
		{
		Assert(not "".Alpha?())
		Assert("a".Alpha?())
		Assert("Abracadabra".Alpha?())
		Assert(not "123xyz".Alpha?())
		}
	Test_AlphaNum?()
		{
		Assert(not "".AlphaNum?())
		Assert("a".AlphaNum?())
		Assert("1".AlphaNum?())
		Assert("123xyZ".AlphaNum?())
		Assert(not "123 xyz".AlphaNum?())
		}
	Test_Lower?()
		{
		Assert(not "".Lower?())
		Assert(not "?".Lower?())
		Assert(not "X".Lower?())
		Assert("a".Lower?())
		Assert(not "Abracadabra".Lower?())
		Assert("123xyz".Lower?())
		}
	Test_Upper?()
		{
		Assert(not "".Upper?())
		Assert(not "?".Upper?())
		Assert(not "x".Upper?())
		Assert("X".Upper?())
		Assert(not "Abracadabra".Upper?())
		Assert("123XYZ".Upper?())
		}
	Test_LeftFill()
		{
		Assert("abc".LeftFill(0) is: "abc")
		Assert("abc".LeftFill(2) is: "abc")
		Assert("abc".LeftFill(3) is: "abc")
		Assert("abc".LeftFill(4) is: " abc")
		Assert("abc".LeftFill(5, '*') is: "**abc")
		}
	Test_TruncateLeftFill()
		{
		Assert("abc".TruncateLeftFill(0) is: "")
		Assert("abc".TruncateLeftFill(2) is: "ab")
		Assert("abc".TruncateLeftFill(3) is: "abc")
		Assert("abc".TruncateLeftFill(4) is: " abc")
		Assert("abc".TruncateLeftFill(5, '*') is: "**abc")
		}
	Test_RightFill()
		{
		Assert("abc".RightFill(0) is: "abc")
		Assert("abc".RightFill(2) is: "abc")
		Assert("abc".RightFill(3) is: "abc")
		Assert("abc".RightFill(4) is: "abc ")
		Assert("abc".RightFill(5, '*') is: "abc**")
		}
	Test_TruncateRightFill()
		{
		Assert("abc".TruncateRightFill(0) is: "")
		Assert("abc".TruncateRightFill(2) is: "ab")
		Assert("abc".TruncateRightFill(3) is: "abc")
		Assert("abc".TruncateRightFill(4) is: "abc ")
		Assert("abc".TruncateRightFill(5, '*') is: "abc**")
		}
	Test_Center()
		{
		Assert("abc".RightFill(0) is: "abc")
		Assert("abc".RightFill(2) is: "abc")
		Assert("abc".RightFill(3) is: "abc")
		Assert("abc".Center(4) is: "abc ")
		Assert("abc".Center(5) is: " abc ")
		Assert("abc".Center(7, '*') is: "**abc**")
		}
	Test_Iter()
		{
		ob = Object()
		for (c in "abc")
			ob.Add(c)
		Assert(ob is: #(a, b, c))
		}
	Test_Count()
		{
		Assert("abcd".Count('a') is: 1)
		Assert("abcd".Count('ab') is: 1)
		Assert("abcd".Count('abc') is: 1)
		Assert("abcd".Count('d') is: 1)
		Assert("abcd".Count('cd') is: 1)
		Assert("abcdabcd".Count('abcd') is: 2)
		Assert("abcdabcd".Count('a') is: 2)
		Assert("abcdabcd".Count('ab') is: 2)
		Assert("abcdabcd".Count('abc') is: 2)
		Assert("abcdabcd".Count('d') is: 2)
		Assert("abcdabcd".Count('cd') is: 2)
		Assert("abcdabcd".Count('abcd') is: 2)
		Assert("abcdabcd".Count('da') is: 1)
		Assert("abcdabcd".Count('abcda') is: 1)
		Assert("abcdabcd".Count('dabcd') is: 1)
		}
	Test_Capitalize()
		{
		Assert('joe'.Capitalize() is: 'Joe')
		Assert('Joe'.Capitalize() is: 'Joe')
		Assert('joE'.Capitalize() is: 'Joe')
		Assert('JOE'.Capitalize() is: 'Joe')
		}
	Test_Capitalized?()
		{
		Assert(not ''.Capitalized?())
		Assert(not '@#$'.Capitalized?())
		Assert(not 'aBc'.Capitalized?())
		Assert('Abc'.Capitalized?())
		}
	Test_CapitalizeWords()
		{
		Assert('FRED FLINSTONE'.CapitalizeWords() is: 'Fred Flinstone')
		Assert('PO BOX 773'.CapitalizeWords() is: 'PO Box 773')
		Assert('123 1ST AVE NE'.CapitalizeWords() is: '123 1st Ave NE')
		Assert('A & R LOGISTICS'.CapitalizeWords() is: 'A & R Logistics')
		Assert('A-R LOGISTICS'.CapitalizeWords() is: 'A-R Logistics')
		Assert('A.R LOGISTICS'.CapitalizeWords() is: 'A.R Logistics')
		Assert('FOOD-4-LESS-LOS BANOS'.CapitalizeWords() is: 'Food-4-Less-Los Banos')
		Assert('ATTN: A/P'.CapitalizeWords() is: 'Attn: A/P')
		Assert('FRED FLINSTONE'.CapitalizeWords(false) is: 'FRED FLINSTONE')
		Assert('fred flintStone'.CapitalizeWords(false) is: 'Fred FlintStone')
		Assert("fred flintStone's book".CapitalizeWords() is: "Fred Flintstone's Book")
		}
	Test_WrapLines()
		{
		Assert("Test".WrapLines(10) is: #(Test))
		Assert("Test\nLine2".WrapLines(10) is: #(Test Line2))
		Assert("This is a test".WrapLines(5) is: #("This" "is a" "test"))
		Assert("This is another test.\nMultiline wrapping test".WrapLines(10)
			is: #("This is" "another" "test." "Multiline" "wrapping" "test"))
		Assert("Testingwrapwithnospaces".WrapLines(5) is: #(Testi ngwra pwith nospa ces))
		}
	Test_Divide()
		{
		Assert("".Divide(1) is: #())
		Assert("".Divide(5) is: #())
		Assert("test".Divide(5) is: #(test))
		Assert("testing1234567890".Divide(5) is: #('testi' 'ng123' '45678' '90'))
		Assert("_".Divide(1) is: #('_'))
		Assert("abc".Divide(1) is: #(a, b, c))
		}
	Test_White?()
		{
		Assert(not "".White?())
		Assert(not "x".White?())
		Assert(not " x ".White?())
		Assert(" ".White?())
		Assert("\n\r\t ".White?())
		}
	Test_Blank?()
		{
		Assert("".Blank?())
		Assert(" ".Blank?())
		Assert("\n\r\t ".Blank?())
		Assert(not "x".Blank?())
		Assert(not " x ".Blank?())
		}
	Test_LineCount()
		{
		Assert("".LineCount() is: 0)
		Assert("hello".LineCount() is: 1)
		Assert("hello\r\n".LineCount() is: 1)
		Assert("hello\r\nworld".LineCount() is: 2)
		Assert("hello\r\nworld\n".LineCount() is: 2)
		}
	Test_ToHex()
		{
		Assert(''.ToHex() is: '')
		Assert('abc'.ToHex() is: '616263')
		Assert('\x00\xff'.ToHex() is: '00ff')
		}
	Test_FromHex()
		{
		// SuJsWebTest Excluded
		Assert(''.FromHex() is: '')
		Assert('00ff'.FromHex() is: '\x00\xff')
		Assert({ 'hello'.FromHex() } throws:)
		Assert('abcd'.FromHex() is: '\xab\xcd')
		Assert('ABCD'.FromHex() is: '\xab\xcd')
		}
	Test_GlobalName?()
		{
		Assert("X".GlobalName?())
		Assert("Fred".GlobalName?())
		Assert("Fred_Flintstone".GlobalName?())
		Assert("Fred?".GlobalName?())
		Assert("Fred!".GlobalName?())
		Assert(not "fred".GlobalName?())
		Assert(not "Fred Flintstone".GlobalName?())
		}
	Test_Identifier?()
		{
		Assert("X".Identifier?())
		Assert("x".Identifier?())
		Assert("Fred".Identifier?())
		Assert("fred".Identifier?())
		Assert("Fr3d".Identifier?())
		Assert("fr3d".Identifier?())
		Assert("Fred_Flintstone".Identifier?())
		Assert("fred_flintstone".Identifier?())
		Assert("Fred?".Identifier?())
		Assert("fred?".Identifier?())
		Assert("Fred!".Identifier?())
		Assert("fred!".Identifier?())
		Assert(not "?".Identifier?())
		Assert(not "_".Identifier?())
		Assert(not "0".Identifier?())
		Assert(not "0xyz".Identifier?())
		Assert(not "Fred Flintstone".Identifier?())
		}
	Test_DynamicName?()
		{
		Assert("_x".DynamicName?())
		Assert("_xyz".DynamicName?())
		Assert("_xyz?".DynamicName?())
		Assert("_xyz!".DynamicName?())
		Assert("x".DynamicName?() is: false)
		Assert("xyz".DynamicName?() is: false)
		Assert("Xyz".DynamicName?() is: false)
		}
	Test_ChangeEol()
		{
		Assert("".ChangeEol('\n') is: '')
		Assert("one".ChangeEol('\n') is: 'one')
		Assert("one\ntwo\r\n".ChangeEol('!') is: 'one!two!')
		}
	Test_ForEachMatch()
		{
		x = []
		"ab cab abc def".ForEachMatch("ab") {|m| x.Add(m) }
		Assert(x is: #(((0, 2)), ((4, 2)), ((7, 2))))
		x = []
		i = 0
		"ab cab abc def".ForEachMatch("ab")
			{|m|
			++i
			if i is 1
				continue
			x.Add(m)
			if i is 2
				break
			}
		Assert(x is: #(((4, 2))))
		}
	Test_SafeEval()
		{
		// SuJsWebTest Excluded
		Assert(''.SafeEval() is: '')
		Assert('123456789.987654321'.SafeEval() is: 123456789.987654321)
		Assert('true'.SafeEval())
		Assert('false'.SafeEval() is: false)
		Assert('#20110101'.SafeEval() is: #20110101)
		Assert('#20110101.112233444'.SafeEval() is: #20110101.112233444)
		Assert('#(a: 1, b: "Qq")'.SafeEval() is: #(a: 1, b: "Qq"))
		Assert('#(a: 1, b: "Qq") '.SafeEval() is: #(a: 1, b: "Qq"))
		Assert('#{a: 1, b: "Qq"}'.SafeEval() is: #{a: 1, b: "Qq"})
		Assert('[a: 1, b: "Qq"]'.SafeEval() is: [a: 1, b: "Qq"])
		Assert('"qwert"'.SafeEval() is: "qwert")
		Assert("'qwert'".SafeEval() is: "qwert")
		Assert('Date'.SafeEval() is: 'Date')
		Assert({'1+1'.SafeEval() } throws: 'invalid SafeEval')
		Assert({'Date()'.SafeEval() } throws: 'invalid SafeEval')
		}
	Test_ToUtf8()
		{
		// SuJsWebTest Excluded
		Assert('hello'.ToUtf8() is: 'hello')
		Assert('abc\x9a123\x80'.ToUtf8() is: 'abc\xc5\xa1123\xe2\x82\xac')
		}
	Test_FromUtf8()
		{
		// SuJsWebTest Excluded
		Assert('hello'.FromUtf8() is: 'hello')
		Assert('abc\xc5\xa1123\xe2\x82\xac'.FromUtf8() is: 'abc\x9a123\x80')
		}
	Test_Map()
		{
		Assert(''.Map(#Upper) is: '')
		Assert('Hello'.Map(#Upper) is: 'HELLO')
		Assert('abc'.Map(function (c) { '(' $ c $ ')' }) is: '(a)(b)(c)')
		Assert('abc'.Map({ '(' $ it $ ')' }) is: '(a)(b)(c)')
		Assert('abc'.Map({ it is 'b' ? '' : it }) is: 'ac')
		}
	Test_Escape()
		{
		Assert(''.Escape() is: '')
		Assert('hello'.Escape() is: 'hello')
		Assert('\t\r\n\x00\x03\xff\u"'.Escape() is: `\t\r\n\x00\x03\xff\\u\"`)
		}
	Test_Ellipsis()
		{
		Assert("".Ellipsis(100) is: "")
		Assert("hello world".Ellipsis(100) is: "hello world")
		Assert("hello world".Ellipsis(11) is: "hello world")
		Assert("hello world!".Ellipsis(11) is: "hello...orld!")
		Assert("hello cruel world".Ellipsis(10) is: "hello...world")
		Assert("hello cruel world".Ellipsis(10, atEnd:) is: "hello crue...")
		}
	Test_Has1of?()
		{
		Assert("".Has1of?("abc") is: false)
		Assert("foo".Has1of?("abc") is: false)
		for s in #(axy, xay, xya)
			for c in #(abc, bac, bca)
				Assert(s.Has1of?(c))
		}
	Test_In?()
		{
		Assert('b'.In?('abc'))
		Assert('x'.In?('abc') is: false)
		Assert('two'.In?('one two three'))
		Assert('Two'.In?('one two three') is: false)
		Assert('a'.In?(#(a, b, c)))
		Assert('x'.In?(#(a, b, c)) is: false)
		}
	Test_UniqueChars()
		{
		Assert("".UniqueChars() is: "")
		Assert("a".UniqueChars() is: "a")
		Assert("aa".UniqueChars() is: "a")
		Assert("abc".UniqueChars() is: "abc")
		Assert("aabc".UniqueChars() is: "abc")
		Assert("axbxcxdxexfx".UniqueChars() is: "axbcdef")
		}
	Test_LineAtPosition()
		{
		s = 'one\ntwo\nthreex'
		Assert(s.LineAtPosition(0) is: 'one')
		Assert(s.LineAtPosition(s.Find('w')) is: 'two')
		Assert(s.LineAtPosition(s.Find('x')) is: 'threex')
		Assert(s.LineAtPosition(99) is: '')
		}
	Test_FirstLine()
		{
		Assert("".FirstLine() is: "")
		Assert("hello".FirstLine() is: "hello")
		Assert("hello\n".FirstLine() is: "hello")
		Assert("hello\r\n".FirstLine() is: "hello")
		Assert("hello\nworld".FirstLine() is: "hello")
		Assert("hello\r\nworld".FirstLine() is: "hello")
		}
	Test_StartPositionOfLine()
		{
		Assert("".StartPositionOfLine(0) is: 0)
		Assert("".StartPositionOfLine(1) is: 0)
		s = "\r\n\taaa\nbbb\r\n"
		Assert(s.StartPositionOfLine(0) is: 0)
		Assert(s.StartPositionOfLine(1) is: 2)
		Assert(s.StartPositionOfLine(2) is: 7)
		Assert(s.StartPositionOfLine(3) is: 12)
		Assert(s.StartPositionOfLine(3) is: 12)
		Assert(s.StartPositionOfLine(4) is: 12)
		}
	Test_ExtractAll()
		{
		s = ""
		rx = ".*"
		Assert(s.ExtractAll(rx) is: #(''))
		s = 'Hello'
		Assert(s.ExtractAll(rx) is: #('Hello'))
		s = "Hello World"
		rx = "(H.*)(W.*)"
		Assert(s.ExtractAll(rx) is: #('Hello World', 'Hello ', 'World'))
		}
	Test_StripInvalidChars()
		{
		s = 'Hello World'
		Assert(s.StripInvalidChars() is: 'Hello World')
		s = 'Hello\x02 World'
		Assert(s.StripInvalidChars() is: 'Hello World')
		s = 'H\x1aello\x02 World'
		Assert(s.StripInvalidChars() is: 'Hello World')
		}
	Test_ForEach1of()
		{
		s = "now is the time for all good men"
		c = "abcde"
		ob = Object()
		s.ForEach1of(c, { ob.Add(it) })
		Assert(ob is: #(9, 14, 20, 27, 30))
		}
	}
