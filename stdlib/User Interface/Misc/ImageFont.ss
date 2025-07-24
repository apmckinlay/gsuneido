// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
class
	{
	New(.ch, .font)
		{
		.suneidoFont? = .font.Prefix?('suneido')
		}

	IsRasterImage()
		{
		return false
		}

	Width(hdc)
		{
		.suneidoFont?
			? .getDimension(hdc).x
			: Max(.getDimension(hdc).x, .getDimension(hdc).y)
		}

	Height(hdc)
		{
		.getDimension(hdc).y
		}

	dimension: false
	getDimension(hdc)
		{
		if .dimension isnt false
			return .dimension
		DoWithHdcObjects(hdc, [.getHfont(hdc)])
			{
			GetTextExtentPoint32(hdc, .ch, 1, ex = Object())
			return .dimension = ex
			}
		}

	getHfont(hdc, size = '', orientation = 0)
		{
		if size is ''
			{
			factor = GetDeviceCaps(hdc, GDC.LOGPIXELSY) / PointsPerInch
			size = StdFonts.FontSize('') * factor
			}

		fonts = Suneido.GetInit('ImageFonts', Object().Set_default(Object()))
		if fonts[.font].Member?(key = size $ ':' $ orientation)
			return fonts[.font][key]

		lf = Object(
			lfFaceName: .font,
			lfHeight: -size,
			lfWeight: .suneidoFont? ? FW.NORMAL : FW.BOLD,
			lfUnderline: false,
			lfItalic: false,
			lfStrikeOut: false,
			lfOrientation: orientation,
			lfEscapement: orientation,
			lfCharSet: CHARSET[GetLanguage().charset])
		hfont = CreateFontIndirect(lf)
		Assert(hfont isnt: 0)
		return fonts[.font][key] = hfont
		}

	Draw(hdc, x, y, w = 0, h = 0, brushImage = false, brushBackground = false,
		orientation = 0)
		{
		font = .getHfont(hdc, h, orientation)
		WithHdcSettings(hdc, .hdcSettings(brushImage, brushBackground, font))
			{
			Assert(GetTextFace(hdc) is: .font)
			if not .suneidoFont?
				{
				DrawText(hdc, .ch, 1, dimen = Object(left: 0, top: 0), DT.CALCRECT)
				.textOut(hdc, x + (w - dimen.right) / 2, y + (h - dimen.bottom) / 2)
				}
			else
				.textOut(hdc, x, y)
			}
		}

	hdcSettings(brushImage, brushBackground, font)
		{
		settings = Object(font, SetBkMode: TRANSPARENT)
		if false isnt backgroundColor = .brushColor(brushBackground)
			{
			settings.SetBkColor = backgroundColor
			settings.SetBkMode = OPAQUE
			}
		if false isnt imageColor = .brushColor(brushImage)
			settings.SetTextColor = imageColor
		return settings
		}

	brushColor(brush)
		{
		if brush is false
			return false
		lb = GetObjectBrush(brush)
		return lb isnt false
			? lb.lbColor
			: false
		}

	textOut(hdc, x, y)
		{
		TextOut(hdc, x, y, .ch, 1)
		}

	Close() {}

	DestroyAllImageFonts() // Manually called for testing / development purposes
		{
		Suneido.
			Extract('ImageFonts', Object().Set_default(Object())).
			Each({ it.Each(DeleteObject) })
		}
	}
