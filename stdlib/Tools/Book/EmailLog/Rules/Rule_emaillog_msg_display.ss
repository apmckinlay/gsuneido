// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
function ()
	{
	if .emaillog_view_all? is true or .emaillog_cur_user is .emaillog_user
		return .emaillog_msg

	return 'Only users with the admin role and the email sender can see the contents'
	}