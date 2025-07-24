// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass (folder, fileBasename, fileExistsFn = false)
		{
		Assert(fileBasename isnt: '')
		dest = folder $ fileBasename
		base = fileBasename.Has?(`.`) ? dest.BeforeLast('.') : dest
		ext = Opt('.', fileBasename.AfterLast('.'))
		if .fileExists?(dest, :fileExistsFn)
			dest = base $ '(' $ .uniqueName()  $ ')' $ ext

		return Object(:dest, :base, :ext)
		}

	// extracted for testing
	fileExists?(dest, fileExistsFn = false)
		{
		if fileExistsFn is false
			return FileExists?(dest)
		return fileExistsFn(dest)
		}
	uniqueName()
		{
		return Display(Timestamp()).Replace('\.', '_')[1 ..]
		}
	}
