// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
#(
	schedulerHeartBeat: (sched_when: 'every 1 minute', sched_func: 'ServerSuneido.Set',
		sched_args: '#SchedulerHeartBeat, Timestamp()')
	threadChecker: (sched_when: 'every 30 minutes', sched_func: 'ThreadChecker',
		sched_args: '')
)
