// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Name: 'Notes'
	New(title = false)
		{
		super(#('EnhancedButton', command: 'Notes', image: 'notes.emf',
			mouseEffect:, imagePadding: .1))
		.set_title(title)
		.image = .EnhancedButton
		UserNotes.EnsureTable()
		.UpdateImage()
		}
	set_title(title)
		{
		if title is false
			{
			title = .Send('Notes_Title')
			title = String?(title) ? title : ''
			}
		.title = title
		}
	UpdateImage()
		{
		opt = .Send('AccessGoTo_CurrentBookOption')
		.cur_book_option = String?(opt)
			? opt
			: Suneido.Member?("CurrentBookOption")
				? Suneido.CurrentBookOption
				: ''
		x = Query1(.query())
		no_notes = x is false or x.text is ""
		if no_notes
			.image.SetImageColor(CLR.Inactive, CLR.Inactive)
		else
			.image.SetImageColor(CLR.Active, CLR.Active)
		.image.ToolTip((no_notes ? "Add notes for " : "Notes for ") $ .title)
		}
	query()
		{
		return UserNotes.Query(.title, .cur_book_option)
		}
	UpdateStatus(ctrl)
		{
		if false isnt noteBtn = ctrl.FindControl('Notes')
			noteBtn.UpdateImage()
		}
	On_Notes()
		{
		result = NotesEditControl(.Window.Hwnd, .title, .query(),
			book_option: .cur_book_option)
		if result is true
			.UpdateImage()
		}
	}
