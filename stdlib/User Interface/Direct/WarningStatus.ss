// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
/* Warning Contributions should be defined as:
class
	{
	// Lower number means it will be processed sooner
	CallClass(data)
		{
		// This will need to return a object of warning objects.
		// Each warning object will need a priority number, and a message to display (msg)
		// IE:
		return Object(
			[priority: <number>, msg: <warning>],
			[priority: <number>, msg: <warning>],
			...
			)
		}
	}
*/
function (data, warningPrefix, contributionName, lastJoin = 'and')
	{
	warnings = Object()
	for cont in Contributions(contributionName)
		warnings.Add(@cont(data))
	warnings.RemoveIf({ it.Empty?() }).Sort!(By(#priority)).Map!({ it.msg })
	if warnings.Size() > 1
		warnings.Add(lastJoin $ ' ' $ warnings.PopLast())
	return Opt(warningPrefix $ ' ', warnings.Join('; '))
	}