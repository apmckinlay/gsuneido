// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Memoize
	{
	Func(eventName)
		{
		return not QueryEmpty?('event_condition_actions', eca_event: eventName)
		}
	}