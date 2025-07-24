// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.
class
	{
	// can either create an instance or inherit and define Remote:
	New(remote = false)
		{
		if remote isnt false
			.Remote = remote
		}
	Default(@args)
		{
		args.Add('ServerEvalProxy.Forward', .Remote, at: 0)
		return ServerEval(@args)
		}
	Forward(@args)
		{
		remote = args[0].Eval() // needs to use Eval
		method = args[1]
		remote[method](@args[2..])
		}
	}