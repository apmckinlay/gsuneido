// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.

// should include TEMPINDEX but there are too many of them
class
	{
	CallClass(observer, block)
		{
		if Sys.Client?()
			{
			block()
			return 0
			}
		.ensure_logfile_exists()
		File("trace.log", "r")
			{|f|
			f.Seek(0, 'end')

			d = Date()
			Trace(TRACE.LOGFILE | TRACE.SLOWQUERIES | TRACE.ALLINDEX)
				{
				Trace("*before* " $ Display(d))
				block()
				Trace("*after* " $ Display(d))
				}

			return .count_slow(observer, f, d)
			}
		}
	ensure_logfile_exists()
		{
		// delete will only work if trace.log isnt open yet
		// but it should work very first time, so file won't grow forever
		DeleteFileApi("trace.log")
		if not FileExists?("trace.log")
			PutFile("trace.log", "")
			// only do this if file doesn't exist
			// causes problems if this is done when exe has trace.log open
		}
	count_slow(observer, f, d)
		{
		n = 0
		Retry()
			{ Assert(f.Readline() is: "*before* " $ Display(d)) }
		for (last = false; false isnt line = f.Readline(); last = line)
			if line.Prefix?("SLOWQ") or line.Prefix?("ALLINDEX")
				{
				++n
				observer(line)
				}
		Assert(last is: "*after* " $ Display(d))
		return n
		}
	}