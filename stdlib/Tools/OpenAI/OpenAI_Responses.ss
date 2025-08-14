// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
OpenAI
	{
	New(token, .usageTracker = false)
		{
		super(token)
		}

	Responses(input, format = false, tools = false, model = 'gpt-4.1', extra = #())
		{
		body = Object(:input, :model).Merge(extra)
		if tools isnt false
			body.tools = tools
		if format isnt false
			body.text = Object(:format)
		res = .Send(.BaseUrl $ 'responses' body)
		if .usageTracker isnt false
			(.usageTracker)(res['usage'])
		return res
		}

	ProcessOpenAIError(res)
		{
		if res.Member?(#usage) and .usageTracker isnt false
			(.usageTracker)(res['usage'])
		}

	// static
	/*
	name: string
	desc: string
	params: [
		paramName: [
			type: string,
			description: string,
			required (optional): boolean ]]
	*/
	DefineFunction(name, desc, params)
		{
		parameters = Object(type: 'object', properties: Object(),
			required: Object(), additionalProperties: false)
		for param in params.Members()
			{
			if params[param].GetDefault(#required, false)
				parameters.required.Add(param)
			parameters.properties[param] = params[param].Copy().Delete(#required)
			}
		return Object(type: 'function',
			:name,
			:parameters,
			description: desc,
			strict:)
		}
	}