// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	clNoThrow: CopyFileAndAttach
		{
		EnsureDirExists(copyfolder, copyto /*unused*/, quiet? /*unused*/)
			{
			return copyfolder
			}
		}
	clThrow: CopyFileAndAttach
		{
		EnsureDirExists(@unused)
			{
			return Object(msg: "Test Error")
			}
		}
	Test_buildFileOb_windows()
		{
		if not Sys.Windows?()
			return

		.SpyOn(MakeUniqueFileName.MakeUniqueFileName_fileExists?).
			Return(false, true, false)
		.SpyOn(MakeUniqueFileName.MakeUniqueFileName_uniqueName).
			Return('19821226.132615128')

		subfolder = CopyFileAndAttach.SubFolder()

		result = .clNoThrow.CopyFileAndAttach_buildFileOb(
			'c:\work\fred.txt', `c:\temp\`, true)
		Assert(result.ext is: '.txt')
		Assert(result.dest is: `c:\temp\` $ subfolder $ 'fred.txt')
		Assert(result.base is: `c:\temp\` $ subfolder $ 'fred')
		Assert(result.copyto is: `c:\temp\`)
		Assert(result.fileBaseName is: subfolder $ 'fred.txt')

		result = .clNoThrow.CopyFileAndAttach_buildFileOb(
			'c:\work\fred.txt', `c:\temp\`, true)
		Assert(result.ext is: '.txt')
		Assert(result.dest is: `c:\temp\` $ subfolder $ 'fred(19821226.132615128).txt')
		Assert(result.base is: `c:\temp\` $ subfolder $ 'fred')
		Assert(result.copyto is: `c:\temp\`)
		Assert(result.fileBaseName is: subfolder $ 'fred(19821226.132615128).txt')

		// copy to destination and give new file name
		result = .clNoThrow.CopyFileAndAttach_buildFileOb(
			'fred.txt', `c:\temp\barney.txt`, true)
		Assert(result.ext is: '.txt')
		Assert(result.dest is: `c:\temp/` $ subfolder $ 'barney.txt')
		Assert(result.base is: `c:\temp/` $ subfolder $ 'barney')
		Assert(result.copyto is: `c:\temp/`)
		Assert(result.fileBaseName is: subfolder $ 'barney.txt')

		result = .clNoThrow.CopyFileAndAttach_buildFileOb(
			'fred.txt', `c:\temp\CON.txt`, true)
		Assert(result is: Object(error: CheckFileName.ReservedNameDisplay))

		result = .clNoThrow.CopyFileAndAttach_buildFileOb(
			'fred.txt', `c:\temp\CON:f.txt`, true)
		Assert(result is: Object(error: CheckFileName.InvalidCharsDisplay))

		result = .clThrow.CopyFileAndAttach_buildFileOb(
			'fred.txt', `c:\temp\barney.txt`, true)
		Assert(result is: Object(error: "Test Error"))
		}

	Test_buildFileOb_linux()
		{
		if Sys.Windows?()
			return

		.SpyOn(MakeUniqueFileName.MakeUniqueFileName_fileExists?).
			Return(false, true, false)
		.SpyOn(MakeUniqueFileName.MakeUniqueFileName_uniqueName).
			Return('19821226.132615128')

		subfolder = CopyFileAndAttach.SubFolder()
		result = .clNoThrow.CopyFileAndAttach_buildFileOb(
			'/work/fred.txt', `/temp/`, true)
		Assert(result.ext is: '.txt')
		Assert(result.dest is: `/temp/` $ subfolder $ 'fred.txt')
		Assert(result.base is: `/temp/` $ subfolder $ 'fred')
		Assert(result.copyto is: `/temp/`)
		Assert(result.fileBaseName is: subfolder $ 'fred.txt')

		result = .clNoThrow.CopyFileAndAttach_buildFileOb(
			'/work/fred.txt', `/temp/`, true)
		Assert(result.ext is: '.txt')
		Assert(result.dest is: `/temp/` $ subfolder $ 'fred(19821226.132615128).txt')
		Assert(result.base is: `/temp/` $ subfolder $ 'fred')
		Assert(result.copyto is: `/temp/`)
		Assert(result.fileBaseName is: subfolder $ 'fred(19821226.132615128).txt')

		// copy to destination and give new file name
		result = .clNoThrow.CopyFileAndAttach_buildFileOb(
			'fred.txt', `/temp/barney.txt`, true)
		Assert(result.ext is: '.txt')
		Assert(result.dest is: `/temp/` $ subfolder $ 'barney.txt')
		Assert(result.base is: `/temp/` $ subfolder $ 'barney')
		Assert(result.copyto is: `/temp/`)
		Assert(result.fileBaseName is: subfolder $ 'barney.txt')

		result = .clNoThrow.CopyFileAndAttach_buildFileOb(
			'fred.txt', `/temp/CON.txt`, true)
		Assert(result is: Object(error: CheckFileName.ReservedNameDisplay))

		result = .clNoThrow.CopyFileAndAttach_buildFileOb(
			'fred.txt', `/temp/CON:f.txt`, true)
		Assert(result is: Object(error: CheckFileName.InvalidCharsDisplay))

		result = .clThrow.CopyFileAndAttach_buildFileOb(
			'fred.txt', `/temp/barney.txt`, true)
		Assert(result is: Object(error: "Test Error"))
		}

	Test_emptyCopyTo()
		{
		// copy to is empty, won't copy the files BUT will still build
		// attachOb so user can see attachment name on the email
		result = .clNoThrow.CopyFileAndAttach_buildFileOb('fred.txt', '', true)
		Assert(result.ext is: '')
		Assert(result.fileBaseName is: 'fred.txt')
		Assert(result.dest is: '')
		Assert(result.base is: '')
		Assert(result.copyto is: '')
		}

	Test_buildAttachOb()
		{
		fn = CopyFileAndAttach.BuildAttachOb
		ob = Object()
		fn(ob, 'fred.txt')
		Assert(ob isSize: 1)
		Assert(ob[0] hasMember: 'attachment0')

		fn(ob, 'barney.pdf')
		Assert(ob isSize: 1)
		Assert(ob[0] hasMember: 'attachment1')

		fn(ob, 'wilma.txt')
		fn(ob, 'betty.txt')
		fn(ob, 'pebbles.txt')
		fn(ob, 'bambam.txt')
		Assert(ob isSize: 2)
		Assert(ob[0] isSize: 5)
		Assert(ob[1] isSize: 1)

		fn(ob, 'dino.txt.gpg')
		Assert(ob[1].attachment1 is: 'dino.txt')

		fn = CopyFileAndAttach.BuildAttachOb
		ob = Object()
		fn(ob, '//server/dir1/dir2/fred.txt', true)
		Assert(ob isSize: 1)
		Assert(ob[0] hasMember: 'attachment0')
		Assert(ob[0]['attachment0'] is: 'fred.txt')

		fn = CopyFileAndAttach.BuildAttachOb
		ob = Object()
		fn(ob, '//server/dir1/dir2/fred.txt')
		Assert(ob isSize: 1)
		Assert(ob[0] hasMember: 'attachment0')
		Assert(ob[0]['attachment0'] is: '//server/dir1/dir2/fred.txt')
		}
	}
