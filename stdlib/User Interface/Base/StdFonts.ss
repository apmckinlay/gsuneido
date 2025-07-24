// Copyright (C) 2012 Suneido Software Corp. All rights reserved worldwide.
class
	{
	New()
		{ throw "not meant to be instantiated" }
	PtSize(lfsize = false)
		{
		return StdFontsSize.PtSize(lfsize, .getLogPixel())
		}
	LfSize(ptsize = false)
		{
		return StdFontsSize.LfSize(ptsize, .getLogPixel())
		}
	getLogPixel()
		{
		if Sys.SuneidoJs?()
			return WinDefaultDpi /* CSS uses 96 px per inch */
		WithDC(NULL)
			{|hdc|
			logPixel = GetDeviceCaps(hdc, GDC.LOGPIXELSY)
			return logPixel <= 0 ? WinDefaultDpi : logPixel
			}
		}
	Ui()
		{
		return "Segoe UI"
		}
	Mono()
		{
		return "Consolas"
		}
	Serif()
		{
		return "Georgia"
		}
	Sans()
		{
		return "Verdana"
		}
	// helper methods
	Font(font) // handle e.g. @mono
		{
		if font[0] is '@'
			{
			name = font[1..].Capitalize()
			if .Member?(name)
				font = StdFonts[name]()
			}
		return font
		}
	Weight(weight) // handle e.g. 'bold'
		{
		if String?(weight)
			weight = FW.GetDefault(weight.Upper(), weight)
		if weight not in ("", false) and (weight < 0 or 999/*= max weight*/ < weight)
			{
			ProgrammerError("invalid font weight: " $ Display(weight))
			return FW.NORMAL
			}
		return weight
		}
	FontSize(size, def = false) // handle e.g. '+2'
		{
		if String?(size)
			size = Number(size) + (def is false ? .PtSize() : def)
		return size
		}

	// At smaller sizes, scintilla fonts seem smaller then standard
	SciSize(size, def = false)
		{
		fontSize = .FontSize(size, def)
		if String?(size) and fontSize < 12 /*= Scintilla Fonts match LogFont size*/
			fontSize++
		return fontSize
		}

	// returns the CSS font specification
	// based on Business > Customize Font
	GetCSSFont(sizeFactor = 1)
		{
		family = .Ui()
		if Suneido.logfont.lfFaceName isnt Suneido.stdfont.lfFaceName
			family = Suneido.logfont.lfFaceName $ ', ' $ family
		logfont = Suneido.logfont
		stdfont = Suneido.stdfont
		baseSize = 16
		size = ((baseSize + stdfont.lfHeight - logfont.lfHeight) * sizeFactor).Round(0)
		return 'font: ' $ size $ 'px ' $ family $ ';'
		}
	}
