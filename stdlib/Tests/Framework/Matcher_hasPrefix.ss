// Copyright (C) 2011 Suneido Software Corp. All rights reserved worldwide.
class
	{
	Match(value, args)
		{
		return String?(value) and value.Prefix?(args)
		}
	Expected(args)
		{
		return "a string starting with " $ Display(args)
		}
	Actual(value, args)
		{
		return String?(value) or Number?(value)
			? "started with " $ Display(value[.. args.Size()])
			: "was a " $ Type(value)
		}
	}