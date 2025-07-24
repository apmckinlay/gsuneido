// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(field)
		{
		ddVals = Datadict(field, #(Prompt))
		return ddVals.Member?("Prompt") ? ddVals.Prompt : field
		}

	WithInfo(field, fallbackMethod = 'Prompt')
		{
		dd = Datadict(field)
		if dd.Member?(#PromptInfo)
			return .promptFromInfo(dd.PromptInfo)

		return Global(fallbackMethod)(field)
		}

	promptFromInfo(info)
		{
		return [info.prefix, Global(info.promptMethod)(info.baseField), info.suffix].
			Each(#Trim).
			Remove('').
			Join(' ')
		}
	}