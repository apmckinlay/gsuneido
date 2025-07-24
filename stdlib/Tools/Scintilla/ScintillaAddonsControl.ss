// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
// NOTE: This control is designed to handle xmin/ymin or height/width.
//		 It is preferential to use height and width as the control then uses the scintilla
//		 font dimensions to calculate the required space (height being the number of lines
//		 and width being the number of characters per line). Whereas, xmin/ymin will
//		 simply set the control's minimum dimensions. You cannot use a combination of
//		 xmin and width or ymin and height.
ScintillaControl
	{
	Height: 		7		// Determines how many lines before scrolling is required
	Width:			60		// Determines how many characters can fit on a line
	MarginLeft: 	false
	LONG_IDLE: 		3000 	// 3 sec
	CHANGE_IDLE: 	500 	// .5 sec
	addons:			false
	IDE: 			false
	New(@args)
		{
		super(@.processArgs(args))
		.Map[SCN.UPDATEUI] = 'UPDATEUI'
		.Map[SCN.PAINTED] = 'PAINTED'
		.Map[SCN.STYLENEEDED] = 'STYLENEEDED'
		.Map[SCN.CHARADDED] = 'CHARADDED'
		.Map[SCN.MARGINCLICK] = 'SCN_MARGINCLICK'
		.Map[SCN.AUTOCSELECTION] = 'SCN_AUTOCSELECTION'

		.styleManager = ScintillaAddonsLineStyles(this)
		.scheme = args.GetDefault(#scheme, '')

		.setupAddons(args)
		.setupContextMenu()
		.setupFont(args)

		.SetModEventmask(SC.MOD_INSERTTEXT | SC.MOD_DELETETEXT | SC.MOD_BEFOREDELETE)
		}

	processArgs(args)
		{
		.commandManager = new ScintillaAddonsCommandManager
		.IDE = args.GetDefault(#IDE, .IDE)
		return args
		}

	setupAddons(args)
		{
		.addons = AddonManager(this, args)
		.styleManager.DefineStyles(.addons)
		.addons.Send(#Init)
		if args.GetDefault(#readonly, false) is true and
			args.GetDefault(#readonlyEnableIdle, false) isnt true
			.longIdle = .changeIdle = class { Reset(){} Kill(){} }
		else
			{
			.longIdle = IdleTimer(.LONG_IDLE, { .addons.Send(#LongIdle) }).Reset()
			.changeIdle = IdleTimer(.CHANGE_IDLE, { .addons.Send(#IdleAfterChange) })
			}
		}

	InitMarkers()
		{
		// ScintillaAddonsControl uses ScintillaAddonsLineStyles to setup its marks
		}

	setupFont(args)
		{
		fontArgs = .ideEnv?() ? .sciArgs() : args
		.setFont(
			fontArgs.GetDefault(#font, StdFonts.Mono()),
			fontArgs.GetDefault(#fontSize, .DefaultFontSize),
			fontArgs.GetDefault(#weight, FW.NORMAL),
			fontArgs.GetDefault(#italic, false))

		if .MarginLeft isnt false
			.SetMarginLeft(0, .MarginLeft)

		.calcYmin(args.GetDefault(#height, false), args.GetDefault(#ymin, false))
		.calcXmin(args.GetDefault(#width,  false), args.GetDefault(#xmin, false))
		}

	InitFont()
		{
		if .ideEnv?()
			.setFont(@.sciArgs())
		else
			super.InitFont()
		}

	ideEnv?()
		{ return .IDE and '' isnt IDESettings.Get(#ide_scifont) }

	sciArgs()
		{
		ide_scifont = IDESettings.Get(#ide_scifont)
		return Object(
			font: ide_scifont.lfFaceName,
			fontSize: .fontSize(ide_scifont),
			weight: ide_scifont.lfWeight,
			italic: ide_scifont.lfItalic is -1)
		}

	fontSize(font)
		{ return font.Member?(#fontPtSize) ? font.fontPtSize : -1 * font.lfHeight }

	setFont(font, fontSize, weight, italic)
		{
		.SendMessageTextIn(SCI.STYLESETFONT, SC.STYLE_DEFAULT, StdFonts.Font(font))
		.StyleSetSize(SC.STYLE_DEFAULT, .sciSize(fontSize))
		// Scintilla supports a font weight of 1 - 999, treat 0 and '' as FW.NORMAL (400)
		.StyleSetWeight(SC.STYLE_DEFAULT,
			StdFonts.Weight(weight in (0, '') ? FW.NORMAL : weight))
		.StyleSetItalic(SC.STYLE_DEFAULT, italic)
		}

	sciSize(fontSize)
		{
		fontOb = .ideEnv?() ? IDESettings.Get(#ide_scifont) : Suneido.logfont
		return StdFonts.SciSize(fontSize, .fontSize(fontOb))
		}

	scheme: ''
	GetSchemeColor(colorName)
		{
		return not Object?(.scheme)
			? IDE_ColorScheme.GetColor(colorName, .scheme)
			: .scheme.Member?(colorName)
				? .scheme[colorName]
				: IDE_ColorScheme.DefaultStyle.GetDefault(colorName, false)
		}

	SCEN_SETFOCUS()
		{
		.addons.Send('SetFocus')
		return super.SCEN_SETFOCUS()
		}

	MOUSEWHEEL(wParam, lParam)
		{
		zoom = .addons.Collect('Scroll_Zoom')
		if Object?(zoom) and zoom.Size() is 1 and zoom[0] is true
			return 'callsuper'

		return super.MOUSEWHEEL(wParam, lParam)
		}

	setupContextMenu()
		{
		.Context_Menu = ["&Undo\tCtrl+Z", "&Redo\tCtrl+Y", "", "Cu&t\tCtrl+X",
			"&Copy\tCtrl+C", "&Paste\tCtrl+V", "&Delete", ""]
		misc = ["Select &All\tCtrl+A", "Find...\tCtrl+F"]
		addonMenu = .addons.Collect(#ContextMenu)
		addonMenu.Filter({ it.Size() is 1 }).Each({ misc.Add( it[0]) })
		.Context_Menu.Add(@misc.SortWith!({ it.Replace('\.\.\.') }))
		addonMenu.Filter({ it.Size() > 1 }).Each({ .Context_Menu.Add("").Add(@it) })
		.commandManager.Set(.Context_Menu)
		}

	calcXmin(width, xmin)
		{
		if width isnt false and xmin isnt false
			ProgrammerError(`Cannot utilize xmin and width at the same time. `  $
				`Please verify the controls arguments and use either width or xmin only`,
				params: [:width, :xmin, name: this.GetDefault(#Name, false)])

		.Xmin = xmin isnt false
			? ScaleWithDpiFactor(xmin)
			: .SendMessageTextIn(SCI.TEXTWIDTH, SC.STYLE_DEFAULT,
				'M'.Repeat(width is false ? .Width : width))
		}

	calcYmin(height, ymin)
		{
		if height isnt false and ymin isnt false
			ProgrammerError(`Cannot utilize ymin and height at the same time. `  $
				`Please verify the controls arguments and use either height or ymin only`,
				params: [:height, :ymin, name: this.GetDefault(#Name, false)])
		borderPadding = 2
		topAdjustment = 1
		if ymin isnt false
			.Ymin = ScaleWithDpiFactor(ymin)
		else
			{
			txtHeight = .TextHeight(0)
			.Ymin = txtHeight * (height is false ? .Height : height) + borderPadding * 2
			topAdjustment = Max(txtHeight / 7 /*= % percent to reduce*/, 1)
			}
		.Top = (.Ymin / Max(height, 1) - borderPadding * topAdjustment).Round(0)
		}

	Addon(@args)
		{
		.addons.Addon(@args)
		}
	ContextMenu(x, y)
		{
		i = ContextMenu(.Context_Menu).ShowCall(this, x, y)
		if i > 0
			.addons.Send("ContextMenuChoice", i - 1)
		return 0
		}
	Default(@args) // used by context menu
		{
		if args[0].Prefix?('On_') and .addons isnt false
			{
			args[0] = args[0].Replace("^On_Context_", "On_")
			return .addons.Send(@args)
			}
		else
			return super.Default(@args)
		}
	Recv(@args) // used by redirected accelerator keys
		{
		if .addons is false
			return false
		if args[0].Prefix?('On_')
			{
			args[0] = args[0].Replace("^On_Context_", "On_")
			return .addons.Send(@args)
			}
		}

	KEYDOWN(wParam)
		{
		command = .commandManager.BuildCommandOb(VK.Find(wParam))
		if false isnt method = .commandManager.GetMethod(command)
			{
			this[method]()
			return 0
			}
		return 'callsuper'
		}

	Set(text)
		{
		super.Set(text)
		.addons.Send(#Set)
		}
	UPDATEUI(lParam = false)
		{
		if lParam isnt false
			{
			scn = SCNotification(lParam)
			// .GetDefault is needed because gSuneido built before 09/09/2021 doesn't
			// put #updated value to the SCNotification object
			updated = scn.GetDefault(#updated, 0)
			if 0 isnt (updated & SC.UPDATE_V_SCROLL)
				.addons.Send(#Scintilla_VScroll)
			if 0 isnt (updated & SC.UPDATE_SELECTION)
				.addons.Send(#Scintilla_Selection)
			}
		.addons.Send(#UpdateUI)
		return 0
		}
	PAINTED()
		{
		.addons.Send(#Painted)
		return 0
		}
	SetReadOnly(readonly = true)
		{
		super.SetReadOnly(readonly)
		.UPDATEUI()
		return // no return value
		}
	CHARADDED(lParam)
		{
		scn = SCNotification(lParam)
		.addons.Send(#CharAdded, scn.ch.Chr())
		return 0
		}
	CHAR(wParam, lParam)
		{
		// disable entering ascii control characters (CTRL + Alphabet) in scintilla
		if wParam < 32 and KeyPressed?(VK.CONTROL) /*=smallest printable char code*/
			return 0

		if wParam is VK.RETURN and .Send('Enter_Pressed') is false // disabling enter key
			return 0

		.Callsuper(.Hwnd, WM.CHAR, wParam, lParam)
		if wParam is VK.BACK
			.addons.Send(#Backspace)

		// should be able to use: .AssignCmdKey(SCK.ESCAPE, SCI.CANCEL)
		// but couldn't get it to work ???
		if wParam is VK.ESCAPE
			.CANCEL()

		return 0
		}
	GETDLGCODE(wParam, lParam)
		{
		if wParam is VK.ESCAPE and not .Window.Base?(Dialog)
			return DLGC.WANTALLKEYS
		return super.GETDLGCODE(wParam, lParam)
		}
	longIdle: false
	EN_CHANGE()
		{
		.ResetTimers()
		.Send('EditorChange')
		return super.EN_CHANGE()
		}
	ResetTimers()
		{
		.longIdle.Reset()
		.changeIdle.Reset()
		}
	SCN_MODIFIED(lParam)
		{
		scn = SCNotification(lParam)
		if 0 isnt (scn.modificationType & SC.MOD_BEFOREDELETE)
			.addons.Send(#BeforeDelete, scn.position, scn.length)
		if 0 isnt (scn.modificationType &
			(SC.MOD_DELETETEXT | SC.MOD_INSERTTEXT))
			.addons.Send(#Modified, scn)
		return super.SCN_MODIFIED(lParam)
		}
	SCN_MARGINCLICK(lParam)
		{
		scn = SCNotification(lParam)
		.addons.Send(#MarginClick, scn)
		return 0
		}
	SCN_DOUBLECLICK()
		{
		.addons.Send(#DoubleClick)
		return super.SCN_DOUBLECLICK()
		}
	SCN_AUTOCSELECTION(lParam)
		{
		scn = SCNotificationText(lParam)
		.addons.Send(#AutocSelection, scn)
		return 0
		}
	SCN_ZOOM()
		{
		super.SCN_ZOOM()
		.addons.Send(#Scintilla_Zoom)
		return 0
		}
	STYLENEEDED(lParam)
		{
		.StyleTo(SCNotification(lParam).position)
		return 0
		}
	StyleTo(pos)
		{
		.addons.Send(#Style, .GetEndStyled(), pos)
		}
	StyleToEnd()
		{
		.StyleTo(.GetTextLength())
		}

	// START: Redirects to ScintillaAddonsLineStyles
	ForEachMarkerByLevel(type, block)
		{ .styleManager.ForEachMarkerByLevel(type, block) }

	GetMarkerTypes()
		{ return .styleManager.MarkerTypes }

	GetMarkerColor(i)
		{ return .styleManager.MarkerColor(i) }

	MarkerIdx(level, type = 'default')
		{ return .styleManager.MarkerIdx(level, type) }

	IndicatorIdx(level)
		{ return .styleManager.IndicatorIdx(level) }

	IndicatorAtPos?(pos)
		{ return .styleManager.IndicatorAtPos?(pos) }
	// END: Redirects to ScintillaAddonsLineStyles

	MakeSummary()
		{
		text = .Get()
		summary = text.BeforeFirst('\n').Trim()[.. 60] /*= summary size*/
		if text.Size() > summary.Size()
			summary $= '...'
		return summary
		}

	MarkersChanged()
		{
		super.MarkersChanged()
		.addons.Send(#Scintilla_MarkersChanged)
		}

	SendToAddons(@args)
		{
		.addons.Send(@args)
		}

	CollectFromAddons(@args)
		{
		return .addons.Collect(@args)
		}

	ConditionalSendToAddons(block, args)
		{
		.addons.ConditionalSend(block, args)
		}

	Valid?()
		{
		return not .addons.Collect(#Valid?).Has?(false)
		}

	ScrollToBottom(noFocus? = false)
		{
		lines = .GetLineCount()
		.GotoLine(lines, :noFocus?)
		}

	GetCurrentWord()
		{
		org = end = .GetCurrentPos()
		while .wordChars.Has?(.GetAt(org - 1))
			--org
		while .wordChars.Has?(.GetAt(end))
			++end
		return org < end ? .GetRange(org, end) : ""
		}
	getter_wordChars()
		{
		return .wordChars = .GetWordChars() // once only
		}

	On_Copy()
		{
		super.On_Copy()
		.addons.Send(#On_Copy)
		}

	ResetAddons()
		{
		.markerCount = 1	// Set to 1 because of the Marker inside of this class
		.addons.Send(#Init)
		}

	Destroy()
		{
		.commandManager.Destroy()
		.longIdle.Kill()
		.changeIdle.Kill()
		.addons.Send(#Destroy)
		super.Destroy()
		}
	}
