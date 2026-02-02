// Copyright (C) 2024 Axon Development Corporation All rights reserved worldwide.
function()
	{
	x = SuRenderBackend().Status()
	dimens = x.GetDefault('dimension', Object(width: 600, height: 400))
	return Object(right: dimens.width, top: 0, left: 0, bottom: dimens.height)
	}