// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
WndPaneControl
	{
	ComponentName: 'TabWndPane'

	SetVisible(visible)
		{
		Assert(Boolean?(visible))
		.Act('SetVisible', visible)
		}
	}