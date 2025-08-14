// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_ResponsesWithFunctionCalls()
		{
		_logs = Object(fnCalls: 0, inputs: Object())
		_fakeResponses = #(
			(
				id: 1
				output: (
					(
						type: "function_call",
						name: "test_func",
						arguments: '{"x":1}',
						call_id: "call1"
					),
					(
						type: "function_call",
						name: "test_func",
						arguments: '{"x":2}',
						call_id: "call2"
					)
				)
			),
			(
				id: 2
				output: (
					(
						type: "message",
						content: (
							(type: "output_text", text: "One more call!")
						)
					),
					(
						type: "function_call",
						name: "test_func",
						arguments: '{"x":3}',
						call_id: "call3"
					)
				)
			),
			(
				id: 3
				output: (
					(
						type: "message",
						content: (
							(type: "output_text", text: "Success!")
						)
					)
				)
			),
		)

		cl = OpenAI_Responses
			{
			// Function to replace .Responses
			Responses(input, format/*unused*/, tools/*unused*/, model/*unused*/, extra)
				{
				_logs.inputs.Add(input)
				id = extra.GetDefault('previous_response_id', 0)
				return _fakeResponses[id]
				}
			}

		// Define the function the AI can call
		testFunc = Object(
			desc: "Test function",
			params: Object(x: Object(type: "number", description: "A number")),
			fn: function(x) { _logs.fnCalls++; return x + 1 } // returns x+1
		)
		funcs = Object(test_func: testFunc)

		messages = Object()
		out = { |text, type/*unused*/| messages.Add(text) }

		agent = OpenAI_Agent("dummy-token", out, :funcs)
		agent.OpenAI_Agent_responses_api = cl("dummy-token")

		agent.Send("test prompt")

		// Validate messages and responses
		Assert(messages.Size() is: 2)
		Assert(messages[0] is: "One more call!")
		Assert(messages[1] is: "Success!")

		Assert(_logs.inputs is: #(
			#([role: "user", content: "test prompt"]),
			#([type: "function_call_output", call_id: "call1", output: 2],
				[type: "function_call_output", call_id: "call2", output: 3]),
			#([type: "function_call_output", call_id: "call3", output: 4])))
		Assert(_logs.fnCalls is: 3)
		}

	}