// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
// Run functions at scheduled times
// source argument is either a table name (which will be passed to SchedTable)
// or a function that returns a list of tasks with:
// (sched_when, sched_func, sched_args)
/* example with table, runs for 5 minutes
SchedTable.Ensure('schedtest')
QueryEnsure('schedtest', [sched_num: 1,
	sched_when: 'every 1 minute', sched_func: 'Print', sched_args: "'hello'"])
stop = Date().Plus(minutes: 5)
QueryEnsure('schedtest', [sched_num: 2,
	sched_when: 'at ' $ stop.Format('H:mm'), sched_func: 'SchedExit'])
Scheduler('schedtest')
*/
class
	{
	scheds: (SchedEvery, SchedAt, SchedOn, SchedMonthlyOn)
	sourceName: 'scheduler'
	CallClass(source = 'schedtasks', sourceName = 'scheduler')
		{
		(new this(source, sourceName)).Run()
		}
	New(.source, .sourceName = 'scheduler')
		{
		if String?(source)
			.source = SchedTable(.sourceName = source)
		.bad = Object().Set_default(false)
		}
	RunEveryMinute(func)
		{
		taskFn = { Object(Object(sched_when: "every 1 minute",
			sched_func: func, sched_name: func, sched_args: '')) }
		Scheduler(taskFn, sourceName: func)
		}
	Run()
		{
		while not Suneido.GetDefault(#SchedulerExit, false)
			{
			SleepUntil(Date().Plus(minutes: 1))
			.runDue(Date())
			}
		}
	prevcheck: false
	runDue(curtime) // pass in time for testing
		{
		.curtime = curtime
		(.source)().Filter(.due?).Each(.run)
		.prevcheck = .curtime
		}
	due?(task)
		{
		if false is sched = .Sched(task.sched_when)
			{
			SuneidoLog("ERROR: Scheduler: bad when: " $ Display(task.sched_when),
				params: task)
			return false
			}
		return sched.Due?(.prevcheck, .curtime)
		}
	Sched(when) // static method, can also be used for validation
		{
		for s in .scheds
			if false isnt sched = Global(s)(when)
				return sched
		return false
		}
	run(task)
		{
		if not .checktask(task)
			return
		expr = (task.sched_func $ '(' $ task.sched_args $ ')')
		.track('SchedulerLastProcessStarted', expr, task)
		if .runAsThread?(task)
			.schedAsThread(task, expr)
		else
			.tryTask(task, expr)
		}

	tryTask(task, expr)
		{
		try
			{
			expr.Eval() // needs Eval
			.track('SchedulerLastProcessCompleted', expr, task)
			}
		catch (e)
			{
			SuneidoLog('ERROR: (CAUGHT) Scheduler: ' $ e, params: task,
				caughtMsg: 'unattended')
			// Use Opt(e) to convert the except to a string to avoid saving
			// any variables of the except to Suneido object
			.track('SchedulerLastProcessCompleted', expr $ ' ERROR ' $ Opt(e), task)
			}
		}
	runAsThread?(task)
		{
		if task.GetDefault('sched_thread?', false) isnt true
			return false
		if not task.Member?('sched_name')
			throw 'run as thread failed, task needs sched_name member for thread name'
		return true
		}

	schedAsThread(task, expr)
		{
		name = .sourceName $ '-' $ task.sched_name
		if .CurrentlyRunning?(name)
			{
			lastLog = ServerSuneido.Get("SchedulerLastLog_" $ task.sched_name, "unknown")
			SuneidoLog(.loggingPrefix(task) $ ': Scheduler - previous ' $
				name $ ' is still running', params: Object(:lastLog))
			return
			}
		Thread({ .tryTask(task, expr) }, :name)
		}

	loggingPrefix(task)
		{
		return task.sched_when.Prefix?('every')	? 'ERRATIC' :'ERROR'
		}

	CurrentlyRunning?(threadName)
		{
		curThreads = .threadList()
		return curThreads.Any?({ it =~ '^Thread-\d+ ' $ threadName $ '$' })
		}

	threadList()
		{
		return ServerEval('Thread.List')
		}

	checktask(task)
		{
		// only report bad tasks once
		// we compare entire task, so if it is updated we'll try again
		if .bad.Has?(task)
			return false
		cf = .checkfunc(task)
		ca = .checkargs(task)
		if cf is '' and ca is ''
			return true
		.bad.Add(task)
		SuneidoLog('ERROR: Scheduler: ' $ Join(', ', cf, ca), params: task)
		return false
		}
	checkfunc(task)
		{
		return task.sched_func =~ '^[A-Z][_a-zA-Z0-9]+([.][A-Z][_a-zA-Z0-9]+)?$'
			? "" : "invalid func"
		}
	checkargs(task)
		{
		try
			('#(' $ task.sched_args $ ')').SafeEval() // ensure args are constants
		catch
			return "invalid args"
		return ""
		}
	track(mem, msg, task)
		{
		ob = ServerSuneido.Get(mem, Object())
		ob[.sourceName] = Date().ShortDateTime() $ ' ' $ msg
		ServerSuneido.Set(mem, ob)
		if not task.sched_when.Has?('every 1 minute')
			BookLog(mem $ ' ' $ msg, name: .sourceName)
		}
	}