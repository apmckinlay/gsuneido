// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Unsortable: true
	UseSubFolder: false
	Xmin: 90
	Ymin: 90
	New(filter = "", file = "", status = "", .reorderOnly = false, addons = #())
		{
		.File = file
		if Paths.Basename(.File) isnt "" and
			.File.BeforeLast(Paths.Basename(.File)) isnt ""
			.Set(.File)

		if filter is ""
			filter = "Image Files (*.bmp;*.gif;*.jpg;*.jpe;*.jpeg;*.ico;*.emf;*.wmf;" $
				"*.tif;*.tiff;*.png;*.pdf;*.txt;*.doc;*.docx;*.xls;*.xlsx)\x00" $
				"*.bmp;*.gif;*.jpg;*.jpe;*.jpeg;*.ico;*.emf;*.wmf;*.tif;*.tiff;" $
				"*.png;*.pdf;*.txt;*.doc;*.docx;*.xls;*.xlsx\x00" $
				"All Files (*.*)\x00*.*"
		.SetFilter(filter)
		.SetStatus(status)
		.addons = AddonManager(this, addons)

		.Send('Data')
		}

	GetFileName()
		{
		// .Get() returns the File with the Labels
		// to get just the File name use .File
		return String?(.Get())
			? Paths.Basename(.Get())
			: false // not a string but binary
		}

	GetFilePath()
		{
		return String?(.Get())
			? .Get().BeforeLast(.GetFileName())
			: false
		}

	SetFilter(filter)
		{
		.filter = filter
		}

	GetFilter()
		{
		return .filter
		}

	SetStatus(status)
		{
		.status = status
		}

	GetStatus()
		{
		return .status
		}

	GetReorderOnly()
		{
		return .reorderOnly
		}

	GetImageReadOnly()
		{
		.GetImageControl().GetReadOnly()
		}

	SetNewValue(file)
		{
		.Send('Field_SetFocus')
		.Set(Paths.Basename(file) is "" ? "" : file)
		.Send('NewValue', .Get())
		}

	ProcessFile(file)
		{
		process = OptContribution('OpenImageProcessFile', .processFile)
		return process(file, .errorHandler)
		}

	processFile(file, errorHandlerFn)
		{
		return not FileExists?(file) ? errorHandlerFn(file) : file
		}

	errorHandler(file, error = false)
		{
		msg = error is false ? 'Could not find file' : 'Could not access file'
		alert = error is false ? .AlertInfo : .AlertError
		alert('View Attachment', .GetAlertMsg(error, msg $ ': ' $ file))
		SuneidoLog('INFO: View Attachment: ' $ msg, params: [:file, :error])
		return false
		}

	GetAlertMsg(error, msg)
		{
		if error is false
			return msg
		return msg $ Opt(' (', error.AfterLast(':').Trim(), ')')
		}

	Attach(file)
		{
		if .Destroyed?()
			return
		if .reorderOnly
			return
		// file is the source, but after Result file it is the destination (as it copied)
		// if UseSubFolder is true it won't  have that prefix in the path, need to check first
		// This is to prevent errored scans that would attach a 0 sized file
		if false is .fileExist?(file)
			return
		origFile = .FullPath()
		if false isnt file = OpenImageSelect.ResultFile(file, useSubFolder: .UseSubFolder)
			if 0 is .Send("ImageDropFile", file)
				{
				.SetNewValue(.ProcessValue(file))
				action = origFile is '' ? '' : AttachmentsManager.ReplaceAction
				.QueueDeleteAttachment(.FullPath(), origFile, action)
				}
		}

	fileExist?(file)
		{
		if '' is bucket = AttachmentS3Bucket()
			return FileExists?(file)
		region = AmazonS3.GetBucketLocationCached(bucket)
		return AmazonS3.FileExist?(bucket, file, region)
		}

	Highlight?(value)
		{
		return value isnt ""
		}

	Open(file = false, hwnd = false)
		{
		Finally({
			.SetEnabled(false)
			if 0 isnt fullPath = .Send('ImageGetFullPath', file is false ? .File : file)
				file = fullPath
			if file is false
				file = .FullPath()
			if hwnd is false
				hwnd = .Window.Hwnd
			.AttemptOpenFile(file, hwnd)
			}, {
			.SetEnabled(true)
			})
		}

	AttemptOpenFile(file, hwnd)
		{
		if ExecutableExtension?(file)
			{
			.AlertInfo('Executable File', 'The attachment: ' $ file $
				'\r\nAppears to be an executable or zip file and will not open from Axon')
			return false
			}
		if false isnt existFile = .ProcessFile(file)
			ShellExecute(hwnd, NULL, Paths.ToLocal(existFile),
				fMask: SEE_MASK.ASYNCOK)
		}

	Default(@args) // used by context menu
		{
		if args[0].Prefix?('On_')
			{
			args[0] = args[0].Replace("^On_Context_", "On_")
			return .addons.Send(@args)
			}
		}

	// Image Event Handlers
	ImageDoubleClick()
		{
		if .File isnt ""
			.Open()
		else if not .GetImageReadOnly() and not .reorderOnly
			.addons.Send(#On_Select)
		}

	ImageContextMenu(x, y)
		{
		menu = Object()
		for menuList in .addons.Collect("ContextMenu")
			menu.Add(@menuList)
		menu.Sort!(By(#order))
		if menu isnt []
			ContextMenu(menu.Each({ it.Delete(#order) })).ShowCall(this, x, y)
		}

	ImageDropFiles(hDrop)
		{
		.addons.Send(#ImageDropFiles, hDrop)
		}

	ImageStartDrag()
		{
		.Send("ImageStartDrag")
		}

	// from suneido.js ImageComponent
	ImageEndDrag()
		{
		.Send("ImageEndDrag")
		}

	ImageFinishDrag()
		{
		.Send("ImageFinishDrag")
		}

	ImageMouseMove()
		{
		.Send("ImageMouseMove")
		}

	ImageMouseLeave(dragging)
		{
		.Send("ImageMouseLeave", dragging)
		}

	// abstract
	Get(@unused) { .accessAbstract(#Get) }
	Set(@unused) { .accessAbstract(#Set) }
	FullPath(@unused) { .accessAbstract(#FullPath) }
	SetTip(@unused) { .accessAbstract(#SetTip) }
	GetImageControl(@unused) { .accessAbstract(#GetImageControl) }
	ProcessValue(@unused) { .accessAbstract(#ProcessValue) }

	accessAbstract(method)
		{
		throw method $ ' should be defined in a subclass'
		}

	GetCopyTo()
		{
		return ''
		}

	ZoomReadonly(value)
		{
		if value isnt ''
			.AttemptOpenFile(value, 0)
		}

	QueueDeleteAttachment(newFile, oldFile, action = '')
		{
		if not OpenImageSettings.Normally_linkcopy?()
			return
		name = .fieldName()
		if .Send('QueueDeleteAttachmentFile', newFile, oldFile, name, action) is 0
			SuneidoLog("INFO: QueueDeleteAttachmentFile not implemented",
				params: Object(:newFile, :oldFile, :name), calls:)
		}

	fieldName()
		{
		if 0 is name = .Send('AttachmentFieldName')
			name = .Name
		return name
		}

	Destroy()
		{
		.Send('NoData')
		super.Destroy()
		}
	}