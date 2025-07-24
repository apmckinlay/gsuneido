// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
function ()
	{
	if false is failed = GetFile('failedTests.txt')
		return

	dir = Dir('*.*')
	if false is index = dir.FindIf({ it.Prefix?('results_') })
		{
		path = ExeDir()
		testName = false is path.Has?('standalone')
			? 'Standalone'
			: 'Client/Server'
		ContinuousTest_Base.SendResults('Continuous Tests - ' $ testName $ ' Failure',
			'Expected a results_ file but none was found.  All tests failed to complete')
		return
		}

	msg = '\r\n===============================================================\r\n'
	for test in failed.Lines()
		msg $= "ERROR: Following test crashed after starting: " $
			test.BeforeFirst(' - ') $  "\r\n"

	AddFile(dir[index], msg)
	DeleteFile('failedTests.txt')
	}
