// Copyright (C) 2020 Axon Development Corporation All rights reserved worldwide.
HtmlDivComponent
	{
	Name: 'Coord'
	New(@args)
		{
		super(@args)
		.xValue = .FindControl('xValue')
		.yValue = .FindControl('yValue')
		.Left = .xValue.Left
		}

	HandleTab()
		{
		if .xValue.HasFocus?()
			{
			.yValue.SetFocus()
			return true
			}
		return false
		}
	}
