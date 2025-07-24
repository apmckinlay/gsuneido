// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
function (field)
	{
	ddVals = Datadict(field, #(SelectPrompt, Prompt, Heading))
	return ddVals.Member?("SelectPrompt")
		? ddVals.SelectPrompt
		: ddVals.Member?('Prompt') and ddVals.Prompt isnt ""
			? ddVals.Prompt
			: ddVals.Member?('Heading') and ddVals.Heading isnt ""
				? ddVals.Heading
				: field
	}
