// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
HorzFormat
	{
	New(formats, widths, font)
		{
		super(@.makeitems(formats, widths, font))
		}
	makeitems(formats, widths, font)
		{
		skip = widths.Member?(#skip)
			? ['Hskip', widths.skip / 1.InchesInTwips()]
			: 'Hskip'
		items = Object()
		str = ''
		for (i = 0; formats.Member?(i); ++i)
			{
			if i > 0
				items.Add(skip)
			fmt = formats[i]
			heading = fmt.Member?("Heading") ? fmt.Heading :
				fmt.Member?("Field") ? Heading(fmt.Field) : ""
			str $= .CSVExportString(heading) $ ','
			items.Add(Object('Vert',
				Object("Wrap", heading, w: widths[i], justify: "center"),
				Object("Hline", widths[i], xstretch: 0)))
			}
		.heading = .CSVExportLine(str)
		if font isnt false
			{
			font = font.Copy()
			normalFontWeight = 400
			extraFontWeight = 200
			font.weight = font.GetDefault('weight', normalFontWeight) + extraFontWeight
			items.font = font
			}
		return items
		}
	ExportCSV(data /*unused*/= '')
		{
		return .heading
		}
	}
