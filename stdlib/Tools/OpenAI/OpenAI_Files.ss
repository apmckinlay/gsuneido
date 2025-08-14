// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
OpenAI
	{
	Upload(file, purpose = 'user_data')
		{
		url = .BaseUrl $ 'files'
		header = Object('Authorization': .AuthKey())
		options = Object(files: Object(
			'purpose=' $ purpose,
			'file=@' $ file))
		response = Https.PostWithOptions('POSTFILES', url, :header , :options)
		return .HandleResponse(response, url).id
		}

	List(purpose = false, limit = 200)
		{
		params = Object(:limit)
		if purpose isnt false
			params.purpose = purpose

		url = Url.Encode(.BaseUrl $ "files", params)
		return .Send(url, body: false, method: "GET")
		}

	Retrieve(fileId)
		{
		return .Send(.BaseUrl $ "files/" $ fileId, body: false, method: 'GET')
		}

	DeleteFile(fileId)
		{
		return .Send(.BaseUrl $ "files/" $ fileId, body: false, method: 'DELETE')
		}
	}