// Copyright (C) 2019 Axon Development Corporation All rights reserved worldwide.
_TabsControl
	{
	x: 0
	y: 0
	w: 0
	h: 0
	ComponentName: 'Tabs'
	initialConstruct(args, startTab)
		{
		super.TabsControl_initialConstruct(args, startTab)
		.ComponentArgs = Object(.Tab.GetLayout(), vertical: .vertical,
			xstretch: 1, ystretch: 1)
		}

	construct(i)
		{
		.ActWith()
			{
			.being_constructed = i
			.ctrls[i] = .Construct(Object("WndPane"
				Object("Border", .controls[i], .border, borderline: 1,
					xstretch: 1, ystretch: 1),
				windowClass: 'SuBtnfaceArrow'))
			Object('Insert', .alternativePos is true ? 0 : 1, .ctrls[i].GetLayout())
			}
		.being_constructed = false
		return .ctrls[i]
		}

	SelectTab(i, source = false, keepFocus = false)
		{
		if not .allowSelectTab?(source, i)
			return
		.Send('SelectTab', i)
		if not keepFocus
			SetFocus(0)
		.deselectCtrl()
		if (.ctrls[i] is false)
			.ConstructTab(i)
		.ctrls[i].SetVisible(true)
		.ctrl = .ctrls[i]
		.Send('TabsControl_SelectTab')
		if not keepFocus
			.FocusFirst(.ctrl.Hwnd)
		.Act(#SelectTab, .ctrl.UniqueId)
		.Window.Refresh()
		}

	ConstructTab(i)
		// WARNING: this does NOT fill in RecordControl data (use Select)
		{
		if true is res = super.ConstructTab(i)
			.ctrls[i].SetVisible(false)
		return res
		}

	destroyTab(i)
		{
		if i is .GetSelected()
			.ctrl = false
		if false isnt c = .ctrls.Extract(i, false)
			{
			id = c.UniqueId
			c.Destroy()
			.Act('RemoveTab', id)
			}
		}
	}