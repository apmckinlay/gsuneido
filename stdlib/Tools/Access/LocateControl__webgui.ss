// Copyright (C) 2022 Axon Development Corporation All rights reserved worldwide.
_LocateControl
	{
	New(@args)
		{
		super(@args)
		if .Member?(#Horz)
			SuRenderBackend().AddOverrideProc(.Horz.Locate.Field.UniqueId, .enterProc)
		}

	enterProc(uniqueId/*unused*/, event, args)
		{
		if event is 'KEYDOWN' and args[0] is VK.RETURN
			{
			.Horz.Locate.FieldReturn()
			.Send("On_Go")
			return false
			}
		return true
		}

	Destroy()
		{
		if .Member?(#Horz)
			SuRenderBackend().RemoveOverrideProc(.Horz.Locate.Field.UniqueId, .enterProc)
		super.Destroy()
		}
	}