// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
function (libs = false, preGetFn = false)
	{
	if libs is false
		libs = GetContributions('ApplicationLibraries')

	if libs.Empty?()
		return ''
	result = UpdateLibraries(libs, :preGetFn)
	if String?(result)
		return result

	return result is 0 ? false : ''
	}