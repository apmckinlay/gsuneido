// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_Sched()
		{
		f = Scheduler.Sched
		Assert(f('garbage') is: false)
		Assert(f('at 15:30') base: SchedAt)
		Assert(f('every 5 minutes') base: SchedEvery)
		}
	Test_main()
		{
		_log = Object()
		source = function ()
			{
			return #(
				(sched_num: 1, sched_when: 'every 10 minutes',
					sched_func: 'Scheduler_Test.Task1', sched_args: 1)
				(sched_num: 2, sched_when: 'every 15 minutes',
					sched_func: 'Scheduler_Test.Task1', sched_args: 2)
				(sched_num: 3, sched_when: 'at 10:30',
					sched_func: 'Scheduler_Test.Task1', sched_args: 3)
				(sched_num: 4, sched_when: 'at 10:45',
					sched_func: 'Scheduler_Test.Task1', sched_args: 4)
				)
			}
		scheduler = new Scheduler(source)
		runDue = scheduler.Scheduler_runDue
		for (_time = #20160401.1000;
			_time <= #20160401.1050;
			_time = _time.Plus(minutes: 1))
			if _time isnt #20160401.1030 // simulate delay due to long task
				runDue(_time)
		Assert(_log is: #(
			"1: 10:00"
			"2: 10:00"
			"1: 10:10"
			"2: 10:15"
			"1: 10:20"
			"1: 10:31" // ran late due to simulated delay
			"2: 10:31"
			"3: 10:31"
			"1: 10:40"
			"2: 10:45"
			"4: 10:45"
			"1: 10:50"
			))
		}
	Task1(id)
		{
		_log.Add(id $ ': ' $ _time.Format("H:mm"))
		}

	Test_runAsThread?()
		{
		runAsThread? = Scheduler.Scheduler_runAsThread?
		Assert(runAsThread?(#()) is: false)
		Assert({ runAsThread?(#(sched_thread?: true)) } throws:
			"run as thread failed, task needs sched_name member for thread name")
		Assert(runAsThread?(#(sched_name: 'hello')) is: false)
		Assert(runAsThread?(#(sched_name: 'hello', sched_thread?: true)) is: true)
		}

	Test_currentlyRunning?()
		{
		curThreads = #(
			"SocketServer-0-connection-1886 HTTP Server",
			"Thread-4 workqueue extra process",
			"Thread-1 http extra process",
			"SocketServer-thread-pool",
			"SocketServer-0 HTTP Server",
			"Thread-3 scheduler extra process",
			"Thread-2 scheduledreports extra process",
			"SocketServer-0-connection-1887 HTTP Server",
			"Thread-5 scheduler-autoAttachments",
			"Thread-15 scheduler-snapshot")
		.SpyOn(Scheduler.Scheduler_threadList).Return(curThreads)
		currentlyRunning? = Scheduler.CurrentlyRunning?
		Assert(currentlyRunning?('http extra process') is: true)
		Assert(currentlyRunning?('SocketServer-thread-pool') is: false)
		Assert(currentlyRunning?('scheduler extra process') is: true)
		Assert(currentlyRunning?('scheduler-autoAttachments') is: true)
		Assert(currentlyRunning?('scheduler-snapshot') is: true)
		Assert(currentlyRunning?('') is: false)
		Assert(currentlyRunning?('hello world') is: false)
		}

	Test_loggingPrefix()
		{
		loggingPrefix = Scheduler.Scheduler_loggingPrefix
		Assert(loggingPrefix(#(sched_when: 'every 1 minutes')) is: 'ERRATIC')
		Assert(loggingPrefix(#(sched_when: 'every 15 minutes')) is: 'ERRATIC')
		Assert(loggingPrefix(#(sched_when: 'at 1:10')) is: 'ERROR')
		Assert(loggingPrefix(#(sched_when: 'at 15:01')) is: 'ERROR')
		Assert(loggingPrefix(#(sched_when: 'On Sun')) is: 'ERROR')
		Assert(loggingPrefix(#(sched_when: 'on 5 at 06:00')) is: 'ERROR')
		Assert(loggingPrefix(#(sched_when: 'on MidMonth at 09:00')) is: 'ERROR')
		}
	}