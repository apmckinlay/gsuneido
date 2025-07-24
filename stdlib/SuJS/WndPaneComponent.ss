// Copyright (C) 2020 Axon Development Corporation All rights reserved worldwide.
// May need its own implementation later, uses HtmlDiv for now
HtmlDivComponent
	{
	Name: 'WndPane'
	ContextMenu: true
	New(control, bgColor)
		{
		super(control)
		.El.SetStyle('background-color', bgColor)
		}
	}
