// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
Container
	{
	Name: "Center"
	ComponentName: "Center"
	New(control, border = 0)
		{
		.ctrl = .Construct(control)
		.ComponentArgs = Object(.ctrl.GetLayout(), border)
		}
	GetChildren()
		{
		return Object(.ctrl)
		}
	}