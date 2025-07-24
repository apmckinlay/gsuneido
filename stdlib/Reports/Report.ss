// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// TODO: do not call GetSize for exporting csv
ReportBase
	{
	New(@args)
		{
		.args = args
		.t = Transaction(read:)
		Suneido.ReportQueryCache = LruCache(.query1, cache_size: 100)
		.AccessPoints = Object()
		x = Params.Get_devmode(args)
		.initDimens(x, args)
		.Data = Object()
		.Params = args.GetDefault(#paramsdata, Record())
		if args.Member?("title")
			.header = Object('PageHead', args.title)
		if args.Member?("header")
			{
			.header = args.header
			args.Delete("header")
			}
		if args.Member?("footer")
			.footer = args.footer
		if args.Member?("printParams")
			.Params.printParams = args.printParams
		.page = 0
		.npages = 0
		.pageRange = args.GetDefault('pageRange', Object(from: 0, to: 9999))
		.input = InputFormat(@args)
		.toclose = Object()

		rptName = args.GetDefault(#title, '')
		if rptName is ""
			rptName = args.GetDefault('name', '')
		BookLog(rptName $ ':' $ .Params.ReportDestination,
			params: .buildBookLogParams(.Params))
		}

	buildBookLogParams(rptParams)
		{
		params = rptParams.DeepCopy()
		filterWheres = Object()
		for m in params.Members()
			if m.Suffix?('Filters')
				filterWheres[m] = GetParamsWhere(m, data: params)
		params.Delete('printParams')
		params.DeleteIf({ it.Suffix?('Filters') })
		if not filterWheres.Empty?()
			params.filterWheres = filterWheres
		return params
		}

	initDimens(x, args)
		{
		minBorder = Object().Set_default(.25) /*= default min Border*/
		if args.Member?("minBorder") and Object?(args.minBorder)
			for side in args.minBorder.Members()
				minBorder[side] = args.minBorder[side]
		.dimens = Object(
			left: Max(minBorder.left, x.left),
			right: Max(minBorder.right, x.right),
			top: Max(minBorder.top, x.top),
			bottom: Max(minBorder.bottom, x.bottom),
			width: x.width,
			height: x.height)
		.SetDimens(args.GetDefault('margins', false))
		}

	query1(query)
		{
		return Query1(query)
		}
	GetReportArgs()
		{
		return .args
		}
	GetDimens()
		{ return .dimens }
	SetDimens(margins = false)
		{
		if margins isnt false
			{
			.dimens.left = margins.left
			.dimens.right = margins.right
			.dimens.top = margins.top
			.dimens.bottom = margins.bottom
			}
		.dimens.W = (.dimens.width - .dimens.left - .dimens.right).InchesInTwips()
		.dimens.H = (.dimens.height - .dimens.top - .dimens.bottom).InchesInTwips()
		}
	GetParamsWhere(@fields)
		{
		fields.data = .Params
		GetParamsWhere(@fields)
		}

	GetUserDefinedLayout(report)
		{
		if not .Params.Member?('params_report_layout')
			return []
		layout = Query1('report_layout_designs',
			:report, rptdesign_name: .Params.params_report_layout)
		return layout isnt false and Object?(layout.rptdesign_layout)
			? layout.rptdesign_layout
			: []
		}

	Print(dc)
		{
		.SetDriver(GdiPrinterDriver(dc))
		.print(quiet?: false)
		}
	ExportCSV(file, quiet? = false, fileCl = false)
		{
		.SetDriver(CsvDriver(file, fileCl))
		return .processWithoutPage(quiet?)
		}
	processWithoutPage(quiet? = false)
		{
		_tran = .t
		_report = this
		noOutput = true
		if .header isnt false
			.Driver.Process(.Construct(.header))
		if false isnt hdr = .input.PageHeader()
			.Driver.Process(.Construct(hdr))

		status = .doWithGeneratingReportErrors()
			{
			while (false isnt (fmt = .input.Next()))
				{
				.Driver.Process(fmt)
				if not (Instance?(fmt) and
					(fmt.Base?(ColHeadsFormat) or fmt.Base?(CanvasRowFormat)))
					noOutput = false
				}
			.processFooter()
			noOutput is true ? ReportStatus.NODATA : ReportStatus.SUCCESS
			}
		.DisplayAlert(status, quiet?)
		return .Close(status)
		}
	processFooter()
		{
		if false isnt footer = .input.PageFooter()
			.Driver.Process(.Construct(footer))
		if .footer isnt false
			.Driver.Process(.Construct(.footer))
		}
	Run(driver = false, quiet? = false)
		{
		if driver is false
			driver = DummyDriver()
		.SetDriver(driver)
		return .processWithoutPage(:quiet?)
		}
	alerts: #(
		// NODATA
		1: (msg: 'No data found for report - are your options correct?',
			flag: #ICONQUESTION)
		// NOFIT
		2: (msg: 'The report cannot be displayed because one of the items is ' $
			'too big for the page. This is sometimes caused by too much ' $
			'information being entered in a field',
			flag: #ICONINFORMATION)
		3: (msg: 'Invalid matcher',
			flag: #ICONINFORMATION)
		4: (msg: 'The report cannot be displayed because one of the items has ' $
			'too much information to sort on',
			flag: #ICONINFORMATION)
		5: (msg: 'Too much data for the report to process. ' $
			'Consider adjusting the report options to filter the data.',
			flag: #ICONINFORMATION)
		6: (msg: 'There was a problem accessing the file', flag: #ICONINFORMATION)
		7: (msg: 'Report has been aborted', flag: #ICONINFORMATION)
		8: (msg: 'Too many records to summarize. ' $
			'Consider adjusting the report options to filter the data.',
			flag: #ICONINFORMATION)
		)
	DisplayAlert(status, quiet? = false, noDelay? = false)
		{
		if quiet? is true or status is ReportStatus.SUCCESS
			return
		if  Suneido.GetDefault('CheckReports', false) is true
			return
		alert = .alerts[status]
		fn = noDelay? is true ? Alert : AlertDelayed
		fn(alert.msg, 'Report', flags: MB[alert.flag], uniqueId: 'ReportAlert' $ status)
		}
	GetStatusMsg(status)
		{
		if not .IsReportError?(status) and status isnt ReportStatus.ABORT
			return ''
		return .alerts[status].msg
		}
	IsReportError?(status)
		{
		return status is ReportStatus.NODATA or .IsGeneratingReportError?(status)
		}
	PrintPDF(filename, quiet? = false, maxPages = false)
		{
		.SetDriver(PdfDriver(filename))
		return .print(quiet?, maxPages)
		}

	print(quiet?, maxPages = false)
		{
		_tran = .t
		_report = this
		.SelectFont(.Driver.GetDefaultFont())

		_env = Object(vbox: false, no_output?:, :maxPages)
		finished? = DoTaskWithPause('Working...', .printNextPage)

		// if merging pdfs' the base report has no return value.
		if Params.HasIndividualReport?()
			_env.no_output? = false

		status = not finished?
			? ReportStatus.ABORT
			: _env.vbox isnt false
				? _env.vbox 	// ReportStatus.NOFIT or INVALIDMATCHER or LONGTEMPINDEX
						// or DERIVEDTOOLARGE
				: _env.no_output?
					? ReportStatus.NODATA
					: ReportStatus.SUCCESS

		.DisplayAlert(status, quiet?)
		return .Close(status)
		}

	printNextPage(_env)
		{
		if .NextPageSuccess?(env.vbox = .NextPage()) is false
			return false

		if .page < .pageRange.from
			return true

		if .page > .pageRange.to
			{
			env.vbox = false
			return false
			}

		.checkForMaxPages(env.maxPages)
		.Driver.AddPage(.dimens)

		env.vbox.Print(
			.dimens.left.InchesInTwips(), .dimens.top.InchesInTwips(),
			.dimens.W, .dimens.H)
		.Driver.EndPage()
		env.no_output? = false
		return true
		}

	// This MaxPagesMessagePrefix is used to detect these exceptions from CreateReportFile
	MaxPagesMessagePrefix: 'Maximum pages reached'
	checkForMaxPages(maxPages)
		{
		if maxPages isnt false and .page > maxPages
			{
			.Close(ReportStatus.ABORT)
			throw .MaxPagesMessagePrefix $ " (" $ maxPages $ "). " $
				"Consider adjusting the report options to filter the data"
			}
		}

	Paint(vbox)
		{
		_tran = .t
		_report = this
		.SelectFont(.Driver.GetDefaultFont())
		vbox.Print(.dimens.left.InchesInTwips(), .dimens.top.InchesInTwips(),
			(.dimens.width - .dimens.left - .dimens.right).InchesInTwips(),
			(.dimens.height - .dimens.top - .dimens.bottom).InchesInTwips())
		}
	header: false
	footer: false
	skip_page: false
	Abort: false
	inputEmpty?: false
	NextPage()
		{
		if (.skip_page)
			{ .skip_page = false; ++.page; ++.npages; return VertFormat(); }

		.vh = 0
		++.page
		++.npages
		_tran = .t
		_report = this
		.SelectFont(.GetReportDefaultFont())

		vbox = VertFormat()
		.pageHeader(vbox)

		item = .doWithGeneratingReportErrors({ .addNextItems(vbox) })
		if .IsGeneratingReportError?(item)
			return item

		if item is 'stop'
			return false

		vbox.EraseTrailingHeader()
		if vbox.Tally() is 0
			return .inputEmpty? ? false : ReportStatus.NOFIT

		.pageFooter(vbox)

		.skipOnEvenOddPage(item)

		return vbox
		}

	doWithGeneratingReportErrors(block)
		{
		try
			return block()
		catch (err,
'REPORT: NOFIT|*regex|temp index entry size|temp index: derived too large|File|DirExists|summarize')
			{
			if err.Prefix?('temp index: derived too large')
				return ReportStatus.DERIVEDTOOLARGE
			else if err.Prefix?('temp')
				return ReportStatus.LONGTEMPINDEX
			else if err.Prefix?('REPORT: NOFIT')
				return ReportStatus.NOFIT
			else if err.Prefix?('File') or err.Prefix?('DirExists')
				return ReportStatus.FILEERROR
			else if err.Prefix?('summarize')
				return ReportStatus.SUMMARIZETOOLARGE
			else // regex
				return ReportStatus.INVALIDMATCHER
			}
		}

	IsGeneratingReportError?(status)
		{
		return status in (ReportStatus.NOFIT,
			ReportStatus.INVALIDMATCHER,
			ReportStatus.LONGTEMPINDEX,
			ReportStatus.DERIVEDTOOLARGE,
			ReportStatus.FILEERROR,
			ReportStatus.SUMMARIZETOOLARGE)
		}

	NextPageSuccess?(res)
		{
		return res isnt false and not .IsGeneratingReportError?(res)
		}

	addNextItems(vbox)
		{
		if false is item = .getNext()
			return 'stop'
		do
			{
			if (.Abort is true)
				return 'stop'

			if (item is "pg0")
				.page = 0
			if item in ('pg', 'pg0', 'pgo', 'pge')
				if vbox.ContentTally() > 0
					break
				else
					continue // ignore pg at start of page
			if false is .addItem(item, vbox)
				break
			}
			while false isnt item = .getNext()
		return item
		}
	getNext()
		{
		if false is item = .input.Next()
			.inputEmpty? = true
		return item
		}
	addItem(item, vbox)
		{
		h = item.GetSize().h
		if (item.Member?("Y") and item.Y > .vh)
			.vh = item.Y
		if (.vh + h > .dimens.H and vbox.Tally() > 0)
			{
			if (not item.Header?)
				.input.Pushback(item)
			return false
			}
		vbox.AddConstructed(item)
		item.OnPage()
		.vh += h
		return true
		}
	pageHeader(vbox)
		{
		if (.header isnt false)
			{
			hdr = .Construct(.header)
			hdr.Header? = true
			vbox.AddConstructed(hdr)
			.vh += hdr.GetSize().h
			}
		if (false isnt hdr = .input.PageHeader())
			{
			hdr = .Construct(hdr)
			hdr.Header? = true
			vbox.AddConstructed(hdr)
			.vh += hdr.GetSize().h
			}
		}
	pageFooter(vbox)
		{
		if (false isnt footer = .input.PageFooter())
			{
			footer = .Construct(footer)
			footer.Y = .dimens.H
			vbox.AddConstructed(footer)
			}
		if (.footer isnt false)
			{
			footer = .Construct(.footer)
			footer.Y = .dimens.H
			vbox.AddConstructed(footer)
			}
		}
	skipOnEvenOddPage(item)
		{
		if ((item is "pge" and .page % 2 is 0) or
			(item is "pgo" and .page % 2 is 1))
			{ .skip_page = true; }
		}
	Remaining() // vertical amount remaining on page in twips
		{
		return .dimens.H - .vh
		}

	Construct(item)
		{
		field = false
		if String?(item) and item =~ "^[a-z]"
			{
			if (item =~ "^pg[0|e|o]?$")
				return item
			field = item
			item = Datadict(item).Format
			}
		else if .fieldObject?(item)
			{
			x = item
			field = item[0]
			item = Datadict(field).Format.Copy()
			for (i in x.Members())
				if (i isnt 0)
					item[i] = x[i]
			}
		fmt = Construct(item, "Format")
		if (field isnt false)
			fmt.Field = field
		.applyProperies(item, fmt)
		return fmt
		}
	fieldObject?(item)
		{
		return Object?(item) and item.Member?(0) and String?(item[0]) and
			item[0] =~ "^[a-z]"
		}
	fmtMembers: (x, y, xmin, ymin, xstretch, ystretch, field, span, font, heading)
	applyProperies(item, fmt)
		{
		if not Object?(item)
			return
		for m in .fmtMembers
			{
			if not item.Member?(m)
				continue
			member = m.Capitalize()
			fmt[member] = item[m]
			if m in ('x', 'y', 'xmin', 'ymin')
				fmt[member] = fmt[member].InchesInTwips()
			}
		}
	GetPage()
		{ return .page; }
	GetGeneratedPages()
		{ return .npages }
	GetWidth()
		{
		return (.dimens.width - .dimens.left - .dimens.right).InchesInTwips()
		}
	GetFont()
		{ return .curfont; }
	curfont: false
	getCurFontSize()
		{
		defaultFont = .Driver.GetDefaultFont()
		if not Object?(.curfont)
			return defaultFont.size
		return .curfont.GetDefault('size', defaultFont.size)
		}
	SelectFont(font)
		{
		if not Object?(font)
			return false

		font = font.Copy()
		font.size = .GetFontSize(font)
		if font.Member?(#weight)
			font.weight = StdFonts.Weight(font.weight) // convert 'bold' to FW.BOLD

		if .Driver isnt false
			{
			font.MergeNew(.Driver.GetDefaultFont())
			.Driver.RegisterFont(font, .getCurFontSize())
			}

		oldfont = .curfont
		.curfont = font
		return oldfont
		}
	RegisterForClose(ob)
		{
		.toclose.Add(ob)
		}
	UnregisterForClose(ob)
		{
		.toclose.Remove(ob)
		}
	Close(status)
		{
		if .Driver isnt false
			status = .Driver.Finish(status)
		for (ob in .toclose)
			ob.Close()
		.t.Complete()
		Suneido.Delete('ReportQueryCache')
		BookLog('Report Closed')
		return status
		}
	}
