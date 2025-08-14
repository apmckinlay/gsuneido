// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
HtmlDivComponent
	{
	SendEmail(filename, data, from, to, info)
		{
		msg = data.GetDefault('mergeOnly?', false)
			? 'Downloading pdf'
			: 'Sending email...'
		SuRender().Overlay.Show('emailAttachment', msg)
		data.info = info
		data.from = from
		data.to = to
		SuUI.GetCurrentWindow().Eval(OpenFileNameComponent.LoadMagickScript(
			info.externalCDN))
		extraAttachments = data.GetDefault('preSignedAtttachments', #()).Copy().
			Add(filename, at: 0)
		list = Object()
		for attach in extraAttachments
			.download(attach, list, data)
		}

	download(attach, list, data)
		{
		path = .getFileName(attach)
		name =  Paths.Basename(path)
		item = Object(fileContent: '', compressed: '', :name, :path, url: attach)
		list.Add(item)
		.downloadAttch(item, list, data)
		}

	downloadAttch(item, list, data)
		{
		xhr = SuUI.MakeWebObject('XMLHttpRequest')
		item.xhr = xhr
		item.type = item.name.Lower().AfterLast('.')
		xhr.AddEventListener('readystatechange', { |event/*unused*/|
			if xhr.readyState is 4/*=DONE*/
				{
				if xhr.status is HttpResponseCodes.OK
					{
					item.fileContent = xhr.response
					if list.Every?({ it.fileContent isnt '' })
						.allDownloaded(list, data)
					}
				else if String(xhr.status)[0] in ('4','5') or xhr.status is 0
					.alertErr("Cannot find file: " $ item.name)
				}
			})
		xhr.Open('GET', item.url)
		xhr.SetRequestHeader('Content-Type',
			MimeTypes.GetDefault(item.type, 'application/octet-stream'))
		xhr.responseType = "arraybuffer"
		xhr.Send()
		}

	allDownloaded(list, data)
		{
		attachmentsSize = 0
		list.Each()
			{
			attachmentsSize += it.fileContent.byteLength
			sb = SuUI.GetCurrentWindow().Uint8Array(it.fileContent)
			it.fileData = SuUI.GetCurrentWindow().ArrayToString(sb)
			}
		data.listAttachments = list
		.loadLibs({ .mergePdfs(data, list, attachmentsSize) })
		}

	loadLibs(block)
		{
		SuUI.GetCurrentWindow().LoadMagick(
			{
			SuUI.GetCurrentWindow().LoadPako(block,
				{ |err| Print(LoadPakoError: err); throw err })
			},
			{ |err| Print(loadImageMagickError: err); throw err })
		}

	mergePdfs(data, list, attachmentsSize)
		{
		if not data.merge_pdf?
			{
			.compressAndEmail(attachmentsSize, data, list)
			return
			}

		mergeableFiles = Object()
		for f in data.mergeableFiles
			mergeableFiles.Add(list.FindOne({ f.Suffix?(it.path) }))
		mergedFile = Object()
		PdfMerger(data.mergeableFiles, mergedFile,
			compress: data.GetDefault('compress', false)
			filesData: mergeableFiles, afterMergedAsync: { |invalidFiles|
				if not invalidFiles.Empty?()
					.alertErr(PdfMerger.InvalidFilesMsg(invalidFiles))
				else
					{
					list[0].fileData = mergedFile.fileData
					for f in data.mergeableFiles
						{
						m = list.FindIf({ f.Suffix?(it.path) })
						if m not in (0, false)
							list.Delete(m)
						}
					.compressAndEmail(attachmentsSize, data, list)
					}
			})
		}

	compressAndEmail(attachmentsSize, data, list)
		{
		if data.GetDefault('mergeOnly?', false)
			{
			.CloseOverlay()
			SuUI.GetCurrentWindow().DownloadFile(
				data.attachFileName, data.listAttachments[0].fileData)
			return
			}

		compress? = attachmentsSize > Email_CreateMIME.MaxSizeInMb().Mb()
		if compress? is false
			{
			data.listAttachments.Each({ it.compressed = it.fileData })
			.sendEmail(data)
			return
			}
		.totalSize = 0
		.compressOne(list, 0, data)
		}

	compressOne(list, compressIdx, data)
		{
		item = list[compressIdx]
		if item.type is 'pdf'
			{
			.compressIfPDF(item, data, compressIdx, list)
			return
			}

		if false isnt options = ImageMagick.GetOption(item.type)
			{
			.compressIfImage(options, item, list, compressIdx, data)
			return
			}

		.compressOther(item, compressIdx, list, data)
		}

	compressIfPDF(item, data, compressIdx, list)
		{
		compressedFile = Object()
		PdfMerger(Object(item.path), compressedFile, compress:,
			maxCompressedFileSizeInMb: Email_CreateMIME.MaxSizeInMb(),
			filesData: Object(item)
			afterMergedAsync: { |invalidFiles|
				if not .compressionValid(invalidFiles)
					{
					if invalidFiles[0].Has?("compressed file size over maximum") and
						data.merge_pdf? is true
						.sendMergedPdfAsLink(data, item)
					else
						.emailAsLinks(data)
					}
				else
					{
					item.compressed = compressedFile.Empty?()
						? item.fileData
						: compressedFile.fileData
					.compressNext(item, compressIdx, list, data)
					}
			})
		}

	compressionValid(invalidFiles)
		{
		return invalidFiles.Empty?() or invalidFiles[0].Has?('nothing compressible')
		}

	sendMergedPdfAsLink(data, item)
		{
		xhr = SuUI.MakeWebObject('XMLHttpRequest')
		xhr.AddEventListener('readystatechange', { |event/*unused*/|
			if xhr.readyState is 4/*=DONE*/
				{
				if xhr.status is HttpResponseCodes.OK
					{
					data.attachments.RemoveIf({ it.Prefix?('(APPEND) ') })
					data.origFileName = data.mergedPdfName
					.emailAsLinks(data)
					}
				// 0: cors preflight request failed, like permission or credential expired
				else if String(xhr.status)[0] in ('4','5') or xhr.status is 0
					.alertErr("There was a problem sending the e-mail. " $
						"Please try again later.")
				}
			})
		xhr.Open('PUT', data.mergedPdfLink)
		xhr.SetRequestHeader('Content-Type', MimeTypes.pdf)

		arrary = SuUI.GetCurrentWindow().Uint8Array(item.fileData)
		blob = SuUI.MakeWebObject('Blob', [arrary], [type: 'application/pdf'])
		file = SuUI.GetCurrentWindow().File(
			[:blob, name: data.mergedPdfName], MimeTypes.pdf)
		xhr.Send(file)
		}

	compressIfImage(options, item, list, compressIdx, data)
		{
		type = item.type
		cmd = ImageMagick.BuildCompressionCmd(type, options)
		sourceBytes = SuUI.GetCurrentWindow().Uint8Array(item.fileContent)
		files = Object(Object(name: 'src.' $ type, content: sourceBytes))
		SuUI.GetCurrentWindow().Magick(files, cmd).Then(
			{|result|
			if result.exitCode is 0
				{
				.imageCompressed(list, compressIdx, result, item, data)
				}
			else
				{
				item.compressed = item.fileData
				.compressNext(item, compressIdx, list, data)
				}
			}).Catch(
			{ |err|
			Print(compressError: err)
			.emailAsLinks(data)
			})
		}

	// does not actually compress, just named compress for consistancy
	compressOther(item, compressIdx, list, data)
		{
		// Pretend we compressed the file.
		if item.Member?('fileData')
			{
			item.compressed = item.fileData
			}
		else
			{
			sourceBytes = SuUI.GetCurrentWindow().Uint8Array(item.fileContent)
			fileContent = SuUI.GetCurrentWindow().ArrayToString(sourceBytes)
			item.compressed = fileContent
			}

		.compressNext(item, compressIdx, list, data)
		}

	imageCompressed(list, compressIdx, result, item, data)
		{
		output = result.outputFiles[0]
		resultFile = SuUI.GetCurrentWindow().File(output,
			MimeTypes.GetDefault(item.type, 'application/octet-stream'))
		resultFile.ArrayBuffer().Then({|arrayBuffer|
			sourceBytes = SuUI.GetCurrentWindow().Uint8Array(arrayBuffer)
			content = SuUI.GetCurrentWindow().ArrayToString(sourceBytes)
			item.compressed = content
			.compressNext(item, compressIdx, list, data)
			})
		}

	compressNext(item, compressIdx, list, data)
		{
		// track total size of the attachments, if > 7MB, send as links
		if ((.totalSize += item.compressed.Size()) > Email_CreateMIME.MaxSizeInMb().Mb())
			{
			if data.merge_pdf? is true
				.sendMergedPdfAsLink(data, list[0])
			else
				.emailAsLinks(data)
			}
		else
			{
			if compressIdx isnt list.Size() - 1 // process next file
				.compressOne(list, compressIdx + 1, data)
			else // if at end of list, send
				.sendEmailIfAllCompressed(data)
			}
		}

	emailAsLinks(data)
		{
		.CloseOverlay()
		newData = data.Copy()
		newData.Delete('listAttachments')
		newData.Delete('mergeableFiles')
		newData.Delete('preSignedAtttachments')
		newData.Delete('historyLinks')
		if newData.GetDefault('emailSent', false)
			return
		data.emailSent = true
		SuRender().Event(false, 'EmailAttachment.EmailAsLinks',
			Object(newData, data.origFileName, data.listAttachments[0].name, 0))
		}

	sendEmailIfAllCompressed(data)
		{
		if data.listAttachments.Every?({ it.compressed isnt '' })
			.sendEmail(data)
		}

	getFileName(url)
		{
		if url.Prefix?('download?')
			return SuUI.GetCurrentWindow().Eval('decodeURIComponent("' $
				url.AfterLast('&saveName=') $ '")')
		return SuUI.GetCurrentWindow().Eval('decodeURIComponent("' $
			Url.Split(url).basepath.AfterFirst('/').AfterFirst('/') $ '")')
		}

	sendEmail(data)
		{
		filename = data.listAttachments[0].name
		_sendError = Object()

		msg = SuBookSendEmail.CreateMime(data.subject, data.message,
			data.listAttachments[0].compressed, filename)
		for file in data.listAttachments[1..]
			msg.AttachFile(file.name, fileContent: file.compressed)
		if not _sendError.Empty?()
			{
			err = 'Please fix the following issues and try again\n'
			err $= _sendError.Join('\n')
			.alertErr(err)
			SuRender().Event(false, 'BookLog' Object(
				'Emailing Attachment(s) failed: ' $ err))
			return
			}

		if msg is false
			{
			.CloseOverlay()
			return
			}

		SuBookSendEmail(data, msg, history: .updateHistory)
		}

	updateHistory(data, mime)
		{
		uploadHistory = Object()
		for file in data.listAttachments
			.sendHistory(file.name, file.compressed, uploadHistory, data, mime)
		}

	sendHistory(fileName, compressed, uploadHistory, data, mime, retry? = false)
		{
		xhr = SuUI.MakeWebObject('XMLHttpRequest')
		uploadHistory[fileName] = false
		xhr.AddEventListener('readystatechange', { |event/*unused*/|
			if xhr.readyState is 4/*=DONE*/
				{
				tmpName = data.historyLinks[fileName].tmpName
				if xhr.status is HttpResponseCodes.OK
					uploadHistory[fileName] = tmpName
				else if String(xhr.status)[0] in ('4','5') or xhr.status is 0
					{
					if retry?
						uploadHistory[fileName] =
							tmpName $ '_' $ xhr.status $ '_att_failed'
					else
						{
						sendHistory = .sendHistory
						SuUI.GetCurrentWindow().SetTimeout({
							sendHistory(fileName, compressed,
								uploadHistory, data, mime, retry?:)
							}, 100) /*= delay */
						}
					}
				if not uploadHistory.Has?(false)
					.saveEmailLog(uploadHistory, data, mime)
			}})
		xhr.Open('PUT', data.historyLinks[fileName].url)
		ext = fileName.AfterLast('.').Lower()
		mimeType = MimeTypes.GetDefault(ext, 'application/octet-stream')
		xhr.SetRequestHeader('Content-Type', mimeType)
		arrary = SuUI.GetCurrentWindow().Uint8Array(compressed)
		blob = SuUI.MakeWebObject('Blob', [arrary], [type: mimeType])
		file = SuUI.GetCurrentWindow().File([:blob, name: fileName], mimeType)
		xhr.Send(file)
		}

	saveEmailLog(uploadHistory, data, mime)
		{
		extraAttach = uploadHistory.Values()
		failed = extraAttach.Filter({ it.Suffix?('_att_failed') })
		if not failed.Empty?()
			SuneidoLog('ERROR: (CAUGHT) Failed to upload attachments for email history, '$
				'the email was sent successfully', params: uploadHistory)
		succeeded = extraAttach.Filter({ not it.Suffix?('_att_failed') })
		pdfNames = [uploaded?:, extraAttach: succeeded]
		SuRender().Event(false, 'LogEmailString', [data.from, data.to,
			mime.ToString(skipFiles:), :pdfNames, quiet?:])
		}

	alertErr(msg)
		{
		.CloseOverlay()
		SuRender().Event(false, 'Alert',
			Object(msg, 'Email Attachment: Not Sent', flags: MB.ICONERROR))
		}

	CloseOverlay()
		{
		SuRender().Overlay.Close('emailAttachment')
		}
	}