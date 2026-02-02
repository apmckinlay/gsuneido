// Copyright (C) 2022 Axon Development Corporation All rights reserved worldwide.
Control
	{
	Name: 'VirtualListExpandButtons'
	New(.switchToForm)
		{
		}

	GetExpandButtons(showExpandButton)
		{
		buttons = Object()
		if showExpandButton is true
			buttons.Add(Object('ExpandButton_LBUTTONDOWN', Object('data-expanded',
				IconFontHelper.GetCode('next.emf'),
				IconFontHelper.GetCode('forward.emf')),
				alwaysDisplay?:))
		if .switchToForm is true
			buttons.Add(Object('VirtualListExpand_SwitchToForm',
				IconFontHelper.GetCode('view_form.emf')))
		return buttons
		}
	}