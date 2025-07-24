// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
// WARNING: doesn't handle multiple Start's (e.g. from multiple books)
class
	{
	Start()
		{
		if not Suneido.Member?(#TimerManager)
			Suneido.TimerManager = new TimerManagerImpl()
		else
			Print("TimerManager: duplicate start not handled")
		}
	Stop()
		{
		if Suneido.Member?(#TimerManager)
			{
			Suneido.TimerManager.Stop()
			Suneido.Delete(#TimerManager)
			}
		}
	}