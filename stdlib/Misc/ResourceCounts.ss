// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
function ()
	{
	info = #("windows.nDefer", "builtin.nRunPiped", "builtin.nSocketClient",
		"windows.nGdiObject", "windows.nUserObject", "windows.nCallback",
		"windows.nWndProc", "windows.nTimer", "builtin.nThread", "builtin.nFile")
	info = info.Intersect(Suneido.Info())
	result = Object()
	for m in info
		if 0 isnt v = Suneido.Info(m)
			{
			m = m.Replace("(windows|builtin).n")
			result[m] = v
			}
	result["GoRoutines"] = Suneido.GoMetric("/sched/goroutines:goroutines")
	return result
	}
