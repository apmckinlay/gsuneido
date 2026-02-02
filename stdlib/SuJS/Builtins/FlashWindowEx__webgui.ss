// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
function (fi)
	{
	args = [fi.hwnd, show?: fi.dwFlags isnt FLASHW.STOP]
	if fi.Member?(#message)
		args.text = fi.message
	SuRenderBackend().RecordAction(false, #SuFlashWindow, args)
	}