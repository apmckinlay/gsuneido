// Copyright (C) 2012 Suneido Software Corp. All rights reserved worldwide.
Memoize
	{
	Func(name)
		{
		ob = Object()
		for contrib in .contribs(name)
			.getContributions(ob, contrib, name)
		return ob
		}

	getContributions(ob, contrib, name) // recursive
		{
		if Function?(contrib)
			contrib = contrib()

		ob.Add(@contrib.Values(list:))
		for m in contrib.Members(named:)
			if Object?(contrib[m]) and Object?(ob.GetDefault(m, #()))
				{
				if not ob.Member?(m)
					ob[m] = Object()
				.getContributions(ob[m], contrib[m], name)
				}
			else if not ob.Member?(m)
				ob[m] = contrib[m]
			else
				throw "GetContributions " $ name $ " duplicate member " $ m
		}

	contribs(name)
		{
		return Contributions.Func(name) // not cached since we're caching
		}

	Init()
		{
		LibUnload.AddObserver(#GetContributions, {|name|
			if IsContribution?(name)
				.ResetCache()
			})
		}
	}