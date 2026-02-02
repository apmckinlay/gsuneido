// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
HorzControl
	{
	Name: 'HorzEqual'
	ComponentName: "HorzEqual"
	New(@args)
		{
		super(@args)
		.ComponentArgs.pad = args.GetDefault('pad', 20/*=default pad*/)
		}
	}