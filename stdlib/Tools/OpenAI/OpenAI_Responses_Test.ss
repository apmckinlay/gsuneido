// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_DefineFunction()
		{
		name = "add"
		desc = "Adds two numbers"
		params = Object(
			a: Object(type: "integer", description: "First number", required: true),
			b: Object(type: "integer", description: "Second number")
		)

		res = OpenAI_Responses.DefineFunction(name, desc, params)
		Assert(res is: #(
			type: "function",
			name: "add",
			description: "Adds two numbers",
			strict:,
			parameters: #(type: "object",
				properties: #(
					a: #(type: "integer", description: "First number"),
					b: #(type: "integer", description: "Second number")),
				required: #('a'),
				additionalProperties: false)))
		}

	Test_OpenAI_Responses_UsageChecker()
		{
		trackerCalls = Object()
		usageTracker = { |usage| trackerCalls.Add(usage) }
		cl = OpenAI_Responses
			{
			OpenAI_https(@unused)
				{
				return _fakeHttpResponse
				}
			}

		responses = cl("token", usageTracker)

		_fakeHttpResponse = Object(
			header: "HTTP/1.1 200 OK\n",
			content: Json.Encode(#(usage: (prompt_tokens: 10, completion_tokens: 20)))
			)

		responses.Responses("test input")
		Assert(trackerCalls isSize: 1)
		Assert(trackerCalls[0] is: #(prompt_tokens: 10, completion_tokens: 20))

		_fakeHttpResponse = Object(
			header: "HTTP/1.1 400 Bad Request\n",
			content: Json.Encode(#(
				error: (message: 'OpenAI error'),
				usage: (prompt_tokens: 20, completion_tokens: 40)))
			)

		Assert({ responses.Responses("test input") }
			throws: 'OpenAI Responses: OpenAI error')
		Assert(trackerCalls isSize: 2)
		Assert(trackerCalls[1] is: #(prompt_tokens: 20, completion_tokens: 40))

		_fakeHttpResponse = Object(
			header: "HTTP/1.1 400 Bad Request\n",
			content: '')
		Assert({ responses.Responses("test input") }
			throws: 'OpenAI Responses: HTTP/1.1 400 Bad Request')
		Assert(trackerCalls isSize: 2)
		}
	}