// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Title: 'Properties'
	CallClass(fullPath, hwnd = 0)
		{
		ToolDialog(hwnd, Object(this, fullPath), title: .Title)
		}

	New(.fullPath)
		{
		}

	Controls()
		{
		details = .Dir(.fullPath, files:, details:)
		controls = details.NotEmpty?()
			? .propertiesLayout(details[0])
			: Object('Form',
				Object('Static', 'Could not find file: ' $ Paths.Basename(.fullPath)))

		if InternalUser?() or not Sys.SuneidoJs?()
			controls.Add('nl',
				Object('LinkButton', 'Full Path', tip:
					'Click to copy: ' $ .fullPath, group: 0))

		return controls
		}

	Dir(fullPath)
		{
		return FileStorage.Dir(fullPath, files:, details:)
		}

	propertiesLayout(details)
		{
		return Object('Form',
			#(Static, 'Name: ', group: 0),
				Object('Static', details.name, group: 1), 'nl',
			#(Static, 'Size: ', group: 0),
				Object('Static', ReadableSize(details.size), group: 1), 'nl',
			#(Static, 'Last Modified: ', group: 0),
				Object('Static', details.date.ShortDateTime(), group: 1))
		}

	On_Full_Path()
		{
		ClipboardWriteString(.fullPath)
		InfoWindowControl('Copied: ' $ .fullPath,
			titleSize: 0, marginSize: 7, autoClose: 5)
		}
	}