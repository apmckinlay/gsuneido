// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
OpenImageAddon
	{
	ContextMenu(readonly? = "")
		{
		if readonly? is ""
			readonly? = .GetImageReadOnly() is true

		return Object(
			Object(name: "Labels", state: .State(readonly?), order: 100))
		}

	ExtraValidation(labels, file)
		{
		handlers = Contributions('OpenImageLabelsValidation')
		errOb = Object()
		for handler in handlers
			for label in labels.Copy()
				{
				if '' is validOb = handler(label, file)
					continue
				labels.Remove(label)
				errOb.Add(validOb.error)
				}
		return errOb
		}

	On_Labels()
		{
		fieldName = .Parent.Name
		fileOb = .Parent.SplitLabel(.Parent.Get())
		file = .Parent.FullPath(fileOb.file, fileOb.subfolder)
		if false is labels = (.labeldialog)(.Window.Hwnd, .GetLabels(), fieldName,
			.GetType(), file)
			return
		.SetLabels(labels)
		.Send('NewValue', .Get())
		.SetTip()
		}

	labeldialog: Controller
		{
		Title: 'Attachment Labels'
		file : false
		CallClass(hwnd, labels, fieldName, type, file = false)
			{
			return OkCancel(Object(this, labels, fieldName, type, file), .Title, hwnd)
			}
		New(origlabel, fieldName = '', type = '', .file = false)
			{
			super(.controls(fieldName, type))
			.attachmentctrl = .Vert.AttachmentLabels
			.attachmentctrl.Set(origlabel)
			}
		controls(fieldName, type)
			{
			Object('Vert'
				#(Static 'Enter one or more words or phrases describing the attachment,' $
					' separated by commas. e.g. web, pod')
				#(Skip 6)
				#(AttachmentLabels xstretch: 1 ystretch: 1)
				#(Skip 6)
				Object('Horz'
					Object('OpenImageLabelsExtra', :fieldName, :type)))
			}
		OK()
			{
			return .attachmentctrl.Get()
			}

		NewValue(s)
			{
			labels = s.Split(',')
			if not .checkLabels?(labels, .file)
				.attachmentctrl.Set(labels.Join(','))
			}

		checkLabels?(labels, file)
			{
			if false is file
				return false
			errOb = Addon_Labels.ExtraValidation(labels, file)
			if not errOb.Empty?()
				{
				.AlertError(.Title, errOb.Join('\r\n'))
				return false
				}
			return true
			}

		AddLabel(label)
			{
			if not .checkLabels?(Object(label), .file)
				return

			attachlabelsText = .attachmentctrl.Get().Trim()
			attachOb = attachlabelsText.Split(',').Map!(#Trim).Remove('')
			if attachOb.Has?(label)
				return

			.attachmentctrl.Set(attachlabelsText $
				(attachlabelsText.Size() > 0 and not attachlabelsText.Suffix?(',')
				? ','
				: '') $ label)
			}
		}
	}