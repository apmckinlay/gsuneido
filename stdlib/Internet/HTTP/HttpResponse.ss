// Copyright (C) 2008 Suneido Software Corp. All rights reserved worldwide.
class
	{
	New()
		{
		.ResponseCode(200) /*= OK */
		.headers = Object()
		}
	fields: ('Date', 'Cache-Control', 'Last-Modified', 'Content-Type',
		'Content-Disposition', 'Expires', 'ETag', 'Location', 'Server',
		'Content-Encoding', 'Connection', 'Access-Control-Allow-Origin',
		'Access-Control-Allow-Methods', 'Access-Control-Allow-Headers', 'Set-Cookie')

	ResponseHeaderField(field, value, headers)
		{
		field = field.Tr('_', '-')
		if not .fields.Has?(field)
			throw "HttpResponse: method not found: " $ field
		if value isnt false
			{
			if Date?(value)
				value = value.InternetFormat()
			headers[field] = value
			}
		return headers[field]
		}

	Default(field, value = false)
		{
		return .ResponseHeaderField(field, value, .headers)
		}

	ResponseCode(code = false)
		{
		if code isnt false
			.response_code = String?(code) ? code : code $ ' ' $ HttpResponseCodes[code]
		return .response_code
		}

	Getter_Headers()
		{
		return .headers
		}

	WriteHeader(sc)
		{
		// use HTTP/1.0 so one request per connnection
		header = Object()
		header.Add("HTTP/1.0 " $ .ResponseCode())
		for member in .headers.Members().Sort!()
			header.Add(member $ ': ' $ .headers[member])
		sc.Writeline(header.Join('\r\n'))
		}
	}