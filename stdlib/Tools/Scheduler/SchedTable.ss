// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
class
	{
	New(.tablename)
		{
		}
	Call()
		{
		return QueryAll(.tablename $ ' where sched_suspended isnt true')
		}
	Ensure(tablename = 'schedtasks')
		{
		Database('ensure ' $ tablename $
			' (sched_num, sched_when, sched_func, sched_args,
				sched_suspended)
			key (sched_num)
			key (sched_when, sched_func, sched_args)')
		}
	}
