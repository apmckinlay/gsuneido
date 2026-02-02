// Copyright (C) 2019 Axon Development Corporation All rights reserved worldwide.
Container
	{
	Name: "Border"
	ComponentName: "Border"
	Ctrl: false
	New(control, border = 10, borderline = 0)
		{
		if Number?(control) and Object?(border)
			{ tmp = control; control = border; border = tmp }
		.Ctrl = .Construct(control)
		.ComponentArgs = Object(.Ctrl.GetLayout(), border, borderline)
		}

	SetCtrl(control)
		{
		if .Ctrl isnt false
			.Ctrl.Destroy()

		.ActWith()
			{
			.Ctrl = .Construct(control)
			Object('SetCtrl', .Ctrl.GetLayout())
			}
		}

	GetChild()
		{
		return .Ctrl
		}

	GetChildren()
		{
		return Object(.Ctrl)
		}
	}