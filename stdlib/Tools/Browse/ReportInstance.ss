// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
ReportBase
	{
	New(.fontScale = false)
		{
		.curfont = false
		.SetDriver(new GdiPreviewDriver)
		}
	default_font: (name: "Arial", size: 9)
	dc: false
	SetDC(dc)
		{
		.dc = dc
		.SelectFont(.default_font)
		}

	SelectFont(font)
		{
		if font is false
			return false

		weight = font.GetDefault("weight", 400) /*= default to normal size */
		id = font.name $ " " $ font.size $ " " $ weight
		oldfont = .curfont
		if not .Fonts.Member?(id)
			{
			lf = Object(
				lfCharSet: CHARSET[font.GetDefault("charset", "DEFAULT")],
				lfFaceName: font.Member?("name") ? font.name : oldfont.name,
				lfHeight: .fontScale is false
					? -font.size * GetDeviceCaps(.dc, GDC.LOGPIXELSY) / PointsPerInch
					: -font.size * .fontScale,
				lfWeight: weight)
			f = CreateFontIndirect(lf)
			if f is 0
				throw "Report: couldn't create font"
			.Fonts[id] = f
			}
		// Cannot restore the original SelectObject value as the end point is ambiguous.
		// As a result, this cannot be changed to use DoWithHdcObjects or WithHdcSettings.
		SelectObject(.dc, .Fonts[id])
		.curfont = font
		return oldfont
		}
	GetFont()
		{
		return .curfont
		}
	GetDimens()
		{
		return #(
			left: .5,
			right: .5,
			top: .5,
			bottom: .5,
			width: 8.5,
			height: 11,
			W: 7.5,
			H: 10)
		}
	}
