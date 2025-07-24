// Copyright (C) 2013 Axon Development Corporation All rights reserved worldwide.
class
	{
	Table: extra_processes
	Ensure()
		{
		Database('ensure ' $ .Table $ ' (ep_sessionId, ep_process) key(ep_sessionId)')
		}

	EnsureProcess(ep_sessionId, ep_process, start = false)
		{
		QueryEnsure(.Table, [:ep_sessionId, :ep_process])
		if start
			.StartOne(ep_sessionId)
		}

	Start()
		{
		QueryApply(.Table $ ' sort ep_sessionId')
			{ |x|
			.StartOne(x.ep_sessionId)
			}
		}

	StartOne(sessionId)
		{
		if Sys.Client?()
			{
			ServerEval('ExtraProcesses.StartOne', sessionId)
			return
			}
		Thread({ .Run(sessionId) })
		}

	///////////////////////////////////////////////////
	EP_Suffix: ' extra process'
	restartDelayInSecs: 30
	Run(sessionId) // split out for testing
		{
		x = Query1(.Table, ep_sessionId: sessionId)
		Thread.Name(sessionId $ .EP_Suffix)
		Suneido.User = Suneido.User_Loaded = 'none'
		Suneido.user_roles = #('none')

		forever
			{
			try
				{
				x.ep_process.Eval() // Eval should be okay here
				break
				}
			catch (e)
				{
				try
					SuneidoLog('ERROR: (CAUGHT) ExtraProcesses - ' $ sessionId $
						' - ' $ e $ ' - will try to restart function in ' $
						.restartDelayInSecs $ ' seconds.', caughtMsg: 'unattended')
				catch
					throw e // throw original error instead of one from SuneidoLog
				Thread.Sleep(.restartDelayInSecs.SecondsInMs())
				// continue, to call function (e.g. scheduler) again after error
				}
			}
		}
	}
