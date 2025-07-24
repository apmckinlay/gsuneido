// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		cl = MakeUniqueFileName
			{
			MakeUniqueFileName_fileExists?(dest)
				{
				return dest.Has?('exists')
				}
			}

		// no file name; it throws
		Assert({ cl('', '') } throws:)

		// no extension; file not exist
		folder = .TempTableName() $ `/`
		fileBasename = `noFileExtension`
		result = cl(folder, fileBasename)
		Assert(result.ext is: '')
		Assert(result.base is: folder $ fileBasename)
		Assert(result.dest is: folder $ fileBasename)

		// no extension; file exists
		fileBasename = `noFileExtension-exists`
		result = cl(folder, fileBasename)
		Assert(result.ext is: '')
		Assert(result.base is: folder $ fileBasename)
		Assert(result.dest matches: folder $ fileBasename $ '\(\d*_\d*\)')

		result = MakeUniqueFileName(folder, fileBasename, { |dest| dest.Has?('exists') })
		Assert(result.ext is: '')
		Assert(result.base is: folder $ fileBasename)
		Assert(result.dest matches: folder $ fileBasename $ '\(\d*_\d*\)')

		// with extension; no other dots; file not exist
		fileBasename = `withFileExtension.extension`
		result = cl(folder, fileBasename)
		Assert(result.ext is: '.extension')
		Assert(result.base is: folder $ fileBasename.RemoveSuffix(`.extension`))
		Assert(result.dest is: folder $ fileBasename)

		result = MakeUniqueFileName(folder, fileBasename, { |dest| dest.Has?('exists') })
		Assert(result.ext is: '.extension')
		Assert(result.base is: folder $ fileBasename.RemoveSuffix(`.extension`))
		Assert(result.dest is: folder $ fileBasename)

		// with extension; no other dots; file exists
		fileBasename = `withFileExtension-exists.extension`
		result = cl(folder, fileBasename)
		Assert(result.ext is: '.extension')
		Assert(result.base is: folder $ fileBasename.RemoveSuffix(`.extension`))
		Assert(result.dest matches: folder $ fileBasename.RemoveSuffix(`.extension`) $
			'\(\d*_\d*\)' $ '.extension')

		// with extension; with other dots; file not exist
		fileBasename = `with.File.Extension.extension`
		result = cl(folder, fileBasename)
		Assert(result.ext is: '.extension')
		Assert(result.base is: folder $ fileBasename.RemoveSuffix(`.extension`))
		Assert(result.dest is: folder $ fileBasename)

		// with extension; with other dots; file exists
		fileBasename = `with.File.Extension-exists.extension`
		result = cl(folder, fileBasename)
		Assert(result.ext is: '.extension')
		Assert(result.base is: folder $ fileBasename.RemoveSuffix(`.extension`))
		Assert(result.dest matches: folder $ fileBasename.RemoveSuffix(`.extension`) $
			'\(\d*_\d*\)' $ '.extension')
		}
	}
