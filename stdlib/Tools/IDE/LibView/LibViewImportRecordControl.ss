// Copyright (C) 2019 Axon Development Corporation All rights reserved worldwide.
Controller
	{
	Title: "Import Records"

	CallClass(hwnd, defaultLib)
		{
		OkCancel(Object(this, defaultLib), .Title, hwnd)
		}
	New(.defaultLib = false)
		{
		if defaultLib isnt false
			.Data.SetField('lib', defaultLib)
		}

	Controls()
		{
		return Object('Record',
			Object('Vert',
			Object('Pair',
				Object('Static', 'Filename')
				Object('OpenFile' title: 'Import' name: "fileName")), 'Skip',
			Object('Pair',
				Object('Static', 'Library')
				Object('ChooseList', LibTreeModel.Libs(), name: "lib"))
			))
		}

	OK()
		{
		data = .Data.Get()
		lib = data.lib.Tr('()')
		if not lib.Blank?() and not TableExists?(lib)
			{
			.AlertError('Import Record', 'Please specify a valid Library.')
			return false
			}
		return Object(fileName: data.fileName, lib: lib.Blank?() ? false : lib)
		}
	}