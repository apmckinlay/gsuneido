// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
// Use this control to be a placeholder for a field value that is not visible
// for example: use it on Params to retain field data since Params only saves control data
Control
	{
	Name: "DummyField"
	ComponentName: "DummyField"
	ComponentArgs: #()
	Hwnd: 0
	New()
		{
		.Send("Data")
		}

	Set(.value)
		{
		}

	value: ""
	Get()
		{
		return .value
		}

	Destroy()
		{
		.Send("NoData")
		super.Destroy()
		}
	}
