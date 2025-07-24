// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
EnhancedButtonControl
	{
	New(command, .image, book = 'imagebook', tip = "")
		{
		super(:command, :image, :book, :tip, imageColor: 0x737373,
			mouseOverImageColor: 0x00cc00, imagePadding: 0.15)
		.Ymin -= 2 	// make it same size as field control (10 in ButtonControl)
					// to avoid double lines between rows
		}

	readonly: false
	SetReadOnly(readonly)
		{
		.readonly = readonly
		.SetEnabled(not readonly)
		}
	}