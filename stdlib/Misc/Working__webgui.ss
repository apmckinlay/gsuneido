// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
class
	{
	CallClass(message, block, quiet? = false, font/*unused*/ = '', size/*unused*/ = '',
		weight/*unused*/ = '', color/*unused*/ = '')
		{
		showOverlay? = not (quiet? or not Sys.GUI?())
		if showOverlay?
			SuRenderBackend().Overlay(message)
		Finally(block, {
			if showOverlay?
				SuRenderBackend().Overlay(message, hide:) })
		}
	}