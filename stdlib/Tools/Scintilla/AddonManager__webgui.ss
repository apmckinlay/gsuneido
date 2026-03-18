// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
_AddonManager
	{
	getAddonsFrom(ob)
		{
		for m in ob.Members(all:)
			{
			if not m.Prefix?('Addon_')
				continue
			addon = Global(m)
			if addon.GetDefault(#SuJsSupport?, true) is false
				continue
			.construct(m, ob[m])
			}
		}
	}
