// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
_Params
	{
	verifyDevMode(report/*unused*/)
		{
		return true
		}

	preparePrinter(params_window/*unused*/, report/*unused*/)
		{
		return true
		}

	On_Page_Setup(@report)
		{
		if (Object?(.report))
			report = .report

		x = .Get_devmode(report)
		if false is newLayout = SuPageSetupControl(x)
			return

		x.Merge(newLayout)
		.update_pdc(x, report)

		.checkReportSize(report)
		}

	On_Print(@report)
		{
		if IsSafari?()
			.AlertWarn('Print compatibility issue',
				"The Safari browser that you are using does not support " $
				"certain printing features. If you have problems with alignment try " $
				"using a different browser such as Chrome or Firefox or save as PDF " $
				"and then print.")
		super.On_Print(@report)
		}

	deleteDC(unused)
		{
		}

	update_pdc(x, report)
		{
		if report.Member?('name')
			{
			reportname = .devmode_reportname(report)
			RetryTransaction()
				{|t|
				t.QueryDo('delete ' $ .devmode_query(reportname))
				x.report = reportname
				x.computer = .devmode_save_name()
				t.QueryOutput('devmode', x)
				}
			}
		}

	getSaveFileName(report, ext)
		{
		if not .paramsWindowValid?(report)
			return false

		return SuGetTempSaveName(ext)
		}

	afterSaveFile(report, ext, filename)
		{
		JsDownload.Trigger(Paths.Basename(filename), .getDefaultFileName(report, ext))
		}

	// overridden to avoid showing the preview window
	buildPreviewWindowForPageCount(report, pdc /*unused*/)
		{
		w = Window(Object(PreviewControl, report, this, extraButtons: .PreviewButtons),
			title: '', show: false)
		w.SetVisible(false)
		return w
		}

	On_PDF_Download(@report)
		{
		.On_PDF_Save_to_file(@report)
		}

	beforeMerge(report, mergeableFiles, filename, compress?)
		{
		if '' is AttachmentS3Bucket() or not Sys.SuneidoJs?()
			return true
		attachments = mergeableFiles.Map({ '(APPEND) ' $ it })
		merge_pdf? = attachments.Size() > 1
		data = Object(:attachments, :merge_pdf?, :compress?, mergeOnly?:)
		defaultFileName = .getDefaultFileName(report, '.pdf')
		EmailAttachment_Mime.SendJS(false, filename, defaultFileName, :merge_pdf?,
			:data, :mergeableFiles)
		.CloseDialog(report)
		return false
		}
	}
