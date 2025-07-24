// Copyright (C) 2000 Suneido Software Corp.
ReportBase
	{
	fontFixed:		false	// false allows the formats to set the font
							// true when a fontSize is passed to ListControl
	booleansAsBox:	false	// false lets the formats display booleans

	noFormat:		false

	SetBackgroundBrush(brush, selected = false)
		{
		return .Driver.SetBackgroundBrush(brush, darkBackground?: selected)
		}

	New(fontFixed, booleansAsBox, noFormat = false, .extraFmts = false, .tooltip = false)
		{
		.fontFixed = fontFixed
		.booleansAsBox = booleansAsBox
		.noFormat = noFormat
		.defaultFormat = Report.Construct(Field_string.Format)
		.SetDriver(new GdiPreviewDriver)
		}
	SetFormats(columns)
		{ // create the formats for displaying the columns
		.dispFormats = Object().Set_default(false)
		.colLeftJustified = Object().Set_default(true)
		if .noFormat
			return
		for col in columns
			.AddColumnsFmt(col)
		}
	AddColumnsFmt(col)
		{
		if false isnt fmt = .getFormat(col)
			{
			fmt = .setWrapFormatWidth(fmt)
			_report = this
			fmt = .dispFormats[col] = Report.Construct(fmt)
			fmt.Field = col
			.colLeftJustified[col] = false
			if fmt.Method?("GetJustify") and fmt.GetJustify() is "left"
				if .fontFixed or fmt.TextFormat_font is false
					.colLeftJustified[col] = true
			}
		}
	getFormat(col)
		{
		fmt = false
		try
			{
			fmt = Datadict(col).Format
			if .extraFmts isnt false and .extraFmts.Member?(col)
				fmt = fmt.Copy().Merge(.extraFmts[col])
			}
		catch (e)
			{
			if col.Prefix?("total_")
				{
				try
					fmt = Datadict(col["total_".Size() ..]).Format
				catch (err)
					SuneidoLog('ERROR: (CAUGHT) ' $ err, calls:,
						caughtMsg: 'failed to get total format')
				}
			else
				SuneidoLog('ERROR: (CAUGHT) ' $ e, calls:,
					caughtMsg: 'failed to get format')
			}
		return fmt
		}
	setWrapFormatWidth(fmt)
		{
		if not Object?(fmt)
			return fmt
		if fmt[0] is "Wrap"
			{
			fmt = fmt.Copy()
			fmt[0] = "Text"
			fmt.Delete("w")
			fmt.width = fmt.Member?('width') ? fmt.width : 40
			}
		if fmt[0] is 'ScintillaRichWrap'
			{
			fmt = fmt.Copy()
			fmt[0] = "ScintillaRichText"
			fmt.Delete("w")
			fmt.width = fmt.Member?('width') ? fmt.width : 40
			}
		return fmt
		}
	AssumeLeftJust?(col)		// returns true if column is certainly left justified
		{
		return .colLeftJustified[col]
		}
	GetHeaderAlign(col)
		{
		fmt = .dispFormats[col]
		if fmt isnt false and fmt.Method?("GetJustify")
			return HDF.STRING | HDF[fmt.GetJustify().Upper()]
		return HDF.STRING
		}
	CompareRows(col, x, y)
		{
		if false isnt fmt = .dispFormats[col]
			if fmt.Method?("DataToString")
				return fmt.DataToString(x[col], x) < fmt.DataToString(y[col], y)
		if Object?(x[col]) or Object?(y[col])
			return Display(x[col]) < Display(y[col])
		return x[col] < y[col]
		}
	PaintCell(col, x, y, w, h, rec)
		{
		try
			{
			_showFormulaError = true
			value = rec[col]
			}
		catch (err, "SHOW")
			{
			value = err.RemovePrefix("SHOW: ")
			}

		if String?(value)
			{
			value = ScintillaRichStripHTML(value)
			value = value.Ellipsis(400, atEnd:) /*= to speed up */
			if .tooltip
				value = value.Tr('\r\n', ' ')
			}

		if .booleansAsBox and Boolean?(value)
			DrawFrameControl(.Driver.GetDC(),
				Object(left: x - 1, top: y, right: x + 12, bottom: y + 13), /*= margins */
				DFC.BUTTON,
				DFCS.BUTTONCHECK | DFCS.FLAT | (value is true ? DFCS.CHECKED : 0))
		else
			{
			// Cannot restore the original SelectObject value due to code structure
			// and performance issues with SelectObject.
			if not .fontFixed
				SelectObject(.Driver.GetDC(), Suneido.hfont)
			.paint(col, value, Object(x, y, w, h), rec)
			}
		}
	paint(col, value, rect, rec)
		{
		if '' isnt invalidVal = ListControl.GetInvalidFieldData(rec, col)
			value = invalidVal
		_report = this
		fullDisplay = false
		if .tooltip
			{
			fullDisplay = _report.TooltipFullDisplay =
				rec.GetInit('vl_full_display', Object())
			fullDisplay.CurrentCol = col
			fullDisplay[col] = ''
			fullDisplay.CellEllipsized = false
			}
		if false is .dispFormats[col]
			.paintWithDefaultFmt(value, rect, col)
		else
			{
			fmt = .dispFormats[col]
			fmt.GetSize() // format may depend on this being called
			if not fmt.Member?('Field')
				value = rec
			rect.data = value
			rect.rec = rec
			rect.textOnly? = true
			fmt.Print(@rect)
			}
		if .tooltip
			{
			if not fullDisplay.CellEllipsized
				fullDisplay.Delete(col)
			else
				fullDisplay[col] = fullDisplay[col].Trim()
			_report.Delete('TooltipFullDisplay')
			}
		}
	paintWithDefaultFmt(value, rect, col)
		{
		.defaultFormat.TextFormat_justify = Number?(value) ? 'right' : 'left'
		if .noFormat and String?(value)
			value = Display(value)
		rect.data = value
		.defaultFormat.Print(@rect)
		if .colLeftJustified[col] and Number?(value)
			.colLeftJustified[col] = false
		}
	SetDC(dc)
		{
		.Driver.SetDC(dc)
		.curFont = false
		}
	MeasureWidth(col, rec)
		{
		_report = this
		_report.Measuring? = true
		if false is fmt = .dispFormats[col]
			{
			fmt = .defaultFormat
			value = rec[col]
			if String?(value)
				value = Display(value)
			}
		else
			value = not fmt.Member?('Field') ? rec : rec[col]
		sz = fmt.GetSize(data: value)
		_report.Delete('Measuring?')
		return sz
		}
	curFont: false
	SelectFont(font)
		{
		if font is false or .fontFixed
			return false
		fontName = font.GetDefault(#name, Suneido.logfont.lfFaceName)
		fontSize = StdFonts.FontSize(font.GetDefault(#size, '0'))
		charset = font.GetDefault(#charset, 'ANSI')
		weight = StdFonts.Weight(font.GetDefault(#weight,
			Suneido.logfont.GetDefault(#lfWeight, 'NORMAL')))
		italic = font.GetDefault(#italic,
			Suneido.logfont.GetDefault(#lfItalic, 0))
		id = fontName $ charset $ fontSize $ weight $ italic
		if not .Fonts.Member?(id)
			{
			lf = Object(
				lfFaceName: fontName,
				lfCharSet: CHARSET[charset],
				lfHeight: -fontSize *
					GetDeviceCaps(.Driver.GetDC(), GDC.LOGPIXELSY) / PointsPerInch,
				lfWeight: weight,
				lfItalic: italic)
			if 0 is f = CreateFontIndirect(lf)
				throw "Format: couldn't create font"
			.Fonts[id] = f
			}
		// Cannot restore the original SelectObject value as the end point is ambiguous.
		// As a result, this cannot be changed to use DoWithHdcObjects or WithHdcSettings.
		SelectObject(.Driver.GetDC(), .Fonts[id])
		oldfont = .curFont
		.curFont = font
		return oldfont
		}
	GetFont()
		{
		return .curFont
		}
	}
