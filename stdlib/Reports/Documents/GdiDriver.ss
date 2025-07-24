// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
Document_Builder
	{
	New()
		{
		.fonts = Object()
		}

	GetDC()
		{
		throw "not implemented"
		}

	GetCopyDC()
		{
		throw "not implemented"
		}

	params: (iTabLength: 4, cbSize: 20) // DRAWTEXTPARAMS.Size()
	AddText(data, x, y, w, h, font /*unused*/, justify = 'left', ellipsis? = false,
		color = false)
		{
		flags = DT.EXPANDTABS | DT.TABSTOP | DT.NOPREFIX
		if (data.Find1of('\r\n') < data.Size())
			flags |= DT.END_ELLIPSIS
		dc = .GetDC()
		if ellipsis?
			flags |= DT.END_ELLIPSIS
		rect = Object(left: x, right: x + w, top: y, bottom: y + h)
		flags = .handleJustification(justify, rect, data, flags)
		.doWithColor(dc, color)
			{
			if not _report.Member?('TooltipFullDisplay')
				DrawTextEx(dc, data, -1, rect, flags, .params)
			else
				{
				newData = DrawTextExOut(
					dc, data, rect, flags | DT.MODIFYSTRING, .params).text
				.CacheTooltipFullDisplay(data, newData)
				}
			}
		}

	CacheTooltipFullDisplay(data, newData)
		{
		fullDisplay = _report.TooltipFullDisplay
		fullDisplay[fullDisplay.CurrentCol] $= Opt(' ', data.Trim())
		if data isnt newData
			fullDisplay.CellEllipsized = true
		}

	doWithColor(dc, color, block)
		{
		if color isnt false and not .darkBackground?
			WithHdcSettings(dc, [SetTextColor: color], block)
		else
			block()
		}

	RegisterFont(font, defaultSize)
		{
		weight = font.GetDefault(#weight, 400) /*= default font weight */
		italic = font.GetDefault(#italic, false)
		underline = font.GetDefault(#underline, false)
		strikeout = font.GetDefault(#strikeout, false)

		size = font.GetDefault(#size, defaultSize)
		size = (size * 20).Int()  /*= to twips */

		name = font.GetDefault(#name, .GetDefaultFont().name)
		id = name $ " " $ size $ " " $ weight $ " " $
			italic $ " " $ underline $ " " $ strikeout
		if .fonts.Member?(id)
			f = .fonts[id]
		else
			{
			lf = Object(
				lfFaceName: name,
				lfHeight: -size,
				lfWeight: weight,
				lfItalic: italic,
				lfUnderline: underline,
				lfStrikeOut: strikeout,
				lfCharSet: CHARSET[font.Member?("charset") ? font.charset : 'DEFAULT']
				)
			f = CreateFontIndirect(lf)
			if f is 0
				throw "Report: couldn't create font"
			.fonts[id] = f
			}
		// Cannot restore the original SelectObject value as the end point is ambiguous.
		// As a result, this cannot be changed to use DoWithHdcObjects or WithHdcSettings.
		Assert(-1 isnt SelectObject(.GetDC(), f))
		}

	GetLineSpecs(font /*unused*/)
		{
		GetTextMetrics(.GetDC(), tm = Object())
		return Object(height: tm.Height, descent: tm.Descent)
		}

	GetTextWidth(font /*unused*/, data)
		{
		return .getTextRect(data).width
		}

	getTextRect(data)
		{
		flags = DT.CALCRECT | .getFlags(justify: 'left')
		h = DrawTextEx(.GetCopyDC(), data, -1, r = Object(), flags, .params)
		return Object(width: r.right, height: h)
		}

	GetCharWidth(width, font /*unused*/, widthChar)
		{
		if width is false
			return 0
		width = width.Int()
		GetTextExtentPoint32(.GetCopyDC(),
			widthChar.Repeat(width), width, sz = Object())
		return sz.x
		}

	handleJustification(justify, rect, data, flags)
		{
		if justify is "right"
			return .RightJustify(rect, data, flags)
		else if justify is "center"
			return .CenterJustify(rect, data, flags)
		return flags
		}

	GetTextHeight(data, lineHeight /*unused*/)
		{
		return .getTextRect(data).height
		}

	getFlags(justify)
		{
		flags = DT.EXPANDTABS | DT.TABSTOP | DT.EXTERNALLEADING | DT.NOPREFIX
		if justify is "right"
			flags |= DT.RIGHT
		else if justify is "center"
			flags |= DT.CENTER
		return flags
		}

	AddMultiLineText(data, x, y, w, h, font /*unused*/, justify = 'left', color = false)
		{
		dc = .GetDC()
		rect = Object(left: x, right: x + w, top: y, bottom: y + h)
		.doWithColor(dc, color)
			{
			DrawTextEx(dc, data, -1, rect, .getFlags(justify), .params)
			}
		}

	AddLine(x, y, x2, y2, thick, color = 0x00000000)
		{
		DoWithHdcObjects(dc = .GetDC(), [pen = CreatePen(PS.SOLID, thick, color)])
			{
			MoveTo(dc, x, y)
			LineTo(dc, x2, y2)
			}
		DeleteObject(pen)
		}

	AddRect(x, y, w, h, thick, fillColor = false, lineColor = false)
		{
		.doWithColorAndBorder(thick, fillColor, lineColor)
			{
			Rectangle(.GetDC(), x, y, x + w, y + h)
			}
		}

	doWithColorAndBorder(thick, fillColor, lineColor, block)
		{
		lineColor = lineColor is false ? CLR.BLACK : lineColor
		pen = CreatePen(thick is 0 ? PS.NULL : PS.SOLID, thick, lineColor)
		hdcSettings = [pen]
		if fillColor isnt false
			hdcSettings.brush = fillColor
		WithHdcSettings(.GetDC(), hdcSettings, block)
		DeleteObject(pen)
		}

	AddCircle(x, y, radius, thick, fillColor = false, lineColor = false)
		{
		.doWithColorAndBorder(thick, fillColor, lineColor)
			{
			Ellipse(.GetDC(), x - radius, y - radius, x + radius, y + radius)
			}
		}

	GetImageSize(data)
		{
		image = Image(data)
		dpiFactor = GetDpiFactor()
		mmInInch = 2540 // 1 inch to .01 mm
		width = dpiFactor * (image.Width(false).InchesInTwips() / mmInInch)
		height = dpiFactor * (image.Height(false).InchesInTwips() / mmInInch)
		image.Close()
		return Object(:width, :height)
		}
	backgroundBrush: false
	darkBackground?: false
	SetBackgroundBrush(brush, darkBackground? = false)
		{
		orig = Object(.backgroundBrush, .darkBackground?)
		.backgroundBrush = brush
		if darkBackground? isnt ''
			.darkBackground? = darkBackground?
		return orig
		}
	AddImage(x, y, w, h, data)
		{
		if .GetDC() is 0 // for fake printer device context
			return
		image = Image(data)
		image.Draw(.GetDC(), x, y, w, h, brushBackground: .backgroundBrush)
		image.Close()
		}

	GetAcceptedImageExtension()
		{
		return '.gif'
		}

	AddRoundRect(x, y, w, h, width, height, thick, fillColor = false, lineColor = false)
		{
		.doWithColorAndBorder(thick, fillColor, lineColor)
			{
			RoundRect(.GetDC(), x, y, x + w, y + h, width, height)
			}
		}
	AddEllipse(x, y, w, h, thick, fillColor = false, lineColor = false)
		{
		.doWithColorAndBorder(thick, fillColor, lineColor)
			{
			Ellipse(.GetDC(), x, y, x + w, y + h)
			}
		}
	AddArc(left, top, right, bottom, xStartArc, yStartArc, xEndArc, yEndArc, thick,
		lineColor = false)
		{
		.doWithColorAndBorder(thick, false, lineColor)
			{
			Arc(.GetDC(), left, top, right, bottom, xStartArc, yStartArc, xEndArc,
				yEndArc)
			}
		}
	AddPolygon(points, thick, fillColor = false, lineColor = false)
		{
		.doWithColorAndBorder(thick, fillColor, lineColor)
			{
			Polygon(.GetDC(), points, points.Size())
			}
		}
	Finish(status)
		{
		.fonts.Each(DeleteObject)
		return status
		}
	}
