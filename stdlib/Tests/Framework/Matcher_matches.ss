// Copyright (C) 2011 Suneido Software Corp. All rights reserved worldwide.
class
	{
	Match(value, args)
		{
		return String?(value) and value =~ args
		}
	Expected(args)
		{
		return "a string that matches: " $ Display(args) $ '\n'
		}
	Actual(value, args/*unused*/)
		{
		return String?(value) or Boolean?(value)
			? "was: " $ Display(value)
			: "was a " $ Type(value)
		}
	}