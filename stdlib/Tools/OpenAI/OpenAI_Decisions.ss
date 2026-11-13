// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
OpenAI
	{
	New(token, .usageTracker = false)
		{
		super(token)
		}

	/*
	input: text or structured messages/tool calls and outputs.
	questions: array of 1-400 predicate, choice, or score questions.
	extra: optional request parameters, e.g. prompt_cache_key,
		prompt_cache_options, and safety_identifier.
	Returns the full response with answers in question order, model, and usage.
	Example:
		OpenAI_Decisions(token).Decisions('The screen is broken.',
			#((type: 'predicate', name: 'damaged',
				instructions: 'Does the customer report a damaged item?')))
	*/
	Decisions(input, questions, model = "gpt-6-luna", extra = #())
		{
		body = Object(:input, :questions, :model).Merge(extra)
		res = .Send(.BaseUrl $ "decisions", body)
		if .usageTracker isnt false
			(.usageTracker)(res.usage)
		return res
		}

	ProcessOpenAIError(res)
		{
		if res.Member?(#usage) and .usageTracker isnt false
			(.usageTracker)(res.usage)
		}
	}
