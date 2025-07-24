// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
function (libs = false)
	{
	if libs is false
		libs = Libraries()
	result = ""
	for lib in libs
		{
		if "" isnt s = CheckLibrary(lib)
			result $= "\n" $ lib $ ":\n" $ s
		if s.Suffix?('interrupt')
			break
		}
	return result
	}