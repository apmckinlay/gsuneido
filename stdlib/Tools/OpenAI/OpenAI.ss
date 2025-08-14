// Copyright (C) 2023 Axon Development Corporation All rights reserved worldwide.
class
	{
	BaseUrl: 'https://api.openai.com/v1/'
	New(.token)
		{
		}

	Send(url, body, method = 'POST')
		{
		header = Object(
			'Authorization': .AuthKey(),
			'Content-Type': 'application/json')
		content = body is false ? '' : Json.Encode(body).ToUtf8()
		response = .https(method, url, content, :header)
		return .HandleResponse(response, url)
		}

	// for test
	https(@args)
		{
		return Https(@args)
		}

	AuthKey()
		{
		return 'Bearer ' $ .token
		}

	HandleResponse(response, url)
		{
		if Https.ResponseCode(response.header) !~ `^2\d\d$`
			.handleError(response, url.AfterLast('/').Capitalize())
		return Json.Decode(response.GetDefault(#content, ''), handleNull: 'skip')
		}

	handleError(response, op)
		{
		ob = Object()
		try
			{
			ob = Json.Decode(response.GetDefault(#content, ''), handleNull: 'empty')
			.ProcessOpenAIError(ob)
			}
		throw 'OpenAI ' $ op $ ': ' $ ob.GetDefault(#error,
			{ [message: response.header.BeforeFirst('\n')] }).message
		}

	ProcessOpenAIError(@unused) { }
	}
