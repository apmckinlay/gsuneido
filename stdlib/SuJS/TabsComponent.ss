// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
GroupComponent
	{
	Name: "Tabs"
	Dir: "vert"

	New(@elements)
		{
		super(.args(elements))
		}

	args(elements)
		{
		.Dir = elements.GetDefault(#vertical, false) ? 'horz' : 'vert'
		return elements
		}

	ctrl: false
	SelectTab(id)
		{
		for ctrl in .GetChildren()
			if ctrl.UniqueId is id
				{
				.ctrl = ctrl
				break
				}
		}
	Recalc()
		{
		super.Recalc()
		.Left = 0
		if .findSelectedCtrl() is false
			.Xmin = .Ymin = 0
		else
			{
			tab = .FindControl(#Tab)
			.Xmin = Max(.ctrl.Xmin, tab.Xmin)
			.Ymin = .ctrl.Ymin + tab.Ymin
			}
		.SetMinSize()
		}

	findSelectedCtrl()
		{
		if .ctrl isnt false
			return .ctrl
		for ctrl in .GetChildren()
			if ctrl.Name isnt 'Tab' and ctrl.GetVisible() is true
				return .ctrl = ctrl
		return false
		}

	RemoveTab(id)
		{
		if false isnt i = .GetChildren().FindIf({ it.UniqueId is id })
			.Remove(i)
		}
	}
