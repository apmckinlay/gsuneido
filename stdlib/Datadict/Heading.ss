// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
function (field)
	{
	ddVals = Datadict(field, #(Heading, Prompt))
	heading = ddVals.Member?("Heading")
		? ddVals.Heading
		: ddVals.Member?("Prompt")
			? ddVals.Prompt
			: field
	return TranslateLanguage(heading)
	}
