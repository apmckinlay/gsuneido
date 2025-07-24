// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Name: FilterButton
	New()
		{
		super(Object('EnhancedButton', command: 'Filter', image: 'zoom.emf',
			mouseEffect:, imagePadding: .1, tip: 'Show/Hide Select',
			imageColor: CLR.Inactive))
		.image = .EnhancedButton
		}
	On_Filter()
		{
		.Send('ToggleFilter')
		}
	UpdateStatus(ctrl, active? = false)
		{
		if false isnt btn = ctrl.FindControl('FilterButton')
			ctrl.Defer({ btn.UpdateImage(active?) }, uniqueID: 'update_filter_button')
		}
	UpdateImage(active? = false)
		{
		color = active? ? CLR.Active : CLR.Inactive
		.EnhancedButton.SetImageColor(color, color)
		}
	}