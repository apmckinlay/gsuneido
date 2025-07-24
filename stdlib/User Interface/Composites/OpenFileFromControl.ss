// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
OpenFileControl
	{
	Name: 'OpenFileFrom'
	ButtonName: 'Choose...'
	New(.title = 'Open', .filter = "", width = 30, file = "", .flags = false, status = "",
		mandatory = false, disableDrop = false, .rootDir = false, .fileType = false)
		{
		super(:title, :filter, :width, :file, :flags, :status, :mandatory, :disableDrop)
		}

	SetReadOnly(readOnly)
		{
		if false isnt ctrl = .FindControl('Field')
			ctrl.SetReadOnly(readOnly)
		}

	On_Choose()
		{
		folders = ViewFileControl.GetDirList(.rootDir, filter: .filter)
		hwnd = .Window.Hwnd
		filename = folders.Empty?()
			? OpenFileName(:hwnd, title: .title, flags: .flags, initialDir: .InitialDir)
			: ViewFileControl(hwnd, Object(this, rootDir: .rootDir,
				fileType: .fileType, file: .Get(), filter: .filter, :folders))

		if filename in ('', false) or .Destroyed?()
			return

		.SetFileName(filename)
		}
	}