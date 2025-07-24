// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
function (map, prompts)
	{
	if String?(prompts)
		prompts = prompts.Split(',')
	fields = Object()
	for prompt in prompts
		{
		prompt = prompt.Trim()
		if map.Member?(prompt)
			fields.Add(map[prompt])
		}
	return fields
	}