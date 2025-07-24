// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
OpenImageAddon
	{
	ContextMenu()
		{
		readonly? = .GetImageReadOnly() is true
		reorderOnly? = .GetReorderOnly() is true

		return Object(
			Object(name: "Select Attachment", state: .State(readonly? or reorderOnly?)
				order: 20))
		}

	// sent by OpenImageAddonsBase
	On_Select()
		{
		.On_Select_Attachment()
		}

	On_Select_Attachment()
		{
		origFile = .FullPath()
		selectedfile = OpenImageGetSelectedFile(.Window.Hwnd, .Parent.File,
			.GetFilter(), .GetStatus(), .Parent.UseSubFolder)

		try
			if selectedfile isnt false
				{
				.SetNewValue(.ProcessValue(selectedfile))
				action = origFile is '' ? '' : AttachmentsManager.ReplaceAction
				.QueueDeleteAttachment(.FullPath(), origFile, action)
				}
		catch(e, 'member not found')
			SuneidoLog('Attachment control destroyed (' $ e $')')
		}
	}
