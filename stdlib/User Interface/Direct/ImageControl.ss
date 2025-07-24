// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
WndProc
	{
	Name: 'Image'
	Unsortable: true
	Xmin: false
	Ymin: false
	readonly: false
	New(image = '', .message = 'no image', .stretch = false, .acceptDrop = false,
		.color = false, .bgndcolor = false, .alwaysReadOnly? = false)
		{
		.orig_message = message
		.CreateWindow("SuBtnfaceNocursor", "", WS.VISIBLE)
		.SubClass()
		.set(image)
		.sizeControlToImage()
		.setBrush()
		.font = Suneido.Member?(#hfont)
			? Suneido.hfont
			: GetStockObject(SO.DEFAULT_GUI_FONT)
		if .acceptDrop is true
			DragAcceptFiles(.Hwnd, true)
		}
	SetTip(tip)
		{
		.ToolTip(tip)
		}
	sizeControlToImage()
		{
		if .Xstretch isnt false or .Ystretch isnt false
			return // size will be set by parent

		if .Xmin is false and .Ymin is false
			{ // size control to match image
			.Xmin = ScaleWithDpiFactor(.iw)
			.Ymin = ScaleWithDpiFactor(.ih)
			}
		else if .Xmin is false
			{ // set Xmin to maintain proportions
			.Xmin = (.iw / .ih) * .Ymin
			}
		else if .Ymin is false
			{ // set Ymin to maintain proportions
			.Ymin = (.ih / .iw) * .Xmin
			}
		}
	highlight?: false
	PAINT()
		{
		hdc = BeginPaint(.Hwnd, ps = Object())
		if .image is '' // Draw .message text
			WithHdcSettings(hdc, .hdcSettings())
				{
				offset = 16
				left = top = 8
				right = .width - offset
				bottom = .height - offset
				DrawText(hdc, .message, -1, Object(:left, :top, :right, :bottom),
					DT.EDITCONTROL | DT.WORDBREAK | DT.NOPREFIX)
				}
		else
			.image.Draw(hdc, .x, .y, .w, .h, .imagebrush, .bgndbrush)
		EndPaint(.Hwnd, ps)
		return 0
		}

	hdcSettings()
		{
		hdcSettings = Object(.font, SetBkMode: TRANSPARENT)
		if .readonly and not .highlight?
			hdcSettings.SetTextColor = CLR.Inactive
		else if .highlight?
			hdcSettings.SetTextColor = CLR.Highlight
		return hdcSettings
		}

	ERASEBKGND(wParam)
		{
		GetClientRect(.Hwnd, rect = Object())
		FillRect(wParam, rect, .GetReadOnly() ? .bgndbrush : .editbrush)
		return 1
		}

	Dragging: false
	depressed: false
	LBUTTONDOWN()
		{
		if not .GetReadOnly() and not .isEmpty?() and .acceptDrop
			{
			.depressed = true
			.Send("ImageStartDrag")
			}
		return 0
		}
	isEmpty?()
		{
		return .Get() is "" and .message is .orig_message
		}
	MOUSEMOVE(wParam)
		{
		if .GetReadOnly() or wParam isnt MK.LBUTTON
			.Dragging = false
		if .Dragging
			SetCursor(LoadCursor(ResourceModule(), IDC.ARROWSTACK))
		else
			SetCursor(LoadCursor(NULL, IDC.ARROW))
		// Prevents the first click in the control from changing the cursor
		// without moving
		if not .GetReadOnly() and .depressed
			.Dragging = true
		TrackMouseEvent(Object(cbSize: TRACKMOUSEEVENT.Size(),
			dwFlags: TME.LEAVE hwndTrack: .Hwnd))
		.Send("ImageMouseMove")
		return 0
		}
	MOUSELEAVE()
		{
		if not .GetReadOnly()
			.Send("ImageMouseLeave", .Dragging)
		.Dragging = false
		.depressed = false
		}
	LBUTTONUP()
		{
		if not .GetReadOnly() and .Dragging
			{
			SetCursor(LoadCursor(NULL, IDC.ARROW))
			.Dragging = false
			.Send("ImageFinishDrag")
			.depressed = false
			return 0
			}
		.depressed = false
		.Send("ImageClick")
		return 0
		}
	LBUTTONDBLCLK()
		{
		.Send("ImageDoubleClick")
		return 0
		}

	DROPFILES(wParam)
		{
		.Send("ImageDropFiles", wParam)
		return 0
		}

	ContextMenu(x, y)
		{
		.Send("ImageContextMenu", x, y)
		return 0
		}
	Get()
		{ return .image; }
	Set(image, highlight? = false, message = false, color = false)
		{
		if false isnt color
			.swapBrush(color)
		.set(image, message)
		.highlight? = highlight?
		.Repaint()
		}
	image: ''
	iw: 0
	ih: 0
	set(image, message = false)
		{
		image = .loadIfBook(image)
		.message = .orig_message
		if image is .image
			return
		if .image isnt ''
			{
			.image.Close()
			.image = ''
			}
		if image is ''
			return
		open = .openPreview(image)
		if open is '' or (String?(open) and open.Prefix?(.invalidFileToOpen))
			return .canNotOpenForPreview(open, image, message)

		.image = open
		hdc = GetDC(.Hwnd)
		.iw = .image.Width(hdc)
		.ih = .image.Height(hdc)
		ReleaseDC(.Hwnd, hdc)

		.sizeImageToControl()
		}
	canNotOpenForPreview(open, image, message)
		{
		.image = ''
		.message = (message is false ? image : message) $
			Opt(' * ', open.AfterFirst(.invalidFileToOpen $ ': ') , ' * ')
		.iw = .ih = 100
		return
		}
	loaded: false
	loadIfBook(image)
		{
		if String?(image) and
			Paths.IsValid?(image) and image.Has?('%') and
			false isnt x = .load(image)
			{
			.loaded = true
			return x.text
			}
		else
			return image
		}
	load(image)
		{
		book = image.BeforeFirst('%')
		if not TableExists?(book)
			return false
		name = image.AfterFirst('%')
		if name.Has?('/')
			{
			path = '/res/' $ name.BeforeLast('/')
			name = name.AfterLast('/')
			return Query1(book, :path, :name)
			}
		else
			return Query1(book, :name)
		}
	max_size: 1048576
	openPreview(image)
		{
		if not String?(image)
			return image  // already an open image
		if Paths.IsValid?(image) and .loaded is false
			{
			if false isnt code = IconFont().MapToCharCode(image)
				return ImageFont(code.char, code.font)

			result = .validFileToOpen(image)
			if result isnt true
				return result
			}

		Image.RunWithErrorLog({ return Image(image) })
		return ''
		}
	invalidFileToOpen: 'InvalidFile'
	validFileToOpen(image)
		{
		if image !~ "(?i)[.](bmp|gif|jpg|jpe|jpeg|ico|emf|wmf)$"
			return .invalidFileToOpen

		try fileSize = .fileSize(image)
		catch (err)
			{
			if err.Has?('does not exist')
				return .invalidFileToOpen $ ": Can not find file"
			return .invalidFileToOpen $ ': Can not open file'
			}
		if fileSize > .max_size
			return .invalidFileToOpen
		return true
		}
	fileSize(image)
		{
		return FileSize(image)
		}
	sizeImageToControl()
		{
		.x = 0
		.y = 0
		if .stretch is true
			{ // stretch image to fit control
			.w = .width
			.h = .height
			}
		else if .iw / .ih < .width / .height
			{ // center image horizontally to maintain proportions
			.w = ((.iw / .ih) * .height).Round(0)
			.h = .height
			.x = ((.width - .w) / 2).Round(0)
			}
		else // .iw / .ih >= w / h
			{ // center image vertically to maintain proportions
			.w = .width
			.h = ((.ih / .iw) * .width).Round(0)
			.y = ((.height - .h) / 2).Round(0)
			}
		}
	SetReadOnly(ro)
		{
		if .alwaysReadOnly?
			return
		ro = ro is true
		if ro is .readonly
			return
		.readonly = ro
		.Repaint()
		}
	GetReadOnly()
		{
		return .alwaysReadOnly? or .readonly
		}
	width: 0
	height: 0
	Resize(x, y, w, h)
		{
		super.Resize(x, y, w, h)
		.width = w
		.height = h
		.sizeImageToControl()
		}
	imagebrush: false
	bgndbrush: false
	editbrush: false
	swapBrush(color)
		{
		.deleteBrush()
		.color = color
		.setBrush()
		}
	setBrush()
		{
		.imagebrush = .color isnt false
			? CreateSolidBrush(.color)
			: false
		.bgndbrush = .bgndcolor isnt false
			? CreateSolidBrush(.bgndcolor)
			: GetSysColorBrush(COLOR.BTNFACE)
		.editbrush = CreateSolidBrush(CLR.WHITE)
		}
	deleteBrush()
		{
		if .imagebrush isnt false
			DeleteObject(.imagebrush)
		if .bgndbrush isnt false
			DeleteObject(.bgndbrush)
		if .editbrush isnt false
			DeleteObject(.editbrush)
		.imagebrush = .bgndbrush = .editbrush = false
		}
	Destroy()
		{
		if not String?(.image)
			.image.Close()
		.deleteBrush()
		super.Destroy()
		}
	}