// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass()
		{
		.initFont()
		.registerClasses()
		InitCommonControlsEx(Object(dwSize: INITCOMMONCONTROLSEX.Size(),
			dwICC: ICC.WIN95_CLASSES | ICC.DATE_CLASSES))
		EnsureWebBrowserVersion()
		}

	initFont()
		{
		WithDC(NULL)
			{|hdc|
			Suneido.logfont = Suneido.stdfont =
				Object(lfFaceName: StdFonts.Ui(),
					lfHeight: (-StdFontsSize.DefaultSize *
						GetDeviceCaps(hdc, GDC.LOGPIXELSY) / PointsPerInch).Round(0)).
				Set_readonly()
			}
		Suneido.hfont = CreateFontIndirect(Suneido.logfont)
		}

	registerClasses()
		{
		defWndProc = GetDefWindowProc()

		// button face background, arrow cursor
		.registerClass("SuBtnfaceArrow", defWndProc)

		// button face background, arrow cursor, no double clicks
		.registerClass("SuBtnfaceArrowNoDblClks", defWndProc,
			style: CS.REDRAW)

		// button face background, custom cursor
		.registerClass("SuBtnfaceNocursor", defWndProc,
			cursor: NULL)

		// white background, pushing hand cursor
		.registerClass("SuWhitePush", defWndProc,
			cursor: LoadCursor(ResourceModule(), IDC.PUSH),
			background: GetStockObject(SO.WHITE_BRUSH))

		// white background, arrow cursor
		.registerClass("SuWhiteArrow", defWndProc,
			background: GetStockObject(SO.WHITE_BRUSH))

		// NO background, pointing hand cursor
		.registerClass("SuNobgndHand", defWndProc,
			background: NULL,
			cursor: LoadCursor(ResourceModule(), IDC.HAND))

		// tooltip background, arrow cursor
		.registerClass("SuToolArrow", defWndProc,
			background: GetSysColorBrush(COLOR.INFOBK))

		// gray background, arrow cursor
		.registerClass("SuGrayArrow", defWndProc,
			background: GetStockObject(SO.LTGRAY_BRUSH))

		// null background, and explicitly does not redraw itself on resize
		.registerClass("SuLightweight", defWndProc,
			style: CS.DBLCLKS,
			background: GetStockObject(SO.NULL_BRUSH))

		// only for messages
		.registerClass("SuMessageOnly", defWndProc,
			style: 0, background: 0, cursor: 0)
		}

	registerClass(className, defWndProc,
		background = false, style = false, cursor = false)
		{
		if style is false
			style = CS.DBLCLKS | CS.REDRAW
		if background is false
			background = 1 + COLOR.BTNFACE
		if cursor is false
			cursor = LoadCursor(NULL, IDC.ARROW)
		RegisterClass(Object(
			:className,
			wndProc: defWndProc,
			:style,
			instance: Instance(),
			icon: LoadIcon(ResourceModule(), IDI.SUNEIDO),
			:cursor,
			:background
			))
		}
	}