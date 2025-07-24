// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
ScintillaAddon
	{
	Setting: 		''
	Set: 			false
	DefaultValue: 	false
	New(parent, options)
		{
		super(parent, options)
		.SyncPreferences()
		}

	SyncPreferences(init = false)
		{
		prevSet = .Set
		.Set = IDESettings.Get(.Setting, .DefaultValue)
		if init and prevSet isnt .Set
			.Init()
		}
	}
