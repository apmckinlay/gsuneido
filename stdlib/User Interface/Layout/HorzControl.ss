// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
//  Copyright (C) 2000 Suneido Software Corp.
Group
	{
	Dir: "horz"
	Name: "Horz"
	New(@controls)
		{
		super(controls)
		if controls.GetDefault('leftAlign', false)
			.Left = .GetChildren()[0].Left
		}
	}