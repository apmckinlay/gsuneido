// Copyright (C) 2022 Axon Development Corporation All rights reserved worldwide.
CanvasControl
	{
	Name: Canvas
	Title: Canvas
	Xstretch: 1
	Ystretch: 1
	ComponentName: 'DrawCanvas'
	ComponentArgs: #()

	ContextMenu(x, y)
		{
		menu = DrawControl.Menu[0].Copy().Map({ it.Tr('&') })
		ContextMenu(menu).ShowCall(this, x, y)
		}

	tracker: false
	SetTracker(tracker, item)
		{
		if .tracker isnt false
			.tracker.Release()
		.tracker = tracker(.Hwnd, item, canvas: this)
		.Act('SetTracker', .tracker.GetSuJSTracker())
		}

	TrackerMouseUp(args)
		{
		if false isnt item = .mouseUp(args)
			{
			.ClearSelect()
			if not Object?(item)
				item = Object(item)
			item.Each({
				.AddItemAndSelect(it.SetColor(.GetColor()).SetLineColor(.GetLineColor()))
				})
			.Send('Canvas_LButtonUp')
			}
		.Send('CanvasChanged')
		}

	mouseUp(args)
		{
		_report = new HtmlDriver
		return (.tracker[args[0]])(@+1args)
		}

	LBUTTONDBLCLK()
		{
		.Send('WhichDrawCanvasClicked', .Name)
		.edit()
		}

	On_Context_Edit()
		{
		.edit()
		}

	edit()
		{
		selects = .GetSelected()
		if selects.Size() isnt 1
			return
		item = selects[0]
		item.Edit(canvas: this)
		.SyncItem(item)
		.Send('CanvasChanged')
		}

	Default(@args)
		{
		method = args[0]
		if method.Prefix?('On_Context_') and args.Member?('item')
			{
			.Send('On_' $ ToIdentifier(args.item))
			}
		return 0
		}
	}
