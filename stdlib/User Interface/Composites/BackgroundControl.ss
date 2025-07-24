// Copyright (C) 2008 Suneido Software Corp. All rights reserved worldwide.
Control
	{
	// width and height are Device-independent pixel (DIPs)
	New(image, controls, width = 500, height = 250)
		{
		.image = .Construct(['Image', image stretch:])
		.image.SetEnabled(false) // so it doesn't steal mouse events
		.controls = .Construct(['Vert' controls])
		// make sure it's big enough for content
		// NOTE: this may stretch the image proportionally
		.Xmin = Max((width * GetDpiFactor()).Round(0), .controls.Xmin)
		.Ymin = Max((height * GetDpiFactor()).Round(0), .controls.Ymin)
		}
	GetChildren()
		{
		return .controls.GetChildren()
		}
	Resize(x, y, w, h)
		{
		.image.Resize(x, y, w, h)
		.controls.Resize(x, y, w, h)
		}
	Destroy()
		{
		.image.Destroy()
		.controls.Destroy()
		super.Destroy()
		}
	}
