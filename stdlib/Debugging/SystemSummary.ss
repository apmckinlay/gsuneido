// Copyright (C) 2004 Suneido Software Corp. All rights reserved worldwide.
function ()
	{
	if TestRunner.RunningTests?() // to avoid calling slow system info script
		return "Testing OS"
	if Suneido.Member?(#SystemSummary)
		return Suneido.SystemSummary
	m = SystemMemory()
	m = m / 1.Gb()
	m = (m * 2).Round(0) / 2 // round to nearest .5
	osname = SystemInfo().OSName
	return Suneido.SystemSummary = osname.Replace("Windows ", "Win") $ " " $ m $ "gb"
	}
