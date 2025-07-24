// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.
PassthruController
	{
	Name: "ShowHide"
	New(control, showFunction)
		{
		super(.layout(control, showFunction))
		.Top = .GetChildren()[0].Top
		.Left = .GetChildren()[0].Left
		}

	layout(control, showFunction)
		{
		if String?(showFunction)
			showFunction = Global(showFunction)
		show = Suneido.GetDefault("ShowAll", "")
		if show is ""
			show = showFunction()
		return show ? control : #(Skip 1)
		}
	}
