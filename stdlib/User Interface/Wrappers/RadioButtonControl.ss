// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// NOTE: does NOT register with RecordControl
// normally only used by RadioButtonsControl and RadioGroupsControl
// (not standalone)
TextPlus
	{
	Name: 'RadioButton'
	PaintOuter(hdc, rect)
		{
		// circle
		h = rect.GetHeight()
		Ellipse(hdc, 0, 0, h, h)
		}
	PaintInner(hdc, rect)
		{
		// dot
		h = rect.GetHeight()
		d = Max(14, (h / 6).Int())
		Ellipse(hdc, d, d, h - d, h - d)
		}
	// for testing, comment this out so you can toggle freely
	Toggle()
		{
		if .Get() isnt true
			.Send('Picked', .GetText())
		}
	}