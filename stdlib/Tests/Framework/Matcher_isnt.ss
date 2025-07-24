// Copyright (C) 2011 Suneido Software Corp. All rights reserved worldwide.
class
	{
	Match(value, args)
		{
		return args isnt value
		}
	Expected(args)
		{
		return "the value to not be " $ Display(args)
		}
	Actual(value/*unused*/, args/*unused*/)
		{
		return "was"
		}
	}