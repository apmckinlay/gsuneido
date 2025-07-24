// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
function (warningsAllMethods)
	{
	if warningsAllMethods is ""
		return ""

	warningText = ""
	for i in warningsAllMethods.Members(list:)
		{
		individualMethodWarnings = warningsAllMethods[i]
		warningText $= Opt('\n', individualMethodWarnings.desc, ':\n')
		individualMethodWarnings.warnings.Each({ warningText $= Opt(it.name, "\n") })
		}
	return warningText
	}