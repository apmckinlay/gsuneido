// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
HorzFormat
	{
	New(.data = false, .numwidth = 15, .uomwidth = 5, mask = '-###,###,###.##',
		.div = ' ', .printZero = false, .font = false, color = false, width = false)
		{
		super(
			Object(printZero is true ? 'Number' : 'OptionalNumber',
				:mask, width: numwidth * .ProrateBy(numwidth, uomwidth, width),
					field: 0, :font, :color),
			Object('Text', width: uomwidth * .ProrateBy(numwidth, uomwidth, width),
				field: 1, :font, :color))
		}
	ProrateBy(numwidth, uomwidth, width)
		{
		if width is false
			return 1
		return width / (numwidth + uomwidth)
		}

	// used by ChooseColumns
	GetDefaultWidth()
		{
		return .numwidth + .uomwidth + .div.Size()
		}

	GetWidths()
		{
		return Object(numwidth: .numwidth, uomwidth: .uomwidth + .div.Size())
		}

	GetSize(data = '')
		{
		widths = Object()
		if Object?(data)
			size = super.GetSize(data, :widths)
		else if false is ob = .buildValue(data)
			size = super.GetSize(#(0, 0), :widths)
		else
			size = super.GetSize(ob, :widths)
		if _report.GetDefault('Measuring?', false)
			{
			ratio = _report.GetCharWidth(.numwidth, .font, NumberFormat.WidthChar) /
				_report.GetCharWidth(.uomwidth, .font, TextFormat.WidthChar)
			size.w = Max((widths[1] * ratio).Ceiling(), widths[0]) +
				Max(widths[1], (widths[0] / ratio).Ceiling())
			}
		return size
		}
	Print(x, y, w, h, data = '')
		{
		if false is ob = .buildValue(data)
			return
		super.Print(x, y, w, h, ob)
		_report.Driver.SetMultiPartsRatio(Object(.numwidth, .uomwidth + .div.Size()))
		}
	ExportCSV(data = '')
		{
		if false is ob = .buildValue(data)
			return ''
		return .CSVExportString(ob[0] $ ' ' $ ob[1])
		}
	buildValue(data)
		{
		if .data isnt false
			data = .data
		split = Split_UOM(data)
		ob = Object(String(split.value), split.uom)
		if .printZero isnt true and (ob[0] is '0.00' or ob[0] is '0')
			return false
		if ob[1] isnt ''
			ob[1] = .div $ ob[1]
		return ob
		}
	}
