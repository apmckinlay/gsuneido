// Copyright (C) 2020 Axon Development Corporation All rights reserved worldwide.
Control
	{
	Name: 'VirtualListThumbBar'
	ComponentName: 'VirtualListThumbBar'
	New(.disableSelectFilter = false)
		{
		.ComponentArgs = Object(.disableSelectFilter)
		}

	OnHome()
		{
		.Send('On_VirtualListThumb_ArrowHome')
		}

	OnEnd()
		{
		.Send('On_VirtualListThumb_ArrowEnd')
		}

	OnSelect()
		{
		.Send('On_VirtualListThumb_ArrowSelect')
		}

	SetSelectPressed(pressed = false)
		{
		if .disableSelectFilter
			return
		.Act('SetSelectPressed', pressed)
		}

	Default(@unused) { }
	}
