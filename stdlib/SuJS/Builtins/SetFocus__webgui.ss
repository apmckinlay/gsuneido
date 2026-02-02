// Copyright (C) 2019 Axon Development Corporation All rights reserved worldwide.
class
	{
	CallClass(uniqueId, force = false)
		{
		if uniqueId is false
			{
			SuServerPrint("SetFocus on false")
			return
			}

		if .handleKillFocus(uniqueId) is false and force is false
			return
		.handleSetFocus(uniqueId)
		}

	handleKillFocus(newFocus)
		{
		if false is prev = SuRenderBackend().Status().Focus
			return true

		if prev is newFocus
			return false

		SuRenderBackend().CallProc(prev, #KILLFOCUS, [wParam: newFocus])
		return true
		}

	handleSetFocus(uniqueId)
		{
		SuRenderBackend().CancelAction(#ignore, 'SetFocus')
		if uniqueId is NULL
			{
			SuRenderBackend().RecordAction(false, 'SuClearFocus', Object())
			SuRenderBackend().Status().Focus = false
			}
		else
			{
			SuRenderBackend().RecordAction(uniqueId, 'SetFocus', Object())
			SuRenderBackend().Status().Focus = uniqueId
			SuRenderBackend().CallProc(uniqueId, #SETFOCUS)
			}
		}
	}