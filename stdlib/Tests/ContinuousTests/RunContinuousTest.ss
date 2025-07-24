// Copyright (C) 2017 Axon Development Corporation All rights reserved worldwide.
class
	{
	/* WARNING
	This is called by the timeout testers server_timeout.go file and
	the SuJsWeb testers server_sujsweb.go.
	If you change this record, be sure to verify that the timeout testers are properly
	reporting. If you change the passed in arguments, you will need to manually update
	those files.
	*/
	CallClass(type, shutdown = false)
		{
		Test.EnsureLibrary(Test.TestLibName())
		SuneidoLog.Ensure()
		client = "gsuneido"
		if Sys.Linux?()
			client = './' $ ExeName()
		cmd = client $ " -c 127.0.0.1 -p " $ ServerPort() $ ' -u '

		if type is 'TimeoutTester'
			{
			PutFile('AttemptingToAuthorize.txt', '')
			timetester = SoleContribution('TimeoutTester')()
			Global(timetester).Setup()
			// adding timeout sleep minutes could increase amazon continous tests cost
			cmd $= AuthorizationHandler.AddTokenToCmdLine(
				timetester $ "(minutes: 4.6); Exit();")
			}
		else if type is 'SuJsWeb'
			{
			SuJsWebTester.Setup()
			cmd $= "SuJsWebTester();"
			}
		else
			cmd $= " ContinuousTest_Base.Run(`" $ type $ "`);Exit();"

		Thread({
			Thread.Name('RunContinuousTest-thread')
			// small server needs enough time to start serving,
			// and also avoid "FATAL: lost connection" error from client
			Thread.Sleep(1000)

			AmazonContinuousTestRunner.Log('Client/Server Command', cmd)
			result = RunPipedOutput(cmd)
			// retry if the server is not ready yet
			if result.Has?("invalid response from server")
				result = '(retry) ' $ RunPipedOutput(cmd)
			AmazonContinuousTestRunner.Log('Client/Server Result', result)
			if shutdown is false
				shutdown = Shutdown
			shutdown(alsoServer:)
			})
		return
		}
	}
