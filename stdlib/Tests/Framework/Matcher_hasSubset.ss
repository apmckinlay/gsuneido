// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
class
	{
	Match(value, args)
		{
		for mem in args.Members()
			{
			if args[mem] is "" and (not value.Member?(mem) or value[mem] is "")
				continue
			if not value.Member?(mem) or value[mem] isnt args[mem]
				return false
			}
		return true
		}
	Expected(args)
		{
		return "the value to have the subset " $ Display(args)
		}
	Actual(value, args)
		{
		return "was " $
			Display(value.Copy().Set_default("missing").Project(args.Members()))
		}
	}
