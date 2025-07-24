// Copyright (C) 2015 Suneido Software Corp. All rights reserved worldwide.
ScintillaAddon
	{
	DoubleClick()
		{
		if false isnt logNum = .getSuneidoLogNum()
			Inspect.Window(Query1('suneidolog', sulog_timestamp: logNum))
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
		return logNum
		}
	}