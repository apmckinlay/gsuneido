// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Name: 'ChooseColor'
	readonly: false
	New(color = 0, .readonly = false)
		{
		.Send('Data')
		.Top = .Horz.Choose.Top
		.Set(.toRGB(color))
		}
	Controls: (Horz (ColorRect xmin: 50 ymin: 22) (Button 'Choose...'))
	Set(value)
		{ .Horz.ColorRect.Set(.toRGB(value)) }
	Get()
		{ return .Horz.ColorRect.Get() }
	Dirty?(dirty = "")
		{ return .Horz.ColorRect.Dirty?(dirty) }
	NewValue(value)
		{ .Send("NewValue", value) }
	On_Choose()
		{
		if .readonly
			return

		if false isnt result = ChooseColorWrapper(.Get(), .Window.Hwnd,
			custColors: Object())
			{
			.Set(result)
			.NewValue(.Get())
			}
		}
	SetReadOnly(readonly)
		{
		.readonly = readonly
		}
	Destroy()
		{
		.Send('NoData')
		super.Destroy()
		}
	toRGB(color)
		{
		return Object?(color)
			? RGB(color[0], color[1], color[2])
			: color
		}
	}
