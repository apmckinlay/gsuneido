// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		encodeClass = Base64BufferEncodeFile
			{
			Base64BufferEncodeFile_maxRead: 	12
			Base64BufferEncodeFile_lineLength: 	4
			}
		encode = encodeClass.Base64BufferEncodeFile_encode

		inFile = FakeFile(content = '')
		outFile = FakeFile('')
		encode(inFile, outFile)
		Assert(outFile.Get() is: '')

		inFile = FakeFile(content = 'This is test content, should be encoded and output')
		outFile.Reset()
		encode(inFile, outFile)
		expected = Base64.EncodeLines(content, linelen: 4)
		Assert(outFile.Get() is: expected)

		RandomLibraryRecord()
		inFile = FakeFile(content = RandomLibraryRecord().text)
		outFile.Reset()
		encode(inFile, outFile)
		expected = Base64.EncodeLines(content, linelen: 4)
		Assert(outFile.Get() is: expected)
		}
	}
