// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
GdiDriver
	{
	New(.pdc = 0)
		{
		}

	metaFile?: false
	AddPage(dimens)
		{
		filename = GetAppTempFullFileName("su")
		.dc = .createMetaFile(dimens, filename)
		.cdc = .createMetaFile(dimens)
		.metaFile? = true
		return filename
		}

	createMetaFile(dimens, file = 0)
		{
		f = 2540 // inches to .01 mm
		r = Object(right: dimens.width * f, bottom: dimens.height * f)
		emfdc = CreateEnhMetaFile(.pdc, file, r, NULL)

		SetupGdiDeviceScale(emfdc)
		// Cannot switch to WithBkMode as the end point is ambiguous
		SetBkMode(emfdc, TRANSPARENT)
		return emfdc
		}

	EndPage()
		{
		DeleteEnhMetaFile(CloseEnhMetaFile(.cdc))
		.cdc = false
		return CloseEnhMetaFile(.dc)
		}

	GetTextWidth(font /*unused*/, data)
		{
		return .getPreviewSize(data).x
		}

	RightJustify(rect, data, flags)
		{
		rect.left = Max(rect.left, rect.right - .GetTextWidth(false, data))
		return flags
		}

	CenterJustify(rect, data, flags)
		{
		rect.left = rect.left +
			Max(0, (rect.right - rect.left - .GetTextWidth(false, data)) / 2)
		return flags
		}

	getPreviewSize(data)
		{
		// need this since DrawTextEx CALCRECT appears to assume kerning
		// but preview does not kern (don't know why)
		// so measurement is too small and text gets cut off
		// GetTextExtentPoint32 appears to NOT assume kerning
		// so it gives the right measurement for preview
		data = data.Detab()
		GetTextExtentPoint32(.GetCopyDC(), data, data.Size(),
			sz = Object())
		return sz
		}

	AddImage(x, y, w, h, data)
		{
		if String?(data) and Paths.IsValid?(data)
			Format.Hotspot(x, y, w, h, [],
				access: [control: "AttachmentGoTo", file: data])
		super.AddImage(x, y, w, h, data)
		}

	AddText(data, x, y, w, h, font, justify = 'left', ellipsis? = false, color = false)
		{
		if justify isnt 'right'
			{
			super.AddText(data, x, y, w, h, font, justify, ellipsis?, color)
			return
			}
		// Handle right justification specially because the cutoff problem with DrawText
		dc = .GetDC()
		data = data.Detab()
		if ellipsis?
			{
			newData = .gdiEllipsis(font, data, w)
			if _report.Member?('TooltipFullDisplay')
				.CacheTooltipFullDisplay(data, newData)
			data = newData
			}
		rect = Object(left: x, right: x + w, top: y, bottom: y + h)
		.doWithSettings(dc, justify, rect, color)
			{ |refPtr|
			ExtTextOut(dc, refPtr.x, refPtr.y, ETO.CLIPPED, rect, data, data.Size(), NULL)
			}
		}

	doWithSettings(dc, justify, rect, color, block)
		{
		prevTextAlign = false
		refPtr = Object(x: rect.left, y: rect.top)
		if justify is 'right'
			{
			prevTextAlign = SetTextAlign(dc, TA.RIGHT)
			refPtr.x = rect.right
			}
		else if justify is 'center'
			{
			prevTextAlign = SetTextAlign(dc, TA.CENTER)
			refPtr.x = (rect.left + rect.right) / 2
			}
		if color isnt false
			prevLineColor = SetTextColor(dc, color)
		block(refPtr)
		if color isnt false
			SetTextColor(dc, prevLineColor)
		if prevTextAlign isnt false
			SetTextAlign(dc, prevTextAlign)
		}

	gdiEllipsis(font/*unused*/, data, w) // for right justify
		{
		do
			{
			textSize = .getPreviewSize(data).x
			if textSize <= w
				return data
			if not data.Suffix?('...')
				data = data[..-2] $ '...' // remove 2 digits to make up for ...
			data = data[..-4] $ '...' /*= remove 1 more digit to try */
			} while data.Size() > 3 /*= only ... */
		return data
		}

	AddMultiLineText(data, x, y, w, h, font, justify = 'left', color = false)
		{
		if justify isnt 'right'
			{
			super.AddMultiLineText(data, x, y, w, h, font, justify, color)
			return
			}
		// Handle right justification specially because the cutoff problem with DrawText
		externalLeading = .getExternalLeading()
		for line in data.Split('\n')
			{
			lineHeight = .getPreviewSize(line is '' ? 'M' : line).y
			.AddText(line, x, y, w, lineHeight, font, justify, :color)
			y += lineHeight + externalLeading
			}
		}

	getExternalLeading()
		{
		dc = .GetDC()
		GetTextMetrics(dc, ob = Object())
		return ob.GetDefault(#ExternalLeading, 0)
		}

	DrawWithinClip(x, y, w, h, block)
		{
		if .metaFile? // cannot SelectClipRgn on meta file
			{
			block()
			return
			}

		.keepClipRegion()
			{ |dc|
			GetClipBox(dc, r = Object())
			hrgn = CreateRectRgn(Max(r.left, x), y, x + w, y + h)
			SelectClipRgn(dc, hrgn)
			block()
			SelectClipRgn(dc, NULL)
			DeleteObject(hrgn)
			}
		}

	keepClipRegion(block)
		{
		dc = .GetDC()
		oldRegion = CreateRectRgn( 0, 0, 0, 0 )
		if GetClipRgn(dc, oldRegion) isnt 1
			{
			DeleteObject(oldRegion)
			oldRegion = 0
			}
		block(dc)
		SelectClipRgn(dc, oldRegion)
		if oldRegion isnt 0
			DeleteObject(oldRegion)
		}

	cdc: false
	GetCopyDC()
		{
		if .cdc is false
			return .dc

		/* WARNING: the below only copies the font (not pen, brush, etc.)

		NOTE: This is used by WrapFormat and TextFormat to do CALCRECT
			since these increase the size of the preview emf.

		NOTE: While GetCurrentObject will return the value of the current type, the below
			code is using a different DC value to get and set. As the end point is
			ambiguous, this cannot be changed to use DoWithHdcObjects or WithHdcSettings.
		*/
		SelectObject(.cdc, GetCurrentObject(.dc, OBJ.FONT))
		return .cdc
		}
	dc: false
	GetDC()
		{ return .dc }
	SetDC(dc)
		{ .dc = dc }
	}