// Copyright (C) 2019 Axon Development Corporation All rights reserved worldwide.
ReportBase
	{
	New(.fontFixed, .booleansAsBox, .noFormat = false, .extraFmts = false,
		.tooltip = false)
		{
		.defaultFormat = Report.Construct(Field_string.Format)
		.Driver = new SuJsHtmlDriver
		.invalid = ToCssColor(CLR.LIGHTRED)
		}

	Default(@args)
		{
		method = args[0]
		if .Driver.Method?(method)
			return .Driver[method](@+1 args)
		throw 'method not found in Report: ' $ method
		}

	SetFormats(.columns)
		{ // create the formats for displaying the columns
		.dispFormats = Object().Set_default(false)
		if .noFormat
			return
		for col in columns
			.AddColumnsFmt(col)
		}

	AddColumnsFmt(col)
		{
		if false isnt fmt = .getFormat(col)
			{
			fmtConverted = .setWrapFormatWidth(fmt)
			_report = this
			.dispFormats[col] = Report.Construct(fmtConverted)
			.dispFormats[col].Field = col
			if fmt isnt fmtConverted
				.dispFormats[col].PrevFormat = fmt[0]
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
		catch
			if col.Prefix?("total_")
				try
					fmt = Datadict(col["total_".Size() ..]).Format
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

	GetHeaderAlign(col)
		{
		fmt = .dispFormats[col]
		if fmt isnt false and fmt.Method?("GetJustify")
			return fmt.GetJustify().Upper()
		return #LEFT
		}

	PaintRow(dataRow, model = false)
		{
		row = Object()
		for col in .columns
			// Browser header may still has 0 width columns (just shrinked)
			if model is false or model.GetColumnWidth(col) isnt false
				row[col] = .PaintCell(col, 0, 0, 0, 0, dataRow)
		return row
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
		rect = Object(:x, :y, :h, :w)
		if '' isnt invalidVal = ListControl.GetInvalidFieldData(rec, col)
			value = invalidVal
		_report = this
		res = false
		if String?(value)
			{
			value = ScintillaRichStripHTML(value)
			value = XmlEntityEncode(value)
			value = value.Ellipsis(400, atEnd:) /*= match VirtualListColModel and
				ListControl */
			if .tooltip
				value = value.Tr('\r\n', ' ')
			}
		if false is fmt = .dispFormats[col]
			{
			res = .paintWithDefaultFmt(value, rect)
			}
		else
			{
			fmt.GetSize() // format may depend on this being called
			if not fmt.Member?(#Field)
				value = rec
			rect.data = value
			rect.rec = rec
			res = .Driver.Print({ fmt.Print(@rect) })
			.addTip(fmt, value, res)
			}
		if rec.GetDefault('list_invalid_row', false) is true or
			rec.GetDefault('list_invalid_cells', #()).Member?(col)
			res.bkColor = .invalid
		return res
		}

	addTip(fmt, value, res)
		{
		// This is to keep multi-lines in tips
		if fmt.Member?(#PrevFormat)
			res.tip = fmt.PrevFormat is 'Wrap'
				? WrapFormat.Format_data(value)
				: ScintillaRichStripHTML(value)
		}

	paintWithDefaultFmt(value, rect)
		{
		.defaultFormat.TextFormat_justify = Number?(value) ? 'right' : 'left'
		if .noFormat and String?(value)
			value = Display(value)
		rect.data = value
		return .Driver.Print({ .defaultFormat.Print(@rect) })
		}

	curFont: false
	GetFont()
		{
		return .curFont
		}
	SelectFont(font)
		{
		if font is false or .fontFixed
			{
			.curFont = false
			return false
			}
		fontName = font.GetDefault(#name, 'Arial')
		fontSize =.fontSize(font.GetDefault(#size, '0'))
		weight = StdFonts.Weight(font.GetDefault(#weight, #normal))
		italic = font.GetDefault(#italic, false) is true ? 'italic' : ''

		oldfont = .curFont
		.curFont = Object(:fontName, :fontSize, :weight, :italic)
		return oldfont
		}
	fontSize(size)
		{
		if String?(size)
			size = ((10 + Number(size)) / 10/*= tenth*/) $ 'em'
		else
			size = size $ 'pt'
		return size
		}
	EnsureFont(@unused)
		{
		return .curFont
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
	}
