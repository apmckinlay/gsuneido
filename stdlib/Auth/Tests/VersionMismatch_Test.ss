// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_programNamesAndSetDir()
		{
		cl = VersionMismatch
			{
			VersionMismatch_exePath()
				{
				return _versionMismatchPath
				}
			VersionMismatch_setDir(dir)
				{
				Assert(dir is: _versionMismatchParentPath)
				}
			}
		func = cl.VersionMismatch_programNamesAndSetDir
		_versionMismatchPath = 'C:/work/stuff/gsuneido.exe'
		_versionMismatchParentPath = 'C:/work/stuff'
		names = func()
		Assert(names.appName is: 'gsuneido')
		Assert(names.exeName is: 'gsuneido.exe')

		_versionMismatchPath = 'C:\\work\\stuff\\gsuneido.exe'
		_versionMismatchParentPath = 'C:\\work\\stuff'
		names = func()
		Assert(names.appName is: 'gsuneido')
		Assert(names.exeName is: 'gsuneido.exe')

		_versionMismatchPath = 'gsuneido.exe'
		_versionMismatchParentPath = ''
		names = func()
		Assert(names.appName is: 'gsuneido')
		Assert(names.exeName is: 'gsuneido.exe')
		}

	Test_getAppFolder()
		{
		func = VersionMismatch.VersionMismatch_getAppFolder
		Assert(func('') is: '')
		Assert(func('AppFolder=') is: '')
		Assert(func('-c hostserver -p 3147 AppFolder=C:/bob') is: 'C:/bob')
		Assert(func('-c hostserver -p 3147 AppFolder=C:/longer/folder/path/then/bob')
			is: 'C:/longer/folder/path/then/bob')
		Assert(func('-c hostserver -p 3147') is: '')
		}

	Test_getGetLatest()
		{
		_versionMismatchAlert = Object()
		cl = VersionMismatch
			{
			VersionMismatch_dir(unused)
				{
				return _versionMismatchDir
				}
			VersionMismatch_alert(msg, detail, appName)
				{
				_versionMismatchAlert.Add(Object(:msg, :detail, :appName))
				}
			}

		func = cl.VersionMismatch_getLatestExe

		_versionMismatchDir = Object('gsuneido20240625.exe', 'gsuneido20240729.exe')
		Assert(func('gsuneido', '') is: '/gsuneido20240729.exe')
		Assert(func('gsuneido', 'gsuneidoFiles') is: 'gsuneidoFiles/gsuneido20240729.exe')
		_versionMismatchDir = Object('csuneido20180101.exe', 'csuneido20190101.exe')
		Assert(func('csuneido', 'oldFiles') is: 'oldFiles/csuneido20190101.exe')

		_versionMismatchDir = Object('gsuneido20200101.exe', 'gsuneido20180225.exe',
			'gsuneido20190525.exe', 'gsuneido20240730.exe', 'gsuneido20210625.exe',
			'gsuneido20240729.exe')
		Assert(func('gsuneido', 'gsuneidoFiles') is: 'gsuneidoFiles/gsuneido20240730.exe')

		_versionMismatchDir = Object()
		Assert(_versionMismatchAlert isSize: 0)
		func('gsuneido', 'serverShare')

		Assert(_versionMismatchAlert[0].msg
			is: 'Cannot read server shared folder serverShare, ' $
				'this is possibly caused by network issues')
		Assert(_versionMismatchAlert[0].detail is: 'getLatestExe failed - #()')
		Assert(_versionMismatchAlert[0].appName is: 'gsuneido')

		_versionMismatchDir = Object('gsuneido.exe')
		func('gsuneido', 'serverShare')
		Assert(_versionMismatchAlert[1].msg
			is: 'Cannot read server shared folder serverShare, ' $
				'this is possibly caused by network issues')
		Assert(_versionMismatchAlert[1].detail
			is: 'getLatestExe failed - #("gsuneido.exe")')
		Assert(_versionMismatchAlert[1].appName is: 'gsuneido')
		}
	}