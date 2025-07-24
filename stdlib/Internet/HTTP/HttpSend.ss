// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(sc, start, header, content)
		{
		if not String?(content)
			ProgrammerError("HttpSend content is not a string")
		data = Object()
		data.Add(start)

		.handleHeaders(data, header, content)

		data.Add("") // blank line between header and content
		data.Add(content)
		// more efficient to use one Join than multiple concatenations
		sc.Write(data.Join('\r\n'))
		}

	handleHeaders(data, header, content)
		{
		if not header.Member?(#Date)
			data.Add('Date: ' $ Date().InternetFormat())

		contentLength? = header.Member?('Content_Length') or
			header.Member?('Content-Length')
		if content isnt '' and contentLength?
			ProgrammerError("HttpSend Content-Length is only allowed if content is empty")
		if not contentLength?
			data.Add('Content-Length: ' $ content.Size())

		for fld in header.Members().Sort!()
			{
			field = fld.Tr('_', '-')
			if Object?(header[fld])
				data.Append(header[fld].Map({ field $ ': ' $ it}))
			else
				data.Add(field $ ': ' $ header[fld])
			}
		}
	}