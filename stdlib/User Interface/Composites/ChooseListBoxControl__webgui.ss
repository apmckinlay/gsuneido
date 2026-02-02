// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
Control
	{
	ComponentName: 'ChooseListBox'
	New(list, select, .listSeparator, fieldHwnd)
		{
		.ComponentArgs = Object(list, select, listSeparator, fieldHwnd)
		}

	CHOOSE(item)
		{
		.Window.Result(item)
		}
	}