// Copyright (C) 2022 Axon Development Corporation All rights reserved worldwide.
Control
	{
	Name: 'ColorRect'
	ComponentName: 'ColorRect'
	Xmin: 50
	Ymin: 50
	New(color = 0, choose? = false)
		{
		if Object?(color)
		   color = RGB(color[0], color[1], color[2])
		.color = color
		.Send('Data')
		.ComponentArgs = Object(ToCssColor(color), choose?)
		}
	LBUTTONDBLCLK()
		{
		.Send("ColorRect_DoubleClick", .color)
		return 0
		}
	Get()
		{
		return .color
		}
	Set(.color)
		{
		.Act('Set', ToCssColor(color))
		}
	UpdateColor(.color) {}
	OK()
		{
		return .Get()
		}
	Dirty?(unused)
		{
		return false
		}
	Destroy()
		{
		.Send('NoData')
		super.Destroy()
		}
	}