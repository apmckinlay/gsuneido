// Copyright (C) 2025 Axon Development Corporation All rights reserved worldwide.
class
	{
	New()
		{
		.start = Timestamp()
		Suneido.User = 'none'
		.UniqueId = SuRenderBackend().NextId()
		SuRenderBackend().Register(.UniqueId, this)
		SuRenderBackend().RecordAction(false, 'RunSuJsWebTestClient',
			args: [.UniqueId])

		tests = SuJsWebTester.GetTests()
		SuRenderBackend().RecordAction(.UniqueId, 'RunTests', args: [tests])

		SuRenderBackend().RegisterBeforeDisconnectFn(.onDisconnect)
		}

	results: 'NO RESULTS'
	TestsResults(.results)
		{
		.results = ContinuousTest_Base.CheckBook(.results, #())
		.Finish()
		}

	finished?: false
	Finish()
		{
		SuRenderBackend().RecordAction(.UniqueId, 'Finish', args: #())
		SuRenderBackend().UnRegister(.UniqueId)

		.output(.results)

		.finished? = true
		ServerSuneido.Set(#sujswebtest_done, true)
		}

	onDisconnect()
		{
		if .finished? is true
			return

		results = 'ERROR: NOT FINISH\r\n' $ .results $ '\r\n\r\n' $ .sulogs()
		.output(results)
		}

	sulogs()
		{
		sulogs = Object()
		QueryApply(`suneidolog where sulog_timestamp > ` $ Display(.start))
			{
			sulogs.Add(it.sulog_timestamp.StdShortDateTimeSec() $ ' ' $
				it.sulog_message)
			}
		return sulogs.Join('\r\n')
		}

	output(results)
		{
		currentDir = GetCurrentDirectory()
		msg = "Instance: " $ currentDir $ "\r\n\r\n" $ results $
			'\r\n' $ ContinuousTests_ErrorLogs()
		testSetName = 'SuJsWeb Tester - ' $ .exeType(currentDir)
		PutFile('stillRunningResults.txt',
			ContinuousTest_Base.FormatResultContent(msg, :testSetName))
		}

	exeType(exePath)
		{
		switch
			{
		case exePath.Has?('current'):
			return 'Current Exe'
		case exePath.Has?('latest'):
			return 'Latest Exe'
		case exePath.Has?('gdev'):
			return 'Dev Exe'
		default:
			return 'Current Exe'
			}
		}
	}
