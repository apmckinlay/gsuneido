// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
class
	{
	Match(value, args)
		{
		return Type(value) in ('Object', 'Class', 'Instance') and value.Base?(args)
		}
	Expected(args)
		{
		return "a value for which " $ Display(args) $ " is a base class"
		}
	Actual(value, args/*unused*/)
		{
		return Type(value) in ('Object', 'Class', 'Instance')
			? "did not (base was " $ Display(value.Base()) $ ")"
			: "was a " $ Type(value)
		}
	}
