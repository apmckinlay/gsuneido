// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Setup()
		{
		.userAttach = ''
		if Suneido.Member?('currentUserAttachmentSettings')
			{
			.userAttach = Suneido.currentUserAttachmentSettings
			Suneido.Delete('currentUserAttachmentSettings')
			}
		}
	Test_main()
		{
		//Test Copy and link to Attachment Directory
		imageSelect = OpenImageSelect
			{
			OpenImageSelect_s3Bucket()
				{
				return ''
				}
			CopyFile(file, copyto /*unused*/ = false)
				{
				return OpenImageSettings.Copyto() $ Paths.Basename(file)
				}
			OpenImageSelect_fileSize(unused) { return 1 }
			}

		oldFile = '\\\\test\\d\\file.fl'
		newFile = imageSelect.ResultFile(oldFile, false, false)
		Assert(oldFile is: newFile msg: "Link Changed Path")
		newFile = imageSelect.ResultFile(oldFile, true, false)
		Assert(Paths.Basename(oldFile) is: newFile msg: "Copy Link Left Path" )

		//Test Copy and link to separate Directory
		imageSelect = OpenImageSelect
			{
			OpenImageSelect_s3Bucket()
				{
				return ''
				}
			CopyFile(file, copyto /*unused*/ = false)
				{
				return file
				}
			OpenImageSelect_fileSize(unused) { return 1 }
			}
		newFile = imageSelect.ResultFile(oldFile, true, false)
		Assert(oldFile is: newFile msg: "Copy to Not Attachment Directory")
		}

	Test_subFolder()
		{
		imageSelect = OpenImageSelect
			{
			OpenImageSelect_s3Bucket()
				{
				return ''
				}
			CopyFile(file, copyto /*unused*/ = false, subfolder /*unused*/ = '')
				{
				return OpenImageSettings.Copyto() $
					"200001\\" $ Paths.Basename(file)
				}
			OpenImageSelect_subFolder(useSubFolder)
				{
				return useSubFolder isnt false ? "200001\\" : ''
				}
			OpenImageSelect_fileSize(unused) { return 1 }
			}

		oldFile = '\\\\test\\d\\file.fl'
		newFile = imageSelect.ResultFile(oldFile, false, true)
		Assert(oldFile is: newFile msg: "Link Changed Sub Path")
		newFile = imageSelect.ResultFile(oldFile, true, true)
		Assert(newFile is: "200001\\" $ Paths.Basename(oldFile)
			msg: "Copy to Not Sub Directory")
		}

	Test_CopyFile()
		{
		mock = Mock(OpenImageSelect)
		mock.When.s3Bucket().Return('')
		mock.When.copyFile([anyArgs:]).Return(false, false, 'Retry', true)
		mock.When.deleteFile?([anyArgs:]).Return(false)
		mock.When.ensureCopyFolder([anyArgs:]).Return(`\copyFolder\subFolder\`)
		mock.When.AlertWarn([anyArgs:]).Return(true)
		mock.When.CopyFile([anyArgs:]).CallThrough()
		mock.When.fileSize([anyArgs:]).Return(1)

		suneidoLog = .SpyOn(SuneidoLog)
		suneidoLog.Return('')

		// No file to copy, return
		Assert(mock.CopyFile(``) is: false)

		// invalid file name, alert occurs, returns
		file = `testFile><` $ Display(Timestamp())[1 ..] $ `.txt`
		filePath = `\testFolder\` $ file
		Assert(mock.CopyFile(filePath) is: false)
		mock.Verify.AlertWarn([anyArgs:])

		// Copy file fails, alert occurs, returns
		file = `testFile` $ Display(Timestamp())[1 ..] $ `.txt`
		filePath = `\testFolder\` $ file
		Assert(mock.CopyFile(filePath) is: false)
		mock.Verify.Times(2).AlertWarn([anyArgs:])

		// Copy file fails, quiet? flag is present, no alert, returns
		Assert(mock.CopyFile(filePath, quiet?:) is: false)
		mock.Verify.Times(2).OpenImageSelect_copyFile([anyArgs:])
		mock.Verify.Times(2).AlertWarn([anyArgs:])

		// Copy file needs to retry
		Assert(mock.CopyFile(filePath, quiet?:) is: `\copyFolder\subFolder\` $ file)
		// copyFile? should get called twice
		mock.Verify.Times(4).OpenImageSelect_copyFile([anyArgs:])
		mock.Verify.Times(2).AlertWarn([anyArgs:])

		// Copy file succeeds, no deletes required, file path is returned
		Assert(mock.CopyFile(filePath) is: `\copyFolder\subFolder\` $ file)
		mock.Verify.Times(5).OpenImageSelect_copyFile([anyArgs:])
		mock.Verify.Times(2).AlertWarn([anyArgs:])

		// Copy file succeeds, deletesource fails, is logged, deletes copied file, returns
		Assert(mock.CopyFile(filePath, deletesource?:) is: false)
		mock.Verify.deleteFile?(filePath)
		mock.Verify.deleteFile?(`\copyFolder\subFolder\` $ file)
		}

	Test_ensureCopyFolder()
		{
		spy = .SpyOn(CheckDirExists)
		callLogs = spy.CallLogs()
		mock = Mock(OpenImageSelect)
		mock.When.s3Bucket().Return('')
		mock.When.subFolder([anyArgs:]).Return(`\subFolder\`)
		mock.When.createDir?([anyArgs:]).Return(true, false)
		mock.When.ensureCopyFolder([anyArgs:]).CallThrough()
		mock.When.fileSize([anyArgs:]).Return(1)

		suneidoLog = .SpyOn(SuneidoLog)
		suneidoLog.Return('')

		// Not using subfolder, returns passed in path
		Assert(mock.ensureCopyFolder('', 'testFolder') is: 'testFolder')

		// Using subfolder, subfolder already exists, returns new path
		spy.ClearAndReturn(true)
		Assert(mock.ensureCopyFolder(true, 'testFolder')
			is: `testFolder\subFolder\`)
		Assert(callLogs[0] is: #(folder: `testFolder\subFolder\`))
		mock.Verify.Never().createDir([anyArgs:])

		// Using subfolder, subfolder doesn't exist, path gets created, returns new path
		spy.ClearAndReturn(false)
		Assert(mock.ensureCopyFolder(true, 'testFolder')
			is: `testFolder\subFolder\`)
		Assert(callLogs[1] is: #(folder: `testFolder\subFolder\`))
		mock.Verify.createDir?(`testFolder\subFolder\`)

		// Using subfolder, subfolder doesn't exist, path doesnt get created,
		// returns original path
		Assert(mock.ensureCopyFolder(true, 'testFolder') is: `testFolder`)
		Assert(callLogs[2] is: #(folder: `testFolder\subFolder\`))
		mock.Verify.Times(2).createDir?(`testFolder\subFolder\`)
		Assert(logs = suneidoLog.CallLogs() isSize: 1)
		Assert(logs[0].message
			is: `ERRATIC: (CAUGHT) Could not create folder: testFolder\subFolder\`)

		// folder is not accessable and DirExists? throws
		spy.ClearAndReturn('test error')
		Assert(mock.ensureCopyFolder(true, 'testFolder') is: #(msg: 'test error'))
		Assert(callLogs[3] is: #(folder: `testFolder\subFolder\`))
		Assert(callLogs isSize: 4)
		mock.Verify.Times(2).createDir?(`testFolder\subFolder\`)

		mock.When.s3Bucket().Return('test')
		Assert(mock.ensureCopyFolder(true, 'testFolder') is: `testFolder\subFolder\`)
		Assert(callLogs isSize: 4)
		mock.Verify.Times(2).createDir?(`testFolder\subFolder\`) // not being called
		}

	Test_GetCopyToFilename()
		{
		// .MakeDir is cleaning up files
		dir = .MakeDir()
		testCl = OpenImageSelect
			{
			OpenImageSelect_s3Bucket() { return '' }
			OpenImageSelect_fileSize(unused) { return 1 }
			}
		get = testCl.GetCopyToFilename
		path = dir $ '/'
		name = .TempTableName()
		Assert(get(path, name) is: path $ name)

		.PutFile(path $ name, 'hello')
		dest = get(path, name)
		DeleteFile(path $ name)
		Assert(dest startsWith: path $ name)
		Assert(dest endsWith: ')')

		.PutFile(path $ name $ '.txt', 'hello2')
		dest = get(path, name $ '.txt')
		DeleteFile(path $ name)
		Assert(dest startsWith: path $ name)
		Assert(dest endsWith: ').txt')
		}

	Teardown()
		{
		if Object?(.userAttach)
			Suneido.currentUserAttachmentSettings = .userAttach
		super.Teardown()
		}
	}
