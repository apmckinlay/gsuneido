// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass()
		{
		if not Sys.Client?()
			{ // standalone mode
			.handleResponse(EventManager.RunServerFuncs(Suneido.User))
			return
			}

		if not .lockEventManager()
			{
			lastrun = Suneido.GetDefault(#EventManagerLastRun, '')
			startedlastrun = Suneido.GetDefault(#EventManagerStarted, '')
			finishedlastrun = Suneido.GetDefault(#EventManagerFinished, '')
			threadCount = Thread.Count()
			SuneidoLog("INFO: overlapping EventManager call (ignored)" $
				(lastrun is "" ? "" : (", last run handleResponse " $ Display(lastrun))) $
				(startedlastrun is ""
					? "" : (", started last run " $ Display(startedlastrun))) $
				(finishedlastrun is ""
					? "" : (", finished last run " $ Display(finishedlastrun))) $
				" threads: " $ threadCount)
			return
			}
		// run in the background (async) so the user doesn't have to wait if it's slow
		try
			RunOnHttp(HttpPort(), 'EventManagerServerFuncs', [Suneido.User],
				asyncCompletion: .handleResponseOnCompletion)
		catch (e)
			{
			Suneido.EventManagerRunning = false
			throw e
			}
		}
	lockEventManager()
		{
		.Synchronized()
			{
			if Suneido.GetDefault(#EventManagerRunning, false)
				return false
			Suneido.EventManagerRunning = true
			Suneido.EventManagerStarted = Date()
			return true
			}
		}
	RunServerFuncs(user) // run on http server
		{
		IM_MessengerManager.LogConnectionsIfChanged()
		ob = Object()
		Plugins().ForeachContribution('Events', 'eventfunction', showErrors:)
			{ |c|
			LogErrors("EventManager", asErratic: #("can't connect to 127.0.0.1"))
				{
				result = (c.serverfunc)(:user)
				ob.Add(result, at: c.name)
				}
			}
		return ob
		}
	handleResponseOnCompletion(result)
		{
		content = result.content
		if content.Prefix?('ERROR ')
			response = content $ ' (from httpserver)'
		else
			{
			try
				response = Unpack(content)
			catch // handleReponse should get a proper error message so just set response
				response = content
			}
		.handleResponse(response)
		}
	handleResponse(result) // run on client
		{
		Suneido.EventManagerLastRun = Date()
		if not Object?(result)
			{
			if String?(result) and result =~ '(?i)socket|connect|not authorized'
				SuneidoLog.Once("INFO: EventManager: " $ result)
			else // if result is an error, don't use Display, because it removes callstack
				{
				resultStr = String(result)
				// content length checking fails sporadically, see more details on 21277
				errPrefx = resultStr.Prefix?(HttpClient.ContentLengthErrPrefix)
					? 'ERRATIC'
					: 'ERROR'
				errStr = errPrefx $ ": (CAUGHT) Event Manager bad result: " $ resultStr
				try
					SuneidoLog(errStr, caughtMsg: 'unattended; no msg to user')
				catch (err)
					ErrorLog(errStr $ ' (err: ' $ err $ ')')
				}
			Suneido.EventManagerFinished = Date()
			Suneido.EventManagerRunning = false
			return false
			}
		Plugins().ForeachContribution('Events', 'eventfunction', showErrors:)
			{ |c|
			LogErrors("EventManager - " $ c.name)
				{
				if result.Member?(c.name) // not a member if the server side failed
					(c.clientfunc)(result: result[c.name])
				}
			}
		Suneido.EventManagerFinished = Date()
		Suneido.EventManagerRunning = false
		}
	}
