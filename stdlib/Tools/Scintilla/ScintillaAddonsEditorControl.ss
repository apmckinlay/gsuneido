// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
ScintillaAddonsEditorBaseControl
	{
	Width: 10 // should be consistent with EditorControl (100 xmin at 12 ~ 10 characters)
	Height: 2 // should be consistent with EditorControl (50  ymin at 12 ~ 2 lines)

	New(@args)
		{
		super(@.processArgs(args))
		.Map[SCEN.SETFOCUS] = 'SCEN_SETFOCUS'
		.mandatory = args.GetDefault(#mandatory, false)
		if .zoom and not args.GetDefault('readonly', false)
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

	Valid?()
		{
		return .validCheck?(.Get(), .mandatory)
		}
	validCheck?(data, mandatory)
		{
		if data.Size() > EditorTextLimit
			return false
		return (mandatory and data is "") ? false : true
		}

	ValidData?(@args)
		{
		return .validCheck?(args[0], args.GetDefault('mandatory', false))
		}

	On_Print()
		{
		if '' is (text = .Get().Trim())
			return
		Params.On_Print(Object('WrapGen', text),
			title: '', name: 'print_editor', previewWindow: .Window.Hwnd)
		}
	}
