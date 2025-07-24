// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	CallClass(prompt = "", title = "", hwnd = 0, ctrl = 'Field',
		valid = function (unused) { return '' }, noCancel = false)
		{
		prompt = TranslateLanguage(prompt)
		title = TranslateLanguage(title)
		return ToolDialog(hwnd, Object(this, prompt, ctrl, valid, noCancel),
			:title, closeButton?: false, keep_size: false)
		}
	New(prompt, ctrl, .valid, noCancel)
		{
		super(.makecontrols(prompt, ctrl, noCancel))
		}
	makecontrols(prompt, ctrl, noCancel)
		{
		ctrl = String?(ctrl) or Class?(ctrl) ? Object(ctrl) : ctrl.Copy()
		ctrl.name = 'Field'
		buttons = noCancel ? #(Horz Fill OkButton) : "OkCancel"
		return Object('Vert'
			Object('Pair' Object('Static' prompt) ctrl)
			'Skip', buttons)
		}
	On_OK()
		{
		field = .Vert.Pair.Field
		if not field.Valid?()
			return
		value = field.Get()
		if '' is err = (.valid)(value)
			.Window.Result(value)
		else
			.AlertInfo('Invalid Entry', err)
		}
	}