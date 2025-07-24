// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
function (socket, fileName, headers, start = "HTTP/1.1 200 OK", delete? = false)
	{
	headers.Content_Length = FileSize(fileName)
	HttpSend(socket, start, headers, '')
	File(fileName)
		{ |f|
		f.CopyTo(socket)
		}
	if delete? and true isnt result = DeleteFile(fileName)
		SuneidoLog('ERROR: (CAUGHT) SendFileToSocket failed deleting file - ' $ result,
			params: Object(:fileName), caughtMsg: 'Programmer may want to clean it up')
	return -1 // response is already handled
	}