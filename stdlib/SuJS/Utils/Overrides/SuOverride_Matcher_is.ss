// Copyright (C) 2024 Axon Development Corporation All rights reserved worldwide.
MatcherWas
	{
	Match(value, args)
		{
		if args is value
			return true
		return false
		}
	Expected(args)
		{
		return .DisplayValue(args)
		}
	}
