// Copyright (C) 2011 Suneido Software Corp. All rights reserved worldwide.
MatcherWas
	{
	Match(value, args)
		{
		return (String?(value) or Object?(value)) and value.Has?(args)
		}
	Expected(args)
		{
		return "the value to be a string or object containing " $ Display(args)
		}
	}