// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Title: 'Select an Attachment'
	New(value, filter, status, usesubfolder = false)
		{
		super(.controls(value, filter, status))
		.value = value
		.openfile = .Vert.OpenFile
		.mode = .Vert.RadioButtons
		if OpenImageSettings.Copyto() is ""
			.mode.SetEnabledChild(.linkcopy_i, false)
		if Paths.IsValid?(value)
			{
			.openfile.Set(value)
			.mode.Set(.link)
			}
		else
			{
			.mode.Set(OpenImageSettings.Normally_linkcopy?()
				? .linkcopy : .link)
			}
		.useSubFolder = usesubfolder
		}
	controls(value, filter, status)
		{
		return Object('Vert'
			Object('OpenFile' title: 'Image', :filter, file: value, :status)
			'Skip'
			Object('RadioButtons', .link, .linkcopy)
			'Skip',
			.copyAndLinkSettings(),
			.currentDirectory()
			'OkCancel'
			)
		}
	copyAndLinkSettings()
		{
		return #(Horz
			Skip
			(LinkButton 'Copy and Link Settings')
			)
		}
	currentDirectory()
		{
		dir = Suneido.Member?('currentUserAttachmentSettings')
			? Suneido.currentUserAttachmentSettings.copyto
			: OpenImageSettings.Copyto()
		return dir is ""
			? #(Skip 0)
			: Object('Vert',
				Object('Static' 'Copy and Link will go to: ' $ dir)
				'Skip')
		}
	link:		"Link (just stores file location)"
	linkcopy:	"Copy and Link"
	linkcopy_i: 1
	openfile: false
	NewValue(value/*unused*/, source)
		{
		if .openfile is false or // during construction
			source isnt .openfile
			return
		}
	On_Copy_and_Link_Settings()
		{
		ToolDialog(.Window.Hwnd, OpenImageSettings)
		if OpenImageSettings.Copyto() is ""
			{
			if .mode.Get() is .linkcopy
				.mode.Set(.link)
			.mode.SetEnabledChild(.linkcopy_i, false)
			}
		else
			.mode.SetEnabledChild(.linkcopy_i, true)
		}
	On_OK()
		{
		file = .openfile.Get()
		if false isnt file = .ResultFile(file, .mode.Get() is .linkcopy,
			.useSubFolder)
			.Window.Result(file)
		}
	ResultFile(file, linkcopy = "", useSubFolder = false)
		{
		if .s3Bucket() isnt ''
			return file
		if not .validFile?(file)
			return false

		if linkcopy is ""
			linkcopy = OpenImageSettings.Normally_linkcopy?()

		if linkcopy isnt true
			return file

		copyto = Suneido.Member?('currentUserAttachmentSettings')
			? Paths.EnsureTrailingSlash(Suneido.currentUserAttachmentSettings.copyto)
			: false
		deletesource? = OpenImageSettings.DeleteSource?()
		subfolder = Suneido.Member?('currentUserAttachmentSettings')
			? ''
			: .subFolder(useSubFolder)

		if false is file = .CopyFile(file, :copyto, :deletesource?, :subfolder)
			return false

		fileSub = subfolder $ Paths.Basename(file)
		if file.BeforeLast(fileSub) is OpenImageSettings.Copyto()
			file = fileSub

		return file
		}

	// extracted for test
	s3Bucket()
		{
		return AttachmentS3Bucket()
		}

	validFile?(file)
		{
		if file is '' // user clicked cancel
			return false
		if '' isnt msg = CheckFileName.WithErrorMsg(file, withPath?:)
			return .alert('Could not copy ' $ file $ '.\n' $ msg, false)

		try
			fs = .fileSize(file)
		catch (e, 'FileSize:')
			{
			try exists = FileExists?(file)
			catch(err)
				exists = err
			SuneidoLog('INFO: ' $ e, params: exists)
			.AlertWarn(.Title, "Can't access file or path: " $ file)
			return false
			}
		if fs is 0
			{
			// filename with unicode could be changed to ascii through win32 api,
			// which can cause file not accessible
			.AlertWarn(.Title, "The file cannot be attached.\r\nIt may be empty, " $
				"not accessible, or the file name may include non-standard characters.")
			return false
			}
		return true
		}

	// extracted for test
	fileSize(file)
		{
		return FileSize(file)
		}

	copyFileWarning(file, dest)
		{
		extraParams = Object()
		warnings = LastContribution("OpenImageCopyFileWarning")
		try
			{
			msg = FileExists?(file) ? warnings.copyLinkIssue : warnings.fileNotExistIssue
			}
		catch (error)
			{
			msg = [logMsg: 'INFO: Attachments - Could not access file',
				userMsg: OpenImageAddonsBase.GetAlertMsg(error,
				'Could not access file' $ ': ' $ file)]
			extraParams.error = error
			}
		if msg.logMsg isnt false
			SuneidoLog(msg.logMsg, params: Record(:file, :dest).Merge(extraParams))
		return msg.userMsg
		}
	subFolder(useSubFolder) // override in test
		{
		return (useSubFolder isnt false ? .SubFolder() : '')
		}

	SubFolder(date = false)
		{
		if date is false
			date = Date()
		return Paths.EnsureTrailingSlash(date.Format('yyyyMM'))
		}

	GetCopyTo()
		{
		return OpenImageSettings.Copyto()
		}

	CopyFile(file, copyto = false, deletesource? = false, quiet? = false, subfolder = '',
		fromBucket = '')
		{
		if file is ""
			return false

		fileBasename = Paths.Basename(file)
		if '' isnt  msg = CheckFileName.WithErrorMsg(fileBasename)
			return .alert('Could not copy ' $ fileBasename $ '.\n' $ msg, quiet?)

		copyfolder = .ensureCopyFolder(subfolder, copyto)
		if Object?(copyfolder)
			return .alert(copyfolder.msg, quiet?)
		rotated = .autoRotateIfNeeded(file)
		result = .retryCopy(copyfolder, fileBasename, copyto, rotated, fromBucket)
		if String?(result)
			return .alert(result, quiet?)

		if result.copy is false
			return .alert(.copyFileWarning(file, result.dest), quiet?)

		if deletesource? and not .deleteFile?(file)
			{
			.deleteFile?(result.dest)
			SuneidoLog('OpenImageSelect.CopyFile: ' $
				'can not delete source file: ' $ file $ ' err: Retry failed')
			return false
			}
		return result.dest
		}

	retryCopy(copyfolder, fileBasename, copyto, file, fromBucket = '')
		{
		dest = .GetCopyToFilename(copyfolder, fileBasename)
		if dest.Size() > CheckFileName.MaxAllowedFileNameChars
			return "The destination file path is too long: " $ dest
		if copyto is ""
			return Object(copy: false, :dest)

		if '' isnt bucket = .s3Bucket()
			{
			if fromBucket isnt ''
				{
				copy = AmazonS3.CopyFile(
					fromBucket, file, bucket, FormatAttachmentPath(dest))
				return Object(:copy, :dest)
				}
			copy = FileStorage.SaveFile(file, dest) isnt false
			if dest.Prefix?(GetAppTempPath())
				.deleteFile?(dest)
			return Object(:copy, :dest)
			}

		if true is x = .copyFile(file, dest)
			return Object(copy: true, :dest)

		if x is 'Retry'
			{
			dest = .GetCopyToFilename(copyfolder, fileBasename)
			if true is .copyFile(file, dest)
				return Object(copy: true, :dest)
			return Object(copy: false, :dest)
			}
		return Object(copy: false, :dest)
		}

	autoRotateIfNeeded(dest)
		{
		if dest.AfterLast('.').Lower() not in ('jpg', 'jpeg')
			return dest

		if not ImageHandler.Available?()
			return dest

		if false is ImageHandler.Orientation(dest)
			return dest

		ext = dest.AfterLast('.')
		tmp = GetAppTempFullFileName() $ ext
		if true isnt ret = ImageHandler.AutoRotate(dest, tmp)
			{
			deleteTmp = DeleteFile(tmp)
			SuneidoLog('INFO: ImageHandler.AutoRotate failed', calls:,
				params: [:dest, :ret, :tmp, :deleteTmp])
			return dest
			}
		return tmp
		}

	GetCopyToFilename(copyfolder, fileBasename)
		{
		fileExistsFn = false
		if '' isnt bucket = .s3Bucket()
			fileExistsFn = { |file|
				AmazonS3.FileExist?(:bucket, file: FormatAttachmentPath(file))
				}
		return MakeUniqueFileName(copyfolder, fileBasename, :fileExistsFn).dest
		}

	// Want to force the CopyFile to fail if file already exists
	// in the case where the file already exists we want to re-try MakeUniqueFileName
	// and try the copy again. (only once) (Possibility of network hiccup)
	copyFile(file, dest)
		{
		msg = CopyFile(file, dest, true)
		return (String?(msg) and msg =~ "CopyFile: The file exists.")
			? "Retry" : (true is msg)
		}

	deleteFile?(file)
		{ return true is DeleteFile(file) }

	ensureCopyFolder(subfolder, copyto)
		{
		if copyto is false
			copyto = .GetCopyTo()
		if subfolder is true
			subfolder = .subFolder(true)
		if subfolder is ''
			return copyto
		copyfolder = copyto $ subfolder
		if '' isnt .s3Bucket()
			return copyfolder
		if String?(exists? = CheckDirExists(copyfolder))
			return Object(msg: exists?)

		if not exists? and not .createDir?(copyfolder)
			{
			SuneidoLog('ERRATIC: (CAUGHT) Could not create folder: ' $ copyfolder)
			return copyto
			}
		return copyfolder
		}

	createDir?(folder)
		{ return true is CreateDir(folder.RightTrim('\\/')) }

	alert(message, quiet?)
		{
		if quiet? is false
			.AlertWarn(.Title, message)
		return false
		}
	}
