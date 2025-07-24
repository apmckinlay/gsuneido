// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
Addon
	{
	RequirementsMet?()
		{
		return true
		}

	Init()
		{
		.lockKey = .GetLockKey()
		.editBtn = .FindControl(#Edit)
		// We cannot assume the state is false or true as it can be either on initial load
		.prevState = NULL
		.RunAddon()
		}

	RunAddon()
		{
		unlocked? = false is user = LockManager.CheckIfLockedForReadOnlyPurposes(.lockKey)
		if .prevState is unlocked?
			return
		.prevState = unlocked?
		.editBtn.SetTextColor(unlocked? ? CLR.BLACK : CLR.amber)
		.editBtn.ToolTip(unlocked? ? 'Alt+E' : user $ ' is currently editing this')
		}
	}
