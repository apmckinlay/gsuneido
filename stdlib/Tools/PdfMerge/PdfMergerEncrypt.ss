// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
// BuiltDate > 20260303
class
	{
	Setup(userPassword, ownerPassword)
		{
		if userPassword is false
			return false
		Assert(userPassword isString:)
		Assert(ownerPassword isString:)
		// using default industry standard permission: -1028
		return PdfEncrypt.KeyEntries(userPassword, ownerPassword)
		}

	EncryptObj(obj, reader, f, encrypt, cleanupFiles)
		{
		if encrypt is false
			return
		.EncryptHead(obj, encrypt)
		if not obj.head.Trim().Suffix?('stream') or not
			obj.tail.Trim().Prefix?('endstream')
			return
		obj.head $= '\n'
		tmp = GetAppTempPath() $ Display(Timestamp()).Tr('#.')
		cleanupFiles.Add(tmp)
		reader.ExtractStreamToJPG(f, obj, tmp)
		streamData = GetFile(tmp)
		encrypted = PdfEncrypt.Encrypt(
			streamData.ToHex(), encrypt.encryptionKey).FromHex()
		PutFile(tmp, encrypted)
		obj.streamFile = tmp
		obj.streamSize = encrypted.Size()
		obj.Delete('streamStart')
		obj.Delete('streamEnd')
		obj.tail = '\n' $ obj.tail
		.updateLength(obj)
		}

	updateLength(obj)
		{
		headLength = obj.head.FindRx('/Length[' $ .ws $ ']')
		remaining = obj.head[headLength + 1 ..]
		remaining = remaining.Has?('/')
			? '/' $ remaining.AfterFirst('/')
			: '>' $ remaining.AfterFirst('>')
		obj.head = obj.head[..headLength] $ '/Length ' $ obj.streamSize $ remaining
		}

	ws: '\x00\x09\x0A\x0C\x0D\x20' // see pdf reference 1.7 > table 3.1
	EncryptHead(obj, encrypt)
		{
		if encrypt is false
			return obj
		obj.head = obj.head.Replace('[^<]<[[:xdigit:]]*?>',
			{
			it[0] $ .encrytStr(it[2..-1], encrypt)
			}).Replace(`(?q)\(`, `__SU_LEFT_PAREN__`).
			   Replace(`(?q)\)`, `__SU_RIGHT_PAREN__`).Replace('\(.*?\)',
			{
			s = it[1..-1].Replace(`__SU_LEFT_PAREN__`, `\(`).
				Replace(`__SU_RIGHT_PAREN__`, `\)`)
			.encrytStr(s.ToHex(), encrypt)
			}).Replace(`__SU_LEFT_PAREN__`, `\(`).Replace(`__SU_RIGHT_PAREN__`, `\)`)
		return obj
		}

	encrytStr(s, encrypt)
		{
		return '<' $ PdfEncrypt.Encrypt(s, encrypt.encryptionKey) $ '>'
		}
	}
