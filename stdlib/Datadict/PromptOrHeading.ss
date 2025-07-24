// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
function (field)
	{
	ddVals = Datadict(field, #(Prompt, Heading, SelectPrompt))
	s = ddVals.Member?("Prompt") and ddVals.Prompt isnt ""
		? ddVals.Prompt
		: ddVals.Member?("Heading")
			? ddVals.Heading
			: ddVals.Member?('SelectPrompt')
				? ddVals.SelectPrompt
				: field
	return TranslateLanguage(s)
	}
