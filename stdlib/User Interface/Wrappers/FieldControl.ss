// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
// Contributions from Claudio Mascioni
EditControl
	{
	Name:			'Field'
	MaxCharacters: 512

	New(width = 20, status = "", readonly = false,
		font = "", size = "", weight = "", underline = false,
		password = false, justify = "LEFT", style = 0,
		set = false, mandatory = false, .trim = true,
		bgndcolor = "", textcolor = "", hidden = false, tabover = false,
		cue = false, readOnlyBgndColor = false, acceptDrop = false,
		upper = false, lower = false, maxCharacters = false)
		{
		super(mandatory, readonly, .buildStyle(style, justify, password, upper, lower),
			:bgndcolor, :textcolor, :hidden, :tabover, :width, height: 1, :cue,
			:font, :size, :weight, :underline, :readOnlyBgndColor, :status)
		.MaxCharacters = maxCharacters is false
			? .MaxCharacters
			: maxCharacters
		if set isnt false
			{
			.Set(set)
			.Send("NewValue", .Get())
			}

		if acceptDrop
			DragAcceptFiles(.Hwnd, true)
		}

	buildStyle(style, justify, password, upper, lower)
		{
		passwordStyle = password is true ? ES.PASSWORD : 0
		if upper and lower
			throw "upper and lower options cannot be used together"
		caseStyle = upper ? ES.UPPERCASE : lower ? ES.LOWERCASE : 0
		return style | ES.AUTOHSCROLL | ES[justify] | passwordStyle | caseStyle
		}

	SetBackColor(color)
		{
		.SetBgndColor(color)
		}

	// 64k is the gSuneido dll argument limit
	// 30k is the default text limit of windows edit control
	pasteLimit: 30000
	PASTE()
		{
		if .GetReadOnly()
			return 'callsuper'

		text = ClipboardReadString()
		if String?(text)
			{
			newText = text.Tr('\r\n', ' ')
			if newText.Size() > .pasteLimit
				newText = newText[::.pasteLimit]
			if text isnt newText
				{
				.SendMessageTextIn(EM.REPLACESEL, true, newText)
				return 0
				}
			}
		return 'callsuper'
		}
	RBUTTONDOWN()
		{
		// need to make sure all the text is highlighted when users do right-click > copy
		// if they don't already have a selection
		if not .HasFocus?() and .GetSel() is #(0,0)
			.SelectAll()
		return 'callsuper'
		}
	EN_KILLFOCUS()
		{
		dirty? = .Dirty?()
		// has to call super here
		super.EN_KILLFOCUS()
		if dirty?
			.Send("NewValue", .Get())
		return 0
		}
	Get()
		{
		text = GetWindowText(.Hwnd)
		return .trim ? text.Trim() : text
		}

	Set(value)
		{
		if not String?(value)
			value = Display(value)
		.Dirty?(false)
		SetWindowText(.Hwnd, value)
		// select all text, this is useful for Browse and other controls
		// which create fields on the fly
		.SelectAll()
		}

	SetReadOnly(readOnly)
		{
		// Prevent edit control from scrolling to the end when going into edit mode
		.SetSel(0,0)
		super.SetReadOnly(readOnly)
		}
	DROPFILES(wParam)
		{
		.Send('FieldDropFiles', wParam)
		return 0
		}

	changed?: false // killfocus sets dirty back to false
	Valid?()
		{
		if super.Valid?() is false
			return false

		return .ValidLength?(.Get())
		}

	ValidLength?(data)
		{
		if ((.Dirty?() or .changed?) and
			.ValidTextLength?(data, .MaxCharacters) is false)
			{
			if not .changed? // only show alert once
				.AlertInfo("Invalid Text", 'data exceeds the ' $ .MaxCharacters $
					' character limit for this field.')
			.changed? = true
			return false
			}
		return true
		}

	ValidTextLength?(data, maxCharacters)
		{
		if not String?(data)
			return true
		return data.Size() <= maxCharacters
		}

	ValidData?(@args)
		{
		maxCharacters = args.GetDefault('maxCharacters', .MaxCharacters)
		return super.ValidData?(@args) and .ValidTextLength?(args[0], maxCharacters)
		}
	}
