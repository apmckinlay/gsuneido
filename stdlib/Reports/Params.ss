// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// On_Preview, On_Print, On_Page_Setup can be called as static methods
Controller
	{
	Title: 'Report'
	New(@report)
		{
		super(.Controls(report))
		.Ensure()
		.validField = .report.Member?('validField') ? .report.validField : false
		.setParams = .report.Member?('SetParams') ? .report.SetParams : Object()
		.loadParams(.report)
		}
	Ensure(table = 'params')
		{
		Database("ensure " $ table $
			" (user, report, params, report_options, params_TS)
			index(report)
			key (user,report)")
		}
	Controls(report)
		{
		.report = report
		if (not report.Member?('Params'))
			report.Params = 'Skip'
		.disableFieldProtects = report.GetDefault('disableFieldProtectRules', false)
		pageCount = report.GetDefault('pageCount', false)
		return Object('Vert',
			'Skip',
			.body(report),
			'Skip',
			.ParamsButtons(.reporter?(report), .export?(report), pageCount)
			'Skip')
		}
	body(report)
		{
		params = report.Params
		trim = Object?(params)
			? params.GetDefault('trim', true)
			: true
		chooseLayout = .ChooseLayoutControl(.report)
		if Object?(params) and params.Flatten().Has?('ParamsChooseLayout')
			chooseLayout = ''
		paramsVert = ['Vert',
			['Horz', params, 'Fill',
				['Vert', ParamsLogo(report), chooseLayout, xstretch: 0]]
			name: 'paramContainer']

		params = ['Scroll', paramsVert, :trim, xmin: 700, ymin: 350]
		noPresets? = .report.GetDefault('NoPresets', false)
		if .report.Member?('title')
			params = Object('Vert',
				Object('Horz', 'Skip'
					Object('TitleNotes', .report.title,
						extra: .report.Member?("name") and noPresets? is false
							? ['Presets', .report.name, .report.title,
								saveCtrlsOnly:]
							: #(Skip 0))),
				#(Skip 5), params)
		return Object('Record', params,
			disableFieldProtectRules: .disableFieldProtects)
		}

	HelpButton_HelpPage()
		{
		book_option = 0
		if false is currentbook = Suneido.GetDefault('CurrentBook', false)
			return book_option
		option = .Name
		if .report.Member?('HelpOption')
			option = .report.HelpOption
		if false isnt rec = QueryFirstBookOption(currentbook, option)
			book_option = rec.path $ '/' $ rec.name
		return book_option
		}

	ChooseLayoutControl(report, xstretch = 1)
		{
		return .hasCustomLayouts?(report)
			? [#ParamsChooseLayout, report.name, xstretch]
			: ''
		}

	hasCustomLayouts?(report)
		{
		return report.Member?(#name) and ReportLayoutDesign.Customizable?(report.name)
		}

	ReportName()
		{
		return .report.name
		}

	reporter?(report)
		{
		return ((report.Member?(0) and Object?(report[0])) and
			(report[0][0] is 'ReporterFormat' or
				report[0][0] is ReporterFormat or
				report[0][0] is 'ReporterCanvasFormat' or
				report[0][0] is ReporterCanvasFormat))
		}

	reporterForm?(report)
		{
		return report[0][0] in ('ReporterCanvasFormat', ReporterCanvasFormat)
		}

	export?(report)
		{
		if not report.Member?(0)
			return false

		if .reporter?(report)
			return not .reporterForm?(report)

		rpt = .getReportFormat(report)
		return Class?(rpt) and rpt.Member?('Export') and rpt.Export
		}
	getReportFormat(report)
		{
		rpt = report[0]
		if Object?(rpt) and rpt.Member?(0)
			rpt = rpt[0]
		if String?(rpt)
			rpt = Global(rpt.Suffix?('Format') ? rpt : rpt $ 'Format')
		return rpt
		}

	ExtraButtons: ()
	Buttons: #(Print, Preview, PDF)
	buttonPad: false
	PreviewButtons: #(Print, PDF)
	ParamsButtons(reporter_link, export = false, pageCount = false)
		{
		ob = Object(.buttonFormat(), 'Skip')

		.addParamButtons(ob, .Buttons)
		if export
			.addParamButtons(ob, #(Export))
		.addParamButtons(ob, .ExtraButtons)
		if pageCount
			ob.Add(#(Button 'Page Count'), 'Skip')
		ob.Add(Object('Button' 'Page Setup', pad: .buttonPad),
			'Skip',
			#(CheckBox 'Print Lines' name: print_lines,
				tip: 'Print lines between records'),
			'Skip',
			reporter_link is true ?
				#(LinkButton name: 'Reporter' command: 'ReporterReport') : '')
		ob.pad = 10
		return Object('Vert' ob name: 'generate_msg')
		}

	buttonFormat()
		{
		allButtons = .Buttons.Copy().Append(.ExtraButtons)
		char = allButtons.MaxWith({|x| x.Size() }).Size()
		if char> 15 and allButtons.Size() > 5 /*= more than 5 buttons and
			the largest button char > 15, switch to Horz */
			{
			.buttonPad = 40
			return 'Horz'
			}
		return 'HorzEqual'
		}

	addParamButtons(ob, buttons)
		{
		for button in buttons
			{
			buttonCtrl = button.Has?("PDF")
				? Object('PDFButton' button) : Object('Button' button, pad: .buttonPad)
			ob.Add(buttonCtrl, 'Skip')
			.PreviewButtons = .PreviewButtons.Copy().AddUnique(button)
			}
		}

	report: false
	ParamsValid?()
		{
		return .params_valid?() is true ? 0 : 1
		}

	params_valid?(silent = false)
		{
		msg = ""
		if (.validField isnt false)
			{
			params = .Vert.Data.Get()
			msg = params[.validField]
			}
		// set RecordControl dirty so that its control's
		// Valid methods are checked
		.Vert.Data.Dirty?(true)
		if ((invalid_fields = .Vert.Data.Valid()) isnt true)
			{
			if (msg isnt "")
				msg $= "\n\n"
			msg $= invalid_fields
			}
		if (msg isnt "")
			{
			if (not silent)
				.AlertWarn("Invalid Parameter(s)", msg)
			return false
			}
		return true
		}
	DisablePreviewDialog() // called by ReportCheck
		{
		if .report.GetDefault('previewDialog', false)
			.report.previewDialog = false
		}
	On_Preview(@report)
		{
		if (Object?(.report))
			report = .report

		if not .paramsWindowValid?(report)
			return false

		if not PreviewLimiter().BeforeOpen()
			return false

		.checkAndResetParams(report)
		.setPreviewParams(report)

		report.paramsdata.ReportDestination = 'preview'
		.SetExtraParamsData(report.paramsdata)
		.add_print_lines(report.paramsdata)

		if false is .verifyDevMode(report)
			return false

		if true is .addFilterIfSlowQuery(report)
			return false
		if report.GetDefault('pageCountOnly', false)
			return .previewForPageCount(report)
		w = .runPreview(report)
		return w // is this used? could be either dialog result or a Window ???
		}

	setPreviewParams(report)
		{
		if not report.Member?('PreviewParams')
			return
		preParams = report.PreviewParams
		for field in preParams.Members()
			if preParams[field] isnt #()
				report.paramsdata[field] = preParams[field]
		}

	On_Page_Count()
		{
		.report.pageCountOnly = true
		if Number?(pages = .On_Preview(.report))  // won't be number if params invalid
			.AlertInfo('Page Count', 'Estimated Pages: ' $ pages)
		.report.pageCountOnly = false
		}

	previewForPageCount(report)
		{
		w = .buildPreviewWindowForPageCount(report)
		previewCtrl = w.Ctrl
		previewCtrl.On_Last()
		pages = previewCtrl.GetNumPages()
		DestroyWindow(previewCtrl.Window.Hwnd)
		return pages
		}

	buildPreviewWindowForPageCount(report)
		{
		return Window(Object(PreviewControl, report, this, extraButtons: .PreviewButtons),
			show: false)
		}

	runPreview(report)
		{
		extraButtons = .export?(report)
			? .PreviewButtons.Copy().AddUnique(#Export)
			: .PreviewButtons

		if (report.Member?('previewDialog') and report.previewDialog is true)
			{
			keep_size = 'Print Preview Dialog'
			if report.Member?('PreviewParams') and report.Member?('name')
				keep_size $= ' - ' $ report.name
			w = ModalWindow(Object("Preview", report, this, :extraButtons),
				:keep_size, useDefaultSize:,
				onDestroy: .onDestroyModalPreview)
			}
		else
			{
			w = Window(Object(PreviewControl, report, this, :extraButtons),
				keep_placement:, useDefaultSize:, onDestroy: .checkNoPage,
				excludeModalWindow: .closeDialog?() ? .Window : false)
			.closeDialog()
			}
		return w
		}

	addFilterIfSlowQuery(report)
		{
		if report.Member?('suppressSlowQuery')
			report.Delete('suppressSlowQuery')
		if false is fieldPrompt = .getSlowQueryFilterFieldPrompt(report)
			return false

		query = .getQuery(report)
		if not String?(query)
			return false

		queryState = Object(sortCol: false, presets: Object())
		filterFields = fieldPrompt.GetFields()
		if false isnt SlowQuery.Validate(query, filterFields, after: false, :queryState,
			indexes: report.GetDefault('slowQueryIndexes', false))
			return false

		if not queryState.Member?('filter')
			{
			report.suppressSlowQuery = true
			// user didn't pick a filter from suggestion window, keep running report
			return false
			}

		SlowQuery.AddParamsIndexedFilter(
			queryState.filter, .Vert.Data, report.slowQueryFilter)
		return true
		}

	getSlowQueryFilterFieldPrompt(report)
		{
		if false is report.GetDefault('slowQueryFilter', false)
			return false

		if report.GetDefault('previewWindow', false)
			return false
		if false is filters = .FindControl(report.slowQueryFilter)
			return false
		if false is fieldPrompt = filters.FindControl('condition_field')
			return false
		if not .paramsScreen?(this)
			return false
		return fieldPrompt
		}

	getQuery(report)
		{
		if not report.Member?(0)
			return false
		if .reporter?(report)
			return false
		fmt = .getReportFormat(report)
		if not Class?(fmt)
			return false
		if not fmt.Base?(QueryFormat)
			return false
		if fmt.Base?(ObjectFormat)
			return false

		_report = new ParamsReportClassForQuery
		_report.Params = report.paramsdata
		fmtInstance = new ParamsFormatClassForQuery(fmt)
		fmtInstance.Params = report.paramsdata
		return fmtInstance.Query()
		}

	paramsScreen?(params)
		{
		return Instance?(params) and params.Base?(Params)
		}

	checkNoPage()
		{
		if .noPage
			{
			.noPage = false
			.ClearFocus() // prevent multiple alerts if users holds down enter key
			Report.DisplayAlert(ReportStatus.NODATA)
			}
		}

	noPage: false
	SetNoPage(.noPage = true)
		{
		}

	CloseDialog(report)
		{
		if not report.GetDefault(#from_preview, false)
			.closeDialog()
		}

	onDestroyModalPreview()
		{
		if Instance?(this)
			.Defer(uniqueID: #ParamsCloseDialog)
				{
				.checkNoPage()
				.closeDialog()
				}
		}

	closeDialog()
		{
		if .closeDialog?()
			.Window.Result(true)
		}

	closeDialog?()
		{
		return .Member?("Window") and .Window.Method?("Result")
		}

	verifyDevMode(report)
		{
		x = .Get_devmode(report)
		hdm = .globalAllocData(x.devmode)
		hwnd = .hwnd(report)
		if hdm is 0 or x.devnames is ''
			{
			if false is hdm = .useDefaultPrinter(x, hwnd)
				return false
			}
		if not .update_pdc(hdm, x, report, fromPreview?:)
			{
			hdm = .useDefaultPrinter(x, hwnd)
			if hdm is false or not .update_pdc(hdm, x, report)
				return false
			}
		if (hdm isnt 0)
			GlobalFree(hdm)
		return true
		}
	globalAllocData(str)
		{
		if str is ""
			return 0
		return GlobalAllocData(str)
		}
	useDefaultPrinter(x, hwnd)
		{
		pd = Object(
			lStructSize: PRINTDLGEX.Size(),
			hwndOwner: hwnd,
			Flags: PD.RETURNDEFAULT | PD.NOPAGENUMS,
			nCopies: 1,
			nStartPage: 0xFFFFFFFF)

		if PrintDlgEx(pd) isnt 0 or pd.hDevMode is 0 or pd.hDevNames is 0
			{
			.AlertWarn("Print Preview", "Unable to get default printer")
			return false
			}
		x.devmode = GlobalData(pd.hDevMode)
		x.devnames = GlobalData(pd.hDevNames)
		hdm = pd.hDevMode
		pd.hDevMode = 0
		.free(pd)
		return hdm
		}

	add_print_lines(paramsdata)
		{
		if not paramsdata.Member?('PrintLines')
			{
			ctrl = .FindControl('print_lines')
			paramsdata.PrintLines = ctrl is false ? false : ctrl.Get()
			}
		}

	On_Print(@report)
		{
		if (Object?(.report))
			report = .report

		if not .paramsWindowValid?(report)
			return

		fromPreview? = report.GetDefault("from_preview", false)
		.checkAndResetParams(report)
		report.paramsdata.ReportDestination = 'printer'
		.SetExtraParamsData(report.paramsdata)
		.add_print_lines(report.paramsdata)
		if not fromPreview? and true is .addFilterIfSlowQuery(report)
			return

		params_window = .hwnd(report)
		if false is .preparePrinter(params_window, report)
			return

		// print to pdf printer driver does not block params window
		// if the params's parent window has hWndParent
		EnableWindow(params_window, false)
		Finally({ .PrintReport(report) },
			{ EnableWindow(params_window, true) })
		.CloseDialog(report)
		}

	GetParamsControl(ctrlName)
		{
		ctrl = .FindControl(ctrlName)
		return ctrl
		}

	preparePrinter(params_window, report)
		{
		pd = .newPRINTDLG(params_window, report)
		x = .Get_devmode(report)
		pd.hDevMode = .globalAllocData(x.devmode)
		pd.hDevNames = .globalAllocData(x.devnames)
		if false is pd = .openPrintDlg(pd, report)
			return false
		x.devmode = GlobalData(pd.hDevMode)
		x.devnames = GlobalData(pd.hDevNames)
		ok = .update_pdc(pd.hDevMode, x, report)
		.free(pd)
		if not ok
			return false
		if PD.PAGENUMS is (pd.Flags & PD.PAGENUMS)
			report.pageRange = Object(from: pd.lpPageRanges.nFromPage,
				to: pd.lpPageRanges.nToPage)
		return true
		}

	openPrintDlg(pd, report)
		{
		if true isnt result = .printDlgAndValid(pd)
			{
			.free(pd)
			if result is false and CommDlgExtendedError() is 0
				return false// user cancelled

			QueryApply1(.devmode_query(.devmode_reportname(report)))
				{
				it.devnames = it.devmode = ''
				it.Update()
				}
			.AlertError("Print", .printerError)
			return false
			}
		return pd
		}

	printerError: "There is a problem using the chosen printer with existing Page Setup" $
		"\nThis might be caused by printer driver or connection issue" $
		"\nThe invalid printer settings are cleaned up for this report" $
		"\nPlease choose the Page Setup and print again"
	invalidDc: 'invalid dc'
	printDlgAndValid(pd)
		{
		try
			{
			Dialog.DoWithWindowsDisabled({ result = PrintDlgEx(pd) }, pd.hwndOwner)
			if result isnt 0 or pd.dwResultAction isnt 1 or pd.hDevNames is 0
				return false
			}
		catch(err, "win32 exception: ACCESS_VIOLATION")
			{
			SuneidoLog('ERRATIC: ' $ err)
			.AlertError("Print", .printerError)
			return false
			}

		// validate pd generated by printer dialog
		devnames = GlobalData(pd.hDevNames)
		if 0 is tempDC = .createDC(devnames, pd.hDevMode)
			return .invalidDc
		DeleteDC(tempDC)
		return true
		}

	createDC(devnames, hDevMode)
		{
		adr = GlobalLock(hDevMode)
		Assert(adr isnt 0)
		dc = CreateDC(0, .devnamesDevice(devnames), 0, adr)
		GlobalUnlock(hDevMode)
		return dc
		}
	devnamesDevice(s)
		{
		deviceOffset = s[2].Asc() + 256 * s[3].Asc() /*= offsets in DEVNAMES */
		return s[deviceOffset..].BeforeFirst('\x00')
		}

	paramsWindowValid?(report)
		{
		if not this.Member?("Vert") or report.Member?("from_preview")
			return true
		return .params_valid?()
		}

	SetExtraParamsData(paramsData /*unused*/)
		{
		}

	On_Export(@report)
		{
		if Object?(.report)
			report = .report

		if false is filename = Dialog.DoWithWindowsDisabled(
			{ .getSaveFileName(report, 'csv') }, .hwnd(report))
			return

		if true isnt msg = CatchFileAccessErrors(filename, { .SaveCSV(filename, report) })
			.AlertInfo("Export", msg)

		.afterSaveFile(report, 'csv', filename)
		.CloseDialog(report)
		}

	SaveCSV(filename, report = false) // also called by ReportCheck
		{
		if report is false
			report = .report

		fromPreview? = report.GetDefault("from_preview", false)
		.checkAndResetParams(report)
		report.paramsdata.ReportDestination = 'csv'
		.SetExtraParamsData(report.paramsdata)
		if not fromPreview? and true is .addFilterIfSlowQuery(report)
			return ReportStatus.SUCCESS
		Report(@report).ExportCSV(filename)
		}

	PrintReport(report = false)
		{
		if report is false
			report = .report
		Report(@report).Print(Suneido.pdc)
		}

	RunWithNoOutput(@report)
		{
		if (Object?(.report))
			report = .report
		if (this.Member?("Vert"))
			{
			if (not .params_valid?())
				return false
			report.paramsdata = .Vert.Data.GetControlData()
			.update_params(report)
			}
		else if not report.Member?('paramsdata')
			report.paramsdata = Record()
		report.paramsdata.ReportDestination = 'nooutput'
		.SetExtraParamsData(report.paramsdata)
		.add_print_lines(report.paramsdata)
		status = Report(@report).Run()
		report.paramsdata.LastReportStatus = status
		return report.paramsdata
		}

	hwnd(report)
		{
		params_window = this.Member?("Window") ? .Window.Hwnd : GetActiveWindow()
		if (report.Member?("previewWindow"))
			{
			params_window = report.previewWindow
			report.Delete("previewWindow")
			}
		return params_window
		}

	On_PDF_Save_to_file(@report)
		{
		if Object?(.report)
			report = .report

		if false is filename = Dialog.DoWithWindowsDisabled(
			{ .getSaveFileName(report, 'pdf') }, .hwnd(report))
			return

		if false is .CreatePdfWithEmailAttachments(report, filename, compress?: 'auto')
			return
		.afterSaveFile(report, 'pdf', filename)

		.CloseDialog(report)
		}

	maxSizeInMb: 10
	CreatePdfWithEmailAttachments(report, filename, alwaysMerge? = false, quiet? = false,
		compress? = false)
		{
		if true isnt result = .pdf(report, filename, :quiet?)
			return result

		unmergeableFiles = #()
		if .mergePdf?(alwaysMerge?, report)
			{
			result = .orderAttachments(report, quiet?)
			if not Object?(result) // false or string
				return result
			unmergeableFiles = result.unmergeableFiles

			if not Boolean?(compress?)
				{
				if String?(fileSize = EmailAttachment.CalculateTotal(
					result.mergeableFiles.AddUnique(filename)))
					{
					.deleteFile(filename, 'Could not clean up PDF temp file')
					.alertError(quiet?, fileSize)
					return .returnVal(quiet?, fileSize)
					}

				compress? = fileSize > .maxSizeInMb.Mb()
				}

			if false is .beforeMerge(report, result.mergeableFiles, filename, compress?)
				return false

			if true isnt msg = .mergePdf(filename, result.mergeableFiles, compress?)
				{
				.alertError(quiet?, msg)
				return .returnVal(quiet?, msg)
				}
			}
		return unmergeableFiles
		}

	beforeMerge(report /*unused*/, mergeableFiles /*unused*/, filename /*unused*/,
		compress? /*unused*/)
		{
		return true
		}

	deleteFile(filename, msg)
		{
		if true isnt result = DeleteFile(filename)
			SuneidoLog('ERRATIC: ' $ msg, params: Object(:result, :filename), calls:)
		}

	alertError(quiet?, msg)
		{
		if not quiet?
			.AlertError("PDF Save To File", msg)
		}

	returnVal(quiet?, msg)
		{
		return quiet? ? msg : false
		}

	mergePdf?(alwaysMerge?, report)
		{
		return alwaysMerge? or report.paramsdata.GetDefault("merge_pdf?", false) is true
		}

	afterSaveFile(report/*unused*/, ext, filename)
		{
		if ext isnt 'pdf'
			return
		ShellExecute(.WindowHwnd(), 'open', filename, fMask: SEE_MASK.ASYNCOK)
		}

	mergePdf(filename, mergeableFiles, compress)
		{
		if mergeableFiles.Has?(filename)
			mergeableFiles.Remove(filename)
		mergeableFiles.Add(filename, at: 0)
		if "" isnt msg = EmailAttachment.AttachmentFilesExist(mergeableFiles)
			{
			.deleteFile(filename, 'Could not clean up PDF temp file')
			return "Append Attachments to PDF Failed, " $ msg
			}

		invalidFiles = #()
		Working('Generating PDF...')
			{
			invalidFiles = PdfMerger(mergeableFiles, filename, :compress)
			}
		if not invalidFiles.Empty?()
			{
			.deleteFile(filename, 'Could not clean up PDF temp file')
			return PdfMerger.InvalidFilesMsg(invalidFiles)
			}
		return true
		}

	orderAttachments(report, quiet? = false)
		{
		files = report.paramsdata.GetDefault('EmailAttachments', #())
		invalidFiles = Object()
		if .hasInvalidEmailAttachment?(files, :invalidFiles)
			{
			return quiet?
				? "Unable to merge files: " $ invalidFiles.Join(', ')
				: false
			}

		mergeableFiles = PdfMerger.FilterFiles(files)
		unmergeableFiles = files.Difference(mergeableFiles)
		if quiet? or mergeableFiles.Size() <= 1
			return Object(:mergeableFiles, :unmergeableFiles)

		if report.paramsdata.GetDefault("merge_pdf_reorder?", false) and
			false is mergeableFiles = ReorderAttachments(.hwnd(report), mergeableFiles)
			return false

		return Object(:mergeableFiles, :unmergeableFiles)
		}

	InvalidMergeFile: 'Params_InvalidMergeFile: '
	hasInvalidEmailAttachment?(emailAttachments, invalidFiles = false)
		{
		mems = emailAttachments.FindAllIf({ it.Prefix?(.InvalidMergeFile) })
		if mems.Empty?()
			return false

		if Object?(invalidFiles)
			invalidFiles.Merge(emailAttachments.Project(@mems).Values())

		return true
		}

	getSaveFileName(report, ext)
		{
		if not .paramsWindowValid?(report)
			return false
		hwnd = .hwnd(report)
		if "" is filename = SaveFileName(
			:hwnd,
			title: "Save " $ ext.Upper() $ " file as",
			filter: ext.Upper() $ " Files (*." $ ext $ ")\x00*." $ ext $
				"\x00All Files (*.*)\x00*.*",
			ext: "." $ ext,
			file: .getDefaultFileName(report, ext))
			return false
		return filename
		}
	getDefaultFileName(report, ext)
		{
		default_name = report.GetDefault('title', 'report') $ ' ' $
			Date().Format("yyyy-MM-dd HHmmss") $ "." $ ext
		return default_name.Tr(CheckFileName.InvalidChars).
			Tr(CheckFileName.InvalidFileChars)
		}
	On_PDF_Email_as_attachment(@report)
		{
		if Object?(.report)
			report = .report
		if not .paramsWindowValid?(report)
			return
		hwnd = .hwnd(report)
		if "" isnt subject = report.GetDefault(#title, "")
			subject = .emailPdfSubject(subject)
		EmailAttachment(hwnd, :subject)
			{
			filename = .pdfName(GetAppTempFullFileName("su"))
			if false is .pdf(report, filename) or
				.hasInvalidEmailAttachment?(.EmailAttachments(report))
				filename = false

			rptResult = Object(:filename,
				attachFileName: .getDefaultFileName(report, 'pdf')
				attachments: .EmailAttachments(report)
				merge_pdf?: report.paramsdata.GetDefault("merge_pdf?", false) is true)
			if report.paramsdata.Member?('EmailSubject')
				rptResult.EmailSubject = report.paramsdata.EmailSubject
			rptResult
			}
		.CloseDialog(report)
		}
	emailPdfSubject(subject)
		{
		prefixPattern = "^(Print|Generate) "
		return PageHeadName() $ " - " $ subject.Replace(prefixPattern, "")
		}

	pdfName(filename)
		{
		extension = filename.AfterLast('.')
		return filename.RemoveSuffix('.' $ extension) $ '.pdf'
		}

	SavePDF(filename) // called by ReportCheck
		{
		.pdf(false, filename)
		}
	pdf(report, filename, quiet? = false)
		{
		if (Object?(.report))
			report = .report

		fromPreview? = report.GetDefault("from_preview", false)
		.checkAndResetParams(report)
		report.paramsdata.ReportDestination = 'pdf'
		.SetExtraParamsData(report.paramsdata)
		.add_print_lines(report.paramsdata)

		report.paramsdata.EmailAttachments = .EmailAttachments(report)

		if not fromPreview? and true is .addFilterIfSlowQuery(report)
			return false

		result = msg = false
		Working('Creating PDF', :quiet?)
			{
			msg = CatchFileAccessErrors(filename)
				{
				result = .PrintPDF(report, filename, quiet?)
				}
			}
		if result is true
			return true

		errMsg = msg isnt true ? msg : result

		if not quiet? and String?(errMsg)
			{
			.AlertInfo("PDF", errMsg)
			return false
			}
		return errMsg
		}

	checkAndResetParams(report)
		{
		if report.Member?('EmailAttachments')
			report.EmailAttachments = Object()
		if (this.Member?("Vert"))
			{
			// don't reset params if from preview
			if (not report.Member?("from_preview"))
				{
				report.paramsdata = .Vert.Data.GetControlData()
				report.paramsdata.lastRan = Timestamp()
				.update_params(report)
				}
			else
				report.Delete("from_preview")
			}
		else if not report.Member?('paramsdata')
			report.paramsdata = Record()
		}

	EmailAttachments(report)
		{
		return report.Member?('EmailAttachments')
			? report.EmailAttachments
			: Object()
		}
	HasIndividualReport?()
		{
		return not _report.Params.GetDefault('EmailAttachments', #()).Empty?()
		}

	PrintPDF(report, filename, quiet? = false)
		{
		if quiet?
			{
			if '' is msg = Report.GetStatusMsg(Report(@report).PrintPDF(filename, true))
				return true
			return msg
			}
		return Report(@report).PrintPDF(filename) is ReportStatus.SUCCESS
		}

	newPRINTDLG(hwnd, report)
		{
		pagenums = report.GetDefault("noPageRange", false) is true
			? PD.NOPAGENUMS : 0
		return Object(
			lStructSize: PRINTDLGEX.Size(),
			hwndOwner: hwnd,
			Flags: PD.USEDEVMODECOPIESANDCOLLATE | PD.NOSELECTION | PD.NOCURRENTPAGE |
				pagenums,
			nMinPage: 1,
			nMaxPage: 9999,
			nPageRanges: 1, nMaxPageRanges: 1,
			lpPageRanges: Object(nToPage: 1, nFromPage: 1),
			nCopies: 1,
			nStartPage: 0xFFFFFFFF)
		}
	uom: 1000 // defined by PSD.INTHOUSANDTHSOFINCHES flag
	On_Page_Setup(@report)
		{
		if (Object?(.report))
			report = .report
		hwndOwner = this.Member?('Window')
			? .Window.Hwnd
			: report.Member?('hwndOwner') ? report.hwndOwner : 0
		psd = .initPageSetupDlg(hwndOwner)
		x = .Get_devmode(report)
		psd.hDevMode = .globalAllocData(x.devmode)
		psd.hDevNames = .globalAllocData(x.devnames)
		psd.rtMargin = Object(
			left: x.left * .uom,
			right: x.right * .uom,
			top: x.top * .uom,
			bottom: x.bottom * .uom)
		if not PageSetupDlg(psd) or psd.hDevNames is 0
			{
			.free(psd)
			if (CommDlgExtendedError() is 0)
				return // user cancelled
			if false is psd = .tryWithoutDevmode(hwndOwner, psd)
				return
			}
		x.devmode = GlobalData(psd.hDevMode)
		x.devnames = GlobalData(psd.hDevNames)
		x.left = psd.rtMargin.left / .uom
		x.right = psd.rtMargin.right / .uom
		x.bottom = psd.rtMargin.bottom / .uom
		x.top = psd.rtMargin.top / .uom
		x.width = psd.ptPaperSize.x / .uom
		x.height = psd.ptPaperSize.y / .uom
		.update_pdc(psd.hDevMode, x, report)
		.free(psd)

		.checkReportSize(report)
		}
	initPageSetupDlg(hwndOwner)
		{
		return Object(
			lStructSize: PAGESETUPDLG.Size(),
			:hwndOwner,
			Flags: PSD.INTHOUSANDTHSOFINCHES + PSD.MARGINS)
		}
	tryWithoutDevmode(hwndOwner, psd)
		{
		// failed so try again without devmode
		psd = .initPageSetupDlg(hwndOwner)
		if not PageSetupDlg(psd) or psd.hDevNames is 0
			{
			.free(psd)
			if (CommDlgExtendedError() isnt 0)
				.AlertError("Page Setup Error",
					"Unable to start Print Dialog, " $
					"please ensure you have a printer set up correctly.")
			return false
			}
		return psd
		}
	On_ReporterReport()
		{
		reporterMode = 'simple'
		rpt_name = .report.name.Replace('Reporter Report - ', 'Reporter - ')
		if .report.Member?(0) and Object?(.report[0])
			{
			if .report[0][0] in ('ReporterFormat', 'ReporterCanvasFormat')
				{
				rpt_name = .report[0][1]
				reporterMode = .report[0][0] is 'ReporterCanvasFormat'
					? 'form'
					: 'simple'
				}
			}

		hwnd = GetActiveWindow()
		Reporter(printMode: true, addButtons: true,
			rpt: rpt_name, title: .report.title,
			onDestroy: { .refresh_reporter_params(rpt_name, hwnd) },
			:reporterMode)
		}
	refresh_reporter_params(rpt_name, hwnd)
		{
		if false is rec = Query1("params", report: rpt_name)
			{
			PubSub.Publish('BrowserRedir_' $ hwnd, 'GoBack')
			return
			}
		if false is .Member?('Vert')
			return
		parent = .Vert
		children = parent.GetChildren()
		i = children.FindIf() {|c| c.Name is 'Data' }
		parent.Remove(i)
		new_report_text = ReporterModel(rec.report).BuildReportText()
		parent.Insert(i, .body(new_report_text))

		// refresh the title, header and printParams
		.report.title = new_report_text.title
		.report.printParams = new_report_text.printParams
		.report.header = new_report_text.header
		if new_report_text.Member?(#footer)
			.report.footer = new_report_text.footer
		.loadParams(.report)
		}

	On_sensitive_file_content_warning()
		{
		Alert('Generated file may contain sensitive information.\nPlease ensure ' $
			'it is getting saved to a secure location.\n' $
			'You may want to consider removing the file once it has been processed',
			'Warning', flags: MB.ICONWARNING)
		}

	free(psd)
		{
		if (psd.hDevNames isnt 0)
			GlobalFree(psd.hDevNames)
		if (psd.hDevMode isnt 0)
			GlobalFree(psd.hDevMode)
		psd.hDevNames = psd.hDevMode = 0 // prevent double free
		}
	Get_devmode(report)
		{
		x = false
		if (report.Member?("name"))
			try x = Query1(.devmode_query(.devmode_reportname(report)))

		if x isnt false
			return x

		return report.GetDefault('default_orientation', 'Portrait') is 'Landscape'
			? Record(width: 11, height: 8.5, left: .5, right: .5, top: .5, bottom: .5)
			: Record(width: 8.5, height: 11, left: .5, right: .5, top: .5, bottom: .5)
		}
	devmode_query(reportname)
		{
		.Ensure_devmode()
		return "devmode
			where computer = " $ Display(.devmode_save_name()) $
			" and report = " $ Display(reportname)
		}
	Ensure_devmode()
		{
		Database("ensure devmode
			(computer, report, bottom, devmode, devnames, height, left, right,
				top, width, devmode_TS)
			key (computer,report)")
		}
	update_params(report)
		{
		if not report.Member?("name") or (report.Member?('NoSaveLoadParams') and
			report.NoSaveLoadParams is true)
			return
		x = Query1(.params_query(report.name))
		report_options = x is false ? #() : x.report_options
		lastRan = x is false ? '' : x.params.GetDefault('lastRan', '')
		if false is report.paramsdata.Member?('lastRan')
			report.paramsdata.lastRan = lastRan

		params = report.paramsdata.Copy()
		.RemoveIgnoreFields(params)
		RetryTransaction()
			{|t|
			t.QueryDo("delete " $ .params_query(report.name))
			t.QueryOutput("params", Object(
				user: Suneido.User,
				report: report.name,
				:params,
				:report_options))
			}
		}
	RemoveIgnoreFields(rec)
		{
		for f in rec.Copy().Members()
			if rec[f] isnt '' and true is
				Datadict(f, getMembers: #(ParamsNoSave)).GetDefault('ParamsNoSave', false)
				rec.Delete(f)
		}
	loadParams(report)
		{
		if not report.Member?("name")
			return

		// get saved params and merge with report params
		if ((report.Member?('NoSaveLoadParams') and
			report.NoSaveLoadParams is true) or
			(false is x = Query1(.params_query(report.name))))
			x = Object(params: Object())

		.setFilterDefaults(report, x)

		.setLayoutDefault(x)

		.setParams(x)

		// have to do RecordControl.Set so that observer gets set up
		.Vert.Data.Set(Record())
		for (field in x.params.Members())
			.Vert.Data.SetField(field, x.params[field])

		if x.params.GetDefault('PrintLines', false) is true and
			false isnt ctrl = .FindControl('print_lines')
			ctrl.Set(true)

		.checkReportSize(report)
		}

	checkReportSize(report)
		{
		if not .reporter?(report)
			return

		rpt_name = report.name.Replace('Reporter Report - ', 'Reporter - ')
		.setMsg(Opt('Alignment Warning: ',
			CheckReportPageSize(rpt_name, .Get_devmode(report))))
		}

	setMsg(msg)
		{
		if false is paramsContainer = .FindControl('paramContainer')
			return

		msgCtrl = paramsContainer.FindControl('paramMsg')
		if msg isnt ''
			{
			if msgCtrl is false
				msgCtrl = paramsContainer.Append(
					Object('Static', textStyle: 'warn', name: 'paramMsg'))
			msgCtrl.Set(msg)
			}
		else
			paramsContainer.Remove(
				paramsContainer.GetChildren().FindIf({|c| c.Name is 'paramMsg' }))
		}

	Record_NewValue(field, value)
		{
		if .report.Member?('AfterField')
			(.report.AfterField)(field, value, data: .Vert.Data.Get())

		if field is 'params_report_layout'
			{
			if false isnt rec = Query1('report_layout_designs', rptdesign_name: value,
				report: .report.name)
				if Object?(rec.rptdesign_layout)
					for m, v in rec.rptdesign_layout
						if .Vert.Data.HasControl?(m)
							.Vert.Data.SetField(m, v)
			}
		}

	setFilterDefaults(report, x)
		{
		if report.Member?('SetFilterDefaults')
			.BuildFilters(report.SetFilterDefaults, x.params)
		}

	BuildFilters(filterDefaults, params)
		{
		for filterName in filterDefaults.Members()
			for field in filterDefaults[filterName]
				{
				if not params.Member?(filterName)
					params[filterName] = Object()
				if false is .repeatConditionExists(filterName, field, Object(:params))
					{
					pos = filterDefaults[filterName].Find(field)
					ob = Record(condition_field: field)
					ob[field] = #(operation: "", value: '')
					params[filterName].Add(ob, at: pos)
					}
				}
		}

	setLayoutDefault(x)
		{
		if false is layoutCtrl = .FindControl('ParamsChooseLayout')
			return
		rptMem = 'params_report_layout'
		rptName = layoutCtrl.ReportName
		if ReportLayoutDesign.Customizable?(rptName)
			x.params[rptMem] = ReportLayoutDesign.DefaultValue(
				rptName, x.params.GetDefault(rptMem, ''))
		}

	setParams(x)
		{
		for (field in .setParams.Members())
			if (.setParams[field] isnt #())
				{
				if String?(field) and field.Suffix?('Filters')
					.addRepeatCondition(field, .setParams[field], x)
				else
					x.params[field] = .setParams[field]
				}
		if .report.Member?('AfterSet')
			(.report.AfterSet)(x.params)
		}
	addRepeatCondition(filterName, repeatOb, x)
		{
		if not x.params.Member?(filterName) or not Object?(x.params[filterName])
			x.params[filterName] = Object()

		for ob in repeatOb
			if false is pos = .repeatConditionExists(filterName, ob.condition_field, x)
				x.params[filterName].Add(ob)
			else
				x.params[filterName][pos] = ob
		}

	repeatConditionExists(filterName, condition_field, x)
		{
		return x.params[filterName].FindIf({ it.condition_field is condition_field })
		}

	params_query(reportname)
		{
		return "params
			where user is " $ Display(Suneido.User) $
			" and report is " $ Display(reportname)
		}
	update_pdc(hdm, x, report, fromPreview? = false)
		{
		if report.Member?('name')
			{
			reportname = .devmode_reportname(report)
			devmode_save_name = .devmode_save_name()
			RetryTransaction()
				{|t|
				t.QueryDo('delete ' $ .devmode_query(reportname))
				x.report = reportname
				x.computer = devmode_save_name
				t.QueryOutput('devmode', x)
				}
			}

		if (hdm is 0)
			{
			.invalidPrinterError()
			return false
			}

		try
			{
			if (Suneido.Member?("pdc"))
				DeleteDC(Suneido.pdc)

			if 0 is Suneido.pdc = .createDC(x.devnames, hdm)
				{
				if not fromPreview?
					.invalidPrinterError()
				return false
				}
			}
		catch(err /*unused*/)
			{
			if not fromPreview?
				.invalidPrinterError()
			return false
			}
		return true
		}
	invalidPrinterError()
		{
		// Windows Vista & 7 Do not have a Select Printer Option
		// from Page Setup. The Validation on the Page Setup screen
		// is resetting the default printer however, (Instead of an error message
		// could we not just replicate that behavior?)
		.AlertWarn("Invalid Printer", "Please open Page Setup and click OK")
		}
	devmode_save_name()
		{
		return Sys.SuneidoJs?() is false and WTS_GetSessionId() is 0
			? GetComputerName()
			: Suneido.User
		}
	devmode_reportname(report)
		{
		return report.Member?('devmode_name') ? report.devmode_name : report.name
		}
	SetNullPdc()
		{
		if Suneido.Member?("pdc")
			DeleteDC(Suneido.pdc)
		Suneido.pdc = NULL
		}
	Destroy()
		{
		.ClearFocus()
		if (this.Member?("Vert") and this.Member?("Params_report") and
			.params_valid?(silent:))
			{
			.report.paramsdata = .Vert.Data.GetControlData()
			.add_print_lines(.report.paramsdata)
			.update_params(.report)
			}
		.report.GetDefault('onDestroy', function(){})()
		super.Destroy()
		}
	}
