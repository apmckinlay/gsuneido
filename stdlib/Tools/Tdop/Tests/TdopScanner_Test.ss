// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
TdopTests
	{
	Test_main()
		{
		src = 'test1 + \r\n1.2 "abc\r\ndefg" true\r\n'
		symbols = TdopSymbols()
		scanner = TdopScanner(src, symbols)
		.checkAndAdvance(scanner, 'IDENTIFIER(test1)', 1, 5, false)
		.checkAndAdvance(scanner, 'ADD', 7, 1, false, 'NUMBER(1.2)')
		.checkAndAdvance(scanner, 'NUMBER(1.2)', 11, 3, true, 'STRING(abc\r\ndefg)')
		.checkAndAdvance(scanner,
			'STRING(abc\r\ndefg)', 15, 11, false, 'IDENTIFIER(true)')
		.checkAndAdvance(scanner, 'IDENTIFIER(true)', 27, 4, false, '<end>')
		.checkAndAdvance(scanner, '<end>', 33, 0, true, '<end>')
		.checkAndAdvance(scanner, '<end>', 33, 0, false, '<end>')

		src = ''
		scanner = TdopScanner(src, symbols)
		.checkAndAdvance(scanner, '<end>', 1, 0, false, '<end>')
		}

	checkAndAdvance(scanner, token, position, length, newline, ahead = false)
		{
		Assert(Display(scanner.Token()) is: token)
		Assert(scanner.Position() is: position)
		Assert(scanner.Length() is: length)
		Assert(scanner.IsNewline() is: newline)
		if ahead isnt false
			Assert(Display(scanner.Ahead()) is: ahead)
		scanner.Advance()
		}
	}