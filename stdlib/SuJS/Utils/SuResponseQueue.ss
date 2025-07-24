// Copyright (C) 2023 Axon Development Corporation All rights reserved worldwide.
Queue
	{
	Lookup(eventId)
		{
		while .Count() > 0 and .Front().eventId < eventId
			.Dequeue()
		if .Count() > 0 and .Front().eventId is eventId
			return .Front()
		return false
		}
	}