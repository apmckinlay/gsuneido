// Copyright (C) 2012 Suneido Software Corp. All rights reserved worldwide.
class
	{
	Match(value, args)
		{
		if args is true
			args = ''
		return Type(value) is 'Except' and value.Has?(args)
		}
	Expected(args)
		{
		return "an exception matching " $ Display(args)
		}
	Actual(value, args/*unused*/)
		{
		return Type(value) is 'Except'
			? "was " $ Display(value)
			: "did not throw an exception"
		}
	}