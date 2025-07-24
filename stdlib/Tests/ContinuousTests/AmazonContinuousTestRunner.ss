// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
class
	{
	RunJob()
		{
		try
			.runJob()
		catch (err)
			{
			AddFile('stillRunningResults.txt', '\nERROR: ' $ err)
			Shutdown(true)
			}
		}

	runJob()
		{
		.Log('AmazonContinuousTestRunner', 'starting')
		job = Json.Decode(GetFile('job.json'))
		.Log('AmazonContinuousTestRunner', job)
		Suneido.NoCredentials = true
		Unuse('axonlib')
		.useLibs(job)
		ServerSuneido.Set('ServerContinuousTestInfo', job)
		.Log('RunningTests', job)
		if job.testGroup is 'Standalone'
			{
			ContinuousTest_Base.Run('ServerContinuousTestInfo')
			Exit()
			}
		else // Client/Server
			RunContinuousTest('ServerContinuousTestInfo')
		}

	useLibs(job)
		{
		for lib in job.libs
			{
			if not TableExists?(lib)
				Database.Load(lib)
			Use(lib)
			}
		LibraryTags.Reset()
		}

	Log(type, text = '')
		{
		try ServerEval('AddFile', 'running.log',
			String(Timestamp()) $ ': ' $ type $ ' - ' $ String(text) $ '\r\n')
		}
	}
