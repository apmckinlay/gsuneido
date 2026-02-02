// Copyright (C) 2020 Axon Development Corporation All rights reserved worldwide.
PassthruController
	{
	Name: 'SpinnerWndProc'

	New(.rangefrom, .rangeto, width, mask, justify, status, mandatory,
		euro, .increase, .rollover)
		{
		super(.layout(euro, width, mask, status, justify, mandatory))
		.updown = .FindControl('UpDown')
		.number = .FindControl(.Name)
		.val = rangefrom
		.Send("Data")
		if Number?(.rangefrom) and Number?(.rangeto)
			{
			.number.ToolTip('Range: ' $ .rangefrom $ ' to ' $ .rangeto)
			.updown.ToolTip('Range: ' $ .rangefrom $ ' to ' $ .rangeto)
			}
		.number.Act('SetVScroll', increase, rollover, rangefrom, rangeto)
		}

	layout(euro, width, mask, status, justify, mandatory)
		{
		return Object('Horz'
			Object((euro ? 'EuroNumber' : 'Number'), name: .Name, rangefrom: .rangefrom,
				rangeto: .rangeto, :width, :mask, :status, :justify, :mandatory),
			'UpDown')
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
		.updown.SetEnabled(readonly isnt true)
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

	Destroy()
		{
		.Send("NoData")
		super.Destroy()
		}
	}