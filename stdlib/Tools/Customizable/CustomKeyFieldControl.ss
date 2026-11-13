// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
KeyFieldControl
	{
	SetReadOnly(readOnly)
		{
		// Required in order to allow the Drop Down button to be interactive when
		// unauthorized (giving the user an alert when interacted with)
		if .Parent.QueryUnauthorized?()
			readOnly = true
		super.SetReadOnly(readOnly)
		}
	}
