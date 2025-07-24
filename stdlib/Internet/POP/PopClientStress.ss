// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
function (seconds = 10, minutes = 0,
	addr = '127.0.0.1', user = 'user', password = 'password')
	{
	i = 0
	seconds += minutes * 60
	start_time = Date()
	stop_time = start_time.Plus(seconds: seconds)
	while Date() < stop_time
		{
		AddFile('log', ++i $ '\n')
		pc = PopClient(addr, user, password)
		pc.List()
		pc.GetMessage(1)
		pc.Close()
		Thread.Sleep(100)
		}
	}