// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
MimeBase
	{
	FileSizeErrorPrefix: 'MimeMulti: AttachFile: file size ('
	FileNotFoundErrorPrefix: 'MimeMulti: AttachFile: can\'t get: '
	New(subtype = 'mixed')
		{
		super('multipart', subtype)
		.parts = Object()
		}

	ValidateFile(filename, fileContent = false)
		{
		if fileContent is false and not FileExists?(filename)
			throw .FileNotFoundErrorPrefix $ filename
		size = fileContent is false ? FileSize(filename) : fileContent.Size()
		max_size = Email_CreateMIME.MaxSizeInMb().Mb()
		if size > max_size
			throw .FileSizeErrorPrefix $ ReadableSize(size) $ ') ' $
				'exceeds maximum (' $ ReadableSize(max_size) $ ')'
		}

	Type(filename, type = false)
		{
		if type is false
			{
			ext = filename.AfterLast('.')
			type = MimeTypes.GetDefault(ext, 'application/octet-stream')
			}
		return type.Split('/')
		}

	boundaryLen: 20
	boundaryKey: 16
	Boundary()
		{
		s = '-'.Repeat(.boundaryLen)
		for (i = 0; i < .boundaryKey; ++i)
			s $= Random(.boundaryKey).Hex()
		return s
		}
	}
