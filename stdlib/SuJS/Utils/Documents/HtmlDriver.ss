// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
Document_Builder
	{
	Getter_Page()
		{
		throw "not implemented"
		}

	pageWidth: 999_999_999_999
	AddPage(dimens)
		{
		.pageWidth = dimens.width * 96 /*=pixelPerInch*/ * 15 /*=pixelToTwip*/
		}

	addToPage(item, x)
		{
		if x > .pageWidth
			return

		try
			{
			.Page.Add(item)
			}
		catch (unused, 'object too large')
			.logTooLarge()
		}

	logTooLarge(_report = false)
		{
		SuneidoLog('ERROR: (CAUGHT) HtmlDriver - .Page object too large', params: report)
		throw 'SHOW: Too many prints on a page'
		}

	AddText(data, x, y, w, h, font, justify = 'left', ellipsis? = false,
		color = false, html = false)
		{
		y = y + PdfFonts.GetFontHeight(font) * .afmToTwip
		font = PdfFonts.GetCompatibleFont(font)
		data = PdfFonts.StripInvalidChars(data)
		.addToPage([#AddText, data, x, y, w, h, font, justify, ellipsis?, color, html], x)
		}

	AddMultiLineText(data, x, y, w, h, font, justify = 'left', color = false)
		{
		font = PdfFonts.GetCompatibleFont(font)
		lineHeight = .GetLineSpecs(font).height
		for line in data.Split('\n')
			{
			.AddText(line, x, y, w, lineHeight, font, justify, :color)
			y += lineHeight
			h -= lineHeight
			if h < lineHeight
				break
			}
		}

	AddLine(x, y, x2, y2, thick, color = 0x00000000)
		{
		.addToPage([#AddLine, x, y, x2, y2, thick, color], Min(x, x2))
		}

	AddRect(x, y, w, h, thick, fillColor = false, lineColor = false)
		{
		.addToPage([#AddRect, x, y, w, h, thick, fillColor, lineColor], x)
		}

	AddImage(x, y, w, h, data)
		{
		if Paths.IsValid?(data)
			{
			if not Jpeg.ValidExtension?(data)
				throw Jpeg.InvalidExtension
			if false is data = ImageHandler.GenerateThumbnail(data)
				throw "Image: Couldn't generate thumbnail"
			}
		new Jpeg(data) // verify that the image is a valid Jpeg
		.addToPage([#AddImage, x, y, w, h,
			'data:image/jpeg;base64,' $ Base64.Encode(data)], x)
		}

	GetImageSize(data, jpeg = false)
		{
		w = h = 0
		if jpeg is false
			{
			size = ImageHandler.GetWidthHeight(data)
			w = size.w
			h = size.h
			}
		else
			{
			w = jpeg.GetWidth()
			h = jpeg.GetHeight()
			}
		return Object(height: h, width: w)
		}

	AddArc(left, top, right, bottom,
		xStartArc = 0, yStartArc = 0, xEndArc = 0, yEndArc = 0,
		thick = 1, lineColor = false)
		{
		.addToPage([#AddArc, left, top, right, bottom,
			xStartArc, yStartArc, xEndArc, yEndArc, thick, lineColor], left)
		}

	AddEllipse(x, y, w, h, thick = 1, fillColor = false, lineColor = false)
		{
		.addToPage([#AddEllipse, x, y, w, h, thick, fillColor, lineColor], x)
		}

	AddCircle(x, y, radius, thick, fillColor = false, lineColor = false)
		{
		.AddEllipse(x - radius, y - radius, radius * 2, radius * 2,
			thick, fillColor, lineColor)
		}

	AddRoundRect(x, y, w, h, width = 0, height = 0, thick = 1,
		fillColor = false, lineColor = false)
		{
		.addToPage([#AddRoundRect, x, y, w, h, width, height, thick, fillColor,
			lineColor], x)
		}

	EnsureFont(font, oldfont)
		{
		return EnsurePDFFont(font, oldfont)
		}

	PixelToUnit: 15 /*=twipToPixel = 1440 / 99*/
	afmToTwip: .02 // afm Spec uses measurements in 1 1000th PSP. This number
				   // Converts from that measurement to Twips.
	GetLineSpecs(font)
		{
		// calculating line height
		descender = Abs(PdfFonts.GetFontDescender(font)) * .afmToTwip
		ascender = PdfFonts.GetFontHeight(font) * .afmToTwip
		return Object(height: descender + ascender, descent: descender)
		}

	GetDefaultFont()
		{
		return Object(angle: 0, size: 10, weight: FW.NORMAL,
			name: "Helvetica", italic: false)
		}

	GetCharWidth(width, font, widthChar)
		{
		if width is false
			return 0
		font = PdfFonts.GetCompatibleFont(font)
		return width * PdfFonts.GetCharWidth(font, widthChar) *
			font.size * .afmToTwip
		}

	GetTextHeight(data, lineHeight)
		{
		numrows = data.Lines().Size()
		return lineHeight * numrows
		}

	GetTextWidth(font, text)
		{
		font = PdfFonts.GetCompatibleFont(font)
		totalWidth = 0
		for c in text
			totalWidth += PdfFonts.GetCharWidth(font, c)
		return totalWidth * font.size * .afmToTwip
		}

	Finish(status)
		{
		return status
		}
	}