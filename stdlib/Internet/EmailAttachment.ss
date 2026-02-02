// Copyright (C) 2008 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Title: 'Email Attachment'
	ComponentName: 'EmailAttachment'
	FileNotFoundErrorPrefix: 'Could not find the following file(s):\n\n'
	SendAttachmentAsLink: false
	CallClass(hwnd, block, subject = "")
		{
		if false is attachInfo = .attachInfo(block)
			return
		result = attachInfo.result
		filename = attachInfo.filename
		attachFileName = attachInfo.attachFileName
		attachments = attachInfo.attachments

		if String?(valid = .checkFileName(attachments.Copy().Add(filename)))
			{
			.AlertInfo(.Title, valid)
			return
			}

		merge_pdf? = result.GetDefault("merge_pdf?", false)
		emailSubject = result.GetDefault('EmailSubject', false)
		if false is x = ToolDialog(hwnd, [this, subject, attachments,
			merge_pdf?, filename, emailSubject, :attachFileName], closeButton?: false)
			return

		if not Object?(x.attachments)
			x.attachments = Object()

		if '' isnt AttachmentS3Bucket()
			return

		_sendError = Object()
		msg = .CreateMIME(x.attachments, filename, x, attachFileName, hwnd)

		if not _sendError.Empty?() and Sys.Client?()
			{
			err = 'Please fix the following issues and try again\n'
			err $= _sendError.Join('\n')
			.AlertError(.Title $ ': Not Sent', err)
			BookLog('Emailing Attachment(s) failed: ' $ err)
			return
			}

		if msg is false
			return

		BookSendEmail(hwnd, x.email_from, x.to, msg,
			pdfNames: Record(orig: filename, rename: attachFileName, :merge_pdf?,
				extraAttach: attachments))
		.DeleteCompressed(msg)
		}

	attachInfo(block)
		{
		result = String?(block) or Object?(block) ? block : block()
		filename = Object?(result) ? result.filename : result
		attachFileName = Object?(result)
			? result.GetDefault('attachFileName', filename)
			: filename
		if filename is false
			return false
		attachments = Object?(result) ? result.attachments : #()
		return Object(:result, :filename, :attachFileName, :attachments)
		}

	CreateMIME(attachments, filename, x, attachFileName, hwnd, quiet? = false)
		{
		attachOb = .attachmentList(Object(:filename, :attachments))

		if .SendAttachmentAsLink
			return .addAttachmentLinks(x, filename, attachFileName, hwnd, quiet?)

		originalSize = .CalculateTotal(attachOb.Flatten())
		if String?(originalSize)
			return false

		attachHandler = .getAttachmentHandler()
		result = .checkSize(originalSize)
		if .attachOriginal?(result, attachHandler)
			return .addAttachmentsToEmail(x, filename, attachFileName, quiet?)
		else if .invalidSize?(result)
			return false

		// over size limit
		if false is result = EmailAttachment_Mime.Compress(attachOb, quiet?, originalSize)
			return .addAttachmentLinks(x, filename, attachFileName, hwnd, quiet?)
		else if String?(result)
			return false

		// successfully compressed
		EmailAttachment_Mime.BuildCompressedAttachments(x, attachOb)
		filename = attachOb.filename

		return .addAttachmentsToEmail(x, filename, attachFileName, quiet?)
		}

	attachmentList(attachOb)
		{
		attachList = Object(attachments: Object(), filename: '')
		for attach in attachOb.GetDefault('attachments', #())
			{
			if Object?(attach) and attach.Member?('fileName')
				attachList.attachments.Add(attach.fileName)
			else if String?(attach)
				attachList.attachments.Add(attach)
			else
				SuneidoLog('ERROR: Unexpected attachOb', params: attachOb, calls:)
			}
		if attachOb.Member?('filename')
			attachList.filename = attachOb.filename
		return attachList
		}

	invalidSize?(result)
		{
		return String?(result) and not result.Has?('exceeds maximum')
		}

	attachOriginal?(result, attachHandler)
		{
		return result is true or attachHandler is false
		}

	CalculateTotal(attachments)
		{
		totalSize = 0
		err = Object()
		for attachment in attachments
			{
			if Number?(result = EmailAttachment_Mime.FileSize(attachment))
				totalSize += result
			else
				err.AddUnique(result)
			}

		if not err.Empty?()
			return err.Join('\n')

		return totalSize
		}

	getAttachmentHandler()
		{
		contrib = OptContribution("AddEmailAttachmentLinks", false)
		if not Object?(contrib) or contrib.Empty?()
			return false
		return Global(contrib[0])
		}

	DeleteCompressed(mime)
		{
		if not EmailAttachment_Mime.Buffered(mime)
			return
		for attachment in mime.GetAttachedFiles()
			if Object?(attachment) and .CompressedFile?(file = attachment.filename)
				DeleteFile(file)
		}

	CompressedFile?(filename)
		{
		return EmailAttachment_Mime.CompressedFile?(filename)
		}

	// Called Externally
	ValidateAttachments(attachments, maxSizeInMb = false)
		{
		maxSizeInMb = maxSizeInMb isnt false
			? maxSizeInMb
			: EmailMimeMaxSizeInMb()
		checkAttachments = attachments.Copy()
		if true isnt result = .checkFileName(checkAttachments)
			return result
		if true isnt result = .checkSize(checkAttachments, maxSizeInMb)
			return result
		return true
		}

	checkFileName(attachments)
		{
		for attach in attachments
			{
			fileExt = attach.AfterLast(".")
			if AmazonSES.InvalidAttachmentExtensions.Has?(fileExt)
				return "Attachments with extension " $ Display(fileExt) $
					" can not be sent through the Axon email service"
			}
		return true
		}
	// attachments can be either a list of attachments or total size of them (in bytes)
	checkSize(attachments, maxSizeInMb = false)
		{
		maxSize = maxSizeInMb isnt false
			? maxSizeInMb.Mb()
			: EmailMimeMaxSizeInMb().Mb()

		// check size of all attachments before loading/ email is written
		totalSize = 0
		if Number?(attachments)
			totalSize = attachments
		else if Object?(attachments)
			totalSize = .CalculateTotal(attachments)

		if String?(totalSize)
			return totalSize

		succeeded? = totalSize <= maxSize
		if not succeeded?
			{
			return "File size (" $	ReadableSize(totalSize) $ ") " $
				"exceeds maximum (" $ ReadableSize(maxSize) $ ")\n" $
				"Email will not be sent"
			}
		return true
		}

	addAttachmentLinks(x, filename, attachFileName, hwnd, quiet? = false)
		{
		addAttach = .getAttachmentHandler()
		return addAttach(x, filename, attachFileName, hwnd, quiet?)
		}
	addAttachmentsToEmail(x, filename, attachFileName, quiet? = false)
		{
		msg = false
		try
			{
			msg = BookSendEmail.CreateMime(x.subject, x.message, filename, attachFileName)
			for file in x.attachments
				{
				if Object?(file) is true
					msg.AttachFile(
						FileStorage.GetAccessibleFilePath(file.fileName),
						attachFileName: file.attachFileName)
				else
					msg.AttachFile(FileStorage.GetAccessibleFilePath(file))
				}
			}
		catch (err, "MimeMulti: AttachFile:")
			{
			if quiet? is false
				.AlertError(.Title, err.Replace("MimeMulti: AttachFile: ", ""))
			return false
			}
		return msg
		}

	New(subject, attachments = #(), .merge_pdf? = false, .filename = "",
		.emailSubject = false, .attachFileName = "")
		{
		super(.layout(attachments, merge_pdf?))
		.Data.AddObserver(.RecordChanged)
		if .emailSubject isnt false and
			false isnt saved = Email_DefaultSubject.GetSavedSubject(.emailSubject.type)
			subject = Email_DefaultSubject.BuildSubject(saved,
				.emailSubject.cols, .emailSubject.data)
		.Data.Set([:subject]) // so rules kick in
		.Data.Dirty?(true) // ensure validation
		if not attachments.Empty?()
			{
			if merge_pdf?
				{
				mergeableFiles = PdfMerger.FilterFiles(attachments)
				formattedList = mergeableFiles.Map({"(APPEND) " $ it}).
					Append(attachments.Difference(mergeableFiles))
				.Data.SetField('attachments', formattedList)
				}
			else
				.Data.SetField('attachments', attachments)
			}
		.Data.SetField('email_signature_enabled',
			UserSettings.Get('email_signature_enabled', true))
		.Data.SetField('email_signature', UserSettings.Get('email_signature', ''))
		}
	layout(attachments = #(), merge_pdf? = false)
		{
		ob = Object('Vert'
			#(Pair (Static 'Reply To')
				(EmailAddress mandatory:, xstretch: 1 name: email_from))
			#(Skip 3)
			#(Pair (Static To)
				(EmailAddresses mandatory:, width: 40, xstretch: 1 name: to))
			#(Skip 3)
			Object('Pair'
				#(Static Subject)
				Object('Horz'
					#(Field mandatory:, xstretch: 1 name: subject),
					.emailSubject isnt false
						? #(Horz #(Skip, small:) #(LinkButton, "Template"))
						: ''))
			#(Skip 3))

		if not attachments.Empty?()
			.attachmentsOption(merge_pdf?, ob, attachments)

		ob.Add(
			#(ScintillaAddonsEditor, name: preview)
			#(Skip 3)
			#(CheckBox "Include Signature", name: "email_signature_enabled")
			#(ScintillaAddonsEditor, height: 4, ystretch: 0,
				name: 'email_signature'))

		return Object('Record', Object('Vert',
			Object('Scroll', ob),
			#(Skip 5),
			#('SendCancel')
			) xmin: 700, ymin: 500)
		}

	attachmentsOption(merge_pdf?, ob, attachments)
		{
		if merge_pdf?
			{
			mergeableFiles = PdfMerger.FilterFiles(attachments)
			formattedList = mergeableFiles.Map({"(APPEND) " $ it}).
				Append(attachments.Difference(mergeableFiles))
			list = .toChooseAsObject(formattedList)

			ob.Add(Object('Pair'
				#(Static 'Attachments')
				Object('ChooseManyAsObject', name: 'attachments', idField: 'File',
					cols: #(File) xstretch: 1, displayField: 'File', :list)))

			ob.Add(#(Pair (Static '')
				#(Horz (Static 'PDF and JPG will be appended to generated PDF, ')
					#(LinkButton "adjust the order of appended files", "Reorder"))))
			ob.Add(#(Skip 3))
			}
		else
			ob.Add(Object('Pair'
				#(Static 'Attachments')
				Object('ChooseManyAsObject', name: 'attachments', idField: 'File',
					cols: #(File) xstretch: 1, displayField: 'File',
					list: .toChooseAsObject(attachments))))
		}

	On_Template()
		{
		if false is subject = Email_DefaultSubject(.Window.Hwnd, .emailSubject)
			return
		.Data.SetField('subject', subject)
		}

	RecordChanged(member)
		{
		if member is 'email_signature_enabled'
			.FindControl('email_signature').SetReadOnly(
				.Data.Get().email_signature_enabled is false)
		}

	AttachmentFilesExist(attchOb)
		{
		if attchOb is ''
			return ''
		return Opt(.FileNotFoundErrorPrefix,
			attchOb.Filter(.fileNotAccessible?).UniqueValues().Join('\n'),
			'\n\nPlease ensure that they exist and ',
			'that you have permission to access them.')
		}

	fileNotAccessible?(file)
		{
		try
			return not FileStorage.Exists?(file)
		return true
		}

	On_Send()
		{
		if .Data.Valid() isnt true
			return
		// There is a hard limit of 998 characters in a subject line
		// 256 is consistent with other email software
		charLimit = 256
		if .Data.GetField('subject').Size() >= charLimit
			{
			.AlertWarn(.Title, 'Subject cannot be longer than ' $ charLimit $
				' characters')
			return
			}
		data = .Data.Get()
		mergeableFiles = unMergeableFiles = false
		if .merge_pdf? and Object?(data.attachments)
			{
			if data.attachments.Empty?()
				{
				.AlertWarn(.Title, "Please select at least one attachment")
				return
				}
			files = .getAppendable(data).Add(.filename, at: 0)
			mergeableFiles = PdfMerger.FilterFiles(files)
			unMergeableFiles = data.attachments.Filter({not it.Prefix?("(APPEND) ")}).
				Append(files.Difference(mergeableFiles))
			allAttachments = unMergeableFiles.Copy().Append(mergeableFiles)
			}
		else
			allAttachments = data.attachments
		if not .mergeAttachments(allAttachments, data, mergeableFiles, unMergeableFiles)
			return
		for field in #(email_from, email_signature_enabled, email_signature)
			UserSettings.Put(field, data[field])

		.setMessage(data)
		merge_pdf? = .merge_pdf? and data.attachments.Size() > 0
		EmailAttachment_Mime.SendJS(this, .filename, .attachFileName, merge_pdf?,
			data, mergeableFiles, unMergeableFiles)
		.Window.Result(data)
		}

	EmailAsLinks(x, filename, attachFileName, hwnd, quiet? = false)
		{
		attachHandler = .getAttachmentHandler()
		if false is msg = attachHandler(x, filename, attachFileName, hwnd, quiet?)
			return
		BookSendEmail(hwnd, x.email_from, x.to, msg,
			pdfNames: Record(orig: filename, rename: attachFileName,
				extraAttach: x.attachments))
		}

	mergeAttachments(allAttachments, data, mergeableFiles, unMergeableFiles)
		{
		if '' isnt AttachmentS3Bucket()
			return true
		if "" isnt msg = .AttachmentFilesExist(allAttachments)
			{
			.AlertError(.Title, msg)
			return false
			}
		if .merge_pdf? and Object?(data.attachments)
			{
			invalidFiles = PdfMerger(mergeableFiles, .filename)
			if not invalidFiles.Empty?()
				{
				.AlertError(.Title, PdfMerger.InvalidFilesMsg(invalidFiles))
				return false
				}
			data.attachments = unMergeableFiles
			}
		return true
		}

	setMessage(data)
		{
		data.message = data.preview.Blank?() ? "see attached" : data.preview
		if data.email_signature_enabled and not data.email_signature.Blank?()
			data.message $= '\r\n--\r\n' $ data.email_signature
		}

	On_Reorder()
		{
		data = .Data.Get()
		mergeableFiles = .getAppendable(data)
		if false is newAppend = ReorderAttachments(.Parent.Hwnd, mergeableFiles)
			return
		.updateAttachmentsList(mergeableFiles, newAppend, data)
		}

	updateAttachmentsList(mergeableFiles, newAppend, data)
		{
		// Used as the replacement for entries that need (APPEND) removed
		appendRemoved = mergeableFiles.Difference(newAppend)
		// Used to find the entries that need (APPEND) removed
		removeAppend = appendRemoved.Map({ "(APPEND) " $ it })
		newAppend.Map!({ "(APPEND) " $ it })
		oldListData = data.attachments
		newListData = .buildNewList(oldListData, newAppend, appendRemoved, removeAppend)
		data.attachments = newListData
		}

	getAppendable(data)
		{
		return  data.attachments.
			Filter({ it.Prefix?("(APPEND) ") }).
			Map!({ it.Replace("\(APPEND\) ", "") })
		}

	buildNewList(oldList, append, appendRemoved, removeAppend)
		{
		return append.Copy().Append(appendRemoved).
			Append(oldList.Difference(removeAppend).Difference(append))
		}

	toChooseAsObject(list)
		{
		return list.Map({ Object(File: it) })
		}
	}
