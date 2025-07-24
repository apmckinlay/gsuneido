// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// passes unhandled messages to next controller
Controller
	{
	Msg(args)
		{
		msg = args[0]
		if .Method?(msg)
			return this[msg](@+1 args)
		if .Member?(#Controller)
			return .Controller.Msg(args)
		return 0
		}
	}