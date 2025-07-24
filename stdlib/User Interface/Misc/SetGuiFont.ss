// Copyright (C) 2004 Suneido Software Corp. All rights reserved worldwide.
function (logfont)
	{
	// handle old saved values
	if logfont.lfFaceName is 'MS Sans Serif'
		logfont.lfFaceName = StdFonts.Ui()

	if logfont is Suneido.logfont
		return
	hfont = CreateFontIndirect(logfont)
	if hfont is 0
		return false
	Suneido.logfont = logfont.Set_readonly()
	Suneido.hfont = hfont
	return true
	}