// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
class
	{
	errors: #(
		File: #(
			'save file'
			'read file': #(Read, Readline, Seek, Tell),
			'save file': #(Write, Writeline)
			),
		DirExists: #('access location'))

	// not handled yet:
	// Seek, Tell,
	CallClass(target, block)
		{
		try
			{
			block()
			return true
			}
		catch(err, "File|DirExists")
			{
			prefix = err.BeforeFirst(":")
			errMsg = err.AfterLast(":").Trim()
			action = .getAction(prefix.Tr('?'), err)
			// ensure the first character of the sentance is uppercase
			// we can't use .Capitalize() here as it will convert the filename to all
			// lowercase letters. the following regex ONLY converts the very first
			// character
			return 'Can not ' $ action $ ": " $ target $ "\n\n" $
				errMsg.Replace('\w', '\u&', 1)
			}
		}
		getAction(prefix, err)
			{
			errType = .errors[prefix]
			errMsg = errType.FindIf(
				{ it.Has?(err.AfterFirst(":").BeforeLast(":").Trim()) })
			return errMsg is false or errMsg is 0 ? errType[0] : errMsg
			}
	}