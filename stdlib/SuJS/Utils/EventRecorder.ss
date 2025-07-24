// Copyright (C) 2023 Axon Development Corporation All rights reserved worldwide.
class
	{
	n: 200
	i: 0
	New()
		{
		.events = Object()
		}

	Add(type, event)
		{
		.events[.i] = Object(date: Date(), :type, :event)
		.i = (.i + 1) % .n
		}

	Get()
		{
		return .events[.i..].Append(.events[...i])
		}

	Format(event)
		{
		return Display(event.date)[1..] $ '\t' $ event.type $ '\t' $
			(String?(event.event) ? event.event : Display(LogFormatEntry(event.event)))
		}
	}
