// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
function(env)
	{
	IM_MessengerManager.LogConnectionsIfChanged()
	args = IM_MessengerManager.ParseQuery(env.body)
	openedOb = Suneido.GetDefault('IM_Opened', #())
	openedOb.Remove(args.user)
	return Json.Encode(#(true))
	}