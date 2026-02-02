// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
function (images, image, color, w/*unused*/, h/*unused*/, dark?/*unused*/ = false,
	padding/*unused*/ = 0)
	{
	images.Add([IconFont().MapToCharCode(image), color])
	}