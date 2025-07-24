// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Title: 'Adjust the order of attachments that are being appended'
	CallClass(hwnd, attachments)
		{
		tmpPath = GetAppTempPath()
		if false is x = ToolDialog(
			hwnd, [this, attachments, tmpPath], closeButton?: false)
			return false

		orderedAttach = Object()
		for level in x
			for field in level.Members().Sort!()
				{
				file = level[field]
				if .specialFile?(file)
					file = Paths.Combine(tmpPath, file)
				orderedAttach.Add(file)
				}
		return orderedAttach
		}

	New(attachments, .tmpPath)
		{
		super(.layout())
		.attachControl = .FindControl("Data").GetControl("order_attachments")
		attach = Object([])
		col = row = 0
		for(i = 0; i < attachments.Size(); i++)
			{
			if col is .attachControl.PerRow
				{
				col = 0
				row++
				attach.Add([])
				}
			attachName = attachments[i]
			if attachName.Prefix?(.tmpPath)
				attachName = Paths.Basename(attachName)
			attach[row].Add(attachName, at: 'attachment' $ col)
			col++
			}
		.attachControl.Set(attach)
		}

	layout()
		{
		return Object('Record', Object('Vert',
			#(Static "You can drag and drop to adjust the order of " $
				"attachments that are being appended" textStyle: 'note')
			#Skip
			Object('Scroll',
				Object(AttachmentsRepeatControl
					name: 'order_attachments', type: '', reorderOnly:)),
			#(Skip 5),
			#(Horz Fill (Button 'OK', width: 8))
			), ymin: 300)
		}

	ImageDisableRemoveAttachment(source /*unused*/)
		{
		return true
		}

	ImageGetFullPath(file)
		{
		return .specialFile?(file) ? Paths.Combine(.tmpPath, file) : file
		}

	// file generated temporarily
	specialFile?(file)
		{
		return file is Paths.Basename(file)
		}

	On_OK()
		{
		.Window.Result(.attachControl.Get())
		}
	}