// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	// when setting up test cases please note that \x0D, \x0A, and \x09 are line feed,
	// carriage return, and tab (respectively) - and will match as valid chars.
	testCases()
		{
		return Object(
			#(text: "This Should, Flag as a, Proper CSV, with columns\r\n" $
				"field1, field2, field3,",
				result: true),
			#(text:
				"this file, should fail, \x17, not valid chars, \x19 include \r\nnewline",
				result: false),
			// a PNG Header
			#(text: '\x89\x50\x4e\x47\x0d\x0a\x1a\x0a\x00\x00\x00\x0d\x49\x48\x44\x52',
				result: false),
			#(text: "",
				result:  false),
			#(text: "ABCDEFGHIJ\r\nKLMNOPQRST",
				result: true),
			#(text: "this has \t a tab\r\nfield1, field2, field3",
				result: true),
			Object(text: "abc" $ 160.Chr() $ "def", // non breaking white-space (ASC 160)
				result: true)
			#(text: ` !"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\]^` $
				"_`abcdefghijklmnopqrstuvwxyz{|}\r\n\t",
				result: true)
			// test UTF-8 marker - made up of extended Ascii so should be okay
			Object(text: "\xEF\xBB\xBF" $
				` !"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\]^` $
				"_`abcdefghijklmnopqrstuvwxyz{|}\r\n\t"
				result: true)
			// only extended Ascii chars - considered not valid (open for discussion)
			Object(text: "\xEF\xBB\xBF\xEF\xBB\xBF\xEF\xBB\xBF\xEF\xBB\xBF"
				result: false)
			// some systems may have inserted an EOF char
			Object(text: "a".Repeat(99) $ '\x1a',
				result: true)
			Object(text: 'abcdefghi234234234234^%$&^(*&)(*&# \x1a',
				result: true)
			// ignore null characthers
			Object(text: '\x00test\x00\x00testned \x00', result: true)
			)
		}

	Test_main()
		{
		for test in .testCases()
			Assert(IsFileText.IsFileText_firstChunkValid?(test.text, chunkIsWholeFile:)
				is: test.result)
		}

	Test_error_from_FileSize()
		{
		testCl = IsFileText
			{
			IsFileText_getFile(@unused) { return false }
			IsFileText_fileSize(@unused) { return 1.Mb() }
			}
		Assert({ testCl(`fakeName`) } throws: "SHOW: Unable to access file")
		}
	}