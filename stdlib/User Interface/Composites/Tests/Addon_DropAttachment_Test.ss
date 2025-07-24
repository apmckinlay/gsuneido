// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_checkFile()
		{
		mock = .mock()

		Assert(mock.checkFile('invalid_???.txt') is: false)
		mock.Verify.alert("The following file contains unsupported characters:\n\n" $
			'invalid_???.txt\n\nPlease reselect files')

		// Does not exist
		mock.When.dir([anyArgs:]).Return([])
		Assert(mock.checkFile('doesNotExist') is: false)
		mock.Verify.alert("The following file does not exist:\n\ndoesNotExist" $
			'\n\nPlease reselect files or ' $
			'rename the file without non-standard characters')

		// Duplicates
		mock.When.dir([anyArgs:]).Return([#file1, #file1])
		Assert(mock.checkFile(#file1) is: false)
		mock.Verify.alert('The following file matched to more than 1 file:\n\nfile1')

		// Standard processing with a valid file
		mock.When.dir([anyArgs:]).
			Return(['folder/'], ['file.exe'], ['file.msp'], ['file.png'])

		Assert(mock.checkFile('folder/') is: false)
		mock.Verify.alert('This field only accepts files')

		Assert(mock.checkFile('file.exe') is: false)
		mock.Verify.alert(ExecutableExtension?.InvalidTypeMsg)

		Assert(mock.checkFile('file.msp') is: false)
		mock.Verify.Times(2).alert(ExecutableExtension?.InvalidTypeMsg)

		Assert(mock.checkFile('file.png'))

		// Dir failures
		mock = .mock()
		mock.When.dir([anyArgs:]).
			Throw('ERROR: Readdir \\path\: The system cannot find the path specified')
		Assert(mock.checkFile('file.png') is: false)
		mock.Verify.alert('The system cannot find the path specified: file.png')
		mock.Verify.Never().log([anyArgs:])

		mock.When.dir([anyArgs:]).Throw(e = 'Unexpected Error')
		Assert(mock.checkFile('file.png') is: false)
		mock.Verify.alert('There was a problem attaching the file: file.png')
		mock.Verify.log(e)
		}

	mock()
		{
		mock = Mock(Addon_DropAttachment)
		mock.When.log([anyArgs:]).Do({ })
		mock.When.alert([anyArgs:]).Do({ })
		mock.When.checkFile([anyArgs:]).CallThrough()
		mock.When.checkInvalidChar([anyArgs:]).CallThrough()
		return mock
		}
	}
