// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
VirtualListThumbImageButtonControl
	{
	Name: 'VirtualListExpandSwitchToFormButton'
	New()
		{
		super(@.layout())
		}

	layout()
		{
		return Object(tip: "switch to form view",
			image: 'view_form.emf', mouseEffect:, imagePadding: 0.1)
		}

	LBUTTONDOWN()
		{
		super.LBUTTONDOWN()
		.Send('VirtualListExpand_SwitchToForm')
		return 0
		}
	}
