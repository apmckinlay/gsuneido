// Copyright (C) 2023 Axon Development Corporation All rights reserved worldwide.
class
	{
	New(.token)
		{
		}

	Responses(input, format = false, tools = false, model = 'gpt-4.1')
		{
		body = Object(:input, :model)
		if tools isnt false
			body.tools = tools
		if format isnt false
			body.text = Object(:format)
		res = .send('https://api.openai.com/v1/responses', body)
		return res
		}

	// Each input must not exceed 8192 tokens in length.
	MaxEmbeddingInputSize: 32000 // = 8000 tokens * 4

	// input can be a string or an string of array
	Embeddings(input, model = 'text-embedding-ada-002', user = 'axon')
		{
		if String?(input)
			input = Object(input)
		res = .send('https://api.openai.com/v1/embeddings', Object(:input, :model, :user))
		return res.data.Map({ it.embedding })
		}

	send(url, body)
		{
		header = Object(
			'Authorization': 'Bearer ' $ .token,
			'Content-Type': 'application/json')
		content = Json.Encode(body).ToUtf8()
		response = Https('POST', url, content, :header)
		if Https.ResponseCode(response.header) !~ `^2\d\d$`
			.handleError(response, url.AfterLast('/').Capitalize())
		return Json.Decode(response.GetDefault(#content, ''), handleNull: 'skip')
		}

	handleError(response, op)
		{
		ob = Object()
		try ob = Json.Decode(response.GetDefault(#content, ''), handleNull: 'empty')
		err = ob.GetDefault(#error, [message: response.header.BeforeFirst('\n')]).message
		if err.Lower() isnt "bad gateway"
			throw 'OpenAI ' $ op $ ': ' $ err
		}
	}
