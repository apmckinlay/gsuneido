// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Controls:
		(Record
			(Vert
				boolean
				number
				string
				Skip
				(Horz (Button One) (Button Two) (Button Clear))
				)
			)
	On_One()
		{ .Data.Set(Record(boolean: false, number: 123, string: "hello")); }
	On_Two()
		{ .Data.Set(Record(boolean: true, number: 456, string: "world")); }
	On_Clear()
		{ .Data.Set(Record()); }
	}