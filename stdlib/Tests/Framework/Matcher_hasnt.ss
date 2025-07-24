// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
class
	{
	Match(value, args)
		{
		return not value.Has?(args)
		}
	Expected(args)
		{
		return "the value to NOT contain " $ Display(args)
		}
	Actual(value/*unused*/, args/*unused*/)
		{
		return "did"
		}
	}