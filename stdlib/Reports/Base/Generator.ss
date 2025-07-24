// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Format
	{
	Generator?()
		{ return true }
	Header()
		{ return false } // no header
	PageFooter()
		{ return false }
	Next() // default implementation, can be overridden
		{
		if (.queue.Empty?())
			.More()
		if (.queue.Empty?())
			return false
		return .queue.PopFirst()
		}
	// handle queuing
	getter_queue()
		{ .queue = Object() } // create on first use
	More()
		{ }	// Generators can define More which calls Output
	Output(item)
		{ .queue.Add(item) }
	}
