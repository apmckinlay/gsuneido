// Copyright (C) 2019 Axon Development Corporation All rights reserved worldwide.
Group
	{
	Dir: "horz"
	Name: "Horz"
	ComponentName: "Horz"
	New(@controls)
		{
		super(controls)
		.ComponentArgs.leftAlign = controls.GetDefault('leftAlign', false)
		}
	}
