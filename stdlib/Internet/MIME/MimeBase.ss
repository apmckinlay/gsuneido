// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
class
	{
	New(.maintype = 'text', .subtype = 'plain')
		{
		.hdr = Object().Set_default(false)
		.extra = Object().Set_default('')
		.Fields = .fields.Copy()
		}
	payload: ''
	SetPayload(.payload)
		{
		return this
		}
	fields: ('From', 'To', 'Subject', 'Date', 'Message-ID', 'Sender', 'Return-Path',
		'Content-Transfer-Encoding', 'Reply-To')
	Default(@args)
		{
		field = args[0].Tr('_', '-')
		if .Fields.Has?(field)
			{
			value = ''
			if args.Size() is 1
				{
				if field is 'Date'
					value = Date()
				else if field is 'Message-ID'
					value = .message_id()
				}
			else
				value = args[1].ChangeEol('\r\n\t').Trim()
			.AddHeader(field, value)
			}
		else
			throw "MimeBase: method not found: " $ args[0]
		return this
		}
	AddHeader(@args)
		{
		name = args[0]
		value = args[1]
		.Fields.AddUnique(name)
		.hdr[name] = Date?(value) ? .date(value) : value
		if args.Size() is 3 /*= extra header value*/
			{
			m = args.Members()[2]
			.extra[name] = '; ' $ m $ '=' $ Display(args[m])
			}
		return this
		}
	AddExtra(@args)
		{
		name = args.Members()[1]
		value = args[name]
		.extra[args[0]] = '; ' $ name $ '=' $ Display(value)
		return this
		}
	date(date)
		{
		return date.Format("ddd, d MMM yyyy HH:mm:ss") $
			' -' $ (date.GetLocalGMTBias() / 60).Pad(2) $ '00' /*= in hours*/
		}
	message_id()
		{
		return '<' $ UuidString().Tr('-', '.') $ '@' $
			LastContribution('MimeBase_MessageIdDomain') $ '>'
		}
	encode: function (s) { s }
	Base64()
		{
		.AddHeader('Content-Transfer-Encoding', 'base64')
		.encode = Base64.EncodeLines
		return this
		}
	MimeVersion: 'MIME-Version: 1.0\r\n'
	ToString()
		{
		return .MimeVersion $ .ExtraHeader(.Fields) $ .ContentType() $ .MessageContent()
		}
	MultiPartToString()
		{
		return .ContentType() $ .ExtraHeader(.Fields) $ .MessageContent()
		}
	ContentType()
		{
		return 'Content-Type: ' $ .maintype $ '/' $ .subtype $
			.extra['Content-Type'] $ '\r\n'
		}
	ExtraHeader(fields)
		{
		s = ''
		for f in fields
			if .hdr.Member?(f)
				s $= f $ ': ' $ .hdr[f] $ .extra[f] $ '\r\n'
		return s
		}
	GetHeaderValue(field)
		{
		return .hdr[field]
		}
	MessageContent()
		{
		s = '\r\n'
		s $= (.encode)(.payload)
		if s[-2 ..] isnt '\r\n'
			s $= '\r\n'
		return s
		}
	}
