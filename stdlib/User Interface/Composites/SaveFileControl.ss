// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
// contributions by Claudio Mascioni
// TODO: add extension if user types in basename (without browse)
OpenFileControl
	{
	Name: 'SaveFile'
	// file may contain only file name
	// or a complete path with file name
	// or a complete path without file name
	New(.title = 'Save', width = 30, .filter = "", .ext = "", file = "",
		.flags = false, status = '', .mandatory = false)
		{
		super(:width, :status, :mandatory, disableDrop:, :filter)
		.Set(file)
		}
	SetFilter(filter)
		{ .filter = filter	}
	SetDefExt(ext)
		{ .ext = ext }
	Valid?()
		{
		if .GetReadOnly()
			return true

		fileName = .Get()
		if not fileName.Blank?() and not CheckDirectory.ValidFileName?(fileName)
			return false

		return super.Valid?()
		}
	On_Browse()
		{
		filename = SaveFileName(hwnd: .Window.Hwnd, title: .title,
			flags: .flags, filter: .filter, ext: .ext,
			file: .GetFileName(), initialDir: .InitialDir)
		if filename is ""
			filename = .Get()
		.Set(Paths.Basename(filename) is "" ? "" : filename)
		.Send("NewValue", .Get())
		}
	}
