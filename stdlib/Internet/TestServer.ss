// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
SocketServer
	{
	Name: 'TestServer'
	Port: 1234
	Run()
		{
		.Writeline("hello")
		while false isnt req = .Readline()
			if req is "quit"
				break
			else
				.Writeline("don't know how to " $ req)
		}
	}
