// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
MimeMultiBase
	{
	New(subtype = 'mixed')
		{
		super(subtype)
		.parts = Object()
		}
	Attach(mime)
		{
		.parts.Add(mime)
		return this
		}

	AttachFile(filename, type = false, attachFileName = false, fileContent = false)
		{
		type = .Type(filename, type)
		.ValidateFile(filename, :fileContent)
		s = fileContent is false ? GetFile(filename) : fileContent
		if type[0] is 'text'
			m = MimeText(s, type[1])
		else
			{
			m = MimeBase(type[0], type[1])
			m.SetPayload(s).Base64()
			}
		if attachFileName is false
			attachFileName = filename
		args = Object('Content-Type')
		nameMember = '\r\n\tname'
		args[nameMember] = Paths.Basename(attachFileName)
		m.AddExtra(@args)
		args = Object('Content-Disposition', 'attachment')
		nameMember = '\r\n\tfilename'
		args[nameMember] = Paths.Basename(attachFileName)
		m.AddHeader(@args)
		.Attach(m)
		}

	ToString(skipFiles = false)
		{
		boundary = .Boundary()
		args = Object('Content-Type')
		bMember = '\r\n\tboundary'
		args[bMember] = boundary
		.AddExtra(@args)
		boundary = '--' $ boundary $ '\r\n'
		s = super.ToString() $ boundary
		parts = skipFiles ? [.parts[0]] : .parts
		for p in parts
			s $= p.MultiPartToString() $ '\r\n' $ boundary
		s = s[.. -2] $ '--\r\n'
		return s
		}
	}
