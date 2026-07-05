// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
EditorControl
	{
	Name: 			'Encrypt Editor'
	DefaultHeight: 	2
	Get()
		{
		return super.Get().Xor(EncryptControlKey())
		}

	Set(value)
		{
		super.Set(value.Xor(EncryptControlKey()))
		}

	ZoomArgs()
		{
		zoomArgs = super.ZoomArgs()
		zoomArgs.zoomDialog = EncryptZoomControl
		return zoomArgs
		}
	}
