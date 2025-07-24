// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
OpenImageAddon
	{
	saveAsPrompt: 'Save As...'
	ContextMenu()
		{
		disabled? = .FileEmpty?() or .GetReorderOnly() is true

		menu = Object(
			Object(name: "Open with Windows", state: .State(disabled?), order: 0)
			Object(name: "Email Attachment", state: .State(disabled?), order: 2)
			Object(name: "Rename Attachment", state: .State(disabled? or
				.GetImageReadOnly() is true ), order: 79)
			Object(name: .saveAsPrompt, state: .State(disabled?), order: 80)
			Object(name: "Properties", state: .State(disabled?), order: 81)
			)
		if not Sys.SuneidoJs?()
			menu.Add(
				Object(name: "Print Attachment", state: .State(disabled?), order: 1))
		return menu
		}

	On_Rename_Attachment(file = false, hwnd = 0)
		{
		if hwnd is 0
			hwnd = .Window.Hwnd
		// .FullPath handles if file is false
		fullPath = .FullPath(file)
		if false isnt file = .ProcessFile(fullPath)
			if OpenImageRename(fullPath, hwnd, this, .GetCopyTo()) is true
				.QueueDeleteAttachment(.FullPath(), fullPath)
		}

	On_Open_with_Windows(file = false, hwnd = false)
		{
		.Open(file, hwnd)
		}

	imageFilePattern: "(?i)[.](bmp|gif|jpg|jpe|jpeg|ico|emf|wmf)$"
	On_Print_Attachment(file = false, hwnd = 0)
		{
		if hwnd is 0
			hwnd = .Window.Hwnd
		file = .FullPath(file)
		if false is file = .ProcessFile(file)
			return
		if file =~ .imageFilePattern
			ToolDialog(hwnd,
				Object('Params',
					Object('Image', file, xstretch: 1, ystretch: 1),
					title: 'Print Image',
					name: 'PrintImage',
					header: false))
		else
			ShellExecute(hwnd, 'print', Paths.ToWindows(file), fMask: SEE_MASK.ASYNCOK)
		}

	On_Email_Attachment()
		{
		file = .FullPath()
		if false isnt existFile = .ProcessFile(file)
			EmailAttachment(.Window.Hwnd, Object(filename: existFile,
				attachFileName: Paths.Basename(existFile), attachments: #()))
		}

	On_Save_As()
		{
		curFilename = .FullPath()
		if false is .ProcessFile(curFilename)
			return
		ext = curFilename.AfterLast('.')
		DoWithSaveFileName(
			hwnd: .Window.Hwnd,
			title: "Save " $ ext.Upper() $ " file as",
			filter: ext.Upper() $ " Files (*." $ ext $ ")\x00*." $ ext $
				"\x00All Files (*.*)\x00*.*",
			ext: "." $ ext,
			file: Paths.Basename(curFilename))
			{ |filename|
			if filename isnt curFilename and
				true isnt CopyFile(curFilename, filename, false)
				throw Paths.Basename(curFilename)
			}
		}

	On_Properties()
		{
		if false isnt .ProcessFile(fullPath = .FullPath())
			FilePropertiesControl(fullPath, .Window.Hwnd)
		}
	}
