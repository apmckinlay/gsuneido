// Copyright (C) 2022 Axon Development Corporation All rights reserved worldwide.
function (ob)
	{
	if false is color = OkCancel(Object('ColorRect', ob.rgbResult, choose?:, xstretch: 1),
		title: 'Choose Color')
		return false
	else if color is true // OkCancelWrapper will return true if value was 0
		color = 0

	ob.rgbResult = color
	return true
	}