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
			// temporarily disable unsupported Scintilla addons
			if addon.Base?(ScintillaAddon) and
				not ScintillaAddonsControl.SupportedAddon?(m)
				continue
			.construct(m, ob[m])
			}
		}
	}
