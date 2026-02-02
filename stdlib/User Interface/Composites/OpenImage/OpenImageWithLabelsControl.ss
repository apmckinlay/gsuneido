// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
OpenImageAddonsBase
	{
	Name: 'OpenImageWithLables'
	UseSubFolder: true
	New(filter = "", file = "", status = "", .type = '', .hideLabel? = false)
		{
		super(filter, file, status, addons: .getAddons())
		.labelctrl = .hideLabel?
			? new class { Set(unused) {} }
			: .FindControl('label')
		.image = .FindControl('image')
		}

	getAddons()
		{
		addons = GetContributions('OpenImageAddons')
		if not .hideLabel?
			addons.Addon_Labels = true
		return addons
		}

	ImageHeight: 126
	TextHeight: 31
	Controls()
		{
		textWidth = imageSize = .ImageHeight
		textHeight = .TextHeight
		image = Object("Image",
			message: 'Double click to select or drag files to here',
			acceptDrop:, xstretch: false, xmin: imageSize, ymin: imageSize,
			name: 'image')
		if .hideLabel? is true
			return image
		return Object("Vert"
			image,
			#(EtchedLine before: 0, after: 0)
			Object(#Horz
				#(Skip small:)
				Object("Static", '', name: 'label', xmin: textWidth, ymin: textHeight)))
		}
	labels: false
	GetLabels()
		{
		return .labels
		}

	GetType()
		{
		return .type
		}

	value: ''
	subfolder: ''
	Get()
		{
		return .subfolder is ''
			? .value
			: Paths.Combine(.subfolder, .value)
		}

	isFirstSet: true

	Set(value)
		{
		if .image is false
			return
		split = .SplitLabel(value)
		.File = split.file
		.subfolder = split.subfolder
		.value = Paths.Basename(.File) is '' ? '' : .File

		errOb = Object()
		if not .isFirstSet and split.labels isnt false and
			.image.GetReadOnly() is false and
			true isnt errOb = .checkFileSize(split.labels, .FullPath())
			{
			.AlertError(.Title, errOb.error)
			split.labels = split.labels.Split(',').Map!(#Trim).
				Remove(errOb.removedLabel).Join(',')
			}
		.isFirstSet = false
		.image.Set(.FullPath(), .Highlight?(.value), message: .value)
		.SetLabels(split.labels)
		.SetTip()
		}

	checkFileSize(labels, file)
		{
		handlers = Contributions('OpenImageLabelsValidation')
		for handler in handlers
			if Object?(err = handler(labels, file))
				return err
		return true
		}

	SetLabels(labels)
		{
		.labels =.CleanupLabels(labels)
		.labelctrl.Set(.labels)
		.attachlabels()
		}

	CleanupLabels(labels)
		{
		return labels.Split(',').Map!(#Trim).Remove('').Sort!().Unique!().Join(', ')
		}

	attachlabels()
		{
		if .labels is "" and .value.Has?(.LabelDelimiter)
			.value = .value.BeforeLast(.LabelDelimiter)
		else if .labels isnt false and .labels isnt ''
			.value = (.value.Has?(.LabelDelimiter)
				? .value.BeforeLast(.LabelDelimiter) : .value) $
				.LabelDelimiter $ .labels
		}

	LabelDelimiter: ' Axon Label: '
	SplitLabel(value)
		{
		if not String?(value)
			value = ''
		file = labels = subfolder = ''
		if not value.Has?(.LabelDelimiter)
			file = value.Trim()
		else
			{
			file = value.BeforeLast(.LabelDelimiter).Trim()
			labels = value.AfterLast(.LabelDelimiter)
			}
		if .CopyAndLinkPath?(file)
			{
			subfolder = Paths.ParentOf(file)
			file = Paths.Basename(file)
			}
		return Object(:file, :labels, :subfolder)
		}

	// a copyAndLink path is like subfolder/filename
	CopyAndLinkPath?(path)
		{
		// It's not likely we will encounter a Linux absolute path here (because
		// we only support Suneido.js on linux and Suneido.js only supports copy & link)
		// Otherwise, we may need to also check CheckDirectory.IsLinuxAbsolutePath?(path)
		return not CheckDirectory.IsDrive?(path) and
			not CheckDirectory.IsFullUNC?(path) and
			path isnt Paths.Basename(path) // path is not a pure file name
		}

	SplitFile(file, label = false)
		{
		pathOb = .SplitLabel(file)
		if not .hasLabel?(pathOb, label)
			return ''
		filename = Paths.Basename(pathOb.file)
		return filename
		}

	SplitFullPath(file, labelFilter = false)
		{
		pathOb = .SplitLabel(file)
		if not .hasLabel?(pathOb, labelFilter)
			return ''
		fullPath = .FullPath(pathOb.file, pathOb.subfolder)
		return fullPath
		}

	hasLabel?(pathOb, labelFilter)
		{
		return labelFilter is false or
			pathOb.labels.Split(',').Map(#Trim).Has?(labelFilter.Trim())
		}

	File: false
	FullPath(file = false, subfolder = false)
		{
		if file is false
			if false is file = .File
				return ''

		if subfolder is false
			subfolder = .subfolder
		copyTo = .getCopyTo(subfolder)
		return Paths.Basename(file) is file and file isnt "" and copyTo isnt ""
			? Paths.Combine(copyTo, subfolder, file)
			: file
		}

	getCopyTo(subfolder)
		{
		Plugins().ForeachContribution('AttachmentPaths', 'subfolder', showErrors:)
			{ |c|
			path = c.getPath
			if false isnt path = path(subfolder)
				return path
			}
		return ""
		}

	GetCopyTo()
		{
		return .getCopyTo(.subfolder)
		}

	SetTip()
		{
		.image.SetTip(.FullPath() $ Opt(' - Labels: ', .labels))
		if '' isnt .labels
			.labelctrl.ToolTip('Labels: ' $ .labels)
		}

	image: false
	GetImageControl()
		{
		return .image
		}

	ProcessValue(value)
		{
		return value $ .labelString()
		}

	labelString()
		{
		return .labels isnt false and .labels isnt ''
			? .LabelDelimiter $ .labels
			: ''
		}

	Static_ContextMenu(x, y)
		{
		menu = Object(Addon_Labels.ContextMenu(.GetImageReadOnly()).Delete(#order))
		ContextMenu(menu).ShowCall(this, x, y)
		return true // to avoid extra context menu from static
		}

	ZoomReadonly(value)
		{
		if value is ''
			return
		split = .SplitLabel(value)
		super.ZoomReadonly(.FullPath(split.file, split.subfolder))
		}
	}
