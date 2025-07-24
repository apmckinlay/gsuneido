// Copyright (C) 2015 Suneido Software Corp. All rights reserved worldwide.
class
	{
	Match(value, args)
		{
		return .sizable?(value) and value.Size() is args
		}
	Actual(value, args/*unused*/)
		{
		if .sizable?(value)
			return "was " $ Display(value.Size())
		return "was a " $ Type(value)
		}
	sizable?(value)
		{
		return String?(value) or Object?(value)
		}
	Expected(args)
		{
		return "a size of " $ Display(args)
		}
	}