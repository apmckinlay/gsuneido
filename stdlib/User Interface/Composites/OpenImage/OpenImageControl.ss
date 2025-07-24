// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
OpenImageAddonsBase
	{
	// file may contain only file name
	// or a complete path with file name
	// or a complete path without file name
	Name: 'OpenImage'
	New(filter = "", file = "", status = "", reorderOnly = false)
		{
		super(filter, file, status, reorderOnly, GetContributions('OpenImageAddons'))
		.image = .FindControl('image')
		}

	Controls()
		{
		return Object(#Image,
			message: .GetReorderOnly()
				? ''
				: 'Double click to select or drag files to here',
			acceptDrop:, xstretch: false, ystretch: false, name: 'image')
		}

	value: ''
	Get()
		{
		return .value
		}

	Set(value)
		{
		if not String?(value)
			value = ''
		.File = value
		.value = Paths.Basename(.File) is '' ? '' : .File
		.image.Set(.FullPath(), .Highlight?(.value), message: .value)
		.SetTip()
		}

	FullPath(file = false)
		{
		if file is false
			file = .File
		copyTo = OpenImageSettings.Copyto()
		return (Paths.Basename(file) is file and file isnt "" and
			copyTo isnt "")
			? copyTo $ file
			: file
		}

	SetTip()
		{
		.image.ToolTip(.value)
		}

	GetImageControl()
		{
		return .image
		}

	ProcessValue(value)
		{
		return value
		}
	}
