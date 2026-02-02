// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
Control
	{
	ComponentName: 'Fill'
	New(min = 0, fill = 1)
		{
		.ComponentArgs = Object(min, fill)
		}
	GetReadOnly() // read-only not applicable to fill
		{
		return true
		}
	}