// Copyright (C) 2020 Axon Development Corporation All rights reserved worldwide.
ScintillaControl
	{
	Height: 		7		// = Determines how many lines before scrolling is required
	Width:			60		// = Determines how many characters can fit on a line
	LONG_IDLE: 		3000 	// 3 sec
	CHANGE_IDLE: 	500 	// .5 sec
	addons:			false
	IDE: false
	ComponentName:	"ScintillaAddons"
	New(@args)
		{
		super(@.processArgs(args))

		.styleManager = ScintillaAddonsLineStyles(this)
		.scheme = args.GetDefault(#scheme, '')

		.setupAddons(args)
		.setupContextMenu()
		.setupFont(args)

		.ComponentArgs.Add(
			args.GetDefault(#width, .Width/*=default width*/),
			args.GetDefault(#xmin, false),
			args.GetDefault(#ymin, false))
		}

	processArgs(args)
		{
		.commandManager = new ScintillaAddonsCommandManager
		.IDE = args.GetDefault(#IDE, .IDE)
		return args
		}

	setupAddons(args)
		{
		.addons = AddonManager(this, args.Copy().DeleteIf({ not .SupportedAddon?(it) }))
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

	supportedAddons: #(Addon_speller, Addon_url, Addon_zoom,
		Addon_highlight_cursor_line,
		Addon_multiple_selection,
		Addon_suneido_style,
		Addon_status, Addon_auto_complete_code, Addon_auto_complete_queries)
	SupportedAddon?(addon)
		{
		return .supportedAddons.Has?(addon)
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
		.Act('SetAddonCommands', .commandManager.GetCommands())
		}

	setupFont(args)
		{
		fontArgs = args
		.SetFont(
			fontArgs.GetDefault(#font, StdFonts.Mono()),
			fontArgs.GetDefault(#fontSize, .DefaultFontSize),
			fontArgs.GetDefault(#weight, FW.NORMAL),
			fontArgs.GetDefault(#italic, false))
		}

	scheme: ''
	GetSchemeColor(colorName)
		{
return IDE_ColorScheme.DefaultStyle.GetDefault(colorName, false)
//		return not Object?(.scheme)
//			? IDE_ColorScheme.GetColor(colorName, .scheme)
//			: .scheme.Member?(colorName)
//				? .scheme[colorName]
//				: IDE_ColorScheme.DefaultStyle.GetDefault(colorName, false)
		}

	Set(text)
		{
		super.Set(text)
		.addons.Send(#Set)
		}

	EN_CHANGE()
		{
		.ResetTimers()
		.Send('EditorChange')
		return super.EN_CHANGE()
		}

	Enter_Pressed(pressed = false)
		{
		.Send('Enter_Pressed', :pressed)
		}

	Backspace_Pressed()
		{
		.addons.Send(#Backspace)
		}

	CHARADDED(c)
		{
		.addons.Send(#CharAdded, c)
		}

	AutocSelection(s)
		{
		if false is .addons.Send(#AutocSelection, [text: s])
			super.AutocSelection(s)
		}

	KEYDOWN(wParam, pressed = false)
		{
		command = .commandManager.BuildCommandOb(VK.Find(wParam), :pressed)
		if false isnt method = .commandManager.GetMethod(command)
			{
			this[method]()
			return 0
			}
		return 'callsuper'
		}

	Scintilla_SetValue()
		{
		.ResetTimers()
		.Send('EditorChange')
		}

	UPDATEUI()
		{
		.addons.Send(#UpdateUI)
		}

	ContextMenu(x, y)
		{
		i = ContextMenu(.Context_Menu).ShowCall(this, x, y)
		if i > 0
			.addons.Send("ContextMenuChoice", i - 1)
		if i is false or i <= 0
			.EnsureSelect()
		return 0
		}

	ResetTimers()
		{
		.longIdle.Reset()
		.changeIdle.Reset()
		}

	SCN_MODIFIED(lParam)
		{
		scn = SCNotification(lParam)
		if 0 isnt (scn.modificationType &
			(SC.MOD_DELETETEXT | SC.MOD_INSERTTEXT))
			.addons.Send(#Modified, scn)
		return super.SCN_MODIFIED(lParam)
		}

	SCN_DOUBLECLICK()
		{
		.addons.Send(#DoubleClick)
		return super.SCN_DOUBLECLICK()
		}

	MarkersChanged()
		{
		if .Destroyed?()
			return
		super.MarkersChanged()
		.addons.Send(#Scintilla_MarkersChanged)
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

	IndicatorAtPos?(pos /*unused*/)
		{ return false }
//		{ return .styleManager.IndicatorAtPos?(pos) }
	// END: Redirects to ScintillaAddonsLineStyles

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

	MakeSummary()
		{
		text = .Get()
		summary = text.BeforeFirst('\n').Trim()[.. 60] /*= summary size*/
		if text.Size() > summary.Size()
			summary $= '...'
		return summary
		}

	SendToAddons(@args)
		{
		.addons.Send(@args)
		}

	ConditionalSendToAddons(block, args)
		{
		.addons.ConditionalSend(block, args)
		}

	CollectFromAddons(@args)
		{
		return .addons.Collect(@args)
		}

	Valid?()
		{
		return not .addons.Collect(#Valid?).Has?(false)
		}

	ScrollToBottom(noFocus? = false)
		{
		.Act(#ScrollToBottom, noFocus?)
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
