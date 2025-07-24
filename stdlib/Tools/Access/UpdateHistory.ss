// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
function (data, historyFields, user = false)
	{
	if not Object?(historyFields)
		return

	for action in historyFields.Members()
		if action is 'Last Modified'
			{
			data[historyFields[action].date] = Timestamp()
			data[historyFields[action].user] = user is false ? Suneido.User : user
			}
	}
