// Copyright (C) 2020 Axon Development Corporation All rights reserved worldwide.
Component
	{
	New()
		{
		.ResetAccels()
		}

	New2()
		{
		}

	GetChildren()
		{
		return [.Ctrl]
		}

	ResetAccels()
		{
		.accels = Object().Set_default(Object())
		}
	SetAccels(accels)
		{
		.accels = Object().Set_default(Object())
		for ac in accels
			{
			.accels[ac.idx][ac.key] = ac.cmd
			}
		}
	GetAccelCmd(ctrl, alt, shift, key)
		{
		idx = (ctrl ? 0x100 : 0) | (alt ? 0x10 : 0) | (shift ? 0x1 : 0)
		return .accels[idx].GetDefault(key.Lower(), false)
		}

	refreshing?: false // combine multiple Refresh
	Refresh()
		{
		if .refreshing?
			return
		.refreshing? = true
		// needed to combine multiple screen refreshes
		// into one to reduce layout reflow
		.refreshTimer = SuDelayed(0, .refresh)
		}
	refreshTimer: false
	refresh()
		{
		.refreshing? = false
		if .Destroyed?()
			return
		.BottomUp(#Recalc)
		}
	// run the delayed refresh before ShowWindow/UpdateWindow to avoid flickering
	RunPendingRefresh()
		{
		if .refreshing?
			{
			.refresh()
			.refreshTimer.Kill()
			return true
			}
		return false
		}

	SetWinSize(w, h)
		{
		if .Ctrl.Xstretch > 0 or .Ctrl.Ystretch > 0
			{
			viewport = SuRender().GetClientRect()
			w = Max(Min(w, viewport.width - 2), .Ctrl.Xmin)
			h = Max(Min(h, viewport.height - 2), .Ctrl.Ymin)
			.SetStyles([ width: w $ 'px', height: h $ 'px' ])
			.TopDownWindowResize()
			}
		}

	TopDownWindowResize()
		{
		if .Ctrl.Method?(#TopDown)
			.Ctrl.TopDown(#WindowResize)
		}

	Center()
		{
		viewportRect = SuRender.GetClientRect()
		midXView = (viewportRect.left + viewportRect.right) / 2
		midYView = (viewportRect.top + viewportRect.bottom) / 2
		containerEl = .GetContainerEl()
		rect = SuRender.GetClientRect(.El)
		.SetStyles([
			top: Max(midYView - rect.height / 2, viewportRect.top) $ 'px',
			left: Max(midXView - rect.width / 2, viewportRect.left) $ 'px'], containerEl)
		}

	GetContainerEl()
		{
		return .El
		}

	RegisterActiveWindow()
		{
		.GetContainerEl().AddEventListener('mousedown', .DoActivate, useCapture:)
		}
	DoActivate(event, target = false)
		{
		if target is false
			target = event.target

		// Need to check this because ListEditWindow is a child of its parent Window
		// In this case, both windows will receive the mousedown event
		if .windowContains?(target)
			{
			if Same?(SuRender().ActiveWindow, this)
				{
				.Event(#DoActivate)
				SuRender().Notification.OnWindowActivated(.UniqueId)
				}
			else
				// need this to block further user actions during ListEditWindow commit
				.EventWithOverlay(#DoActivate)
			}
		}
	windowContains?(target)
		{
		forever
			{
			match = 0
			try
				match = Same?(this, target.Window())
			if match isnt 0
				return match
			newtarget = false
			try
				newtarget = target.parentElement
			if newtarget is false
				return false
			target = newtarget
			}
		}
	UpdateOrder(order, active?)
		{
		.GetContainerEl().SetStyle('z-index', order)
		if active? is true
			{
			SuRender().ActiveWindow = this
			SuRender().Notification.OnWindowActivated(.UniqueId)
			.PlaceActive()
			}
		else if Same?(SuRender().ActiveWindow, this)
			SuRender().ActiveWindow = false
		}

	PlaceActive()
		{
		}

	HighlightDefaultButton(unused)
		{
		}
	}
