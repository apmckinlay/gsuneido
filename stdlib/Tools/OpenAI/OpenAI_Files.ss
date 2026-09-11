// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
OpenAI
	{
	Upload(file, purpose = #user_data, fileName = "")
		{
		file, rename = .renameFile(file, copy?:)
		fileName, unused = .renameFile(fileName)

		url = .BaseUrl $ "files"
		header = Object(Authorization: .AuthKey())
		options = Object(files: [
			"purpose=" $ purpose,
			"file=@" $ file $ Opt(";filename=", fileName)])
		response = Https.PostWithOptions(#POSTFILES, url, :header, :options)
		if rename isnt ""
			DeleteFile(rename)
		return .HandleResponse(response, url).id
		}

	renameFile(file, copy? = false)
		{
		rename = ""
		if file.Has?(',') or file.Has?(';')
			{
			rename = file.Tr(",;")
			if copy?
				CopyFile(file, rename, false)
			file = rename
			}
		return file, rename
		}

	List(purpose = false, after = false, limit = 200)
		{
		params = Object(:limit)
		if purpose isnt false // Only return files with the given purpose
			params.purpose = purpose

		if after isnt false // Pagination: An object ID that defines our place in the list
			params.after = after

		url = Url.Encode(.BaseUrl $ "files", params)
		return .Send(url, body: false, method: #GET)
		}

	Retrieve(fileId)
		{
		return .Send(.BaseUrl $ "files/" $ fileId, body: false, method: #GET)
		}

	DeleteFile(fileId)
		{
		return .Send(.BaseUrl $ "files/" $ fileId, body: false, method: #DELETE)
		}
	}
