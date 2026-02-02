// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
function (logfont)
	{
	if Suneido.Member?(#logfont) and logfont is Suneido.logfont
		return

	Suneido.logfont = logfont.Set_readonly()
	SuRenderBackend().RecordAction(false, #SuSetGuiFont, [logfont])
	}