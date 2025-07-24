// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
// NOTE: To add new icons, see http://appserver.axon:8080/Wiki?UpdateSuneidoIconFont
Singleton
	{

	New()
		{
		.map = Object()
		for font in IconFontHelper.Fonts.Members()
			.initFont(font)
		}
	initFont(fontName)
		{
		try
			{
			data = Query1('imagebook', path: "/res", name: fontName).text
			size = data.Size()
			n = Object()
			hfont = AddFontMemResourceEx(data, size, 0, n)
			if hfont is 0
				throw 'AddFontMemResourceEx failed'
			.initMap(fontName)
			}
		catch (e)
			{
			SuneidoLog('ERROR: IconFont: ' $ e)
			Fatal('Initialize icon resources failed')
			}
		}
	map: #()
	initMap(fontName)
		{
		.map[fontName] = IconFontHelper.InitMap(fontName)
		}
	MapToCharCode(image)
		{
		image = image.RemoveSuffix('.emf')
		if image.Size() is 1 and image =~ '[[:graph:]]'
			return Object(code: image.Asc(), char: image, font: 'Arial')

		for fontName in .map.Members()
			if .map[fontName].Member?(image)
				{
				code = .map[fontName][image]
				return Object(:code, char: code.Chr(),
					font: fontName.RemoveSuffix('.ttf'))
				}
		return false
		}

	// Run Before / After changes and compare the output. Should be similar
	TestIcons()
		{
		controls = Object(#Vert)
		i = 0
		for font in IconFontHelper.Fonts.Members()
			{
			controls.Add(Object('Horz', 'Fill',
				Object('Static', font, size: '+2', weight: 'bold'), 'Fill', xstretch: 0))
			controls.Add('EtchedLine')
			horz = Object(#Horz)
			for image in IconFontHelper.Fonts[font]
				{
				if image.Prefix?('UNUSED')
					continue
				horz.Add(
					Object(#EnhancedButton, :image, imagePadding: .1, tip: image),
					#Skip)
				if ++i >= 10 /*= row size*/
					{
					controls.Add(horz, #Skip)
					horz = Object(#Horz)
					i = 0
					}
				}
			if horz.Size() > 1
				controls.Add(horz)
			}
		Window(controls)
		}

	// override super.Reset to avoid getting cleared
	Reset() {}
	}