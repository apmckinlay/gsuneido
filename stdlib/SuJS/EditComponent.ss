// Copyright (C) 2018 Axon Development Corporation All rights reserved worldwide.
Component
	{
	styles: `
		.su-edit {
			padding: 0;
			position: relative;
			border-width: 1px;
			border-style: solid;
			border-color: #767676;
			box-sizing: border-box;
			margin: 0px;
		}
		.su-edit:focus {
			border-color: deepskyblue;
			outline: none;
		}`
	Name: 'Edit'
	DefaultWidth: 20
	DefaultHeight: 1
	ContextMenu: true
	New(.readonly = false, .bgndcolor = "", textcolor = "", hidden = false,
		tabover = false, font = "", size = "", weight = "", underline = false,
		width = false, height = false, .readOnlyBgndColor = false)
		{
		LoadCssStyles('su-edit.css', .styles)
		.init(height)
		.SetHidden(hidden)
		.SetFontAndSize(font, size, weight, underline, width, height)

		if tabover is true or .readonly is true
			.El.tabIndex = "-1"
		if .readonly is true
			.El.readOnly = true
		styles = Object()
		if textcolor isnt ""
			styles.color = ToCssColor(textcolor)
		// not sure why this is needed. It seems that Chrome uses a default width even
		// if min-width is set
		styles.width = .Xmin $ 'px'
		.SetStyles(styles)
		if .readOnlyBgndColor is false
			.readOnlyBgndColor = ToCssColor(CLR.ButtonFace)
		.SetBgndColor(.readonly is true ? .readOnlyBgndColor : .bgndcolor)

		.SetMinSize()
		.El.AddEventListener('focus', .focus)
		.El.AddEventListener('blur', .blur)
		.El.AddEventListener('input', .EN_CHANGE)
		.El.AddEventListener('mouseup', .onMouseUp)
		}

	init(height)
		{
		.height = height is false ? .DefaultHeight : height
		if .height > 1
			{
			.CreateElement('textarea', className: 'su-edit')
			.SetStyles(#('resize': 'none'))
			// flex align-items: baseline not work for empty fields on WebKit
			.El.placeholder = ' '
			.El.rows = .height
			}
		else
			{
			.CreateElement('input', className: 'su-edit')
			.El.spellcheck = false
			// flex align-items: baseline not work for empty fields on WebKit
			.El.placeholder = ' '
			.El.autocomplete = "su-do-not-autofill"
			}
		.SetStyles(#('padding-left': '4px', 'padding-right': '4px'))
		}
	padding: 10 // (4px padding + 1px border) * 2
	SetFontAndSize(font, size, weight, underline, width = false, height/*unused*/ = false,
		text = "M")
		{
		.SetFont(font, size, weight, underline)
		.OrigYmin = ymin = .Ymin
		.OrigXmin = xmin = .Xmin
		metrics = SuRender().GetTextMetrics(.El, text)
		.Xmin = xmin isnt 0
			? xmin
			: metrics.width * (width is false ? .DefaultWidth : width) + .padding
		lineHeight = metrics.height + .padding - 2/*=border*/
		.Ymin = ymin isnt 0
			? ymin
			: lineHeight * .height + 2/*=border*/
		if .height > 1
			{
			.El.SetStyle('line-height', lineHeight $ 'px')
			if ymin isnt 0
				.El.SetStyle('height', ymin $ 'px')
			}
		}

	SetReadOnly(readOnly)
		{
		if (.readonly)
			return
		readOnly = .enabled is false or readOnly is true
		.El.readOnly = readOnly
		.SetBgndColor(readOnly ? .readOnlyBgndColor : .bgndcolor)
		}

	SetReadOnlyColor(color)
		{
		.readOnlyBgndColor = ToCssColor(color)
		.SetBgndColor(.El.readOnly is true ? .readOnlyBgndColor : .bgndcolor )
		}

	GetReadOnly()
		{
		return .El.readOnly
		}

	enabled: true
	SetEnabled(.enabled)
		{
		.SetReadOnly(not enabled)
		}

	GetEnabled()
		{
		return .enabled
		}

	SetBgndColor(color)
		{
		if color is "" or color is false
			color = CLR.WHITE
		.SetStyles(Object('background-color': ToCssColor(color)))
		}

	Set(value)
		{
		if (not String?(value))
			value = Display(value)
		.El.value = value
		}

	// used by SetWindowText to set value directly
	SetText(text)
		{
		.El.value = text
		}

	GetSel()
		{
		return [.El.selectionStart, .El.selectionEnd]
		}

	SetSel(start, end)
		{
		// .SetSelectionRange tirggers focus on WekKit
		if SuRender().Engine is 'WebKit' and SuGetFocus() isnt .El
			return
		.El.SetSelectionRange(start, end)
		}

	SelectAll()
		{
		.SetSel(0, .El.value.Size())
		}

	ReplaceSel(text)
		{
		value = .El.value
		range = .GetSel()

		tmpTxt = value[..range[0]] $ text $ value[range[1]..]
		if .El.HasAttribute('maxlength')
			{
			maxlength = Number(.El.GetAttribute('maxlength')) - 3 /*=number of ellipsis*/
			.El.value = tmpTxt.Ellipsis(maxlength, atEnd:)
			}
		else
			.El.value = tmpTxt

		.EN_CHANGE()
		}

	SetCue(cue)
		{
		.El.placeholder = cue
		}

	blur()
		{
		.EventWithFreeze('EN_KILLFOCUS')
		return 0
		}

	focus()
		{
		.Event('EN_SETFOCUS')
		.SetValid() // don't color invalid when focused
		return 0
		}

	onMouseUp(event)
		{
		.Event(#UpdateSel, [.El.selectionStart, .El.selectionEnd])
		if event.button isnt 0
			return
		.Event('LBUTTONUP')
		}

	SetValid(valid? = true, force = false)
		{
		if SuGetFocus() is .El and not force
			valid? = true // don't color when we have focus
		.SetBgndColor(.GetReadOnly()
			? .readOnlyBgndColor
			: valid?
				? .bgndcolor
				: CLR.ErrorColor)
		}

	EN_CHANGE()
		{
		if .Destroyed?() or .GetReadOnly()
			return 0
		.Event(#UpdateSel, [.El.selectionStart, .El.selectionEnd])
		.Event('EN_CHANGE', .El.value)
		return 0
		}

	On_Delete()
		{
		.selectAllIfNoSelection()
		SuUI.GetCurrentDocument().ExecCommand('delete')
		}
	On_Cut()
		{
		.selectAllIfNoSelection()
		SuUI.GetCurrentDocument().ExecCommand('cut')
		}
	On_Copy()
		{
		.selectAllIfNoSelection()
		SuUI.GetCurrentDocument().ExecCommand('copy')
		}

	selectAllIfNoSelection()
		{
		sel = .GetSel()
		if sel[0] is sel[1]
			.SelectAll()
		}

	On_Paste()
		{
		SuClipboardPasteString(this, .ReplaceSel)
		}
	On_Undo()
		{
		SuUI.GetCurrentDocument().ExecCommand('undo')
		}
	}
