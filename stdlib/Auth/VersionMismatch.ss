// Copyright (C) 2024 Axon Development Corporation All rights reserved worldwide.
class
	{
	// This class can only access Built-in functions
	CallClass()
		{
		if '' is s = Cmdline()
			{
			Fatal('Version Mismatch')
			return
			}

		if .skip?()
			return

		file = .programNamesAndSetDir()
		if "" is appFolder = .getAppFolder(s)
			{
			.alert('Invalid configuration file, please re-install the program.\r\n' $
				'Run the dated ' $ file.appName $
				' in the Shared axon folder on the server'
				'get AppFolder failed - ' $ s, file.appName)
			return
			}

		if false is copyExe = .getLatestExe(file.appName, appFolder)
			return

		.updateExeAndLaunch(file, copyExe)
		}

	skip?()
		{
		// This can only be ran while running from a local exe
		// slashes are standardized for comparison only so shouldn't matter which slash
		return false is .exePath().Tr('/', '\\').Has?(Getenv('AppData').Tr('/', '\\'))
		}

	programNamesAndSetDir()
		{
		path = .exePath()
		pos = path.FindLast1of("\\/:")
		.setDir(path[.. pos])
		exeName = pos is false ? path : path[pos + 1 ..]
		appName = exeName.Replace('\.exe', '')
		return Object(:exeName, :appName)
		}

	alert(msg, details, title)
		{
		ErrorLog('INFO: Update local exe - ' $ details $ ', from ' $ .exePath())
		Fatal('Updating ' $ title[0].Upper() $ title[1..] $ ':\r\n\r\n' $ msg $ '\r\n' $
			'Please contact your system administrator')
		Exit(true)
		}

	getAppFolder(s)
		{
		delimiter = 'AppFolder='
		i = s.FindLast(delimiter)
		appFolder = i is false ? "" : s[i + delimiter.Size() ..]
		return appFolder
		}

	getLatestExe(appName, appFolder)
		{
		pattern = '^' $ appName $ '\d\d\d\d\d\d\d\d.exe$'
		for exe in dir = .dir(appFolder $ '/' $ appName $ '*.exe').Sort!({ |x,y| x > y })
			if exe =~ pattern
				return appFolder $ `/` $ exe

		return .alert('Cannot read server shared folder ' $ appFolder $
			', this is possibly caused by network issues',
			'getLatestExe failed - ' $ Display(dir), appName)
		}

	updateExeAndLaunch(file, copyExe)
		{
		if true isnt result = .moveFile(file.exeName)
			{
			.alert(file.exeName $ ' is unable to rename on the local drive',
				'MoveFile failed - ' $ result, file.appName)
			return
			}

		if true isnt result = CopyFile(copyExe, file.exeName, false)
			{
			.alert('Cannot copy ' $ file.exeName $ ' from server shared folder' $
				', this is possibly caused by network issues',
				'copy file failed - ' $ result, file.appName)
			return
			}

		if 0 isnt result = System('start ' $ file.exeName)
			{
			.alert('Unable to launch executable',
				'Unable to launch executable - ' $ result, file.appName)
			return
			}
		.cleanUpOld(file.exeName)
		Exit(true)
		}

	moveFile(exe)
		{
		if true is MoveFile(exe, exe $ '.' $ Display(Date()).Tr('#.'))
			return true
		result = ''
		for i in .. 3 /*= retry */
			{
			n = 2 << i
			Thread.Sleep(n + Random(n))
			if true is result = MoveFile(exe, exe $ '.' $ Display(Date()).Tr('#.'))
				return true
			}
		return result
		}

	cleanUpOld(exeName)
		{
		for exe in Dir(exeName $ '.*')
			unused = DeleteFileApi(exe) // to avoid errors thrown
		}

	// below are extracted for tests
	dir(path)
		{
		return Dir(path)
		}

	exePath()
		{
		return ExePath()
		}

	// To stop issues for shortcuts that use 'start in' arg that isn't the exe path
	setDir(dir)
		{
		SetCurrentDirectory(dir)
		}
	}
