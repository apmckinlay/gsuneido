// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
CustomizableFieldDialogPropertiesEditor
	{
	GetControls()
		{
		return Object('Vert', #Skip)
		}
	Get()
		{
		x = super.Get()
		x.control = Object()
		x.format = Object()
		return x
		}
	}