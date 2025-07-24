// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
function (logfont)
	{
	Suneido.logfont = logfont
	font = logfont.lfFaceName
	css = `body, input, button, table, textarea
	{`
	css $= 'font-family: ' $ Component.FontFamily(font) $ ';\r\n'
	css $= 'font-size: ' $
		-StdFontsSize.LfSize(logfont.fontPtSize, WinDefaultDpi) $ 'px;\r\n'
	css $= 'font-weight: ' $ logfont.GetDefault(#lfWeight, FW.NORMAL) $ ';\r\n'
	if logfont.GetDefault(#lfItalic, 0) isnt 0
		css $= 'font-style: italic;\r\n'
	css $= '}'
	LoadCssStyles('su_general_styles', css, override?:)
	}
