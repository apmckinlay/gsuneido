// Copyright (C) 2004 Suneido Software Corp. All rights reserved worldwide.
// contributions by Claudio Mascioni
OpenFileControl
	{
	Name: 'BrowseFolder'
	New(title = 'Browse for Folder', pidlroot = 0, dirpath = false)
		{
		super(disableDrop:)
		.SetOpenDirTitle(title)
		switch dirpath
			{
		case false:
			.opendirpath = ""  // open in computer resources
		case true:
			.opendirpath = GetCurrentDirectory() // open in current directory
		default:
			.opendirpath = dirpath // open in passed path
			}
		.SetPidlRoot(pidlroot)  // this is a pidl format
		.Set(.opendirpath)
		.Send('Data')
		}
	SetOpenDirTitle(text)
		{
		.opendirmsg = text
		}
	SetPidlRoot(value)
		{
		.pidlroot = value
		}
	Set(value)
		{ .Horz.Field.Set(.opendirpath = value)	}
	SetFilter(filter/*unused*/)
		{ } // BrowseFolder does not have filters
	On_Browse()
		{
		selecteddir = BrowseFolderName(
			hwnd: .Window.Hwnd,
			title: .opendirmsg,
			pidlRoot: .pidlroot,
			initialPath: .opendirpath)
		if selecteddir is ""
			return
		.Set(selecteddir)
		.Send("NewValue", .Get())
		}
	Destroy()
		{
		.Send('NoData')
		super.Destroy()
		}
	}