// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Title: 'Select Source Scanner'
	CallClass(hwnd, sources)
		{
		if sources is false
			return .AlertInfo(.Title, 'No Scanners Detected')
		OkCancel(Object(this, sources), .Title, hwnd)
		}
	New(.sources) {	}

	Controls()
		{
		return Object('Record'
			Object('Vert'
				Object('Static', 'Please choose TWAIN (TW) Scanner when ever possible')
				#Skip
				Object('Horz',
					Object('ChooseListControl', width: 20,  list: .sources
						mandatory:, name: 'sourcesWindow', selectFirst:)
				)))
		}

	OK()
		{
		data = .Data.Get()
		ctrl = .FindControl('sourcesWindow')
		if false is .sources.Has?(data.sourcesWindow)
			{
			ctrl.SetValid(false)
			return false
			}
		else
			{
			UserSettings.Put('Clscan - Scanner Source', value: data.sourcesWindow)
			return true
			}
		}
	}
