// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
// FIXME: there is double lines if border is 1 in the grid
Format
	{
	Xstretch: ""
	Ystretch: ""
	New(format, .width = 0, .top = false, font = false, .border = 0, .access = false,
		.lead = 35)
		{
		.data = Object()
		for row in format
			{
			rows = Object()
			for item in row
				{
				if font isnt false
					{
					item = false is Object?(item) ? Object(item) : item.Copy()
					if not item.Member?('font')
						item.font = font
					}
				rows.Add(_report.Construct(item))
				}
			.data.Add(rows)
			}
		.initStretch()
		}

	initStretch()
		{
		.totalXStretch = 0
		.totalYStretch = 0
		.maxColXStretches = Object().Set_default(0)
		.maxRowYStretches = Object().Set_default(0)

		for row in .data.Members()
			{
			pos = 0
			for col in .data[row].Members()
				{
				item = .data[row][col]
				if item.Xstretch > .maxColXStretches[pos]
					.maxColXStretches[pos] = item.Xstretch

				if item.Ystretch > .maxRowYStretches[row]
					.maxRowYStretches[row] = item.Ystretch

				pos += item.Member?('Span') and item.Span >  0 ? item.Span : 1
				}
			}

		for maxColXstretch in .maxColXStretches
			.totalXStretch += maxColXstretch
		for maxRowYStretch in .maxRowYStretches
			.totalYStretch += maxRowYStretch

		if .Xstretch is ""
			.Xstretch = .totalXStretch
		if .Ystretch is ""
			.Ystretch = .totalYStretch
		}

	GetSize(data/*unused*/ = #{})
		{
		w = .width
		h = numCols = 0
		.ascents = Object().Set_default(0)
		.descents = Object().Set_default(0)
		.colwidths = Object().Set_default(0)
		.rowheights = Object().Set_default(0)

		for row in .data
			if (row.Size() > numCols)
				numCols = row.Size()

		for row in .data.Members()
			{
			pos = 0
			for (col in .data[row].Members())
				{
				item = .data[row][col]
				size = item.GetSize()
				if (size.h - size.d > .ascents[row])
					.ascents[row] = size.h - size.d
				if (size.d > .descents[row])
					.descents[row] = size.d
				if (size.w > .colwidths[pos])
					.colwidths[pos] = size.w
				pos += item.Member?('Span') and item.Span >  0 ? item.Span : 1
				}
			.rowheights[row] = .ascents[row] + .descents[row]
			h += .rowheights[row] + .lead + .border
			}

		if .width is 0
			w = .calcTotalItemWidth()

		d = h - .ascents[0]
		return Object(:w, :h, d: .top is true ? d : 0)
		}

	calcTotalItemWidth()
		{
		w = 0
		for colwid in .colwidths
			w += colwid + .lead + .border
		return w + .border
		}

	access: false
	Print(orgx, orgy, w, h, data = #{})
		{
		y = orgy
		size = .GetSize()
		hf = .totalXStretch is 0
			? 0
			: Max(0, (w - .calcTotalItemWidth()) / .totalXStretch)
		vf = .totalYStretch is 0
			? 0
			: Max(0, (h - size.h) / .totalYStretch)
		for (row in .data.Members())
			{
			x = orgx
			height = .rowheights[row] + (vf * .maxRowYStretches[row])
			accumulatedwidth = span = pos = 0
			for (col in .data[row].Members())
				{
				borderWidth = .colwidths[pos] + (hf * .maxColXStretches[pos]) +
					.lead + .border
				item = .data[row][col]
				size = item.GetSize()
				yi = y + .ascents[row] + size.d - size.h
				if (item.Member?('Span'))
					{
					for (span = item.Span, j = 0; j < span - 1; ++j)
						{
						borderWidth += .colwidths[pos + j + 1]  +
							(hf * .maxColXStretches[pos + j + 1]) + .lead + .border
						}
					pos += span - 1
					}
				accumulatedwidth += borderWidth
				if .Xstretch is 0 and .width isnt 0 and accumulatedwidth > .width
					{
					borderWidth -= accumulatedwidth - .width
					}
				if borderWidth <= 0
					continue
				if .border isnt 0
					{
					h = height + .lead + .border
					_report.AddRect(x, y, borderWidth, h, thick: .border)
					}
				++pos
				item.Print(x + .lead / 2, yi + .lead / 2,
					borderWidth - .lead - .border, height)
				if .access isnt false and .access.Size() > row and
					.access[row] isnt false and .access[row].Size() > col and
					.access[row][col] isnt false
					.Hotspot(x + .lead / 2, yi + .lead / 2, borderWidth, height, data,
						access: .access[row][col])
				x += borderWidth
				}
			y += height + .lead + .border
			}
		}
	ExportCSV(data /*unused*/= #())
		{
		csv = ''
		for (row in .data.Members())
			{
			line = ''
			for (col in .data[row].Members())
				{
				item = .data[row][col]
				line $= item.ExportCSV() $ ','
				}
			csv $= .CSVExportLine(line)
			}
		return csv
		}
	}
