// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
HtmlDivComponent
	{
	program:	false
	Name:		'BookBrowser'
	Open(ctrl)
		{
		.GetChild().SetVisible(false)
		.program = .Construct(ctrl)
		DoStartup(.program)
		.SetupCtrl(.program)
		.WindowRefresh()
		}

	Close()
		{
		.program.Destroy()
		.program = false
		.GetChild().SetVisible(true)
		.WindowRefresh()
		}

	Recalc()
		{
		super.Recalc(.program)
		}

	GetChildren()
		{
		// Have to return the superclass controller's children as well, or you
		// will leak when Destroy() is called -- because the superclass
		// Destroy() destroys only the children returned by .GetChildren()
		controllerChildren = super.GetChildren()
		return .program is false
			? controllerChildren
			: Object(.program).Union(controllerChildren)
		}
	}
