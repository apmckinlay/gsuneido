// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
/*
to stop the server:
	SocketClient('127.0.0.1', TestSocketServer.Port)
		{ it.Writeline("exit") }
*/
SocketServer
	{
	Name: 'TestSocketServer'
	Port: 1234
	before: 0
	during: 0
	after: 0
	close: false
	Run()
		{
		.init()
		try
			.run()
		catch (e)
			Print(error: e)
		}
	Killer(killer)
		{
		Suneido.TestSocketServerKiller = killer
		}
	init()
		{
		Print('--------')
		.small = "helloworld".Repeat(10) // 100
		.large = "helloworld".Repeat(10000) // 100,000
		.line = .small $ '\n'
		.longline = "helloworld".Repeat(500) $ '\n' // 5000 > max of 4000
		}
	run()
		{
		while (false isnt (req = .Readline()))
			{
			Print(:req)
			switch req.BeforeFirst(' ')
				{
			case 'line':
				.output(.line)
			case 'longline':
				.output(.longline)
			case 'small':
				.output(.small)
			case 'large':
				.output(.large)
			case 'before':
				.before = Number(req.AfterFirst(' '))
			case 'during':
				.during = Number(req.AfterFirst(' '))
			case 'after':
				.after = Number(req.AfterFirst(' '))
			case 'close':
				.close = true
			case 'exit':
				Suneido.TestSocketServerKiller.Kill()
			default:
				Alert("invalid request: " $ req)
				Suneido.TestSocketServerKiller.Kill()
				}
			}
		}
	output(data)
		{
		if .before > 0
			Thread.Sleep(.before)
		if .during <= 0
			.Write(data)
		else
			{
			half = (data.Size() / 2).Int()
			.Write(data[.. half])
			Thread.Sleep(.during)
			.Write(data[half ..])
			}
		if .after > 0
			Thread.Sleep(.after)
		if .close
			throw "close"
		}
	}
