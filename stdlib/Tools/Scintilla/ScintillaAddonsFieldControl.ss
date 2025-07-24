// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
ScintillaAddonsEditorControl
	{
	ComponentName: 'ScintillaAddonsField'
	Xmin: 100
	Xstretch: false
	Ystretch: false
	MarginLeft: 6 // to match Field
	DefaultFontSize: 9 // should match Init
	New(@args)
		{
		super(@args.Merge(#(height: 1)))
		.trim = args.GetDefault(#trim, true)
		.SetExtraAscent(2) // to align with Field
		.Top += 2 // adjust for extra ascent
		}
	Addons: (Addon_speller:)

	// allow tabbing through field and ignore return key
	GETDLGCODE(lParam)
		{
		if (false isnt (m = MSG(lParam)) and
			(m.message is WM.CHAR or m.message is WM.KEYDOWN) and
			((m.wParam is VK.TAB or m.wParam is VK.ESCAPE or m.wParam is VK.RETURN) and
				.AutoCActive() is 0))
			return DLGC.WANTCHARS // ie. do NOT want to absorb tab or return
		return 'callsuper'
		}

	EN_CHANGE()
		{
		.Defer(.stripNewLines, uniqueID: 'scintillaStripNewLines')
		super.EN_CHANGE()
		}

	stripNewLines()
		{
		s = .Get()
		if s.Has?('\n')
			.PasteOverAll(s.Tr("\r\n", " "))
		}

	Get()
		{
		s = super.Get()
		return .trim ? s.Trim() : s
		}
	//TODO clear select when focus lost (to match Field)
	}
