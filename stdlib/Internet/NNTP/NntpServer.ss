// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
SocketServer
	{
	Name: "NNTP Server"
	Port: 119

	Impl: NntpSampleImpl
	logging: true
	Run()
		{
		user = .RemoteUser()
		.log(user, 'SERVING NNTP SERVER')
		.State = 'Start'
		try
			{
			.Writeline('200 suneido news server ready')
			while (.State isnt 'Closed' and false isnt request = .Readline())
				{
				.log(user, 'REQUEST: ' $ request)
				command = request.Extract("^[a-zA-Z]*").Upper().Tr(' ', '_')
				if command is ''
					{
					.State = 'Closed'
					.Writeline('400 invalid command')
					return
					}
				impl = Global(.Impl)
				if not impl.Member?(command)
					{
					.State = 'Closed'
					SuneidoLog('ERROR: NNTP command not implemented - ' $ request)
					.Writeline('501 command not implemented')
					return
					}
				args = request[command.Size()..].Trim()
				result = Global(.Impl)[command](:args, server: this)
				.writeResponse(result, user)
				}
			}
		catch (e)
			.errorHandler(e, user)
		.log(user, 'DONE SERVING NNTP SERVER')
		}

	writeResponse(result, user)
		{
		if String?(result)
			{
			.Writeline(result)
			.log(user, 'RESPONSE: ' $ result)
			}
		else
			{
			for s in result
				.Writeline(s)
			.log(user, 'RESPONSE: ' $ result.GetDefault(0, '') $
				', size: ' $ result.Size())
			}
		}

	errorHandler(e, user)
		{
		.log(user, 'ERROR: ' $ e)
		if e.Has?('lost connection') or e.Prefix?('socket') or .State is 'Closed'
			return
		SuneidoLog("ERROR: NntpServer: " $ e)
		try .Writeline('500 Internal Server Error')
		}

	log(user, s)
		{
		if .logging
			Rlog('logs/nntp/nntp', user $ ' - ' $ Thread.Name() $ ' - ' $ s)
		}
	}
