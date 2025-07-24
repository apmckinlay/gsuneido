// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
//TODO: when using Json, should using boolean value not the 'true'/'false' string.
//TODO: try to use Json instead of concatenating strings for return events
class
	{
	CallClass(env)
		{
		return .HandleRequest(env.queryvalues.Set_default(''))
		}

	HandleRequest(args)
		{
		if false is user = .ValidUser(args)
			return 'Please Login as valid user.'

		if args.Has?('types')
			return .FormatType(user)

		from = Date(args.from).NoTime()
		to = from.Plus(days: args.weeks * .weekdays - 1)
		evt_src = .parse_src(args.req)
		strReq = '"from":"' $ from.Format('yyyy-M-d') $ '",' $
			'"to":"' $ to.Format('yyyy-M-d') $ '",' $
			'"src":"' $ args.req $ '",'
		return .format(strReq, .combine_events(from, to, evt_src, user))
		}

	ValidUser(args /*unused*/)
		{
		return ''
		}

	GetEventSources(user /*unused*/)
		{
		return .EventSources().Map(Global)
		}
	EventSources()
		{
		sources = Object()
		Plugins().ForeachContribution('CalendarEvents', 'eventSources')
			{ |x|
			if x.Member?('source') and .ValidSource?(x)
				sources.Add(x.source)
			}
		return sources
		}
	ValidSource?(x /*unused*/)
		{ return true }

	FormatType(user)
		{
		ts = Object()
		i = 0
		for c in .GetEventSources(user)
			{
			t = Object()
			t.Index = i++

			t.Inverted = c.Member?('Inverted') and c.Inverted ? 'true' : 'false'
			t.Display = not c.Member?('Display') or c.Display ? 'block' : 'none'
			if c.Member?('Tooltip')
				t.Tooltip = 'true'
			if c.Member?('Color')
				t.Color = c.Color
			.formatSubTypes(c, user, t)
			ts[c.EventType] = t
			}
		return Json.Encode(ts)
		}

	formatSubTypes(c, user, t)
		{
		if c.Method?('SubTypes')
			{
			subtypes = c.SubTypes(:user).Copy()
			subtypes.Map!({|subtype|
				subtype = subtype.Copy()
				subtype.Display = (not subtype.Member?('Display') or
					subtype.Display is true)
					? 'block' : 'none'
				subtype
				})
			t.SubTypes = subtypes
			}
		}

	parse_src(srcStr)
		{
		evt_src = Object()
		src = srcStr.Split('$*$')
		for s in src
			if not s.Has?('__')
				evt_src[s] = Object()
			else
				{
				t = s.BeforeFirst('__')
				st = s.AfterLast('__')
				if not evt_src.Member?(t)
					evt_src[t] = Object()
				evt_src[t].Add(st)
				}
		return evt_src
		}

	combine_events(from, to, src, user)
		{
		events = Object()
		for ev_source in .GetEventSources(user)
			if src.Member?(ev_source.EventType)
				{
				sourceEvents = Object()
				try
					sourceEvents = ev_source(from, to, src[ev_source.EventType], :user)
				catch(err, 'win32 exception: ACCESS_VIOLATION')
					SuneidoLog('ERRATIC: ' $ err)

				for e in sourceEvents
					.add_events(e, from, to, events)
				}
		events.Sort!(By(#date))
		return events
		}

	add_events(e, from, to, events)
		{
		if not e.Member?("subtype")
			e.subtype = ""

		if not e.Member?("multidays")
			{
			e.start_date = e.date
			events.Add(e)
			return
			}

		if e.multidays is 1
			{
			e.start_date = e.date
			e.Delete('multidays')
			events.Add(e)
			return
			}

		sd = e.date	//start date
		if sd > from
			{
			e.completed = 0
			e.start_date = e.date
			}
		else
			{
			e.completed = from.MinusDays(sd)
			e.start_date = e.date
			e.date = from
			}
		e.start_date = e.start_date.NoTime()
		e.span = .span(e.date, e.multidays - e.completed)
		e.end = (e.completed + e.span >= e.multidays)
		events.Add(e)

		.extendEvents(e, to, events)
		}
	weekdays: 7
	span(d, lastDays)
		{
		return Min(lastDays, .weekdays - d.WeekDay())
		}
	extendEvents(e, to, events)
		{
		if e.span isnt e.multidays - e.completed
			{
			event = .copy_event(e)
			while not event.end and event.date.Plus(days: .weekdays) < to
				{
				event.date = event.date.Plus(days: event.span)
				event.completed = event.completed + event.span
				event.span = .span(event.date, event.multidays - event.completed)
				event.end = (event.completed + event.span >= event.multidays)
				events.Add(event)
				event = .copy_event(event)
				}
			}
		}
	format(strReq, events)
		{
		s = '({' $ strReq $ '"events":['
		for e in events
			{
			s $= '{"type":"' $ e.type $ '",' $
				'"subtype":"' $ e.subtype $ '",' $
				'"start_date":"' $ e.start_date.Format('yyyy-M-d') $ '",'

			if not e.Member?("multidays")
				s $= '"multidays":1,"span":1,'
			else
				{
				end_date = e.start_date.Plus(days: e.multidays - 1)
				s $= '"multidays":' $ String(e.multidays) $ ',' $
					'"span":' $ String(e.span) $ ',' $
					'"end_date":"' $ end_date.Format('yyyy-M-d') $ '",' $
					'"completed":' $ String(e.completed) $ ',' $
					'"end":' $ String(e.end) $ ','
				}

			title = .handleInvalidChar(e.title)
			s $= '"title":"' $ title $ '"'
			if e.Member?('id')
				s $= ',"id":"' $ .handleInvalidChar(e.id) $ '"'
			if e.Member?('desc')
				s $= ',"desc":"' $ .handleInvalidChar(e.desc) $ '"'
			s $= '},'
			}
		return s.RemoveSuffix(',') $ ']})'
		}
	handleInvalidChar(desc)
		{
		// special char "  \
		if desc.Find('\\') < desc.Size()
			desc = desc.Replace('\\', '\\\\\\')
		if desc.Find('"') < desc.Size()
			desc = desc.Replace('"', '\\\\"')
		desc = desc.Replace('\r\n', ' ').Replace('\n', ' ')
		return desc
		}
	copy_event(e)
		{
		evt = e.Project(#(type, subtype, start_date, date, title, multidays,
			completed, span, end))
		if e.Member?('desc')
			evt.desc = e.desc
		if e.Member?('id')
			evt.id = e.id
		return evt
		}
	}