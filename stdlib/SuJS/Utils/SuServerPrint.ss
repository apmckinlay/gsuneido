// Copyright (C) 2019 Axon Development Corporation All rights reserved worldwide.
function (@args)
	{
	if not Sys.SuneidoJs?()
		{
		Print(@args)
		return
		}
	SuRenderBackend().RecordAction(false, 'Print', args)
	}