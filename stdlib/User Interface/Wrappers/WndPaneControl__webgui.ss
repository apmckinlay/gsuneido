// Copyright (C) 2020 Axon Development Corporation All rights reserved worldwide.
Control
	{
	Name: 'WndPane'
	ComponentName: 'WndPane'
	New(control, windowClass = "SuWhiteArrow")
		{
		.ctrl = .Construct(control)
		.ComponentArgs = Object(.ctrl.GetLayout(), .getBgColor(windowClass))
		}

	getBgColor(windowClass)
		{
		if windowClass.Has?('White')
			return 'white'
		return ToCssColor(CLR.ButtonFace)
		}

	SetEnabled(enabled)
		{
		Assert(Boolean?(enabled))
		.ctrl.SetEnabled(enabled)
		}
	GetEnabled()
		{ return .ctrl.GetEnabled() }
	SetReadOnly(readOnly)
		{
		Assert(Boolean?(readOnly))
		.ctrl.SetReadOnly(readOnly)
		}
	GetReadOnly()
		{ return .ctrl.GetReadOnly() }
	SetVisible(visible)
		{
		Assert(Boolean?(visible))
		.ctrl.SetVisible(visible)
		super.SetVisible(visible)
		}
	SetFocus()
		{
		.ctrl.SetFocus()
		}
	Update()
		{
		.ctrl.Update()
		super.Update()
		}
	GetControl()
		{
		return .ctrl
		}
	GetChildren()
		{
		return Object(.ctrl)
		}
	ContextMenu(x, y)
		{
		.Send("WndPane_ContextMenu", x, y)
		return 0
		}
	HasFocus?()
		{
		return .ctrl.HasFocus?()
		}
	Destroy()
		{
		.ctrl.Destroy()
		super.Destroy()
		}
	}