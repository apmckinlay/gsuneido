// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
MatcherWas
	{
	Match(value, args)
		{
		return value.Members().Sort!() is args.Copy().Sort!()
		}
	Expected(args)
		{
		return "the object members to be " $ Display(args)
		}
	}