// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
OpenImageAddon
	{
	ContextMenu()
		{
		disabled = .Send('ImageDisableRemoveAttachment') is true or .FileEmpty?() or
			.GetImageReadOnly() is true
		return Object(
			Object(name: "Remove Attachment", state: .State(disabled), order: 21))
		}

	On_Remove_Attachment()
		{
		if .Parent.Method?('SetLabels')
			.SetLabels('')

		origFile = .FullPath()
		.SetNewValue(.ProcessValue(""))
		.QueueDeleteAttachment('', origFile, AttachmentsManager.RemoveAction)
		}
	}
