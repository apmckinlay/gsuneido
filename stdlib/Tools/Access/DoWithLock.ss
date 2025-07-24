// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
function(key, block)
	{
	if true isnt sid = LockManager.Lock(key)
		throw "This record is already being edited by " $
			LockManager.UserFromSessionId(sid) $ ". " $
			"None of the updates were applied"
	Finally(block,
		{ LockManager.Unlock(key) })
	}