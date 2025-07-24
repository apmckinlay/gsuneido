// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Title: #Font
	AlertMsg: 'Font will be changed for each new applicable window'
	New(noTitle? = false, name = #logfont, subHeading = '', .fontFaceName = false,
		.fontHeight = false, initFont = false)
		{
		super(Object(#Record, Object(#Border, .Controls(noTitle?, subHeading, name))))
		.fontStatic = .FindControl(#font)
		.Font = initFont isnt false
			? initFont
			: .InitFont()
		.showFont()
		}

	InitFont()
		{ return .DefaultLogFont() }

	DefaultLogFont()
		{
		lfFaceName = .fontFaceName is false ? StdFonts.Mono() : .fontFaceName
		lfHeight = StdFonts.LfSize(.fontHeight)
		fontPtSize = StdFonts.PtSize(lfHeight)
		return Object(lfItalic: 0, :lfHeight, lfOrientation: 0, lfStrikeOut: 0,
			:lfFaceName, lfEscapement: 0, lfUnderline: 0, lfWidth: 0,
			lfCharSet: 0, lfQuality: 1, lfPitchAndFamily: 8, lfClipPrecision: 1,
			lfOutPrecision: 1, lfWeight: 400, :fontPtSize)
		}

	Controls(noTitle?, subHeading, name)
		{
		.Name = name
		ctrl = Object(#Vert)
		if not noTitle?
			ctrl.Add(Object(#TitleNotes, .Title), #Skip)
		if subHeading isnt ''
			ctrl.Add(Object(#Static, subHeading), #Skip)
		return ctrl.Add(
			#(Horz,
				#(Button 'Customize Font' pad: 20),
				#Skip,
				#(Button 'Restore Standard Font' pad: 20)),
			#Skip,
			#(Static '' name: font))
		}

	Getter_Font()
		{ return .Font }

	showFont()
		{
		.fontStatic.SetFont(
			.Font.lfFaceName,
			.Font.GetDefault(#fontPtSize, StdFonts.PtSize(.Font.lfHeight)),
			.Font.GetDefault(#lfWeight, #NORMAL))
		.fontStatic.Set(.FontMsg())
		}

	FontMsg()
		{ return 'Selected font: ' $ .Font.lfFaceName }

	On_Customize_Font()
		{
		cf = Object()
		cf.lStructSize = CHOOSEFONT.Size()
		cf.hwndOwner = .Window.Hwnd
		cf.lpLogFont = .Font.Copy()
		cf.nSizeMin = 9
		cf.nSizeMax = 28
		cf.Flags = CFO.LIMITSIZE | CFO.TTONLY | CFO.SCRIPTSONLY | CFO.INITTOLOGFONTSTRUCT
		if false isnt CenterDialog(.Window.Hwnd, { ChooseFont(cf) })
			.ChangeFont(cf.lpLogFont)
		}

	On_Restore_Standard_Font()
		{ .changeFont(.DefaultLogFont()) }

	ChangeFont(font)
		{
		if .Font isnt font
			.changeFont(font)
		}

	changeFont(font)
		{
		.SetFont(font)
		.showFont()
		.AlertInfo(.Title, .AlertMsg)
		}

	SetFont(.Font)
		{ .Font.fontPtSize = StdFonts.PtSize(.Font.lfHeight) }

	Get()
		{ return .Font }
	}
