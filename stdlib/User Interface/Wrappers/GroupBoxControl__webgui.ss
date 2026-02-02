// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
Container
	{
	Name: "GroupBox"
	ComponentName: "GroupBox"
	New(.text, control)
		{
		.ctrl = .Construct(control)
		.ComponentArgs = Object(.text, .ctrl.GetLayout())
		}

	GetChildren()
		{
		return Object(.ctrl)
		}

	SetCaption(caption)
		{
		.text = caption
		.Act(#SetCaption, caption)
		}

	GetCaption()
		{
		return .text
		}
	}