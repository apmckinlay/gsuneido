// Copyright (C) 2020 Axon Development Corporation All rights reserved worldwide.
function (el)
	{
	target = SuRenderBackend().GetRegisteredControl(el)
	if not target.Method?(#GetWindowText)
		{
		SuServerPrint("target not support GetWindowText", :target)
		return ''
		}
	return target.GetWindowText()
	}