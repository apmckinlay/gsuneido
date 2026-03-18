// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
ScintillaAddon
	{
	DoubleClick()
		{
		if false is Suneido.Member?('SuLogs')
			return
		if false isnt rec = .getSuneidoLogNum()
			Inspect.Window(rec)
		}

	prefix: 'SuneidoLog ['
	getSuneidoLogNum()
		{
		line = .GetLine()
		if not line.Prefix?(.prefix)
			return false
		log = line.AfterFirst(.prefix).BeforeFirst(']')
		if false is logNum = Date(log)
			return false

		for rec in Suneido.SuLogs
			if rec.sulog_timestamp is logNum
				return rec
		return false
		}
	}