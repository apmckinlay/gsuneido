// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	// tasks: [filename: [t: Date(), :saveName]]
	CallClass(tasks)
		{
		if not Sys.SuneidoJs?()
			return

		Suneido.JsDownloadAlertWindow = ModalWindow(Object(this, tasks),
			title: 'Outstanding Downloads', keep_size: false, closeButton?: false)
		FlashWindowEx(Object(hwnd: Suneido.JsDownloadAlertWindow.Hwnd,
			dwFlags: FLASHW.ALL, message: 'You have outstanding downloads'))
		}

	New(.tasks)
		{
		.initTaskCmdMap()
		.container = .FindControl(#container)
		}

	Controls()
		{
		return Object('Vert',
			Object('AlertText', 'Your browser appears to be blocking ' $
				'the download of multiple files. Please enable the setting ' $
				'(automatic downloads of multiple files), ' $
				'then click ' $ (.tasks.Size() > 1 ? 'each link' : 'the link') $
				' to retry the download:')
			#Skip,
			Object('Vert', .buildTaskLinks(.tasks), name: 'container')
			#(Horz, Fill, (Button 'Cancel Download')))
		}

	buildTaskLinks(tasks)
		{
		vert = Object('Vert')
		for filename, ob in tasks
			vert.Add(Object('LinkButton',
				ob.saveName.Ellipsis(50/*=len*/) $
					' - created at ' $ ob.t.ShortDateTime(),
				command: filename))
		return vert
		}

	Update(newTasks)
		{
		if .Destroyed?()
			return

		if newTasks.Empty?()
			{
			.Window.CLOSE()
			return
			}

		if .tasks is newTasks
			return

		.tasks = newTasks
		.initTaskCmdMap()
		.container.Remove(0)
		.container.Append(.buildTaskLinks(.tasks))
		}

	Refresh()
		{
		.Update(ServerEval('JsDownload.CheckTask', Suneido.User))
		}

	tasksMap: #()
	initTaskCmdMap()
		{
		.tasksMap = Object()
		for filename in .tasks.Members()
			.tasksMap[ToIdentifier(filename.Trim())] = filename
		}

	On_Cancel_Download()
		{
		if false is YesNo('Are you sure you want to cancel all downloads?',
			title: 'Cancel Download', flags: MB.ICONWARNING)
			return
		.tasks.Members().Each(JsDownload.DeleteTask)
		.Refresh()
		}

	Recv(@args)
		{
		if false is filename = .tasksMap.GetDefault(args[0].RemovePrefix('On_'), false)
			return 0
		SuRenderBackend().RecordAction(false, 'SuDownloadFile', [
			target: Base64.Encode(filename.Xor(EncryptControlKey())),
			saveName: .tasks[filename].saveName])
		.Delay(1000/*=1 sec*/, uniqueID: 'refresh')
			{
			.Refresh()
			}
		return 0
		}

	Destroy()
		{
		Suneido.Delete(#JsDownloadAlertWindow)
		super.Destroy()
		}
	}