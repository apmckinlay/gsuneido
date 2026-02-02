// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
function (action, params, allComponents)
	{
	msg = 'action ' $ action.action $ ' not found in the component'
	id = params.id = action.uniqueId
	controls = SuRenderBackend().SuRenderBackend_controls

	nearbys = Object()
	for (i = -5; i <= 5/*=upper*/; i++)
		{
		if controls.Member?(id + i)
			nearbys.Add(Display(controls[id + i]) $ ' - ' $ (id + i))
		}
	params.controlSize = controls.Size()
	params.controlMax = controls.Members().Max()
	params.controlNearbys = nearbys

	LogErrors('Dump debug36105.log')
		{
		allControls = Object()
		for m in controls.Members()
			allControls[m] = Display(controls[m])
		File('debug36105.log', 'a')
			{ |f|
			f.Writeline(Display(Date()).RightFill(80/*=w*/, '*'))
			mems = allComponents.Members().MergeUnion(allControls.Members()).Sort!()
			for m in mems
				f.Writeline(m $ '\t' $ allControls.GetDefault(m, '').RightFill(40/*=w*/) $
					'\t' $ allComponents.GetDefault(m, '').RightFill(40/*=w*/))
			}
		}

	if false isnt control = SuRenderBackend().GetRegisteredControl(id)
		{
		params.control = Display(control)
		params.controlId = control.UniqueId
		}

	SuneidoLog('ERROR: (CAUGHT) ' $ msg, :params, caughtMsg: 'for debug 36105')
	SuRenderBackend().DumpStatus(msg)
	}