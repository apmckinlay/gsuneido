// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
/*
pos = 26
//       123456789 123456789 123456789
code = "Foo { bar(a,b,c) { 'hello' } }"
Print(' '.Repeat(pos-1) $ '\\/', pos)
Print(code)
if false is range = CodeNest(code, pos)
	return false
Print(' '.Repeat(range[0]) $ '^', range[0])
Print(' '.Repeat(range[1]) $ '^', range[1])
*/
Test
	{
	Test_one()
		{
		code = "Foo { bar(a,b,c) { 'hello' } }"
		test = {|pos, expected| Assert(CodeNest(code, pos) is: expected, msg: pos) }
		test(0, false)
		test(4, false)
		test(5, [4, 29])
		test(7, [4, 29])
		test(11, [9, 15])
		test(18, [17, 27])
		test(19, [17, 27])
		test(20, [19, 25])
		test(22, [19, 25])
		test(25, [19, 25])
		test(26, [17, 27])
		test(27, [17, 27])
		test(28, [4, 29])
		test(29, [4, 29])
		test(30, false)
		test(99, false)
		}
	}