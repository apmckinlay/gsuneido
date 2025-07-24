// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
function (env, newest)
	{
	if '' is v = env.GetDefault('if_modified_since', '')
		return false

	if newest is Date.Begin()
		return false

	if false is date = Date(v.Replace(' GMT', ''))
		return false

	return newest.GMTime() <= date
	}