// Copyright (C) 2022 Axon Development Corporation All rights reserved worldwide.
Test
	{
	Test_OK()
		{
		classMock = Mock(OpenImageRename)
		classMock.When.OK([anyArgs:]).CallThrough()
		classMock.When.AlertInfo([anyArgs:]).Do({ })
		classMock.When.getCopyToFilename([anyArgs:]).Return('unused for this test')
		classMock.When.attemptToRename([anyArgs:]).Return(true)

		ctrlMock = Mock()
		classMock.OpenImageRename_fileNameCtrl = ctrlMock
		classMock.OpenImageRename_oldFileName = 'oldFileName.txt'

		ctrlMock.When.Get().Return('')
		Assert(classMock.OK() is: false)
		classMock.Verify.AlertInfo('Invalid file name', 'Please enter a File Name')

		ctrlMock.When.Get().Return('/')
		Assert(classMock.OK() is: false)
		classMock.Verify.AlertInfo('Invalid file name',
			CheckFileName.InvalidCharsDisplay $ '\nPlease enter another File Name.')

		ctrlMock.When.Get().Return('.txt')
		Assert(classMock.OK() is: false)
		classMock.Verify.AlertInfo('Invalid file name',
			'File Name cannot be blank or just an extension\n' $
			'Please enter another File Name.')

		ctrlMock.When.Get().Return('testFile.txt')

		// resulting path over 255 chars should result in warning message
		classMock.When.getCopyToFilename([anyArgs:]).Return('a'.Repeat(260))
		Assert(classMock.OK() is: false)
		classMock.Verify.AlertInfo('Invalid file name',
			'Destination path is too long. Please choose a shorter file name')

		// back to valid file path
		classMock.When.getCopyToFilename([anyArgs:]).Return('unused for this test')
		Assert(classMock.OK())

		ctrlMock.When.Get().Return('testFile')
		Assert(classMock.OK())

		ctrlMock.When.Get().Return('testFile.log')
		Assert(classMock.OK())

		ctrlMock.When.Get().Return('testFile.csv')
		Assert(classMock.OK())

		ctrlMock.When.Get().Return('testFile.pdf')
		Assert(classMock.OK() is: false)
		classMock.Verify.AlertInfo('Invalid file name', 'Cannot convert pdfs')

		ctrlMock.When.Get().Return('oldFileName')
		Assert(classMock.OK() is: 'File not renamed')

		ctrlMock.When.Get().Return('oldFileName.pdf')
		Assert(classMock.OK() is: false)
		classMock.Verify.Times(2).AlertInfo('Invalid file name', 'Cannot convert pdfs')

		ctrlMock.When.Get().Return('oldFileName.txt')
		Assert(classMock.OK() is: 'File not renamed')
		}
	}