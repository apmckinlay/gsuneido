// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
Addon_status
	{
	StatusBarControl: EditorStatusbar
	Status(@status)
		{
		if not status.Member?(#invalid) and status.Member?(#normal)
			.AddonControl.ClearError()
		.AddonControl.Status(@status)
		}
	}
