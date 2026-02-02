// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
WndProc
	{
	Name: "PreviewPage"
	New(report, scale = 1, dimens = false, pdc = 0)
		{
		.CreateWindow("SuWhitePush", "", WS.VISIBLE)
		.SubClass()

		.report = Report(@report)
		.report.SetDriver(new GdiPreviewDriver(:pdc))

		.dimens = dimens is false ? .report.GetDimens() : dimens
		.origWidth = .dimens.width * .ctrlScaleFactor
		.origHeight = .dimens.height * .ctrlScaleFactor
		.pages = Object()
		.pg = 0; // the page currently displayed
		.npages = 0; // the number of pages generated
		.NextPage()
		.Scale(scale)
		.Defer({ if .Member?('Hwnd') SetFocus(.Hwnd) })
		.linePen = CreatePen(PS.DOT, 1, CLR.SpotBlue)
		.lineBrush = GetStockObject(SO.NULL_BRUSH)
		}
	scale: 1
	ctrlScaleFactor: 80
	minZoom: .25
	maxZoom: 4
	Scale(by) // 2 = double, .5 = half
		{
		x = .scale * by
		if x > .minZoom and x < .maxZoom
			.scale = x
		.refreshZoom()
		}
	refreshZoom()
		{
		.Xmin = .origWidth * .scale
		.Ymin = .origHeight * .scale
		InvalidateRect(.Hwnd, NULL, true)
		}
	ResetScale()
		{
		.scale = 1
		.refreshZoom()
		}
	GetScale()
		{
		return .scale
		}
	FirstPage()
		{
		.pg = 0
		InvalidateRect(.Hwnd, NULL, true)
		}
	LastPage()
		{
		while (.nextpage())
			{ }
		.pg = Max(0, .npages - 1)
		InvalidateRect(.Hwnd, NULL, true)
		return true
		}
	NextPage()
		{
		if (.pg + 1 < .npages)
			++.pg
		else if (not .nextpage())
			{
			Beep()
			return false
			}
		InvalidateRect(.Hwnd, NULL, true)
		.sortAccessPoints()
		return true
		}
	max_rect_height: false
	sortAccessPoints()
		{
		if .noAccessPointsOnPage?(.pg)
			return
		.report.AccessPoints[.pg].Sort!(
			{|a, b|	a.bottom <= b.bottom and a.left <= b.left })
		.max_rect_height = 0
		for point in .report.AccessPoints[.pg]
			{
			height = point.bottom - point.top
			if height > .max_rect_height
				.max_rect_height = height
			}
		}
	Empty?()
		{
		return .npages is 0
		}
	eof: false
	status: ''
	nextpage()
		{
		if (.eof)
			return false
		filename = .report.AddPage(.dimens)
		if not .report.NextPageSuccess?(vbox = .report.NextPage())
			{
			.eof = true
			DeleteEnhMetaFile(.report.EndPage())
			DeleteFile(filename)
			if Report.IsGeneratingReportError?(vbox)
				{
				if .pg isnt 0 // first page handled by PreviewControl
					.report.DisplayAlert(vbox)
				.status = vbox
				}
			return false
			}
		.report.Paint(vbox)
		.set_emf(.report.EndPage())

		.emfpg = .pg = .npages++
		.pages[.pg] = filename
		return true
		}
	GetStatus()
		{
		return .status
		}
	PrevPage()
		{
		if (.pg is 0)
			{
			Beep()
			return false
			}
		else
			{
			--.pg
			InvalidateRect(.Hwnd, NULL, true)
			return true
			}
		}
	GetPageNum()
		{
		return .pg
		}

	GetNumPages()
		{
		return .npages
		}

	emfpg: false
	PAINT()
		{
		hdc = BeginPaint(.Hwnd, ps = Object())
		if (not .pages.Member?(.pg))
			{
			EndPaint(.Hwnd, ps)
			return 0
			}
		if .emfpg isnt .pg
			{
			.set_emf(GetEnhMetaFile(.pages[.pg]))
			.emfpg = .pg
			}

		w = .origWidth * .maxZoom
		h = .origHeight * .maxZoom

		WithCompatibleDC(hdc, w, h)
			{|hdcBmp|
			r = Object(right: w, bottom: h)
			brush = CreateSolidBrush(CLR.WHITE)
			FillRect(hdcBmp, r, brush)
			DeleteObject(brush)

			PlayEnhMetaFile(hdcBmp, .emf, r)
			SetStretchBltMode(hdc, STRETCH.HALFTONE)
			StretchBlt(hdc,
				0, 0, .Xmin, .Ymin,
				hdcBmp, 0, 0, w, h,
				ROP.SRCCOPY)
			}
		.paintHotspot(hdc)
		EndPaint(.Hwnd, ps)
		return 0
		}
	lineRect: false
	linePen: false
	lineBrush: false
	paintHotspot(hdc)
		{
		if .lineRect isnt false
			DoWithHdcObjects(hdc, [.linePen, .lineBrush])
				{
				Rectangle(hdc,
					.lineRect.left, .lineRect.top,
					.lineRect.right, .lineRect.bottom)
				}
		}
	emf: false
	set_emf(emf)
		{
		if .emf isnt false
			DeleteEnhMetaFile(.emf)
		.emf = emf
		}

	// scroll by dragging
	dragging: false
	LBUTTONDOWN(lParam)
		{
		SetFocus(.Hwnd) // to get mousewheel
		.dragging = true
		.dragx = LOSWORD(lParam)
		.dragy = HISWORD(lParam)
		SetCapture(.Hwnd)
		return 0
		}
	LBUTTONDBLCLK(lParam)
		{
		.locate_target(.pg, LOSWORD(lParam), HISWORD(lParam))
		return 0
		}
	noAccessPointsOnPage?(pg)
		{
		return not .report.Member?('AccessPoints') or not .report.AccessPoints.Member?(pg)
		}
	locate_target(page, x, y, find? = false)
		{
		if .noAccessPointsOnPage?(page)
			return

		x = x / .scale
		y = y / .scale
		rect = false

		page_points = .report.AccessPoints[page]
		low_bound = page_points.BinarySearch(Object(bottom: y),
			{|a, b| a.bottom < b.bottom})
		for (i = low_bound; i < page_points.Size(); i++)
			{
			target = page_points[i]
			if POINTinRECT(target, [:x, :y])
				{
				if find? is true
					{
					rect = .scaleTarget(target)
					SetCursor(LoadCursor(ResourceModule(), IDC.HAND))
					}
				else
					.access(target)
				break
				}
			else if y + .max_rect_height < target.top
				break
			}
		.repaintLineRect(rect)
		}
	access(target)
		{
		access = target.access
		if access.control is 'AccessGoTo'
			AccessGoTo(access.access, access.goto_field, access.goto_value, .Hwnd)
		else if access.control is 'ReportGoTo'
			ReportGoTo(access.report, access.params, .Hwnd)
		else if access.control is 'AttachmentGoTo'
			AttachmentGoTo(access.file, .Hwnd)
		else if access.control is 'DynamicGoTo'
			DynamicGoTo(access.data, access.func, .Hwnd)
		}
	scaleTarget(target)
		{
		return Object(left: target.left * .scale, top: target.top* .scale,
			right: target.right * .scale, bottom: target.bottom * .scale)
		}
	repaintLineRect(rect)
		{
		if rect is false
			{
			if .lineRect isnt false
				.repaintRect(.lineRect)
			.lineRect = false
			}
		else if rect isnt .lineRect
			{
			if .lineRect isnt false
				.repaintRect(.lineRect)
			.lineRect = rect
			.repaintRect(.lineRect)
			}
		}
	rect_border: 2
	repaintRect(rc)
		{
		rc = rc.Copy()
		rc.left = rc.left - .rect_border
		rc.right = rc.right + .rect_border
		rc.top = rc.top - .rect_border
		rc.bottom = rc.bottom + .rect_border
		InvalidateRect(.Hwnd, rc, true)
		}
	last_x: 0
	last_y: 0
	MOUSEMOVE(lParam)
		{
		x = LOSWORD(lParam)
		y = HISWORD(lParam)
		if (not .dragging)
			{
			if .last_x isnt x or .last_y isnt y
				{
				if .report.Member?('AccessPoints') and .report.AccessPoints.Member?(.pg)
					.locate_target(.pg, x, y, find?:)
				.last_x = x
				.last_y = y
				}
			return 0
			}
		if (x is .dragx and y is .dragy)
			return 0
		.Send("Scroll", x - .dragx, y - .dragy)
		return 0
		}
	LBUTTONUP()
		{
		ReleaseCapture()
		.dragging = false
		return 0
		}

	MOUSEWHEEL(wParam)
		{
		scroll = GetWheelScrollInfo(wParam)
		if KeyPressed?(VK.CONTROL)
			{
			.Send('Zoom', scroll.clicks > 0 ? 'In' : 'Out')
			return 0
			}
		amount = scroll.page? // mousewheel Set to Scroll Page at Time
			? scroll.clicks.Sign() * .getClientHeight()
			: scroll.lines * 20  /*= speed up scrolling */
		.Send("Scroll", 0, amount)
		return 0
		}
	getClientHeight()
		{
		GetClientRect(.Hwnd, r = Object())
		return r.bottom - r.top
		}

	DESTROY()
		{
		if .linePen isnt false
			DeleteObject(.linePen)
		if .lineBrush isnt false
			DeleteObject(.lineBrush)
		.set_emf(false)
		for file in .pages
			if true isnt deleteResult = DeleteFile(file)
				SuneidoLog('ERRATIC: Could not clean up report preview temp file',
					params: Object(:deleteResult, :file), calls:)
		.report.Close(false)
		return 0
		}
	}
