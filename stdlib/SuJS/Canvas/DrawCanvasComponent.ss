// Copyright (C) 2022 Axon Development Corporation All rights reserved worldwide.
CanvasComponent
	{
	ContextMenu: true
	New()
		{
		.El.AddEventListener('mousedown', .mousedown)
		.El.AddEventListener('keydown', .keydown)
		.El.AddEventListener('dblclick', .dblclick)
		}

	tracker: false
	dragging: false
	resizing: false
	SetTracker(tracker)
		{
		.tracker = Construct(Object(tracker, canvas: this))
		}

	OnContextMenu(event)
		{
		if .GetReadOnly() is true or .tracker is false or .dragging is true
			return
		_event = event
		pos = .getPos(event)
		.tracker.MouseDown(pos.x, pos.y)
		.tracker.MouseUp(pos.x, pos.y)
		super.OnContextMenu(event)
		}

	keydown(event)
		{
		if .GetReadOnly() is true
			return
		if .arrowKeys(event)
			return
		if event.key is 'Delete'
			.DeleteSelected()
		else if event.ctrlKey is true
			.ctrlKeys(event)
		}

	arrowKeys(event)
		{
		x = y = 0
		switch event.key
			{
			case 'ArrowUp': y = -1
			case 'ArrowDown': y = 1
			case 'ArrowLeft': x = -1
			case 'ArrowRight': x = 1
			default: return false
			}
		.MoveSelected(x, y)
		return true
		}

	ctrlKeys(event)
		{
		// CTRL + A, X, V do not work inside a Book unless we handle them here
		//  Ctrl+C / On_Copy will be redirected to CanvasControl
		// so we don't handle it here
		switch event.key
			{
			case 'a' :
				.Event('Send', #On_Select_All)
			case 'x':
				.Event('Send', #On_Cut)
			case 'v':
				.Event('Send', #On_Paste)
			default :
			}
		}

	mousemoved?: false
	mousedown(event)
		{
		if .GetReadOnly() is true or .tracker is false or event.button isnt 0
			return
		_event = event
		.dragging = true
		pos = .getPos(event)
		if .tracker.Base?(SuDrawRectTracker)
			{
			for item in .GetAllItems()
				if item.IsHandle?(pos.x, pos.y)
					{
					.resizing = item
					.prevtracker = .tracker
					.origx = pos.x
					.origy = pos.y
					.tracker.ResizeDown(item, pos.x, pos.y)
					.StartMouseTracking(.mouseup, .mousemove)
					return 0
					}
			}
		.mousemoved? = false
		.tracker.MouseDown(@pos)
		.StartMouseTracking(.mouseup, .mousemove)
		}

	mousemove(event)
		{
		_event = event
		pos = .getPos(event)
		.mousemoved? = true
		if .resizing isnt false
			.tracker.ResizeMove(.resizing, pos.x, pos.y)
		else
			.tracker.MouseMove(@pos)
		}

	mouseup(event)
		{
		_event = event
		.dragging = false
		pos = .getPos(event)
		.selectPoint(pos)
		if .resizing isnt false
			{
			.tracker.ResizeUp(.resizing, pos.x, pos.y)
			.resizing.Resize(.origx, .origy, pos.x, pos.y)
			.tracker = .prevtracker
			.resizing = false
			.StopMouseTracking()
			return
			}
		if false isnt args = .tracker.MouseUp(@pos)
			.Event('TrackerMouseUp', args)
		.StopMouseTracking()
		}

	selectPoint(pos)
		{
		if .mousemoved? is false and .MaybeClearSelect() and
			false isnt i = .ItemAtPoint(pos.x, pos.y)
			.Select(i)
		}

	getPos(event)
		{
		rect = SuRender.GetClientRect(.El)
		return Object(x: event.clientX - rect.left, y: event.clientY - rect.top)
		}

	dblclick(event)
		{
		if .GetReadOnly() is true or event.button isnt 0
			return
		.Event('LBUTTONDBLCLK')
		}
	}
