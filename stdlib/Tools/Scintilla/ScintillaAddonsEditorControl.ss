// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
ScintillaAddonsControl
	{
	Width: 10 // should be consistent with EditorControl (100 xmin at 12 ~ 10 characters)
	Height: 2 // should be consistent with EditorControl (50  ymin at 12 ~ 2 lines)

	New(@args)
		{
		super(@.processArgs(args))
		.Map[SCEN.SETFOCUS] = 'SCEN_SETFOCUS'
		.mandatory = args.GetDefault(#mandatory, false)
		if .readonly = args.GetDefault('readonly', false)
			.setBackground(GetSysColor(COLOR.BTNFACE))
		if .zoom and not .readonly
			.SetReadOnly(false)
		}

	Addons: (Addon_speller:, Addon_url:)
	processArgs(args)
		{
		.zoom = args.GetDefault('zoom', false)
		addons = .Addons.Copy()
		addons.Addon_zoom = [zoom: .zoom, zoom_ctrl: ScintillaZoomControl]
		defaults = [
			wrap:,
			margin: 	0,
			exStyle: 	WS_EX.STATICEDGE,
			font: 		Suneido.logfont.lfFaceName,
			fontSize: 	'+0',
			weight: 	Suneido.logfont.GetDefault(#lfWeight, FW.NORMAL),
			italic: 	Suneido.logfont.GetDefault(#lfItalic, 0) is -1
			]
		if not args.Member?(#ymin) // Cannot specify both ymin and height
			defaults.height = .Height
		return args.MergeNew(defaults).MergeNew(addons)
		}

	KEYDOWN(wParam, pressed = false)
		{
		if super.KEYDOWN(wParam, :pressed) is 0
			return 0
		return .Eval(EditorKeyDownHandler, wParam,
			zoomArgs: [this, .zoom, ScintillaZoomControl],
			:pressed)
		}

	ZoomReadonly(value)
		{
		ScintillaZoomControl(0, value, readonly:)
		}

	Hasfocus?: false
	HasFocus?()
		{
		return .Hasfocus? or super.HasFocus?()
		}

	SetReadOnly(readOnly)
		{
		if .readonly
			return

		super.SetReadOnly(readOnly)
		.setBackground(readOnly is true ? GetSysColor(COLOR.BTNFACE) : CLR.WHITE)
		}
	setBackground(bgnd)
		{
		.StyleSetBack(0, bgnd)
		.StyleSetBack(SC.STYLE_DEFAULT, bgnd)
		}
	Valid?()
		{
		return .validCheck?(.Get(), .mandatory)
		}
	validCheck?(data, mandatory)
		{
		return (mandatory and data is "") ? false : true
		}

	ValidData?(@args)
		{
		return .validCheck?(args[0], args.GetDefault('mandatory', false))
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
		if (.Send("Dialog?") isnt true and not .Valid?() and GetFocus() isnt .Hwnd)
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
	On_Print()
		{
		if '' is (text = .Get().Trim())
			return
		Params.On_Print(Object('WrapGen', text),
			title: '', name: 'print_editor', previewWindow: .Window.Hwnd)
		}
	}
