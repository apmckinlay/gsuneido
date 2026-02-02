// Copyright (C) 2020 Axon Development Corporation All rights reserved worldwide.
ScintillaAddonsControl
	{
	ComponentName: 'ScintillaAddonsRichEditor'
	New(@args)
		{
		super(@.processArgs(args))
		.textLimit = args.GetDefault('textLimit', false)
		if .readonly = args.GetDefault(#readonly, false)
			.setBackground(GetSysColor(COLOR.BTNFACE))
		}

	processArgs(args)
		{
		.zoom = args.GetDefault('zoom', false)
		zoomArgs = [zoom: .zoom, zoom_ctrl: ScintillaRichWordZoomControl]
		return args.MergeNew([wrap:, margin: 8, exStyle: WS_EX.STATICEDGE,
			Addon_url:, Addon_multiple_selection:, Addon_zoom: zoomArgs,
			font: 		Suneido.logfont.lfFaceName,
			fontSize: 	'+0',
			weight: 	FW.NORMAL,
			italic: 	false])
		}

	setBackground(bgnd)
		{
		.StyleSetBack(0, bgnd)
		}

	On_Bold()
		{
		if .GetReadOnly() is true
			return

		.Act(#On_Bold)
		.Dirty?(true)
		}

	On_Italic()
		{
		if .GetReadOnly() is true
			return

		.Act(#On_Italic)
		.Dirty?(true)
		}

	On_Underline()
		{
		if .GetReadOnly() is true
			return

		.Act(#On_Underline)
		.Dirty?(true)
		}

	On_Strikeout()
		{
		if .GetReadOnly() is true
			return

		.Act(#On_Strikeout)
		.Dirty?(true)
		}

	On_ResetFont()
		{
		// Might be able to use SCI_CLEARDOCUMENTSTYLE here
		if .GetReadOnly() is true
			return

		.Act(#On_ResetFont)
		.Dirty?(true)
		}

	UpdateButtons(states)
		{
		.Send('ToggleFontButton', 'bold', states.bold)
		.Send('ToggleFontButton', 'italic', states.italic)
		.Send('ToggleFontButton', 'underline', states.underline)
		.Send('ToggleFontButton', 'strikeout', states.strikeout)
		}

	EN_CHANGE()
		{
		.getCached = false
		return super.EN_CHANGE()
		}

	Dirty?(dirty = "")
		{
		if dirty is true
			.getCached = false
		return super.Dirty?(dirty)
		}

	Get()
		{
		if .getCached isnt false
			return .getCached

		s = super.Get()
		str = ScintillaRichEditorHelper.Build(s, .styles)
		.getCached = str
		return str
		}

	GetText(includeStrikeThrough = false)
		{
		return .getText(0, .GetLength(), includeStrikeThrough)
		}

	// TODO: implement includeStrikeThrough
	getText(start, end, includeStrikeThrough/*unused*/ = false)
		{
		return .GetRange(start, end)
		}

	SearchText()
		{
		return .GetText(includeStrikeThrough:)
		}

	ScintillaRichEditor_UpdateStyleObject(.styles)
		{
		.getCached = false
		}

	styles: #()
	getCached: false
	Set(s)
		{
		.getCached = s
		parsed = ScintillaRichEditorHelper.Parse(s)
		.superSet(parsed.s)
		.styles = parsed.styles
		if .styles.NotEmpty?()
			.Act(#SetStyleObject, .styles)

		.Dirty?(false)
		}

	// factored out so we can test
	superSet(s)
		{
		super.Set(s)
		}

	AppendText(s)
		{
		.getCached = false
		oldSize = .GetLength()
		super.AppendText(s)
		newSize = .GetLength()
		if .styles.Empty?()
			return
		from = .styles.Last().to.Copy()
		ScintillaRichEditorHelper.NextPos(.GetRange(oldSize, newSize),
			from, to = Object())
		.styles = .styles.Copy()
		.styles.Add(Object(:from, :to,
			style: ScintillaRichEditorHelper.DefaultStyle.Copy()))
		}

	Trim()
		{
		.getCached = false
		s = super.Get()
		if s.Blank?() or .styles.Empty?()
			{
			super.Trim()
			return
			}

		trimmed = ScintillaRichEditorHelper.TrimStyledText(s, .styles)
		if trimmed.s.Size() isnt s.Size()
			{
			super.Trim()
			.Act(#SetStyleObject, .styles = trimmed.styles)
			}
		}

	SetReadOnly(readOnly)
		{
		if .readonly
			return
		super.SetReadOnly(readOnly)
		.setBackground(readOnly is true ? GetSysColor(COLOR.BTNFACE) : 0xffffff)
		}

	KEYDOWN(wParam, pressed = false)
		{
		if super.KEYDOWN(wParam, pressed) is 0
			return 0
		return .Eval(EditorKeyDownHandler, wParam,
			zoomArgs: [this, .zoom, ScintillaRichWordZoomControl],
			:pressed)
		}
	Valid?()
		{
		return ScintillaEditorValid(.Get(), .textLimit)
		}

	ValidData?(@args)
		{
		return ScintillaEditorValid(args[0], args.GetDefault('textLimit', false))
		}

	SetValid(valid? = true)
		{
		if (GetFocus() is .Hwnd)
			valid? = true
		// have to check .readonly as well because if we are in a block running from the
		// ignoring_readonly method (like from Set), then the readonly flag will actually
		// be 0 when this is called resulting in the background color incorrectly
		// switching to white
		.setBackground(.GETREADONLY() is 1 or .readonly
			? GetSysColor(COLOR.BTNFACE)
			: valid? is false ? CLR.ErrorColor : CLR.WHITE)
		}

	SCEN_KILLFOCUS()
		{
		if not .Valid?() and GetFocus() isnt .Hwnd
			{
			.SetValid(false)
			Beep()
			}
		return super.SCEN_KILLFOCUS()
		}

	SCEN_SETFOCUS()
		{
		super.SCEN_SETFOCUS()
		.SetValid() // don't color invalid when focused
		return 0
		}
	}