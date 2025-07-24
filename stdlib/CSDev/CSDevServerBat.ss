// Copyright (C) 2019 Axon Development Corporation All rights reserved worldwide.
class
	{
	CallClass(lineEnd)
		{
		str = ''
		sections = Object(.main(), .runOptions, .run(), .runStandalone,	.options())
		.setupCopyFiles(sections)
		sections.Add(.loadSaved, .loadFile(), .loadPort, .saveOptions,
			.randomizeport)

		for section in sections
			str $= section.Join(lineEnd) $ lineEnd $ lineEnd
		return str $ ':end'
		}

	main()
		{
		return Object('@echo off',
			'setlocal',
			'title Start G Server',
			'set port=3147',
			'set lastport=3147',
			'echo:',
			'echo \tPress:',
			'echo \t0 to run with last option',
			'echo \t1 to run with defaults',
			'echo \t2 to change settings',
			'echo \t3 to copy latest and run with last option',
			'echo \t4 to copy current and run with last option',
			'echo:',
			'call :loadSaved',
			'echo:',
			'CHOICE /C 01234 /N /D 0 /T 5 /M ""',
			`REM echo '"%ERRORLEVEL%"'`,
			'IF %ERRORLEVEL%==1 GOTO RunOptions',
			'IF %ERRORLEVEL%==2 GOTO Run',
			'IF %ERRORLEVEL%==3 GOTO Options',
			'IF %ERRORLEVEL%==4 GOTO CopyLatest',
			'IF %ERRORLEVEL%==5 GOTO CopyCurrent',
			'echo Invalid Entry',
			'GOTO end')
		}

	runOptions: #(
		':RunOptions'
		'IF %OPTION%==0 GOTO Run'
		'IF %OPTION%==1 call :randomizeport'
		'IF %OPTION%==2 GOTO RunGStandalone'
		'GOTO Run')

	run()
		{
		return Object(
			':Run'
			'echo Starting...'
			'IF NOT %port%==%lastport% call :loadPort'
			'call :saveOptions'
			'gsport.exe -s -p %port% CSDevServerWindow() >> golang.log'
			'del golang.log'
			'GOTO end'
			)
		}

	runStandalone: #(
		':RunGStandalone'
		'gsuneido.exe'
		'GOTO end')

	options()
		{
		return Object(
			':Options',
			'cls',
			'CHOICE /C 012 /N /M "0 for defaults, 1 for Random Port, ' $
				'2 for GSuneido Standalone',
			'IF %ERRORLEVEL%==1 set OPTION=0',
			'IF %ERRORLEVEL%==2 set OPTION=1',
			'IF %ERRORLEVEL%==3 set OPTION=2',
			'GOTO RunOptions')
		}

	setupCopyFiles(sections)
		{
		defFunc = function () { return false }
		paths = Object()
		files = Object('gsuneido.exe')
		paths.Add(OptContribution('CSDevStaffExePath', defFunc)() at: 'CopyLatest')
		paths.Add(OptContribution('CSDevLatestExePath', defFunc)() at: 'CopyCurrent')
		files.Add('gsport.exe')
		for label in paths.Members()
			sections.Add(.copyFiles(label, paths[label], files))
		}

	copyFiles(label, path, files)
		{
		ret = Object(':' $ label)
		for file in files
			ret.Add('copy /Y ' $ path $ file $ ` .\` $ file)
		ret.Add('GOTO RunOptions')
		return ret
		}

	loadSaved: #(
		':loadSaved'
		'if exist csDevOption.sv ('
		'call :loadfile'
		') else ('
		'echo \tNo Saved Options found, running with defaults'
		'set OPTION=0'
		')'
		'exit /b')

	loadFile()
		{
		return Object(':loadfile',
			'(',
			'set /p "OPTION="',
			'set /p "lastport="',
			') <csDevOption.sv',
			'IF %OPTION%==0 Echo \tLast Option was: Defaults',
			'IF %OPTION%==1 Echo \tLast Option was: Random Port',
			'IF %OPTION%==2 Echo \tLast Option was: Standalone',
			'Echo \tLast Port was: %lastport%',
			'exit /b')
		}

	loadPort: #(
		':loadPort'
		'set file1=%AppData%\suneido%lastport%.err'
		'set file2=suneido%port%.err'
		'IF EXIST %file1% ('
		'REN "%file1%" "%file2%"'
		') else ('
		'Echo \t%file1% could not be found, creating new error log'
		')'
		'exit /b'
		)

	saveOptions: #(
		':saveOptions'
		'REM Save Selected Option'
		'('
		'echo %OPTION%'
		'echo %port%'
		') > csDevOption.sv'
		'exit /b')

	randomizeport: #(
		':randomizeport'
		'REM RANDOM returns value between 1 and 32738 - this converts it into'
		'REM a value between 1000 and 9999'
		'set /a port=%RANDOM%*(9999 - 1000 + 1)/32768 + 1000'
		'exit /b')
	}
