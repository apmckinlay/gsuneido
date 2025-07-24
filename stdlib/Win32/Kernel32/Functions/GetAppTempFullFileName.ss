// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
function (prefix = '')
	{
	path = GetAppTempPath()
	return GetTempFileName(path, prefix)
	}