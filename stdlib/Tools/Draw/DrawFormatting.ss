// Copyright (C) 2015 Suneido Software Corp. All rights reserved worldwide.
ReportBase
	{
	New()
		{
		.SetDriver(new GdiPreviewDriver)
		}

	Paint(items, rc)
		{
		_report = this
		for item in items
			{
			rc1 = Object()
			rcItem = item.BoundingRect()
			rc1.top = rcItem.y1 <= rcItem.y2 ? rcItem.y1 : rcItem.y2
			rc1.left = rcItem.x1 <= rcItem.x2 ? rcItem.x1 : rcItem.x2
			rc1.bottom = rcItem.y1 > rcItem.y2 ? rcItem.y1 : rcItem.y2
			rc1.right = rcItem.x1 > rcItem.x2 ? rcItem.x1 : rcItem.x2
			if .rectIntersection?(rc, rc1)
				item.Paint()
			}
		}

	rectIntersection?(rc1, rc2)
		{
		if rc2.top > rc1.bottom or rc2.bottom < rc1.top or rc2.right < rc1.left or
			rc2.left > rc1.right
			return false
		else
			return true
		}

	RegisterFont(font, defaultSize)
		{
		weight = font.GetDefault(#weight, FW.NORMAL)
		italic = font.GetDefault(#italic, false)
		underline = font.GetDefault(#underline, false)
		strikeout = font.GetDefault(#strikeout, false)

		size = font.GetDefault(#size, defaultSize)
		size = size * WinDefaultDpi / PointsPerInch

		name = font.GetDefault(#name, .GetDefaultFont().name)
		id = name $ " " $ size $ " " $ weight $ " " $
			italic $ " " $ underline $ " " $ strikeout
		if .Fonts.Member?(id) and strikeout is false
			f = .Fonts[id]
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
			if (f is 0)
				throw "Report: couldn't create font"
			.Fonts[id] = f
			}
		// Cannot restore the original SelectObject value as the end point is ambiguous.
		// As a result, this cannot be changed to use DoWithHdcObjects or WithHdcSettings.
		Assert(-1 isnt SelectObject(.GetDC(), f))
		}
	}