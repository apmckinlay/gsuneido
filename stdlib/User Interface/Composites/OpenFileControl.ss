// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
// contributions by Claudio Mascioni
// TODO: add multiple file option
Controller
	{
	Name: 'OpenFile'
	InitialDir: ""
	// file may contain only file name
	// or a complete path with file name
	// or a complete path without file name
	New(.title = 'Open', filter = "", width = 30, file = "", .flags = false,
		status = "", mandatory = false, disableDrop = false, .attachment? = false)
		{
		super(.controls(width, status, mandatory, disableDrop))
		.field = .Horz.Field
		.Top = .Horz.Top
		.SetFilter(filter)
		.Set(file)
		.Send('Data')
		}

	ButtonName: 'Browse...'
	controls(width, status, mandatory, disableDrop)
		{
		field = Object('Field', :width, :status, :mandatory, readonly: Sys.SuneidoJs?())
		if not disableDrop
			{
			field.cue = 'Browse or drag a file to here'
			field.acceptDrop = true
			}
		return Object('Horz'
			field
			#(Skip 2)
			Object('Button', .ButtonName))
		}

	GetFileName()
		{
		return Paths.Basename(.Get())
		}

	GetFilePath()
		{
		return .Get().BeforeLast(.GetFileName())
		}

	SetFilter(.filter)
		{
		}

	SetStatus(status)
		{
		.field.SetStatus(status)
		}

	Set(value)
		{
		if value isnt ""
			.InitialDir = value.BeforeLast(Paths.Basename(value))
		// if is gived only the path, displayed value is ""
		value = Paths.Basename(value) is "" ? "" : value
		.field.Set(value)
		}

	FieldDropFiles(wParam)
		{
		if false isnt file = DragQueryFile(wParam, 0)
			.SetFileName(file)
		}

	Get()
		{
		return .field.Get()
		}

	Valid?()
		{
		return .field.Valid?()
		}

	Dirty?(dirty = "")
		{
		return .field.Dirty?(dirty)
		}

	NewValue(value)
		{
		.Send("NewValue", value)
		}

	On_Browse()
		{
		filename = OpenFileName(hwnd: .Window.Hwnd, title: .title,
			filter: .filter, flags: .flags, file: .GetFileName(),
			initialDir: .InitialDir, attachment?: .attachment?)
		 if filename is "" or .Destroyed?()
			return
		.SetFileName(filename)
		}

	SetFileName(filename)
		{
		.Set(Paths.Basename(filename) is "" ? "" : filename)
		.Send("NewValue", .Get())
		.field.SetValid(true) // in case it was invalid from mandatory
		}

	// have to redefine SetReadOnly because Button's SetReadOnly is only a
	// stub, have to use SetEnabled method
	SetReadOnly(readOnly)
		{
		.field.SetReadOnly(readOnly)
		.Horz.Browse.SetEnabled(not readOnly)
		}

	SetFont(font,size)
		{
		.field.SetFont(font, size)
		}

	Destroy()
		{
		.Send('NoData')
		super.Destroy()
		}
	}