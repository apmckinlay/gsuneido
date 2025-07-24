// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
CanvasItem
	{
	defaultFontSize: 9
	New(text, x1, y1, x2, y2, font = false, .name = '', encoded = false,
		.justify = 'left', fromDraw = false)
		{
		.text = encoded ? Base64.Decode(text) : text
		.posLocked? = false
		.font = font is false ? Object() : font

		if fromDraw is true
			{
			_report.RegisterFont(.font, .defaultFontSize)
			s = .parseText()
			.x1 = x1
			.y1 = y1
			.x2 = .x1 + .GetWidth(s, .font)
			.y2 = .y1 + .GetHeight(s, .font)
			}
		else
			{
			.x1 = Min(x1, x2)
			.y1 = Min(y1, y2)
			.x2 = Max(x1, x2)
			.y2 = Max(y1, y2)
			}
		.rect = Object(left: .x1, top: .y1, right: .x2, bottom: .y2)
		}

	GetWidth(s, font)
		{
		lines = s.Lines()
		if lines.Empty?()
			return 20/*=default empty width*/
		return lines.Map({ _report.GetTextWidth(font, it) }).Max() /
			_report.GetDefault(#PixelToUnit, 1)
		}

	GetHeight(s, font)
		{
		if s.Blank?()
			return 20/*=default empty width*/
		lineSpecs = _report.GetLineSpecs(font)
		return _report.GetTextHeight(s, lineSpecs.height) /
			_report.GetDefault(#PixelToUnit, 1)
		}

	Paint()
		{
		try
			if .font is false or .font is #()
				.font = _report.GetFont()
		catch (unused, 'method not found in Report: GetFont')
			.font = #()

		_report.RegisterFont(.font, .defaultFontSize)
		_report.DrawWithinClip(.x1, .y1, .x2 - .x1, .y2 - .y1)
			{
			_report.AddMultiLineText(.getDisplayText(),
				.x1, .y1, .x2 - .x1, .y2 - .y1, .font,
				justify: .justify, color: .GetLineColor())
			}
		}

	getDisplayText()
		{
		displayText = .parseText()
		if _report.Base?(Report)
			displayText = displayText.
				Replace('(?q)<Page#>', {|unused| String(_report.GetPage()) }).
				Replace('(?q)<Short Date>', {|unused| Date().ShortDate() }).
				Replace('(?q)<Long Date>', {|unused| Date().LongDate() })
		return displayText
		}

	parseText()
		{
		extraLabels = GetContributions('DrawText_ExtraLabels')
		text = .text.Replace('<[\w \/]+?>')
			{|s|
			label = s[1..-1].Trim()
			replace = s
			if false isnt item = extraLabels.FindIf({ it.label is label})
				replace = (extraLabels[item].getter)()
			replace
			}
		return text.RemoveBlankLines()
		}

	BoundingRect()
		{
		return Object(x1: .rect.left, y1: .rect.top, x2: .rect.right, y2: .rect.bottom)
		}

	SetSize(x1, y1, x2, y2)
		{
		.x1 = x1
		.y1 = y1
		.x2 = x2
		.y2 = y2
		}
	ResetSize()
		{
		result = ResetSizeControl(0,
			Object(left: .rect.left, top: .rect.top, right: .rect.right,
				bottom: .rect.bottom))
		if (result is false)
			return
		.x1 = Number(result.left)
		.y1 = Number(result.top)
		.x2 = Number(result.right)
		.y2 = Number(result.bottom)
		.rect = Object(left: .x1, top: .y1, right: .x2, bottom: .y2)
		}

	Resize(origx, origy, x, y)
		{
		x1 = .Resizing?(.x1, origx) ? x : .x1
		y1 = .Resizing?(.y1, origy) ? y : .y1
		x2 = .Resizing?(.x2, origx) ? x : .x2
		y2 = .Resizing?(.y2, origy) ? y : .y2

		.x1 = Min(x1, x2)
		.y1 = Min(y1, y2)
		.x2 = Max(x1, x2)
		.y2 = Max(y1, y2)
		.rect = Object(left: .x1, top: .y1, right: .x2, bottom: .y2)
		}

	DisplayFont()
		{
		return "Object" $ Display(.font)[1..]
		}

	StringToSave()
		{
		'CanvasText(text: ' $ Display(Base64.Encode(.text)) $ ', x1: ' $ Display(.x1) $
			', y1: ' $ Display(.y1) $ ', x2: ' $ Display(.x2) $
			', y2: ' $ Display(.y2) $
			', font: ' $ .DisplayFont() $
			', encoded:, justify: ' $ Display(.justify) $ ')'
		}

	ObToSave()
		{
		return Object('CanvasText', Base64.Encode(.text), .x1, .y1, .x2, .y2,
			.font.Copy(), encoded:, justify: .justify)
		}

	Edit()
		{
		font = .font.Copy()
		font.size *= .reverseBy(.ScaleBy)
		result = DrawTextAsk(title: 'Edit Text', content: .text, :font, justify: .justify)
		if result isnt false
			{
			.text = result.text
			result.font.size *= .ScaleBy
			.font = result.font
			.justify = result.justify
			}
		.rect = Object(left: .x1, top: .y1, right: .x2, bottom: .y2)
		}

	reverseBy(by)
		{
		return 1 / by
		}

	SetText(text)
		{
		.text = text
		}

	GetText()
		{
		return .text
		}

	GetFont()
		{
		return .font
		}

	GetJustify()
		{
		return .justify
		}

	Move(dx, dy)
		{
		if .posLocked?
			return
		.x1 += dx
		.y1 += dy
		.x2 += dx
		.y2 += dy
		.rect = Object(left: .x1, top: .y1, right: .x2, bottom: .y2)
		}

	GetName()
		{
		return .name
		}

	scaleOffset: 20
	Scale(by, print? = false)
		{
		.x1 *= by
		.x2 *= by
		.y1 *= by
		.y2 *= by
		.rect = Object(left: .x1, top: .y1, right: .x2, bottom: .y2)
		if .font.Member?(#size)
			.font.size *= print?
				? (by / (PointsPerInch / WinDefaultDpi * .scaleOffset))
				: by
		}

	ReverseScale(by, print? = false)
		{
		reverseBy = .reverseBy(by)
		.x1 *= reverseBy
		.x2 *= reverseBy
		.y1 *= reverseBy
		.y2 *= reverseBy
		.rect = Object(left: .x1, top: .y1, right: .x2, bottom: .y2)

		if .font.Member?(#size)
			.font.size *= print?
				? ((PointsPerInch / WinDefaultDpi * .scaleOffset) / by)
				: reverseBy
		}

	GetSuJSObject()
		{
		return Object('SuCanvasText', .x1, .y1, .x2, .y2, .parseText(), .font, .justify,
			id: .Id)
		}

	ToggleLock()
		{
		.posLocked? = not .posLocked?
		}
	}
