// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass()
		{
		subfolder = .SubFolderName()
		tmp = GetTempPath()
		path = tmp $ subfolder $ '/'
		EnsureDir(path)
		return path
		}

	SubFolderName()
		{
		return OptContribution('TempFolder', function() { return 'suneido' })()
		}
	}