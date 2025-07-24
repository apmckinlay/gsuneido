// Copyright (C) 2011 Suneido Software Corp. All rights reserved worldwide.
class
	{
	Match(value, args)
		{
		return String?(value) and value.Suffix?(args)
		}
	Expected(args)
		{
		return "a string ending with " $ Display(args)
		}
	Actual(value, args)
		{
		return String?(value) or Number?(value)
			? "ended with " $ Display(value[-args.Size() ..])
			: "was a " $ Type(value)
		}
	}