// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		_report = Mock()
		_report.When.PlainText?().Return(false)
		c = WrapFormat
				{
				Data: false
				WrapFormat_measure(line, font /*unused*/)
					{
					return line.Size()
					}
				}
		wf = c()
		wrapLine = { |data, w| wf.WrapDataLines(data, w, font: false) }
		Assert(wrapLine('', 30) is: '')
		Assert(wrapLine('hello', 30) is: 'hello')
		Assert(wrapLine('hello world', 6) is: 'hello\nworld')
		Assert(wrapLine('hello\nworld', 99) is: 'hello\nworld')
		Assert(wrapLine('hello  \n  world', 99) is: 'hello\n  world')
		Assert(wrapLine('helloworld', 6) is: 'hellow\norld')
		Assert(wrapLine('hello,world', 6) is: 'hello,\nworld')
		Assert(wrapLine('hello,1world', 6) is: 'hello,\n1world')
		Assert(wrapLine('hello1,world', 8) is: 'hello1,\nworld')
		Assert(wrapLine('hello 1,2345', 9) is: 'hello\n1,2345')
		Assert(wrapLine('hello world '.Repeat(10), 25)
			is: 'hello world hello world\n'.Repeat(5)[.. -1])

		// test line limit (default is 30 lines)
		s = '12345678901234567890123456789012345678901234567890'
		wrapped_s = '1\n2\n3\n4\n5\n6\n7\n8\n9\n0\n1\n2\n3\n4\n5\n6\n7\n8\n' $
			'9\n0\n1\n2\n3\n4\n5\n6\n7\n8\n9\n...'
		Assert(wrapLine(s, 1) is: wrapped_s)
		}

	Test_Format_data()
		{
		f = WrapFormat { Data: false }.Format_data
		Assert(f(12345) is: "12345")
		Assert(f(#20070727) is: "#20070727")
		Assert(f(false) is: "false")
		Assert(f("") is: "")
		Assert(f("Test") is: "Test")
		Assert(f("   \tTest   ") is: "   \tTest")
		Assert(f("Test\nNext Line") is: "Test\nNext Line")
		Assert(f("    Test\nNext Line") is: "    Test\nNext Line")
		Assert(f("    \n    Test\nNext Line      \n\t") is: "    Test\nNext Line")
		Assert(f("\t     Test\nTest2\n\n\n\nTest3      \n\t")
			is: "\t     Test\nTest2\n\n\n\nTest3")
		Assert(f(" ".Repeat(1200) $ "\t" $ "Test" $ "\t" $ " ".Repeat(100))
			is: " ".Repeat(1200) $ "\t" $ "Test")
		}
	}
