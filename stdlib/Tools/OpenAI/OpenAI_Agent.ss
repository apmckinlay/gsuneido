// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
class
	{
	sendInstractions?: true
	New(token, .output, .instructions = false, .funcs = false, .format = false,
		tools = false, usageTracker = false, .toolTracker = false)
		{
		.responses_api = OpenAI_Responses(token, :usageTracker)

		.origInstructions = .instructions
		.tools = .prepareFunc(tools, funcs)
		}

	prepareFunc(tools, funcs)
		{
		if tools is false
			tools = Object()

		for name in funcs.Members()
			{
			func = funcs[name]
			tools.Add(OpenAI_Responses.DefineFunction(name, func.desc, func.params))
			}

		return tools
		}

	prevId: false
	Send(input, model = 'gpt-4.1', extra = #())
		{
		if String?(input)
			input = [[role: 'user', content: input]]

		origExtra = extra
		continue? = true
		while continue?
			{
			continue? = false

			extra = origExtra.Copy()
			if .sendInstractions?
				{
				extra['instructions'] = .instructions
				.sendInstractions? = false
				}

			if .prevId isnt false
				extra['previous_response_id'] = .prevId

			response = .responses_api.Responses(input, .format, .tools, :extra, :model)
			.prevId = response.id
			input = Object()

			if Object?(response.output)
				{
				for output in response.output
					if .processOutput(output, .funcs, input) is true
						continue? = true
				}
			}
		}

	processOutput(output, funcs, input)
		{
		switch (output.type)
			{
		case 'message':
			for outputText in output.content
				if outputText.type is 'output_text'
					(.output)(outputText.text, type: #message)
				else // refusal
					(.output)(outputText.text, type: #refusal)
			return false
		case 'function_call':
			args = Json.Decode(output.arguments)
			if .toolTracker isnt false
				(.toolTracker)(output.type, name: output.name, :args)
			result = (funcs[output.name].fn)(@args)
			input.Add([type: 'function_call_output',
				call_id: output.call_id,
				output: result])
			return true
		default:
			if .toolTracker isnt false
				(.toolTracker)(output.type, :output)
			return false
			}
		}

	Clear()
		{
		.prevId = false
		.sendInstractions? = true
		}
	}