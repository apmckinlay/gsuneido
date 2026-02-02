// Copyright (C) 2019 Axon Development Corporation All rights reserved worldwide.
EditControl
	{
	Name: 			"Editor"
	ComponentName:	"Editor"
	Unsortable: true
	Hasfocus?:	false
	DefaultHeight: 4

	New(style = 0, .readonly = false, .font = "", .size = "", .zoom = false,
		mandatory = false, set = "", height =  false, .tabthrough = false,
		hidden = false, tabover = false, width = false, weight = false,
		readOnlyBgndColor = false, status = '')
		{
		super(mandatory, readonly, style, :hidden, :tabover, :font, :size, :weight,
			:width, :height, :readOnlyBgndColor, :status)
		.Set(set)
		.findreplacedata = Record()
		.AddContextMenuItem("Find...\tCtrl+F", .On_Find)
		.AddContextMenuItem("Print...\tCtrl+P", .On_Print)
		if .zoom is false
			.AddContextMenuItem("Zoom...\tF6", .On_Zoom)
		}

	EN_KILLFOCUS()
		{
		super.EN_KILLFOCUS()
		if .Dirty?()
			.Send("NewValue", .Get())
		}

	KEYDOWN(wParam, pressed = false)
		{
		return .Eval(EditorKeyDownHandler, wParam, zoomArgs: .zoomArgs(), :pressed)
		}

	zoomArgs()
		{
		return [this, .zoom, font: .font, size: .size]
		}

	On_Print()
		{
		if '' is (text = .Get().Trim())
			return
		Params.On_Print(Object('WrapGen', text),
			title: '', name: 'print_editor', previewWindow: .Window.Hwnd)
		}

	On_Find()
		{
		s = .GetSelText()
		if s > "" and not s.Has?('\n')
			.findreplacedata.find = s
		_hwnd = .WindowHwnd()
		x = FindDialog(.findreplacedata)
		if x is #next
			.On_Find_Next()
		else if x is #prev
			.On_Find_Previous()
		}

	On_Find_Next()
		{
		return .findNextPrev()
		}

	On_Find_Previous()
		{
		return .findNextPrev(prev:)
		}

	findNextPrev(prev = false)
		{
		if .findreplacedata.find.Blank?()
			return false
		from = .GetSel()[prev ? 0 : 1]
		if false is match = Find.DoFind(super.Get(), from, .findreplacedata, :prev)
			return false
		.SetSel(match[0], match[0] + match[1])
		return true
		}

	On_Zoom()
		{
		EditorZoom(@.zoomArgs())
		}

	SetFontAndSize(@args)
		{
		args.Add('SetFontAndSize', at: 0)
		.Act(@args)
		}

	MakeSummary()
		{
		if .GetHidden() is true
			return ''

		text = .Get().Trim()
		summary = text.BeforeFirst('\n').Trim()[::60/*=summary length*/]
		if text.Size() > summary.Size()
			summary $= '...'
		return summary
		}
	}
