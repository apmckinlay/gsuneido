// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	client(usageTracker = false)
		{
		cl = OpenAI_Decisions
			{
			OpenAI_https(method, url, content,
				header)
				{
				_decisionRequests.Add(Object(:method, :url, :content, :header))
				return _decisionResponse
				}
			}
		return cl(#token, usageTracker)
		}

	Test_Request()
		{
		_decisionRequests = Object()
		_decisionResponse = Object(header: "HTTP/1.1 200 OK\n",
			content: Json.Encode(#(answers: (), model: "gpt-6-luna", usage: ())))
		questions = #((type: predicate, name: damaged,
				instructions: "Does the customer report a damaged item?"))
		client = .client()
		client.Decisions("The screen is broken.", questions)
		Assert(_decisionRequests isSize: 1)
		request = _decisionRequests[0]
		Assert(request.method is: #POST)
		Assert(request.url is: "https://api.openai.com/v1/decisions")
		Assert(request.header
			is: #(Authorization: "Bearer token",
				"Content-Type": "application/json"))
		Assert(Json.Decode(request.content)
			is: Object(input: "The screen is broken.", :questions, model: "gpt-6-luna"))
		}

	Test_StructuredInputAndOptions()
		{
		_decisionRequests = Object()
		_decisionResponse = Object(header: "HTTP/1.1 200 OK\n",
			content: Json.Encode(#(answers: (), usage: ())))
		input = #(
			(role: user,
				content: (
					(type: input_text, text: "Evaluate this image.",
						prompt_cache_breakpoint: (mode: explicit)),
					(type: input_image, image_url: "https://example.com/image.png"),
					(type: input_json, data: (item: screen))
					)),
			(type: function_call, call_id: call1, name: inspect, arguments: "{}"),
			(type: function_call_output, call_id: call1, output: broken)
			)
		questions = #(
			(type: choice, instructions: "Choose a category.",
				choices: ((value:), (value: "true"))),
			(type: score, instructions: ((type: input_text, text: "Rate damage.")),
				levels: ((label: low), (label: high)))
			)
		extra = #(prompt_cache_key: customer,
			prompt_cache_options: (mode: explicit, ttl: "30m"),
			safety_identifier: user1)
		client = .client()
		client.Decisions(input, questions, model: "gpt-6-luna", :extra)
		Assert(Json.Decode(_decisionRequests[0].content)
			is: Object(:input, :questions, model: "gpt-6-luna").Merge(extra))
		}

	Test_AnswersAndUsage()
		{
		_decisionRequests = Object()
		usage = #(input_tokens: 42,
			input_tokens_details: (cached_tokens: 0, cache_write_tokens: 0),
			output_tokens: 0, output_tokens_details: (reasoning_tokens: 0),
			total_tokens: 42)
		expected = Object(model: "gpt-6-luna", :usage,
			answers: #(
				(type: predicate, name: damaged, probability: 0.95),
				(type: choice, name: category, choice: false, confidence: 0.8,
					probabilities: ((value: false, probability: 0.8),
						(value: "false", probability: 0.2))),
				(type: score, name: severity, score: 0.75, confidence: 0.75,
					probabilities: ((label: low, value: 0, probability: 0.25),
						(label: high, value: 1, probability: 0.75))),
				(type: refusal, name: restricted)
				))
		_decisionResponse = Object(header: "HTTP/1.1 200 OK\n",
			content: Json.Encode(expected))
		trackerCalls = Object()
		client = .client({|value| trackerCalls.Add(value) })
		Assert(client.Decisions(#input, #()) is: expected)
		Assert(trackerCalls is: [usage])
		}

	Test_Errors()
		{
		_decisionRequests = Object()
		trackerCalls = Object()
		client = .client({|usage| trackerCalls.Add(usage) })
		_decisionResponse = Object(header: "HTTP/1.1 400 Bad Request\n",
			content: Json.Encode(#(error: (message: "Invalid question"),
				usage: (input_tokens: 10, total_tokens: 10))))
		Assert({ client.Decisions(#input, #()) }
			throws: "OpenAI Decisions: Invalid question")
		Assert(trackerCalls is: #((input_tokens: 10, total_tokens: 10)))

		_decisionResponse = Object(header: "HTTP/1.1 400 Bad Request\n",
			content: Json.Encode(#(error: (message: "No usage"))))
		Assert({ client.Decisions(#input, #()) }
			throws: "OpenAI Decisions: No usage")
		Assert(trackerCalls isSize: 1)

		_decisionResponse = Object(header: "HTTP/1.1 400 Bad Request\n", content: "")
		Assert({ client.Decisions(#input, #()) }
			throws: "OpenAI Decisions: HTTP/1.1 400 Bad Request")
		Assert(trackerCalls isSize: 1)

		client = .client()
		_decisionResponse = Object(header: "HTTP/1.1 400 Bad Request\n",
			content: Json.Encode(#(error: (message: "Invalid question"), usage: ())))
		Assert({ client.Decisions(#input, #()) }
			throws: "OpenAI Decisions: Invalid question")
		}
	}
