// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
HtmlContainer
	{
	Name: "Tabs"
	Xstretch: 1
	Ystretch: 1

	styles: `
		.su-tabs-container {
			position: relative;
			flex-grow: 1;
			align-self: stretch;
		}`
	New(tab, .vertical = false, alternativePos = false)
		{
		LoadCssStyles('su-tabs.css', .styles)
		.CreateElement('div')
		.SetStyles(Object(
			'display': 'inline-flex',
			'flex-direction': .vertical is true ? 'row' : 'column',
			'align-items': 'baseline'))

		.tab = .Construct(tab)
		.tab.SetStyles(#('flex-shrink': '0', 'align-self': 'stretch'))

		.tabs = Object()

		.container = CreateElement('div', .El, className: 'su-tabs-container',
			at: alternativePos is true ? 0 : 1)
		.Recalc()
		}

	ctrl: false
	SelectTab(id)
		{
		if false isnt i = .tabs.FindIf({ it.UniqueId is id })
			.ctrl = .tabs[i]
		else
			.ctrl = false
		}

	Recalc()
		{
		.Left = 0
		if .findSelectedCtrl() is false
			.Xmin = .Ymin = 0
		else
			{
			.ctrl.Recalc()
			.Xmin = Max(.ctrl.Xmin, .tab.Xmin)
			.Ymin = .ctrl.Ymin + .tab.Ymin
			}
		.SetMinSize()
		}

	Insert(control)
		{
		_at = [parent: this, parentEl: .container, at: false]
		.tabs.Add(el = .Construct(control))
		DoStartup(el)
		.WindowRefresh()
		}

	RemoveAll()
		{
		.tabs.Each(#Destroy)
		.tabs = Object()
		.WindowRefresh()
		}

	GetChildren()
		{
		return Object(.tab).Append(.tabs)
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
		if false isnt i = .tabs.FindIf({ it.UniqueId is id })
			{
			.tabs[i].Destroy()
			.tabs.Delete(i)
			.WindowRefresh()
			}
		}
	}
