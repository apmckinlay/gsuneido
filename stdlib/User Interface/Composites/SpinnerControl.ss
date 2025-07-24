// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// Contributions by Luis Alfredo Barrasa
PassthruController
	{
	Name: 'Spinner'
	New(rangefrom = 0, rangeto = 99999, width = false, set = false,
		mask = "##,###", justify = 'RIGHT', status = "", mandatory = false,
		euro = false, increase = 1, rollover = false)
		{
		super(Object(SpinnerWndProc rangefrom, rangeto, width, mask,
			justify, status, mandatory, euro, increase, rollover))
		.Send('Data')
		.mandatory = mandatory
		if (set isnt false)
			{
			.Set(set)
			.Send("NewValue", .Get())
			}
		.Top = .GetChild().Top
		}
	// block Data message from NumberControl
	Data()
		{
		}
	// resend NewValue so this becomes responsible for Get
	NewValue(val, source)
		{
		// have to set SpinnerWndProc's value since it does not get
		// NewValue from the NumberControl
		if (source is .SpinnerWndProc.GetNumberControl())
			.SpinnerWndProc.SetValue(val)
		.Send('NewValue', val)
		}
	Get()
		{
		return .SpinnerWndProc.Get()
		}
	Set(val)
		{
		.SpinnerWndProc.Set(val)
		}
	EditHwnd()
		{
		return .SpinnerWndProc.EditHwnd()
		}
	Valid?()
		{
		return .SpinnerWndProc.GetNumberControl().Valid?()
		}
	OverrideRanges(low, high, setLow = false)
		{
		.SpinnerWndProc.OverrideRanges(low, high, setLow)
		}
	Destroy()
		{
		.Send("NoData")
		super.Destroy()
		}
	}
