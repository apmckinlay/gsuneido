// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// NOTE: unlike Pane, contents has this hwnd as its parent
WndProc
	{
	Name:		WndPane

	New(control, windowClass = "SuWhiteArrow", hidden = false)
		{
		style = 0
		if hidden is false
			style |= WS.VISIBLE
		.CreateWindow(windowClass, "", style, exStyle: WS_EX.CONTROLPARENT)
		.SubClass() // FIXME: Why is this subclassed?
					// TODO: Does this even need to be a WndProc (the only thing
					//       it handles is ContextMenu, and that gets reflected
					//       by the parent WndProc)
					// My guess is the only reason this is a WndProc at all is
					// that before the following change:
					//     2007.05.25 12:29 apm - use standard ContextMenu method
					// It was directly handling WM_CONTEXTMENU.
		.ctrl = .Construct(control)
		.recalc()
		}
	recalc()
		{
		.Xmin = .ctrl.Xmin
		.Ymin = .ctrl.Ymin
		.Top = .ctrl.Top
		.Left = .ctrl.Left
		.Xstretch = .ctrl.Xstretch
		.Ystretch = .ctrl.Ystretch
		}
	Recalc()
		{
		.recalc()
		}
	Resize(x, y, w, h)
		{
		super.Resize(x, y, w, h)
		.ctrl.Resize(0, 0, w, h)
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
	Destroy()
		{
		.ctrl.Destroy()
		super.Destroy()
		}
	}
