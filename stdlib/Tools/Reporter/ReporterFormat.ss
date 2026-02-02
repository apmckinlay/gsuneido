// Copyright (C) 2004 Suneido Software Corp. All rights reserved worldwide.
QueryFormat
	{
	Export: true
	field_to_page_break_on: false
	pageHeader: false
	sectionHeader: false
	sectionFooter: false
	New(@args)
		{
		super(@.setup(args))
		.loadCanvasHeadersAndFooters(.Data.reporterCanvas)
		.setColToPageBreakOn(.Data)
		}
	setup(args)
		{
		if not args.Member?(#Sf)
			{
			r = ReporterModel(args[0]).Report(_report.Params)
			r[0].Delete(0) // 'ReporterFormat' member
			return r[0]
			}
		return args
		}

	loadCanvasHeadersAndFooters(multiCanvasData)
		{
		.pageHeader = .LoadCanvasFromSavedData(multiCanvasData, 'page_header')
		.sectionHeader = .LoadCanvasFromSavedData(multiCanvasData, 'section_header')
		.sectionFooter = .LoadCanvasFromSavedData(multiCanvasData, 'section_footer')
		}

	LoadCanvasFromSavedData(multiCanvasData, sectionOfCanvas)
		{
		if multiCanvasData is '' or multiCanvasData is #()
			return false

		canvas_member_center = sectionOfCanvas

		return Object('CanvasRowFormat',
			multiCanvasData.GetDefault(canvas_member_center, #()),
			printParams?: sectionOfCanvas is 'page_header')
		}

	getMembersByPrefix(data, prefix)
		{
		return data.Members().Filter({ String?(it) and it.Prefix?(prefix) })
		}
	setColToPageBreakOn(data)
		{
		if .sectionHeader is false and .sectionFooter is false
			return
		// go through arg.Data's members, finding those that begin with page, and doing a
		// logical or with the values of them, checking if either page0, page1 ... pageN
		// is true, as we only need to know if 'New Page For Each' is true for any one of
		// them
		sort_cols_opts = .getMembersByPrefix(data, 'sort').Map!(
			{ |x|
				Object(
					field: data.selectFields.PromptToField(data[x]),
					prompt: data[x],
					pageBreak: data[ 'page' $ x.RemovePrefix('sort')]
					rank: Number(x.RemovePrefix('sort'))
				)
			}).Filter({ it.field isnt false and it.prompt isnt '' })

		max_rank = 0
		for i in sort_cols_opts
			{
			if i.rank >= max_rank and (i.pageBreak isnt '' or i.pageBreak isnt false)
				{
				max_rank = i.rank
				.field_to_page_break_on = i.field
				}
			}
		}
	Header()
		{
		ob = super.Header()
		if this.Member?('Headertext') and _report.GetPage() is 1
			ob = Object('Vert'
				Object('Wrap', data: .Headertext, w: _report.GetWidth())
				ob)
		if .pageHeader isnt false
			ob = Object('Vert', .pageHeader, ob)
		return ob
		}
	Output()
		{
		return .format
		}
	calcRegex: '^calc\d|^total_calc\d|^max_calc\d|^min_calc\d|^average_calc\d'
	getter_format() // once only
		{
		.format = Object('Row')
		.formats = Object()
		for col in .Columns
			{
			fmt = (col.text =~ .calcRegex
				? .calc_format(col.text)
				: Datadict(col.text).Format).Copy()
			if fmt[0] is 'Image'
				{
				imageMax = 8.5.InchesInTwips()
				fmt.width = (col.width / Reporter.LandscapeChars) * imageMax
				fmt.height = false // to respect image ratio
				}
			else
				fmt.width = col.width
			fmt.field = col.text
			prompt = .Sf.FieldToPrompt(fmt.field)
			if .Data.coloptions.Member?(prompt)
				fmt.heading = .Data.coloptions[prompt].GetDefault(#heading, prompt)
			.format.Add(fmt)
			.formats[col.text] = fmt
			}
		return .format
		}
	calc_format(col)
		{
		idx = .Data.formulas.FindIf({ it.key is col.Extract('\d\d*') })
		return .type_formats[.Data.formulas[idx].type]
		}
	getter_type_formats()
		{
		tf = Object()
		for c in CustomFieldTypes(reporter:)
			tf[c.name] = CustomFieldTypes.GetFormat(c)
		return .type_formats = tf // once only
		}
	Total()
		{
		return .Totalfields
		}
	Count()
		{
		return .Countfields
		}
	Before_(field, data)
		{
		if field is .field_to_page_break_on and .sectionHeader isnt false
			.Append(.sectionHeader)

		prompt = .Sf.FieldToPrompt(field)
		for (i = 0; i < Reporter.SortRows; ++i)
			if .Data['sort' $ i] is prompt
				break
		if .Data['show' $ i] is true
			{
			fmt = (field =~ .calcRegex
				? .calc_format(field)
				: Datadict(field).Format).Copy()
			fmt.justify = 'left'
			if fmt[0] is 'Text'
				fmt.xstretch = 1
			if fmt[0] is 'Wrap'
				fmt.width = 100
			fmt.data = data[field]
			prompt_ob = Object('Text', prompt $ ': ')
			.Append(Object('Vert',
				#(Vskip .05),
				(fmt.data is '' ? prompt_ob : Object('Horz', prompt_ob, fmt)),
				#(Vskip .05)))
			}
		return false
		}
	After_(field, data)
		{
		if field is .field_to_page_break_on and .sectionFooter isnt false
			.Append(.sectionFooter)

		prompt = .Sf.FieldToPrompt(field)
		for (i = 0; i < Reporter.SortRows; ++i)
			if .Data['sort' $ i] is prompt
				break
		if .Data['total' $ i] is true
			.After(data, subtotal:)
		if .Data['page' $ i] is true
			.Append('pg')
		return false
		}
	After(data, subtotal = false)
		{
		if subtotal is false
			{
			summary_rec = .getSummaryRec()
			data.Merge(summary_rec)
			.addDoubleLineDividers()
			}

		// min, max, average do not currently work on sort breaks
		functionParams = Object(
			Object(group: .Totalfields, prefix: 'total_', alwaysCheck: true)
			Object(group: .Minfields, prefix: 'min_', alwaysCheck: false)
			Object(group: .Maxfields, prefix: 'max_', alwaysCheck: false)
			Object(group: .Averagefields, prefix: 'average_', alwaysCheck: false)
			Object(group: .Countfields, prefix: 'count_', alwaysCheck: true))

		for fp in functionParams
			if fp.group isnt #() and (fp.alwaysCheck or subtotal is false)
					.addSortSummary(data, fp.group, fp.prefix, subtotal)

		return false
		}
	addDoubleLineDividers()
		{
		fields = Object(.Totalfields, .Minfields, .Maxfields, .Averagefields,
			.Countfields)
		fields = fields.Flatten().UniqueValues()
		if fields.Empty?()
			return
		dividers = Object('_output')
		for field in fields
			dividers.Add(Object('DoubleLine'), at: field)
		.Append(dividers)
		}
	addSortSummary(data, fields, type, subtotal)
		{
		format = Object('_output')
		for field in fields
			{
			fld_format = .formats[field].Copy()
			fld_format.data = data[type $ field]
			if type is "count_"
				fld_format.data = "Count: " $ fld_format.data
			format.Add(Object(subtotal ? "Total" : "Vert", fld_format),	at: field)
			}
		.addSummaryPromptsToFormat(type, format)
		.Append(format)
		}
	addSummaryPromptsToFormat(type, format)
		{
		if .Printsummaryfields.Empty?() or type is "total_" or type is "count_"
			return
		prompt = type.BeforeFirst("_").Capitalize()
		for f in .Printsummaryfields
			format.Add(Object('Vert',
				Object('Text' prompt, justify: "right", xstretch: 1)), at: f)
		}
	getSummaryRec()
		{
		if .Query.Has?('summarize')
			return []
		summaryOb = Object()
		for field in .Minfields
			summaryOb.Add('min ' $ field)
		for field in .Maxfields
			summaryOb.Add('max ' $ field)
		for field in .Averagefields
			summaryOb.Add('average ' $ field)
		if summaryOb.Empty?()
			return []
		summaryRec = _tran.Query1(QueryStripSort(.Query) $
			' summarize ' $ summaryOb.Join(', '))
		return summaryRec isnt false ? summaryRec : []
		}
	AfterAll()
		{
		return .GetDefault('Footertext', "") isnt ""
			? Object('Vert', 'Vskip'
				Object('Wrap', data: .Footertext, w: _report.GetWidth()))
			: false
		}
	}
