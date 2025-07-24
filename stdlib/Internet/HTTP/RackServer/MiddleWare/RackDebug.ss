// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
RackComposeBase
	{
	Call(env)
		{
		try
			return .App(:env)
		catch (e)
			{
			msg = "EXCEPTION: " $ e $ "\n" $
				FormatCallStack(e.Callstack(), indent:)
			Print(msg)
			return ["500 Internal Server Error", msg]
			}
		}
	}