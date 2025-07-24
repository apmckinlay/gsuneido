// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
class
	{
	Match(value, args)
		{
		return Object?(value) and not value.Member?(args)
		}
	Expected(args)
		{
		return "the value to NOT have member " $ Display(args)
		}
	Actual(value, args/*unused*/)
		{
		return Object?(value) ? "did" : "was a " $ Type(value)
		}
	}