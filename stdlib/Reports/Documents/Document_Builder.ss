// Copyright (C) 2014 Axon Development Corporation All rights reserved worldwide.
class
	{
	PlainText?()
		{
		return false
		}

	AddPage(dimens/*unused*/)
		{
		}

	EndPage()
		{
		}

	AddText(data /*unused*/, x /*unused*/, y /*unused*/, w /*unused*/, h /*unused*/,
		font /*unused*/, justify /*unused*/ = 'left', ellipsis?  /*unused*/= false,
		color /*unused*/ = false)
		{
		}

	AddMultiLineText(data /*unused*/, x /*unused*/, y /*unused*/, w /*unused*/,
		h /*unused*/, font /*unused*/, justify  /*unused*/ = 'left',
			color /*unused*/ = false)
		{
		}

	AddLine(x /*unused*/, y /*unused*/, x2 /*unused*/, y2 /*unused*/,
		thick /*unused*/ = false, color /*unused*/= 0x00000000)
		{
		}

	AddRect(x /*unused*/, y /*unused*/, w /*unused*/, h /*unused*/, thick /*unused*/,
		fillColor /*unused*/ = false, lineColor /*unused*/ = false)
		{
		}

	AddCircle(x /*unused*/, y /*unused*/, radius /*unused*/, thick /*unused*/,
		fillColor /*unused*/ = false, lineColor /*unused*/ = false)
		{
		}

	AddRoundRect(x /*unused*/, y /*unused*/, w /*unused*/, h /*unused*/,
		width /*unused*/, height /*unused*/, thick /*unused*/,
		fillColor /*unused*/ = false, lineColor /*unused*/ = false)
		{
		}

	AddEllipse(x /*unused*/, y /*unused*/, w /*unused*/, h /*unused*/,
		thick /*unused*/, fillColor /*unused*/ = false, lineColor /*unused*/ = false)
		{
		}

	AddArc(left /*unused*/, top /*unused*/, right /*unused*/, bottom /*unused*/,
		xStartArc /*unused*/, yStartArc /*unused*/, xEndArc /*unused*/,
		yEndArc /*unused*/, thick /*unused*/, lineColor /*unused*/ = false)
		{
		}

	AddPolygon(points /*unused*/, thick /*unused*/, fillColor  /*unused*/ = false,
		lineColor /*unused*/ = false)
		{
		}

	RegisterFont(font /*unused*/, defaultSize /*unused*/)
		{
		}

	EnsureFont(font, oldFont /*unused*/)
		{
		return font
		}

	GetDefaultFont()
		{
		return #(name: "Arial", size: 10, weight: 400, angle: 0, italic: false)
		}

	GetLineSpecs(font /*unused*/)
		{
		}

	GetCharWidth(width /*unused*/, font /*unused*/, widthChar /*unused*/)
		{
		}

	GetTextWidth(font /*unused*/, data /*unused*/)
		{
		}

	GetTextHeight(data /*unused*/, lineHeight /*unused*/)
		{
		}

	GetImageSize(data /*unused*/)
		{ }

	GetImageSizeAdjustment()
		{
		return 1
		}

	AddImage(x /*unused*/, y /*unused*/, w /*unused*/, h /*unused*/, data /*unused*/)
		{
		}

	GetAcceptedImageExtension()
		{
		}

	DrawWithinClip(x /*unused*/, y /*unused*/, w /*unused*/, h /*unused*/, block)
		{
		block()
		}

	SetMultiPartsRatio(@unused)
		{
		}

	Finish(status)
		{
		return status
		}
	}