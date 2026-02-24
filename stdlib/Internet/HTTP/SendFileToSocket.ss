// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(socket, fileName, headers, start = "HTTP/1.1 200 OK", delete? = false)
		{
		headers.Content_Length = FileSize(fileName)
		HttpSend(socket, start, headers, '')
		File(fileName)
			{ |f|
			f.CopyTo(socket)
			}
		if delete? and true isnt result = DeleteFile(fileName)
			SuneidoLog('ERROR: (CAUGHT) SendFileToSocket failed deleting file - ' $ result
				params: [:fileName], caughtMsg: 'Programmer may want to clean it up')
		return -1 // response is already handled
		}

	BuildHeaders(fileName)
		{
		return Object(
			Expires: Date().Plus(days: 20).InternetFormat(),
			Content_Type: "application/octetstream",
			Content_Transfer_Encoding: "Binary",
			Content_Disposition: 'attachment; filename="' $ Url.Decode(fileName) $ '"')
		}
	}