// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_ValidateFile()
		{
		valid = MimeMultiBase.ValidateFile

		Assert({ valid(.TempName())	} throws: `MimeMulti: AttachFile: can't get`)

		valid('abc.txt', fileContent: 'hello.world')

		tmpFile = .MakeFile('hello world')
		valid(tmpFile)

		.SpyOn(EmailMimeMaxSizeInMb).Return(0.00001) // 10 bytes

		Assert({ valid('abc.txt', fileContent: 'hello.world')	}
			throws: `MimeMulti: AttachFile: file size`)

		Assert({ valid(tmpFile) } throws: `MimeMulti: AttachFile: file size`)
		}
	}