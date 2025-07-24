// Copyright (C) 2020 Axon Development Corporation All rights reserved worldwide.
HtmlDivComponent
	{
	New(@args)
		{
		super(@args)
		.Left = .FindControl('ChooseList').Xmin - 1
		.typeField = .FindControl('ChooseList')
		}

	HandleTab()
		{
		if .typeField.HasFocus?()
			{
			.FindControl('Field').SetFocus()
			return true
			}
		return false
		}
	}
