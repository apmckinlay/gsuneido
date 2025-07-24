// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
function (users)
	{
	func = OptContribution('BookNotification_UserList',
		function(users){ return users.Split(',') })
	return func(users)
	}