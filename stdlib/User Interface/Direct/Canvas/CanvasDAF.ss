// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
CanvasItem
	{
	posLocked?: false
	New(x1, y1, x2, y2, .field, font = false, .justify = 'left', fromDraw = false)
		{
		.font = font is false ? Object() : font.Copy()
		.sortPoints(x1, y1, x2, y2)
		.updateDisplay()
		if fromDraw
			{
			_report.RegisterFont(.font, .defaultFontSize)
			s = .display.text
			.x2 = .x1 + CanvasText.GetWidth(s, .font)
			.y2 = .y1 + CanvasText.GetHeight(s, .font)
			}
		}

	sortPoints(x1, y1, x2, y2)
		{
		.x1 = Min(x1, x2)
		.y1 = Min(y1, y2)
		.x2 = Max(x1, x2)
		.y2 = Max(y1, y2)
		}

	defaultFontSize: 10
	BoundingRect()
		{
		return Object(x1: .x1, y1: .y1, x2: .x2, y2: .y2)
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
		result = ResetSizeControl(0, Object(x1: .x1, y1: .y1, x2: .x2, y2: .y2))
		if (result is false)
			return
		x1 = Number(result.x1)
		y1 = Number(result.y1)
		x2 = Number(result.x2)
		y2 = Number(result.y2)
		.sortPoints(x1, y1, x2, y2)
		}

	StringToSave()
		{
		return 'CanvasDAF(x1: ' $ Display(.x1) $  ', y1: ' $ Display(.y1) $
			', x2: ' $ Display(.x2) $ ', y2: ' $ Display(.y2) $
			', field: ' $ Display(.field) $
			', font: ' $ Display(.font) $
			', justify: ' $ Display(.justify) $ ')'
		}

	ObToSave()
		{
		return Object('CanvasDAF', .x1, .y1, .x2, .y2, .field, .font.Copy(), .justify)
		}

	Resize(origx, origy, x, y)
		{
		if .Resizing?(.x1, origx)
			.x1 = x
		if .Resizing?(.y1, origy)
			.y1 = y
		if .Resizing?(.x2, origx)
			.x2 = x
		if .Resizing?(.y2, origy)
			.y2 = y
		.sortPoints(.x1, .y1, .x2, .y2)
		}

	defaultMultiplier: 20
	Scale(by, print? = false)
		{
		.x1 *= by
		.x2 *= by
		.y1 *= by
		.y2 *= by
		.rect = Object(left: .x1, top: .y1, right: .x2, bottom: .y2)
		if .font.Member?(#size)
			.font.size *= print? ?
				(by / (PointsPerInch / WinDefaultDpi * .defaultMultiplier)) : by
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
			.font.size *= print? ?
				((PointsPerInch / WinDefaultDpi * .defaultMultiplier) / by) : reverseBy
		}

	reverseBy(by)
		{
		return 1 / by
		}

	Move(dx, dy)
		{
		if .posLocked?
			return
		.x1 += dx
		.x2 += dx
		.y1 += dy
		.y2 += dy
		}

	NoPaint?: false
	Paint(data = false)
		{
		if .NoPaint? is true
			return

		if data isnt false // from report
			{
			format = Datadict(.field).Format
			if String?(format) or Class?(format)
				format = Object(format, font: .font, justify: .justify)
			else
				{
				format = format.Copy()
				format.font = .font
				format.justify = .justify
				}
			fmt = Construct(format, "Format")
			fmt.Print(.x1, .y1, .x2 - .x1, .y2 - .y1, data: data[.field], print?:)
			return
			}
		else
			.print(.display.text, .display.back)
		}

	display: (text: '', back: false)
	updateDisplay(_canvas = false)
		{
		if canvas is false or false is prompt = canvas.Send('FieldToPrompt', .field)
			.display = Object(text: 'NOT FOUND', back: Object(
				fill: CLR.ErrorColor, line: CLR.DARKRED))
		else
			{
			cols = canvas.Send('GetDAFAvailableCols')
			.display = cols.Find(prompt) is false
				? Object(text: 'NOT FOUND', back: Object(
					fill: CLR.ErrorColor, line: CLR.DARKRED))
				: Object(text: '<' $ prompt $ '>', back: Object(
					fill: CLR.ButtonGreen, line: CLR.darkgreen))
			}
		}

	print(displayText, background = false)
		{
		_report.RegisterFont(.font, .defaultFontSize)
		_report.DrawWithinClip(.x1, .y1, .x2 - .x1, .y2 - .y1)
			{
			if background isnt false
				_report.AddRect(.x1, .y1, .x2 - .x1, .y2 - .y1, 1,
					fillColor: background.fill, lineColor: background.line)

			_report.AddMultiLineText(
				displayText,
				.x1, .y1, .x2 - .x1, .y2 - .y1,
				.font, justify: .justify, color: .GetLineColor())
			}
		}

	ToggleLock()
		{
		.posLocked? = not .posLocked?
		}

	Edit(canvas)
		{
		font = .font.Copy()
		font.size *= .reverseBy(.ScaleBy)
		if false isnt result = DrawDAFAsk(canvas, .field, font, .justify)
			{
			.field = result.field
			result.font.size *= .ScaleBy
			.font = result.font
			.justify = result.justify
			.updateDisplay(canvas)
			}
		}

	DesignChanged(canvas)
		{
		oldDisplay = .display
		.updateDisplay(canvas)
		if oldDisplay isnt .display
			canvas.SyncItem(this, recursive?:)
		}

	GetField()
		{
		return .field
		}

	GetSuJSObject()
		{
		return Object('SuCanvasText', .x1, .y1, .x2, .y2, .display.text,
			.font, .justify, .display.back, id: .Id)
		}

	Valid?()
		{
		return .display.back isnt false and .display.back.fill isnt CLR.ErrorColor
		}
	}
