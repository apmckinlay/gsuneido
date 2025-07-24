// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
MimeMultiBase
	{
	New(subtype = 'mixed')
		{
		super(subtype)
		.parts = Object()
		.attachFiles = Object()
		}

	Attach(mime)
		{
		.parts.Add([content: mime, type: #mime])
		return this
		}

	AttachFile(filename, type = false, attachFileName = false)
		{
		.ValidateFile(filename)

		cleanupOriginal? = EmailAttachment.CompressedFile?(filename)
		attachName = attachFileName isnt false ? attachFileName : filename
		.attachFiles[attachName] = Object(:filename, :type, :attachName,
			:cleanupOriginal?)
		.parts.Add([content: attachName, type: #file])
		return this
		}

	GetAttachedFiles()
		{
		return .attachFiles
		}

	GetMimeTextMessageContent()
		{
		s = ''
		for part in .parts.Filter({ it.type is #mime })
			s $= part.content.MessageContent()
		return s.Trim()
		}

	finalFile: false
	ToString()
		{
		if .finalFile isnt false and FileExists?(.finalFile)
			return .finalFile
		boundary = .buildBoundary()
		.finalFile = GetAppTempFullFileName(#fim)
		File(.finalFile, #w)
			{ .outputParts(it, boundary) }
		return .finalFile
		}

	buildBoundary()
		{
		boundary = .Boundary()
		.AddExtra('Content-Type', '\r\n\tboundary': boundary)
		return '--' $ boundary
		}

	outputParts(file, boundary)
		{
		file.Write(super.ToString() $ boundary)

		for part in .parts
			{
			file.Write('\r\n')
			if part.type is #file
				.attachFileToFinalFile(file, .attachFiles[part.content])
			else
				.writeMimeToFinalFile(file, part.content)
			file.Write('\r\n' $ boundary)
			}
		file.Write('--\r\n')
		}

	attachFileToFinalFile(file, attach)
		{
		attachmentType = .Type(attach.filename, attach.type)
		name = Paths.Basename(attach.attachName)
		m = MimeBase(attachmentType[0], attachmentType[1])
		m.AddHeader('Content-Disposition', 'attachment', '\r\n\tfilename': name)
		m.Content_Transfer_Encoding('base64')
		m.AddExtra('Content-Type', '\r\n\tname': name)
		file.Write(m.MultiPartToString())

		Base64BufferEncodeFile(attach.filename, file)
		}

	writeMimeToFinalFile(file, mime)
		{
		for line in mime.MultiPartToString().Lines()
			file.Writeline(line)
		}
	}
