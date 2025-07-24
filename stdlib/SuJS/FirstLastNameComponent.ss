// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
HtmlDivComponent
	{
	New(@args)
		{
		super(@args)
		.first = .Horz.first.CapitalizeWords
		.last = .Horz.last.CapitalizeWords
		.Left = .Horz.first.Left
		}

	HandleTab()
		{
		if .first.HasFocus?()
			{
			.last.SetFocus()
			return true
			}
		return false
		}
	}
