// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.

// why isn't this a method in BookNotification ?
// doesn't seem to be used anywhere else

function (list, user)
	{
	if list is 'ALL'
		return true
	if String?(list)
		{
		if not list.Has?(user)
			return false
		list = list.Split(',')
		}
	return list.Has?(user)
	}