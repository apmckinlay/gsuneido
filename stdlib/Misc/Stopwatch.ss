// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
class
	{
	New()
		{
		.t = .t2 = .date()
		}
	Call()
		{
		t2 = .t2
		.t2 = .date()
		if .t is t2
			return ReadableDuration.Between(.t, .t2)
		return "+ " $ ReadableDuration.Between(t2, .t2) $
			" = " $ ReadableDuration.Between(.t, .t2)
		}
	date() // overridden by test
		{
		return Date()
		}
	}