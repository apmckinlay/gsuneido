// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		test = function (src, dst)
			{ Assert(FormatCode('{\n' $ src $ '\n}') is: '{\n' $ dst $ '\n}\r\n') }
		test('for(i=0;i<10;++i)', 'for (i = 0; i < 10; ++i)')
		test('hello; \nworld', 'hello\nworld')
		test('{ | x | stuff; }', '{|x| stuff }')
		test('xxx;\n\t}', 'xxx\n\t}')
		test('{ xxx;}', '{ xxx }')
		test('\n{ xxx; return; }', '\n{ xxx; return }')
		test('if xx\n\t;', 'if xx\n\t;')
		test('if xx\n\t;\n', 'if xx\n\t;\n')
		test('for (a;\nb;\nc)', 'for (a;\nb;\nc)')
		test('fn(123); // comment', 'fn(123) // comment')
		test('if xx\n\t; // comment', 'if xx\n\t; // comment')
		test('x = 123; // comment', 'x = 123 // comment')
		test('x = 123;\t// comment', 'x = 123\t// comment')
		test('for(x=0; x<5; x++)\r\nPrint(x); // comment',
			'for (x = 0; x < 5; x++)\r\nPrint(x) // comment')
		test('for(x=0; x<5; x++)\r\nPrint(x);\t\t// comment',
			'for (x = 0; x < 5; x++)\r\nPrint(x)\t\t// comment')
		test('for(x=0; x<5; x++)\r\n;\t\t// comment',
			'for (x = 0; x < 5; x++)\r\n;\t\t// comment')
		}
	Test_curlies_on_separate_lines()
		{
		Assert(FormatCode(s = '\t{ }\r\n') is: s)
		Assert(FormatCode('\t{ xxx }') is: '\t{\r\n\txxx\r\n\t}\r\n')
		Assert(FormatCode('\t{|a, b| xxx }') is: '\t{|a, b|\r\n\txxx\r\n\t}\r\n')
		}
	}