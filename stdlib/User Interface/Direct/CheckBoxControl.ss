// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
TextPlus
	{
	Name: 'CheckBox'
	CustomizableOptions: #(hidden, readonly, tabover)
	Antialias: false
	New(@args)
		{
		super(@args)
		.Send('Data')
		if "" isnt set = args.GetDefault(#set, "")
			{
			.Set(set)
			.Send('NewValue', .Get())
			}
		}
	PaintOuter(hdc, rect)
		{
		// box
		Rectangle(hdc, rect.GetX(), rect.GetY(), rect.GetX2(), rect.GetY2())
		}
	PaintInner(hdc, rect, imageBrush)
		{
		checkmark = ImageResource('checkmark.emf')
		checkmark.Draw(hdc, rect.GetX() + 1, rect.GetY() + 1, rect.GetHeight() - 2,
			rect.GetWidth() - 2, imageBrush)
		}
	Toggle()
		{
		super.Toggle()
		if not .GetReadOnly()
			.Send('NewValue', .Get())
		}

	Destroy()
		{
		.Send('NoData')
		super.Destroy()
		}
	}