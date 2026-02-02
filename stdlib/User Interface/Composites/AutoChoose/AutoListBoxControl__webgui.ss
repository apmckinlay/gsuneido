// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
ListBoxControl
	{
	ComponentName: "AutoListBox"
	SELCHANGE(curSel)
		{
		super.SELCHANGE(curSel)
		.Send('AutoListBox_Click', curSel)
		}
	}