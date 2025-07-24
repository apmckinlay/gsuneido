// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
OpenImageAddon
	{
	ContextMenu()
		{
		if Sys.SuneidoJs?()
			return #()

		disabled? = .GetImageReadOnly() is true or .GetReorderOnly() is true or
			 not Scanning().ScanningAllowed?() or not Scanning().ScannerAvailable?()

		return Object(
			Object(name: "Scan Attachment", state: .State(disabled?), order: 30)
			Object(
				Object(name: "Scan Attachment", state: .State(disabled?))
				Object(name: "Select Source...", state: .State(disabled?))
				order: 31
				)
			)
		}

	On_Scan_Attachment()
		{
		if not OpenImageSettings.Normally_linkcopy?()
			{
			.On_Scan_Attachment_As()
			return
			}
		filename = GetAppTempPath() $ GetScanFilename() $ ".pdf"
		.scanAttachment(filename)
		}

	On_Scan_Attachment_As()
		{
		if "" is filename = SaveFileName(
			hwnd: .Window.Hwnd,
			title: "Save Scanned file as",
			filter: "PDF Files (*.pdf)\x00*.pdf"
			ext: ".pdf",
			file: GetScanFilename() $ ".pdf")
			return
		.scanAttachment(filename)
		}

	scanAttachment(filename)
		{
		if .GetImageReadOnly()
			{
			.AlertInfo("Scan Attachment", 'Need to be in edit mode to Scan Attachment')
			return
			}
		Working('Scanning Attachment...')
			{
			if true is resultMsg = Scanning().Scan(filename)
				.Attach(filename) // Will copy the file and attach
			if String?(resultMsg)
				.AlertInfo("Scan Attachment", resultMsg)
			}
		}

	On_Select_Source()
		{
		Scanning().SelectScanner(.Window.Hwnd)
		}
	}
