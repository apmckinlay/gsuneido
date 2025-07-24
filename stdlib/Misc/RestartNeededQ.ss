// Copyright (C) 2004 Suneido Software Corp. All rights reserved worldwide.
function (margin/*unused*/ = 50000000, memory = 48000000, days = 7)
	{
	return ServerEval("MemoryArena") > memory or
		Date().MinusDays(ServerSuneido.Get("start_time")) >= days
	}