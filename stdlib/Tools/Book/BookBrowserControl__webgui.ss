// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
PassthruController
	{
	Name:			'BookBrowser'
	ComponentName: 	'BookBrowser'
	program:		false

	Open(control)
		{
		.ActWith()
			{ |reservation|
			try
				{
				.program = .Construct(.wrapper(control))
				DoStartup(.program)
				}
			catch (err)
				{
				try .program.Destroy()
				SuRenderBackend().CancelAllAfter(reservation.at)
				.program = .Construct(
					.wrapper(#('Center', ('Static', "unable to display page"))))
				.Defer({ throw err })
				}
			Object(#Open, .program.GetLayout())
			}
		}

	wrapper(ctrl)
		{
		// use WndPane to clip application screen,
		// so it does not overlap with other controls, like bookmarks
		return Object('WndPane', ctrl, windowClass: "SuBtnfaceArrow")
		}

	Close()
		{
		if .program isnt false
			{
			.program.Destroy()
			.program = false
			}
		else
			SuneidoLog('WARNING: BookBrowserControl tried to call .Destroy on program' $
				' but program is set to false (boolean.Destroy).  Look for errors prior' $
				' to this message from screen entry')
		Suneido.CurrentBookOption = ''
		.Act(#Close)
		}

	unfrozen: false
	PageFrozen?()
		{
		if false is items = .getValidationItems()
			return false
		frozen = false
		for item in items
			if not item.ConfirmDestroy()
				if (item.Method?('CloseWindowConfirmation') and
					not item.CloseWindowConfirmation())
					return true
				else
					frozen = true
		if frozen and CloseWindowConfirmation(.Window.Hwnd)
			{
			// .unfrozen is to prevent more than one confirmation
			// this is necessary because PageFrozen? can be called multiple times
			.unfrozen = Date()
			frozen = false
			}
		return frozen
		}
	getValidationItems()
		{
		return .program is false or
			(Date?(.unfrozen) and Date().MinusSeconds(.unfrozen) < 1)
			? false
			: .Window.GetValidationItems()
		}

	ProgramPage?()
		{
		return .program isnt false
		}
	GetChildren()
		{
		// Have to return the superclass controller's children as well, or you
		// will leak when Destroy() is called -- because the superclass
		// Destroy() destroys only the children returned by .GetChildren()
		controllerChildren = super.GetChildren()
		return .program is false
			? controllerChildren
			: Object(.program).Union(controllerChildren)
		}

	BookRefresh()
		{
		.Send('BookRefresh')
		}
	}