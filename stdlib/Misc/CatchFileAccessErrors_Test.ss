// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		name = .TempName()
		filename = name $ ":/" $ name
		msg = CatchFileAccessErrors(filename)
			{
			PutFile(filename, 'hello') // invalid, suneido should throw can't open error
			}
		// cannot rely on what the rest of the message might say in this case
		Assert(msg has: 'Can not save file: ' $ filename $ '\n\n')

		error = "File: can't open file in mode 'a'"
		msg = CatchFileAccessErrors(name)
			{
			throw error
			}
		Assert(msg is: 'Can not save file: ' $ name $ '\n\n' $
			"Can't open file in mode 'a'")

		error = "File: Writeline: write file: " $
			"The process cannot access the file because another process has locked a " $
			"portion of the file."
		msg = CatchFileAccessErrors(name)
			{
			throw error
			}
		Assert(msg is: 'Can not save file: ' $ name $ '\n\n' $
			'The process cannot access the file because another process has ' $
				'locked a portion of the file.')

		error = "DirExists?: CreateFile file: " $
			"The system cannot contact a domain controller to service the " $
				"authentication request. Please try again later."
		msg = CatchFileAccessErrors(name)
			{
			throw error
			}
		Assert(msg is: "Can not access location: " $ name $ "\n\n" $
			"The system cannot contact a domain controller to service the " $
				"authentication request. Please try again later.")

		error = 'File: Readline: An unexpected network error occurred'
		msg = CatchFileAccessErrors(name)
			{
			throw error
			}
		Assert(msg is: "Can not read file: " $ name $ '\n\n' $
			'An unexpected network error occurred')

		Assert(
			{ CatchFileAccessErrors(name)
				{
				throw 'ERROR: unknown error'
				}}
			throws: 'ERROR: unknown error')

		Assert(
			{ CatchFileAccessErrors(name)
				{
				throw 'syntax error'
				}}
			throws: 'syntax error')

		msg = CatchFileAccessErrors(name)
			{
			PutFile(name, 'hello')
			}
		DeleteFile(name)
		Assert(msg)
		}
	}