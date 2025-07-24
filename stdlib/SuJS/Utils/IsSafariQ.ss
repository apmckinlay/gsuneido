// Copyright (C) 2024 Axon Development Corporation All rights reserved worldwide.
function ()
	{
	if false is userAgent = Suneido.GetDefault(#userAgent, false)
		return false
	userAgent = userAgent.Lower()
	return userAgent.Has?('safari') and not userAgent.Has?('chrome')
	}
