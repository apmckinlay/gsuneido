// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
function(@unused)
	{
	IM_MessengerManager.LogConnectionsIfChanged()
	Suneido.IM_Contacts = Date()
	return Json.Encode(#())
	}