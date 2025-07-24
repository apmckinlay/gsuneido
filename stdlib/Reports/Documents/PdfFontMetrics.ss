// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
Singleton
	{
	New()
		{
		.charWidths = Object().Set_default(Object())
		.charHeights = Object()
		.fontInfoOb = Object().Set_default(Object())
		}

	GetFontReference(font)
		{
		s = .weightReference(font.GetDefault(#weight, 0))
		if font.Member?(#italic) and font.italic is true
			s $= "I"
		return s
		}

	weightReference(weight)
		{
		if String?(weight)
			weight = StdFonts.Weight(weight)
		return weight > 550 /*= bold*/ ? 'B' : ''
		}

	GetFontName(font)
		{
		fontRef = .initFontInfo(font)
		return .fontInfoOb[fontRef].FontName
		}
	GetFontCharRange(font)
		{
		fontRef = .getFontName(font)
		if .charWidths[fontRef].Empty?()
			.readFontsWidth(fontRef)
		charWidthsIndex = .charWidths[fontRef].Members()
		return Object('FirstChar': charWidthsIndex.Min(),
			'LastChar': charWidthsIndex.Max())
		}
	GetFontWidths(font)
		{
		fontRef = .getFontName(font)
		if .charWidths[fontRef].Empty?()
			.readFontsWidth(fontRef)
		widths = Object()
		charWidthsIndex = .charWidths[fontRef].Members()
		firstIndex = charWidthsIndex.Min()
		lastIndex = charWidthsIndex.Max()
		for (i = firstIndex; i <= lastIndex; i++)
			widths[i - firstIndex] = .charWidths[fontRef].GetDefault(i, 0)
		return widths
		}
	GetFontAscender(font)
		{
		fontRef = .initFontInfo(font)
		if not .charHeights.Member?(fontRef)
			.readFontHeight(fontRef)
		return Number(.fontInfoOb[fontRef].GetDefault('Ascender', .charHeights[fontRef]))
		}
	GetFontDescender(font)
		{
		fontRef = .initFontInfo(font)
		return Number(.fontInfoOb[fontRef].GetDefault('Descender', 0))
		}
	GetFontCapHeight(font)
		{
		fontRef = .initFontInfo(font)
		if not .charHeights.Member?(fontRef)
			.readFontHeight(fontRef)
		return Number(.fontInfoOb[fontRef].GetDefault('CapHeight', .charHeights[fontRef]))
		}
	GetFontBBox(font)
		{
		fontRef = .initFontInfo(font)
		return .fontInfoOb[fontRef]['FontBBox'].Split(' ').Map!(Number)
		}
	GetFontItalicAngle(font)
		{
		fontRef = .initFontInfo(font)
		return Number(.fontInfoOb[fontRef]['ItalicAngle'])
		}
	initFontInfo(font)
		{
		fontRef = .getFontName(font)
		if .fontInfoOb[fontRef].Empty?()
			.readFontInfo(fontRef)
		return fontRef
		}
	readFontInfo(fontRef)
		{
		if false is afm = Query1Cached('imagebook', name: fontRef $ '.afm')
			throw 'PdfFontMetrics cannot find the adobe font metrics for ' $ fontRef

		fontInfo = .fontInfoOb[fontRef]
		for line in afm.text.Lines()
			{
			entry = line.BeforeFirst(' ')
			value = line.AfterFirst(' ')

			if entry is 'C' or entry.Prefix?('Start') or entry.Prefix?('End') or
				entry is 'Comment' or entry is ''
				continue

			fontInfo[entry] = value
			}
		}
	GetCharWidth(font, character)
		{
		charAsc = character.Asc()
		fontRef = .getFontName(font)
		if .charWidths[fontRef].Empty?()
			.readFontsWidth(fontRef)
		fontWidths = .charWidths[fontRef]
		defaultWidth = fontWidths[32] /*= size of space */
		return fontWidths.GetDefault(charAsc, defaultWidth)
		}

	readFontsWidth(fontRef)
		{
		if false is afm = Query1Cached('imagebook', name: fontRef $ '.afm')
			throw 'PdfFontMetrics cannot find the adobe font metrics for ' $ fontRef

		widths = .charWidths[fontRef]
		for line in afm.text.Lines()
			if line.Prefix?('C ')
				if -1 isnt asc = Number(line.AfterFirst('C ').BeforeFirst(';').Trim())
					widths[asc] = Number(line.AfterLast('WX ').BeforeFirst(';').Trim())
		}

	getFontName(font)
		{
		ref = .GetFontReference(font)
		fontName = PdfFonts.GetCompatibleFontName(font.GetDefault("name", ""))
		fontRef = fontName $ Opt('-', ref.Replace('B', 'Bold').Replace('I', 'Oblique'))
		return fontRef
		}

	GetFontHeight(font)
		{
		fontRef = .getFontName(font)
		if not .charHeights.Member?(fontRef)
			.readFontHeight(fontRef)
		return .charHeights[fontRef]
		}

	readFontHeight(fontRef)
		{
		if false is afm = Query1Cached('imagebook', name: fontRef $ '.afm')
			throw 'PdfFontMetrics cannot find the adobe font metrics for ' $ fontRef

		uy = afm.text.AfterFirst('FontBBox ').BeforeFirst('\n').Trim().AfterLast(' ')
		height = Number(uy)
		.charHeights[fontRef] = height
		}
	}
