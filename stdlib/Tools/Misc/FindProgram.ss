// Copyright (C) 2015 Suneido Software Corp. All rights reserved worldwide.
// look for an external program in GetCurrentDirectory() and ExeDir()
// return false if file does not exist or throws an error trying to find it
class
	{
	CallClass(name)
		{
		exeName = Function?(name)
			? name()
			: name
		file = Sys.Windows?() ? exeName $ '.exe' : exeName
		for dir in [ApplicationDir(), ExeDir(), GetCurrentDirectory()]
			if .fileExists?(path = Paths.Combine(dir, file))
				return path
			else if .fileExists?(path = Paths.Combine(dir, exeName, file))
				return path
		return false
		}
	fileExists?(path)
		{
		try return FileExists?(path)
		return false
		}
	}