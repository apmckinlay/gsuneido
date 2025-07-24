// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Title: 'Refactor'
	CallClass(source, control)
		{
		OkCancel(Object("Refactor", new control, source.Ctrl), .Title, source.Hwnd)
		}
	New(aRefactor, libview)
		{
		super(.controls(aRefactor))

		library = libview.CurrentTable()
		name = libview.CurrentName()
		.editor = libview.Editor
		if aRefactor.SelectWord
			.editor.SelectCurrentWord()
		.select = .editor.GetSelect()

		.Data = .Vert.Data
		.data = Record(
			:library,
			:name,
			text: .editor.Get(),
			select: .select,
			ctrl: this,
			editor: .editor,
			:libview
			)
		.Data.Set(.data)
		.Defer(.init)
		.constructed = true
		}

	init()
		{
		if .ref.Init(.data) is false
			.On_Cancel()
		}

	controls(aRefactor)
		{
		.ref = aRefactor
		header = Object('Title', .ref.Name)
		if .ref.Desc isnt ""
			header = Object('Vert', header, #(Skip 4),
				Object('Static', .ref.Desc, size: '+2'))
		controls = Object?(.ref.Controls) ? .ref.Controls : .ref.Controls()
		return Object('Vert',
			Object('Border', header),
			'EtchedLine',
			Object('Record',
				Object('Border', controls)),
			)
		}

	constructed: false
	Edit_Change()
		{
		if .constructed isnt false
			.Data.HandleFocus()
		}

	OK()
		{
		data = .Data.Get()
		if '' isnt errs = .ref.Errors(data)
			{
			.AlertError('Refactor Error', errs)
			return false
			}
		if '' isnt warns = .ref.Warnings(data)
			if not OkCancel(warns $ "\n\nContinue?", .ref.Name, flags: MB.ICONQUESTION)
				return false
		if .ref.Process(data) is true
			{
			line = .editor.GetFirstVisibleLine()
			.editor.PasteOverAll(data.text)
			.editor.SetFirstVisibleLine(line)
			.editor.SetSelect(.select.cpMin)
			return true
			}
		return true
		}
	}