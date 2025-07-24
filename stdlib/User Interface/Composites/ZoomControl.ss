// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// Testing:
//	- ESC should NOT modify text
Controller
	{
	edit: false
	EditorControl: 'Editor' // overridden for Scintilla
	CallClass(hwnd, text, readonly, font = "", size = "")
		{
		ctrl = [this, text, title: 'Zoom', :readonly, :font, :size]
		if readonly
			ModalWindow(ctrl, border: 0)
		else
			text = OkCancel(ctrl, "Zoom", hwnd)
		return text
		}
	New(.text, readonly = false, font = "", size = "")
		{
		super(.Layout(readonly, font, size))
		.edit = .FindControl(.EditorControl)
		.edit.Set(text)
		.Defer(.setfocus)
		}
	Commands:
		(
		(Undo,	"Ctrl+Z",	"Undo the last action")
		(Cut,	"Ctrl+X",	"Cut the selected text to the clipboard")
		(Copy,	"Ctrl+C",	"Copy the selected text to the clipboard")
		(Paste,	"Ctrl+V",	"Insert the contents of the clipboard")
		(Print,	"Ctrl+P",	"Print the contents of the clipboard")
		)
	Layout(.readonly, font, size)
		{
		layout = Object('Vert',
			.readonly
				? #('Toolbar', 'Copy', '', 'Print')
				: #('Toolbar', 'Cut', 'Copy', 'Paste','', 'Undo', 'Print'),
			Object(.EditorControl, :readonly, :font, :size, zoom:,
				name: .EditorControl))
		return layout
		}
	setfocus()
		{
		if .edit isnt false
			.edit.SetFocus()
		}
	Get()
		{
		.edit.Get()
		}
	On_Print()
		{
		.edit.On_Print()
		}
	OK()
		{
		return .edit.Get()
		}
	Cancel()
		{
		return .text // original text
		}
	CloseZoom() // sent by editor control on F6
		{
		.Send('On_OK')
		//needed to make sure the window gets closed when in "ReadOnly"
		if .readonly
			.Window.Result(0)
		}
	}
