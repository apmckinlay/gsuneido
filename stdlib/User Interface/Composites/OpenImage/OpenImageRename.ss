// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Title: 'Rename Attachment'
	CallClass(fullPath, hwnd, parentCtrl, copyTo = '')
		{
		return OkCancel([this, fullPath, copyTo, parentCtrl], .Title, hwnd)
		}

	New(.fullPath, .copyTo, .parentCtrl)
		{
		.fileNameCtrl = .FindControl('renamedFileName')
		endPos = .oldFileName.BeforeLast('.').Size()
		.fileNameCtrl.SetSel(0, endPos)
		}

	Controls()
		{
		.oldFileName = Paths.Basename(.fullPath)
		return Object(#Vert
			Object(#Horz, #(Static 'Rename file to: '),
				Object(#Field, set: .oldFileName, mandatory:, name: 'renamedFileName')))
		}

	OK()
		{
		filename = .fileNameCtrl.Get()
		if true isnt msg = .valid(filename)
			{
			.AlertInfo('Invalid file name', msg)
			return false
			}

		if filename.AfterLast('.').Blank?() and not .oldFileName.AfterLast('.').Blank?()
			filename = filename.RightTrim('.') $ '.' $ .oldFileName.AfterLast('.')
		if .oldFileName is filename
			return 'File not renamed'

		newFilePath = .getCopyToFilename(filename)
		if newFilePath.Size() > CheckFileName.MaxAllowedFileNameChars
			{
			.AlertInfo('Invalid file name',
				'Destination path is too long. Please choose a shorter file name')
			return false
			}

		return .attemptToRename(newFilePath)
		}

	valid(filename)
		{
		if true isnt msg = .basicFileValidation(filename)
			return msg

		newExtension = filename.AfterLast('.')
		oldExtension = .oldFileName.AfterLast('.')
		if ((newExtension is 'pdf' or oldExtension is 'pdf') and
			newExtension isnt oldExtension and not newExtension.Blank?())
			return 'Cannot convert pdfs'
		return true
		}

	basicFileValidation(filename)
		{
		if filename.Blank?()
			return 'Please enter a File Name'
		if '' isnt msg = CheckFileName.WithErrorMsg(filename)
			return msg $ '\nPlease enter another File Name.'
		return true
		}

	getCopyToFilename(filename)
		{
		return OpenImageSelect.GetCopyToFilename(
			Paths.EnsureTrailingSlash(Paths.ParentOf(.fullPath)), filename)
		}

	attemptToRename(newFilePath)
		{
		try
			{
			if true isnt result = .processFileRename(newFilePath)
				throw result
			}
		catch (e)
			{
			if not e.Has?('try again')
				e $= '.\r\nPlease try again'
			.AlertInfo('Rename Attachment', e)
			if AttachmentS3Bucket() isnt ''
				return false

			// Do not delete the file if it's the rename itself that failed
			// due to potential network hiccups, it's possible the rename could have
			// failed because an existnig file with the same name could already exist
			// a delete at this point would delete an unrealed file.
			if not e.Has?('Failed to rename attachment')
				DeleteFile(newFilePath)
			return false
			}
		return true
		}

	processFileRename(newFilePath)
		{
		if '' isnt msg = .copyAndCheck(newFilePath)
			return msg

		copyToRemove = .copyTo isnt ''
			? .copyTo
			: OpenImageSettings.Copyto()
		newFile = OpenImageSettings.Normally_linkcopy?() is true
			? newFilePath.RemovePrefix(copyToRemove.RightTrim(`\/`)).LeftTrim(`\/`)
			: newFilePath
		.parentCtrl.SetNewValue(.parentCtrl.ProcessValue(newFile))
		return true
		}

	copyAndCheck(newFilePath)
		{
		if '' isnt bucket = AttachmentS3Bucket()
			{
			if true isnt AmazonS3.CopyFile(bucket, FormatAttachmentPath(.fullPath),
				bucket, FormatAttachmentPath(newFilePath))
				return 'Failed to rename attachment on S3'
			return ''
			}

		// if the copy itself fails we don't want to delete
		// there is the potential that we could be deleting an already existing
		// attachment that is linked somewhere else
		if true isnt CopyFile(.fullPath, newFilePath, true)
			return 'Failed to rename attachment'

		if not FileExists?(newFilePath)
			return 'Failed to rename attachment'

		return ''
		}
	}
