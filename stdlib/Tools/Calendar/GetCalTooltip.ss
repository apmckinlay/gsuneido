// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(env)
		{
		args = env.queryvalues.Set_default('')

		validFn = OptContribution('ValidCalendarUser', function(@unused) { return '' })
		if false is user = (validFn)(args)
			return 'Please Login as valid user.'

		page = args.page
		if not page.GlobalName?()
			return "Invalid Request"

		ev_combiner = .getCombiner(page)
		if String?(ev_combiner)
			{
			SuneidoLog('ERROR: ' $ ev_combiner, calls:)
			return 'Invalid Event Sources'
			}

		for ev_source in ev_combiner.GetEventSources(user)
			if ev_source.EventType is args.event_type and
				ev_source.Member?('Tooltip')
				return ev_source.Tooltip(args.id)

		return "Not Found!"
		}

	getCombiner(page)
		{
		try
			{
			ev_combiner = GetContributions('RackRoutes').
				FindOne({ it[1].AfterFirst('/') is page })[2]
			if String?(ev_combiner)
				ev_combiner = Global(ev_combiner)
			if not Class?(ev_combiner) or not ev_combiner.Member?('GetEventSources')
				return 'Invalid Event Sources'
			}
		catch (err)
			return 'Invalid Event Sources: ' $ err
		return ev_combiner
		}
	}