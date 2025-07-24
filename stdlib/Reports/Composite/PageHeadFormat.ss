// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
VertFormat
	{
	Header?: true
	Xstretch: 1
	New(@args)
		{
		super(@.format(args))
		}
	params: #(WrapItems)
	format(args)
		{
		header = Object(Object('PageTitle', args.GetDefault('title2', PageHeadName())))
		titleFont = args.GetDefault('titleFont', Object(size: 18, weight: FW.SEMIBOLD))
		for i in args.Members(list:)
			header.Add(Object('Horz',
				'Hfill',
				Object('Text',
					args[i] $ (i is 0 and _report.Params.report_option isnt ""
						? ' - ' $ _report.Params.report_option : ""),
					font: titleFont),
				'Hfill'
				))
		header.Add('Vskip')

		.params  = PrintParams()
		if .params isnt Object('WrapItems')
			header.Add(.params, "Vskip")
		return header
		}

	// only print the params
	ExportCSV(data /*unused*/ = '')
		{
		return Opt(_report.Construct(.params).ExportCSV(), '\n,\n')
		}
	}
