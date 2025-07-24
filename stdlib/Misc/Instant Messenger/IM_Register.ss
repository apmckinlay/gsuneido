// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
function(env)
	{
	IM_MessengerManager.LogConnectionsIfChanged()
	args = IM_MessengerManager.ParseQuery(env.body)
	if not Suneido.Member?('IM_Opened')
		Suneido.IM_Opened = Object()
	// lost connections leave IM_Opened - need to also check # of connections
	if Suneido.IM_Opened.Has?(args.user) and
		Sys.Connections().CountIf({ it.Prefix?(args.user $ '@') }) > 1
		return Json.Encode(#(allowOpen: false))
	Suneido.IM_Opened.AddUnique(args.user)
	return Json.Encode(#(allowOpen: true))
	}