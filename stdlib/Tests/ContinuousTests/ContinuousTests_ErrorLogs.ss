// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
class
	{
	files: #(Error: 'error.log',
		Error2: 'error2.log',
		Output: 'output.log',
		Output2: 'output2.log',
		Client: 'clienterror.log')
	archiveDirectory: './logs'
	logHeading: ' Log since last Continuous Tests'
	CallClass()
		{
		return Check_LogFiles.GetLog(.files, .logHeading, .archiveDirectory,
			.archiveFunction) $ '\r\n' $ Init.ServerInitTimes()
		}

	archiveFunction(file, archiveDirectory, copyFile, cleanPattern)
		{
		if FileExists?(file)
			{
			result = CopyFile(file, archiveDirectory $ copyFile, false)
			if result is true
				try PutFile(file, '')
			}
		DeleteOldestMatchingFiles(archiveDirectory, cleanPattern)
		}
	}
