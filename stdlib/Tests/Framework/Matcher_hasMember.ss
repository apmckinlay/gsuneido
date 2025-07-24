// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
class
	{
	Match(value, args)
		{
		return Type(value) in ('Object', 'Record', 'Class', 'Instance') and
			value.Member?(args)
		}
	Expected(args)
		{
		return "the value to have member " $ Display(args)
		}
	Actual(value, args/*unused*/)
		{
		return Type(value) in ('Object', 'Record', 'Class', 'Instance')
			? "did not"
			: "was a " $ Type(value)
		}
	}