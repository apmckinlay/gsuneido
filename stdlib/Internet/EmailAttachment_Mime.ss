// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
class
	{
	Buffered(mime)
		{
		return Instance?(mime) and mime.Base?(MimeMultiBuffered)
		}

	Compress(attachOb, quiet?, originalSize = false)
		{
		// newAttachments will be an object of attachments, if able to compress, they
		// will be compressed, otherwise will be same as the original
		result = Object()
		try
			{
			Working('Compressing Attachments...', :quiet?)
				{
				result = .compressAttach(attachOb, originalSize)
				}
			}
		catch (err)
			{
			SuneidoLog("EmailAttachment: Unable to compress attachment - " $ err, calls:)
			return false
			}
		if true isnt continue? = .continue?(@result)
			return continue?

		attachOb.attachments = result.newAttachments
		attachOb.filename = result.newFileName
		return true
		}

	compressAttach(attachOb, originalSize)
		{
		result = Object(newAttachments: Object(), newFileName: '', :originalSize,
			uncompressedAttachments: Object(), total: 0, quit?: false, sizeError?: false)
		for filename in attachOb.attachments.Copy().Add(attachOb.filename)
			{
			file = .compressOneFile(filename, result.uncompressedAttachments)
			if Number?(size = .FileSize(file))
				result.total += size
			else
				{
				result.sizeError? = true
				break
				}
			if result.total > EmailMimeMaxSizeInMb().Mb()
				{
				result.quit? = true
				break
				}
			if filename is attachOb.filename
				result.newFileName = file
			else
				result.newAttachments.Add(file)
			}
		return result
		}

	compressSuffix: '_temp_compressed'
	compressOneFile(filename, uncompressedAttachments)
		{
		// POST: compressed file name is returned if successful, otherwise original name
		newFileName = .getTempName(filename)
		compressedFile = filename
		if filename.Lower().Suffix?('pdf')
			{
			result = PdfMerger(Object(filename), newFileName, compress:,
					maxCompressedFileSizeInMb: EmailMimeMaxSizeInMb())

			if result.Empty?()
				compressedFile = newFileName
			else
				uncompressedAttachments.Add(compressedFile $ ' reason: ' $
					result.Join(', '))

			return compressedFile
			}
		filename = FileStorage.GetAccessibleFilePath(filename)
		if true is ImageHandler.Compress(filename, newFileName)
			compressedFile = newFileName
		else
			uncompressedAttachments.Add(
				filename $ ' reason: compressing image file failed')

		return compressedFile
		}

	getTempName(filename)
		{
		return GetAppTempPath() $ Paths.Basename(filename) $ .compressSuffix $
			Display(Timestamp()).Tr('#.', '')
		}

	continue?(quit?, sizeError?, total, originalSize, uncompressedAttachments)
		{
		if quit? or sizeError?
			{
			.outputBookLog(total, originalSize, uncompressedAttachments)

			if sizeError?
				return 'Size Error'

			return false
			}

		return true
		}


	outputBookLog(totalSize, originalSize, uncompressedAttachments)
		{
		msg = 'Building MIME for compressed attachment(s) failed: original size ' $
			ReadableSize(originalSize) $ '; compressed size ' $ ReadableSize(totalSize)
		msg $= Opt('; uncompressed attachments: ', uncompressedAttachments.Join(', '))
		BookLog(msg)
		}

	FileSize(file, _sendError = false)
		{
		try
			{
			fs = FileStorage.FileSize(file)
			}
		catch (e)
			{
			if Object?(sendError)
				sendError.AddUnique('Can not access file or path: ' $ file)
			prefix = OptContribution('Hosted?', function() { return false })()
				? 'ERROR: (CAUGHT) '
				: 'INFO: '
			SuneidoLog(prefix $ e, calls:)
			return 'Can not access file or path: ' $ file
			}
		if fs in (0, false)
			{
			if Object?(sendError)
				sendError.AddUnique('Can not attach file. Either file is empty or ' $
					'there was a problem accessing file: ' $ file)
			return 'Can not attach file. Either file is empty or ' $
				'there was a problem accessing file: ' $ file
			}

		return fs
		}

	BuildCompressedAttachments(x, attachOb)
		{
		x.attachments = Object()
		for attachment in attachOb.attachments
			{
			displayName = attachment
			if displayName.Has?(.compressSuffix)
				displayName = attachment.BeforeLast(.compressSuffix)
			x.attachments.Add(Object(fileName: attachment, attachFileName: displayName))
			}
		}

	CompressedFile?(filename)
		{
		return filename.Prefix?(GetAppTempPath()) and filename.Has?(.compressSuffix)
		}

	SendJS(controller, filename, attachFileName, merge_pdf?, data,
		mergeableFiles = #(), unMergeableFiles = #())
		{
		new this(controller, filename, attachFileName, merge_pdf?,
			data, mergeableFiles, unMergeableFiles)
		}

	New(.controller, .filename, .attachFileName, .merge_pdf?, data,
		mergeableFiles = #(), unMergeableFiles = #())
		{
		if '' is bucket = AttachmentS3Bucket()
			return

		data.attachFileName = .attachFileName
		data.origFileName = .filename
		if .fromLocal?(.filename)
			file = .tmpFileUrl(.filename, .attachFileName)
		else
			file = AmazonS3.PresignUrl(bucket, FormatAttachmentPath(.filename),
				download?:)
		preSignedAttachments = Object()
		historyLinks = Object()
		tmpName = .tmpName(.attachFileName)
		historyLinks[.attachFileName] = Object(:tmpName,
			url: AmazonS3.PresignUrl(bucket, tmpName, method: 'PUT'))
		if .merge_pdf?
			preSignedAttachments.Add(.tmpFileUrl(.filename, Paths.Basename(.filename)))
		.collectPresignedUrls(data, historyLinks, bucket, preSignedAttachments)
		data.preSignedAtttachments = preSignedAttachments
		data.historyLinks = historyLinks
		data.merge_pdf? = .merge_pdf?
		if data.merge_pdf?
			{
			tmpFileName = .monthFolder() $
				'PDF_' $ Date().StdShortDate() $ '_' $ UuidString() $ '.pdf'
			data.mergedPdfName = tmpFileName
			data.mergedPdfLink = AmazonS3.PresignUrl(bucket, tmpFileName, method: 'PUT')
			}

		if Object?(mergeableFiles)
			mergeableFiles = mergeableFiles.Map(Paths.ToStd)
		if Object?(unMergeableFiles)
			unMergeableFiles = unMergeableFiles.Map(Paths.ToStd)
		data.mergeableFiles = mergeableFiles
		if data.GetDefault('mergeOnly?', false)
			{
			info = Object(externalCDN: OptContribution('ExternalJsCDN', ''))
			SuRenderBackend().RecordAction(false, 'EmailAttachmentComponent.SendEmail',
				Object(file, data, 'from', 'to', info))
			}
		else
			.send(data, unMergeableFiles, file)
		}

	collectPresignedUrls(data, historyLinks, bucket, preSignedAttachments)
		{
		for m in data.GetDefault('attachments', #()).Members()
			{
			data.attachments[m] = Paths.ToStd(data.attachments[m])
			attachment = data.attachments[m]
			if not attachment.Prefix?('(APPEND) ')
				{
				tmpName = .tmpName(attachment)
				historyLinks[Paths.Basename(attachment)] = Object(:tmpName,
					url: .preSignedUrl(bucket, tmpName, method: 'PUT'))
				}
			attachmentPath = attachment.RemovePrefix('(APPEND) ')
			if .fromLocal?(attachmentPath)
				preSignedAttachments.Add(
					.tmpFileUrl(attachmentPath, Paths.Basename(attachmentPath)))
			else
				preSignedAttachments.Add(.preSignedUrl(bucket, FormatAttachmentPath(
					attachment.RemovePrefix('(APPEND) ')), download?:))
			}
		}

	preSignedUrl(@args)
		{
		AmazonS3.PresignUrl(@args)
		}

	fromLocal?(filename)
		{
		return filename.Prefix?('temp/') or filename.Lower().Has?(GetTempPath().Lower())
		}

	tmpName(attachment)
		{
		name = Paths.Basename(attachment)
		return MakeUniqueFileName(.monthFolder(), name, {|unused| true }).dest
		}

	monthFolder()
		{
		return Date().Format('yyyyMM') $ '/'
		}

	send(data, unMergeableFiles, file)
		{
		data.unMergeableFiles = unMergeableFiles
		info = BookEmailInfo().Copy()
		from = SuBookSendEmail.CleanupDisplayName(data.email_from)
		to = SuBookSendEmail.FormatAddressesForSend(data.to)
		info.sesFrom = AmazonSES.SourceEmail(from)
		info.service = NetworkService.IPAddress()
		info.serviceOther = NetworkService.OtherIP(info.service)
		info.authentication = SoleContribution('SESAuthentication')()
		info.externalCDN = OptContribution('ExternalJsCDN', '')
		.controller.Act('SendEmail', file, data, from, to, info)
		}

	tmpFileUrl(filename, attachFileName)
		{
		return 'download' $ Url.BuildQuery(Object(Base64.Encode(
			filename.Xor(EncryptControlKey())), token: .token())) $
			'&' $ Url.EncodeValues(Object(saveName: attachFileName))
		}

	token()
		{
		return SuRenderBackend().Token
		}
	}
