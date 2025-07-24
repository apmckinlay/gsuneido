// Copyright (C) 2020 Axon Development Corporation All rights reserved worldwide.
HtmlDivComponent
	{
	New(@args)
		{
		super(@args)
		.Left = .FindControl(#Form).Left
		}
	}
