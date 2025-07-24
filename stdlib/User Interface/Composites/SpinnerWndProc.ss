// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// Contributions by Luis Alfredo Barrasa
WndProc
	{
	Name: 'SpinnerWndProc'

	New(.rangefrom, .rangeto, width, mask, justify, status, mandatory,
		euro, .increase, .rollover)
		{
		.CreateWindow("SuWhiteArrow", "", WS.VISIBLE | WS.TABSTOP, w: 1, h: 1)
		.number = .Construct((euro ? 'EuroNumber' : 'Number'),
			name: .Name,
			:rangefrom, :rangeto, :width, :mask, :status, :justify, :mandatory)
		.Ymin = .number.Ymin
		.updown = .Construct(Object('UpDown', UDS.ALIGNRIGHT))
		.Xmin = .number.Xmin + .updown.Xmin
		.Top = .number.Top
		.SubClass()
		.val = rangefrom
		.Send("Data")
		if Number?(.rangefrom) and Number?(.rangeto)
			{
			.number.ToolTip('Range: ' $ .rangefrom $ ' to ' $ .rangeto)
			.updown.ToolTip('Range: ' $ .rangefrom $ ' to ' $ .rangeto)
			}
		}
	SETFOCUS()
		{
		SetFocus(.number.Hwnd)
		return 0
		}
	last_pos: 0
	VSCROLL(wParam)
		{
		if not .GetEnabled()
			return 0
		.val = .number.Get()
		if LOWORD(wParam) isnt SB.THUMBPOSITION
			return 0
		newpos = HIWORD(wParam)
		.setValFromNewScrollPos(newpos)
		.handleBoundaries()
		.number.Set(.val)
		.number.Dirty?(true)
		.last_pos = newpos
		.Send('NewValue', .val)
		return 0
		}

	setValFromNewScrollPos(newpos)
		{
		if newpos < .last_pos
			.val += .increase
		if newpos > .last_pos
			.val -= .increase
		if newpos is 0 and .last_pos is 0
			.val += .increase
		if newpos is 100 and .last_pos is 100 /* = max scroll pos*/
			.val -= .increase
		}

	handleBoundaries()
		{
		if .rollover
			{
			if .val > .rangeto
				.val = .rangefrom
			if .val < .rangefrom
				.val = .rangeto
			}
		else
			{
			if .val > .rangeto
				.val = .rangeto
			if .val < .rangefrom
				.val = .rangefrom
			}
		}

	Get()
		{
		return .val
		}
	Set(val)
		{
		.val = val
		.number.Set(val)
		}
	SetValue(val)
		{
		.val = val
		}
	EditHwnd()
		{
		return .number.Hwnd
		}
	GetNumberControl()
		{
		return .number
		}

	SetEnabled(enabled)
		{
		.number.SetEnabled(enabled)
		.updown.SetEnabled(enabled)
		}
	GetEnabled()
		{
		return .number.GetEnabled()
		}

	SetReadOnly(readonly)
		{
		.number.SetReadOnly(readonly)
		.updown.SetReadOnly(readonly)
		}
	GetReadOnly()
		{
		return .number.GetReadOnly()
		}

	OverrideRanges(low, high, setLow = false)
		{
		.rangefrom = low
		.rangeto = high
		.number.SetRange(low, high)
		if .val < low or setLow
			.Set(low)
		else if .val > high
			.Set(high)
		.number.ToolTip('Range: ' $ .rangefrom $ ' to ' $ .rangeto)
		.updown.ToolTip('Range: ' $ .rangefrom $ ' to ' $ .rangeto)
		}

	Resize(x, y, w, h)
		{
		super.Resize(x, y, w, h)
		.number.Resize(0, 0, w - .updown.Xmin, h)
		.updown.Resize(w - .updown.Xmin, 0, .updown.Xmin, h)
		}
	Destroy()
		{
		.Send("NoData")
		.number.Destroy()
		.updown.Destroy()
		super.Destroy()
		}
	}
