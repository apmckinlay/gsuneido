// Copyright (C) 2026 Axon Development Corporation All rights reserved worldwide.
HtmlDivComponent
	{
	Name: "Date Number"
	New(@args)
		{
		super(@args)
		.Left = .Vert.Pair.Left
		.date = .Vert.Pair.Date
		.odometer = .Vert.Pair2.Odometer
		.hours = .Vert.Pair3.Hours
		}

	HandleTab()
		{
		if .date.GetChildren()[0].HasFocus?()
			{
			.hours.SetFocus()
			return true
			}
		else if .hours.HasFocus?()
			{
			.odometer.Horz.Value.SetFocus()
			return true
			}
		return .odometer.HandleTab()
		}
	}
