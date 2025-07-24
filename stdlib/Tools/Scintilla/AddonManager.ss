// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
// TODO catch exceptions
// MAYBE let addons subscribe to events so all events aren't sent to all addons
// TODO if value is false then ignore addon (e.g. to override base class)
class
	{
	New(@args)
		{
		.parent = args[0]
		.addons = Object()
		for ob in args
			.getAddonsFrom(ob)
		}
	Addon(@args)
		{
		.getAddonsFrom(args)
		}
	getAddonsFrom(ob)
		{
		for m in ob.Members(all:)
			if m.Prefix?('Addon_')
				.construct(m, ob[m])
		}
	construct(name, options)
		{
		c = Global(name)
		Assert(c.For.Has?(.parent.Base().Name))
		addon = c(.parent, options) // make an instance
		.addons.Add(addon)
		}
	Send(@args)
		{
		addonFound = false
		for addon in .addons
			if addon.Method?(args[0])
				{
				addon[args[0]](@+1args)
				addonFound = true
				}
		return addonFound
		}
	ConditionalSend(block, args)
		{
		for addon in .addons
			if Type(block) is #Block and addon.Method?(args[0]) and block(addon)
				addon[args[0]](@+1args)
		}
	Collect(@args)
		{
		results = []
		for addon in .addons
			if addon.Method?(args[0])
				results.Add(addon[args[0]](@+1args))
		return results
		}
	SendToOneAddon(@args)
		{
		found = false
		targetAddon = false
		for addon in .addons
			if addon.Method?(args[0])
				{
				Assert(not found,
					msg: 'AddonManager: multiple definitions for ' $ args[0])
				targetAddon = addon
				found = true
				}
		if targetAddon is false
			return 0
		targetAddon[args[0]](@+1args)
		}
	}