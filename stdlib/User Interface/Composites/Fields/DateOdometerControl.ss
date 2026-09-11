// Copyright (C) 2026 Axon Development Corporation. All rights reserved worldwide.
Controller
	{
	Name:          "Date Number"
	ComponentName: #DateOdometer
	New(.readonly = false, .mandatory = false)
		{
		.Left = .Vert.Pair.Left
		.Top = .Vert.Pair.Top
		.date = .Vert.Pair.Date
		.odometer = .Vert.Pair2.Odometer
		.hours = .Vert.Pair3.Hours
		.Send(#Data)
		}

	Controls: (Vert,
		(Pair, (Static, Date), (ChooseDate, name: Date)),
		(Pair, (Static, Hours), (Number, mask: "###,###,###", name: Hours), name: Pair3),
		(Pair, (Static, Odometer), (TimeSpan, odometer:, name: Odometer), name: Pair2))

	Data(source/*unused*/) { }

	Valid?()
		{
		if not .odometer.Valid?()
			return false

		if .mandatory is true and .date.Get() is "" and .odometer.Get() is "" and
			.hours.Get() is ""
			return false

		return true
		}

	NewValue(value/*unused*/, source)
		{
		if source is .odometer
			.clear_values([.date, .hours])
		else if source is .hours
			.clear_values([.date, .odometer])
		else
			.clear_values([.hours, .odometer])

		.Send(#NewValue, .Get())
		}

	clear_values(ctrls)
		{
		for ctrl in ctrls
			ctrl.Set("")
		}

	Set(x)
		{
		if Date?(x)
			{
			.date.Set(Date(x))
			.clear_values([.hours, .odometer])
			}
		else if Number?(x)
			{
			.hours.Set(Number(x))
			.clear_values([.date, .odometer])
			}
		else
			{
			.odometer.Set(x is "" ? "0 " $ OptContribution("MileageUOM", { #kms })() : x)
			.clear_values([.date, .hours])
			}
		}

	Get()
		{
		date = .date.Get()
		odometer = .odometer.Get()
		hours = .hours.Get()
		if Date?(date)
			return date
		else if Number?(hours)
			return hours
		else
			return odometer
		}

	Dirty?(dirty = "")
		{
		return .date.Dirty?(dirty) or .odometer.Dirty?(dirty)
		}

	readonly: false
	SetReadOnly(on = true)
		{
		.date.SetReadOnly(on)
		.odometer.SetReadOnly(on)
		.hours.SetReadOnly(on)
		}

	SetFocus()
		{
		.date.SetFocus()
		}

	HandleTab()
		{
		if GetFocus() is .date.Field.Hwnd
			{
			SetFocus(.hours.Hwnd)
			return true
			}
		else if GetFocus() is .hours.Hwnd
			{
			SetFocus(.odometer.Get_Value_Ctrl().Hwnd)
			return true
			}
		return .odometer.HandleTab()
		}

	Destroy()
		{
		.Send(#NoData)
		super.Destroy()
		}
	}
