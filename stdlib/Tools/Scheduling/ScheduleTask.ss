// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
class
	{
	New(args)
		{
		.uid = 0				// unique identiying int
		this['taskTable'] = args.taskTable
		this['taskname'] = args.taskname
		this['task'] = args.task
		this['runinterval'] = args.GetDefault('runinterval', false)	// boolean
		this['rundaily'] = args.GetDefault('rundaily', false)	// boolean
		this['runweekly'] = args.GetDefault('runweekly', false)	// boolean
		this['runmonthly'] = args.GetDefault('runmonthly', false)	// boolean
		this['rundate'] = args.GetDefault('rundate', false)	// boolean
		this['suspended'] = args.GetDefault('suspended', false)	// boolean
		this['run_skipped'] = args.GetDefault('run_skipped', false)
		this['threaded'] = false	// boolean
		this['interval'] = args.GetDefault('interval', 0)	// int
		this['interval_units'] = args.GetDefault('interval_units', 0)	// int 0..4
		this['daily_time'] = args.GetDefault('daily_time',
			Date().Hour() * 100 + Date().Minute())	// int /*=2400 format */
		this['weekly_time'] = args.GetDefault('weekly_time', this['daily_time'])// int
		this['weekly_day'] = args.GetDefault('weekly_day', 0)			// int 0..6
		this['monthly_time'] = args.GetDefault('monthly_time', this['daily_time']) // int
		this['monthly_day'] = args.GetDefault('monthly_day', 0)			// int 0..30
		this['date'] = args.GetDefault('date',
			Date().NoTime().Plus(hours: Date().Hour() minutes: Date().Minute())) // date
		this['prev_event'] = args.GetDefault('prev_event', Date.Begin())	// date
		this['next_event'] = args.GetDefault('next_event', Date.End())		// date
		this['prev_interval'] = args.GetDefault('prev_interval', Date())	// date
		this['time_span'] = args.GetDefault('time_span', false)		// boolean
		this['time_start'] = args.GetDefault('time_start', 0)		// int
		this['time_end'] = args.GetDefault('time_end', 0)		// int
		this['uid'] = .getUID(this['taskTable'])
		}
	getUID(taskTable)
		{
		if taskTable is ""
			return .MinUID()
		uIDMax = Query1(taskTable $ " summarize max uid")
		return uIDMax is false ? .MinUID() : ++uIDMax.max_uid
		}
	UpdateTaskRecord(taskRec)
		{
		for field in taskRec.Members()
			taskRec[field] = this[field]
		// NOTE: if taskRec contains a uid field, it is overwritten with this.uid.
		//		this may not be desirable and it may be necessary to save the old
		//		value of uid and reset it after a call to UpdateTaskRecord
		}
	MinUID()
		{ return 0 }
	}
