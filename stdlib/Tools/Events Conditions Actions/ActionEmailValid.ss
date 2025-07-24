// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
function (action)
	{
	toMem = action.action_setting.FindIf({ it.field is 'action_email_to' })
	toOb = action.action_setting[toMem]
	return toOb.value is ''
		? 'Recipient required for ' $ action.action_name
		: ''
	}
