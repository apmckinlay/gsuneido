// Copyright (C) 2006 Suneido Software Corp. All rights reserved worldwide.
// contributed by Claudio Mascioni
OpenFileControl
	{
	Name: 'BrowseImage'
	// file may contain only file name, or a complete path with file name
	// or a complete path without file name
	New(title = '', filter = "", width = 30, file = false, showimage = true,
		status = "", opendirmsg = "")
		{
		super(width: width, status: status)
		.title = TranslateLanguage(title)
		// filter is in this format: "bmp, gif, jpg, jpe, jpeg, ico, emf, wmf, tif,
		// tiff, png, pdf". i.e. to can select only jpg files, filter is "jpg, jpeg"
		.filter = filter
		.Set(file)
		.showimage = showimage
		.opendirmsg = opendirmsg is ""
			? 'Select an images folder'
			: opendirmsg
		.Send('Data')
		}
	Set(value)
		{
		if value is false
			return
		.file = value
		 // if is gived only the path, displayed value is ""
		.Horz.Field.Set((Paths.Basename(value) is "") ? "" : value)
		}
	On_Browse()
		{
		filename = BrowseImageName(.Window.Hwnd, .title, .filter, .file,
			.showimage, .opendirmsg)
		if filename is false
			filename = .Get()
		.Set((Paths.Basename(filename) is "") ? "" : filename)
		.Send("NewValue", .Get())
		}
	Destroy()
		{
		.Send('NoData')
		super.Destroy()
		}
	}
