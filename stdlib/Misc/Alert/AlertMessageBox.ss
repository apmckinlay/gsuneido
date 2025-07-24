// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
// flags and return values are a subset of MessageBox
// NOTE: should be called via Alert, not directly
// displays message with AlertTextControl
Controller
	{
	CallClass(hwnd, msg, title = "ALERT", flags = 0)
		{
		result = ToolDialog(hwnd, [this, msg, flags, GetWorkArea()], border: 0, :title,
			closeButton?: MB.YESNO isnt (flags & MB.BUTTONBITS), keep_size: false)
		return result is false ? ID.CANCEL : result // false is from window X close button
		}
	New(.msg, flags, rect)
		{
		super(.controls(msg, flags, rect))
		.Window.SetupAccels(#((Copy, "Ctrl+C")))
		.Redir('On_Copy', .FindControl(#AlertText))
		}
	controls(msg, flags, rect)
		{
		ymin = ((rect.bottom-rect.top)*.75 /*=scaleFactor*/)/GetDpiFactor()
		xmin = 9999 // set this very large so it will trim down
		ctrls =
			[#Vert,
				[#WndPane, // for white background
					[#Scroll,
						[#Border,
							horz = [#Horz, [#AlertText msg]],
						]
					// trim: will shrink the Scroll to match the size of the AlertText
					noEdge:, wndclass: 'SuWhiteArrow', trim:, :xmin, :ymin,
					xstretch: 0, ystretch: 0],
				],
				[#Border, .Buttons(flags)]
			]
		if false isnt icon = .icon(flags)
			horz.Add(icon, #Skip, at: 1)
		return ctrls
		}
	Buttons(flags)
		{
		buttons = [#HorzEqual #Fill xstretch: 0]
		switch flags & MB.BUTTONBITS
			{
		case MB.OK: // 0, default
			buttons.Add(#(Button OK width: 8))
		case MB.OKCANCEL:
			buttons.Add(#(Button OK width: 8), #Skip, #(Button Cancel))
		case MB.YESNO:
			.allowCancel = false
			buttons.Add(#(Button Yes width: 8), #Skip, #(Button No))
		case MB.YESNOCANCEL:
			buttons.Add(#(Button Yes width: 8), #Skip, #(Button No), #Skip,
				#(Button Cancel))
			}
		return buttons
		}
	icon(flags)
		{
		icon = false
		switch flags & 0x70
			{
		case 0:
			return false
		case MB.ICONINFORMATION:
			icon = [#Image 'info.emf' color: 0xd57601]
		case MB.ICONQUESTION:
			icon = [#Image 'questionMark_black.emf' color: 0xd57601]
		case MB.ICONWARNING:
			icon = [#Image 'triangle-warning.emf' color: CLR.orange]
		case MB.ICONERROR:
			icon = [#Image 'cross.emf' color: 0x1239ef]
			}
		if icon isnt false
			{
			icon.bgndcolor = 0xffffff
			icon.xmin = 30
			}
		return icon
		}
	On_Yes()
		{
		.Window.Result(ID.YES)
		}
	On_No()
		{
		.Window.Result(ID.NO)
		}
	On_OK()
		{
		.Window.Result(ID.OK)
		}
	allowCancel: true
	On_Cancel()
		{
		if .allowCancel
			.Window.Result(ID.CANCEL)
		}
	}