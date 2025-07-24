// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
// mirror a directory tree with fromDir and toDir
// this handles folders with or without slashes
// returns failed files
class
	{
	CallClass(fromDir, toDir, keepExtra? = false, preCopyFunc = false, skipDirs = #())
		{
		_failedList = Object()
		_keepExtra? = keepExtra?
		_preCopyFunc = preCopyFunc
		_skipDirs = skipDirs.Map(Paths.ToStd)
		.mirror(fromDir, toDir)
		return _failedList
		}

	mirror(fromDir, toDir)
		{
		fromDir = .ensureSlash(fromDir)
		toDir = .ensureSlash(toDir)

		if false is .ensureDir(toDir.BeforeLast(.slash()))
			{
			_failedList.Add(toDir)
			return
			}

		fromDirList = .dirList(fromDir)
		toDirList = .dirList(toDir)
		for item in fromDirList.Members()
			{
			fullFromDir = fromDir $ item
			fullToDir = toDir $ item

			if .skip?(fullFromDir)
				continue

			if item.Suffix?('/') // folder, Dir always return foward slash
				.mirror(fullFromDir, fullToDir)
			else // file
				{

				if fromDirList[item] isnt toDirList.GetDefault(item, false) or
					.different?(fullFromDir, fullToDir)
					{
//Print("copyFile", fullFromDir, fullToDir)
					if not .copyFile(fullFromDir, fullToDir)
						_failedList.Add(fullFromDir)
					}
				}
			}

		.deleteExtra(fromDirList, toDirList, toDir)
		}

	skip?(curDir)
		{
		return _skipDirs.NotEmpty?() and
			_skipDirs.Has?(Paths.ToStd(curDir))
		}

	deleteExtra(fromDirList, toDirList, toDir)
		{
		if _keepExtra?
			return

		for name in toDirList.Members()
			if not fromDirList.Member?(name) and not .skip?(toDir $ name)
				.deleteItem(toDir, name)
		}

	deleteItem(dir, item)
		{
		Assert(dir.Suffix?('/') or dir.Suffix?(`\`))
		fullName = dir $ item
		if item.Suffix?(`\`) or item.Suffix?(`/`)
			{
			if false is .deleteDir(fullName)
				_failedList.Add(fullName)
			}
		else // file
			{
			if false is .deleteFile(fullName)
				_failedList.Add(fullName)
			}
		}

	DeleteFilesAndDirs(dir, listFile)
		{
		dir = .ensureSlash(dir)
		deletionFilePath = dir $ listFile
		if false is deletedListStr = .getFile(deletionFilePath)
			return #()

		_failedList = Object()
		for f in deletedListStr.Lines()
			if '' isnt f = f.Trim()
				.deleteItem(dir, f)
		.deleteItem(dir, listFile)
		return _failedList
		}

	getFile(file)
		{
		return GetFile(file)
		}

	ensureDir(dir)
		{
		return true is EnsureDir(dir)
		}

	different?(from, to)
		{
		try
			return not CompareFiles?(from, to)
		catch
			_failedList.Add(from)
		return false
		}

	copyFile(from, to)
		{
		if Function?(_preCopyFunc) and not _preCopyFunc(:from, :to)
			return false
		try
			return true is CopyFile(from, to, false)
		return false
		}

	deleteFile(file)
		{
		if FileExists?(file)
			return DeleteFile(file)
		return true
		}

	deleteDir(dir)
		{
		return true is DeleteDir(dir)
		}

	dirList(dir)
		{
		list = Object()
		for item in Dir(dir $ '*.*', details:)
			list[item.name] = item.size
		return list
		}

	slash()
		{
		return Paths.ToLocal('/')
		}

	ensureSlash(dir)
		{
		slash = .slash()
		localDir = Paths.ToLocal(dir)
		if not localDir.Suffix?(slash)
			localDir $= slash
		return localDir
		}
	}
