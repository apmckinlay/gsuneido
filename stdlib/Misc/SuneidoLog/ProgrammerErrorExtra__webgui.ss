// Copyright (C) 2024 Axon Development Corporation All rights reserved worldwide.
function (msg)
	{
	if Sys.SuneidoJs?()
		SuRenderBackend().DumpStatus(msg)
	}