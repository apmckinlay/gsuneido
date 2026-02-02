// Copyright (C) 2012 Suneido Software Corp. All rights reserved worldwide.
// NOTE: for keyboard accellerators to work
//		 the top level control (.Window.Ctrl) must be a Controller
Controller
	{
	Name: "WordPad"
	Unsortable: true

	New(@args)
		{
		super(.layout(args))
		.edit = .Vert.Editor
		.Redir('On_Bold', .edit)
		.Redir('On_Italic', .edit)
		.Redir('On_Underline', .edit)
		.Redir('On_Strikeout', .edit)
		.Redir('On_ResetFont', .edit)
		.bold_button = .FindControl('B')
		.italic_button = .FindControl('I')
		.underline_button = .FindControl('U')
		.strikeout_button = .FindControl('S')
		.reset_button = .FindControl('X')
		.Send('Data')
		}
	edit: class // "null" class to avoid checks for .edit existing
		{
		Default(@unused)
			{ return '' }
		}
	Commands:
		(
		(Bold,		"Ctrl+B")
		(Italic,	"Ctrl+I")
		(Underline,	"Ctrl+U")
		(Strikeout, "Ctrl+S")
		(ResetFont,	"Ctrl+Space")
		)
	layout(args)
		{
		textLimit = args.GetDefault('textLimit', 50_000) /*= text limit with html format*/
		args.Set_default(false)
		.readonly = args.readonly
		ignoreOb = Object()
		if args.ignoreList isnt false
			ignoreOb.ignore = args.ignoreList
		layout = Object('Vert')
		if not .readonly
			layout.Add(Object('Horz',
				.button('B', 'Bold', 'Bold (Ctrl+B)', weight: 'BOLD'),
				.button('I', 'Italic', 'Italic (Ctrl+I)', italic:),
				.button('U', 'Underline', 'Underline (Ctrl+U)', underline:),
				.button('S', 'Strikeout', 'Strikeout (Ctrl+S)', strikeout:),
				'Skip',
				.button('X', 'ResetFont', 'Remove formatting (Ctrl+Space)')
				))
		sci = Object('ScintillaAddonsRichEditor',
			height: args.height, width: args.width,	xmin: args.xmin, ymin: args.ymin,
			Addon_speller: ignoreOb, zoom: args.zoom, readonly: .readonly,
			tabthrough: args.tabthrough, :textLimit)
		// If no font or fontSize was specified, don't pass them to
		// ScintillaAddonsRichEditor. This will default the ScintillaAddonsRichEditor
		// to use the Suneido.logfont
		if args.font isnt false
			sci.font = args.font
		if args.fontSize isnt false
			sci.fontSize = args.fontSize
		layout.Add(sci)
		return layout

		}
	button(text, command, tip,
		weight = '', italic = false, underline = false, strikeout = false)
		{
		return Object('EnhancedButton', text, :command, :tip, tabover:, :weight,
			font: 'Verdana', size: 10,
			:italic, :underline, :strikeout, buttonStyle:, mouseEffect:)
		}

	SetReadOnly(readonly = true)
		{
		.edit.SetReadOnly(readonly)
		}

	Dirty?(dirty = "")
		{
		return .edit.Dirty?(dirty)
		}

	Set(s)
		{
		.edit.Set(s.Trim())
		}

	Get()
		{
		return .edit.Get()
		}

	GetText()
		{
		return .edit.GetText()
		}

	Scintilla_DoubleClick()
		{
		.Send('Scintilla_DoubleClick')
		}

	Scintilla_SetFocus()
		{
		.Defer(.setAccels, uniqueID: 'set_accels')
		}

	accelsSet: false
	setAccels()
		{
		// make accelerators work inside book
		if .Destroyed?() or not .Window.Member?(#Ctrl) or
			not .Window.Ctrl.Base?(Controller) or .accelsSet is true
			return

		.accelsSet = true
		.curAccels = .Window.SetupAccels(.Commands)
		.Window.Ctrl.Redir('On_Bold', .edit)
		.Window.Ctrl.Redir('On_Italic', .edit)
		.Window.Ctrl.Redir('On_Underline', .edit)
		.Window.Ctrl.Redir('On_Strikeout', .edit)
		.Window.Ctrl.Redir('On_ResetFont', .edit)
		}

	Scintilla_KillFocus()
		{
		if .Window.Ctrl is this // from list edit window
			return
		.resetAccels()
		}

	resetAccels()
		{
		if .accelsSet is false
			return
		.Window.Ctrl.RemoveRedir(.edit)
		if .curAccels isnt false
			.Window.RestoreAccels(.curAccels)
		.accelsSet = false
		}

	GetSelText()
		{
		return .edit.GetSelText()
		}

	// handles scrolling to the end of the text
	// WITHOUT setting Focus.
	ScrollEnd()
		{
		.edit.SetSelect(.edit.GetLength())
		}

	SetSelect(chrg)
		{
		if chrg is "end"
			chrg = .edit.GetLength()
		.edit.SetSelect(chrg)
		.edit.SetFocus()
		}

	SetText(value)
		{
		.edit.SetText(value)
		}

	SetFocus()
		{
		.edit.SetFocus()
		}

	AppendText(text)
		{
		.edit.AppendText("\r\n" $ text)
		}

	BeforeSave()
		{
		.edit.BeforeSave()
		}

	NewValue(value)
		{
		.Send("NewValue", value)
		}

	MenuSelect(tip)
		{
		.Vert.Status.Set(tip)
		}

	Status(status)
		{
		.Vert.Status.Set(status)
		}

	ToggleFontButton(button, state)
		{
		if not .readonly
			this['ScintillaRichWordAddonsControl_' $ button $ '_button'].Pushed?(state)
		}

	On_Print()
		{
		Params.On_Print(Object('ScintillaRichWrap', .Get(), width: 100),
			title: SelectPrompt(.Name), name: 'print_editor', previewWindow: .Window.Hwnd)
		}

	CloseZoom()
		{
		.Send("CloseZoom")
		}

	ZoomReadonly(value)
		{
		ScintillaRichWordZoomControl(0, value, readonly:)
		}

	ValidData?(@args)
		{
		return ScintillaAddonsRichEditorControl.ValidData?(@args)
		}

	Valid?()
		{
		return .edit.Valid?()
		}

	Default(@args)
		{
		return .edit[args[0]](@+1 args)
		}
	curAccels: false
	Destroy()
		{
		.resetAccels()
		.Send('NoData')
		super.Destroy()
		}
	}
