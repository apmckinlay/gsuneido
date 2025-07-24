// Copyright (C) 2018 Axon Development Corporation All rights reserved worldwide.
GroupComponent
	{
	Dir: "horz"
	Name: "Horz"
	New(@elements)
		{
		super(elements)
		if elements.GetDefault('leftAlign', false)
			.Left = .GetChildren()[0].Left
		}
	}
