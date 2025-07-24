// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
HtmlDivComponent
	{
	New(@args)
		{
		super(@args)
		.value = .Horz.Value
		.uom = .Horz.Uom
		.Left = .Horz.Value.Left
		}

	HandleTab()
		{
		if .value.HasFocus?()
			{
			.uom.SetFocus()
			return true
			}
		return false
		}
	}
