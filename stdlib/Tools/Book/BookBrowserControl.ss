// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// this exists primarily to control program pages
PassthruController
	{
	// data:
	x:			false
	program:	false
	Name:		'BookBrowser'
	// interface:
	Open(control)
		{
		try
			{
			.program = .Construct(.wrapper(control))
			DoStartup(.program)
			if not .Window.RunPendingRefresh()
				.resize_program()
			}
		catch (err)
			{
			try .program.Destroy()
			.program = .Construct(
				.wrapper(Object('Center', Object('Static', "unable to display page"))))
			.resize_program()
			.Defer({ throw err })
			}
		}
	wrapper(ctrl)
		{
		// use WndPane to clip application screen,
		// so it does not overlap with other controls, like bookmarks
		return Object('WndPane', ctrl, windowClass: "SuBtnfaceArrow")
		}
	resize_program()
		{
		if .x isnt false
			.program.Resize(.x, .y, .w, .h)
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
		}
	Resize(.x, .y, .w, .h)
		{
		if (.program isnt false)
			.program.Resize(x, y, w, h)
		super.Resize(x, y, w, h)
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
			.restoreAttachmentFiles(items)
			// .unfrozen is to prevent more than one confirmation
			// this is necessary because PageFrozen? can be called multiple times
			.unfrozen = Date()
			frozen = false
			}
		return frozen
		}
	restoreAttachmentFiles(items)
		{
		for item in items
			if item.Method?('RestoreAttachmentFiles')
				item.RestoreAttachmentFiles()
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
