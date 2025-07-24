// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
// REFACTOR: extract duplicate code in child classes for fonts
class
	{
	New()
		{
		.Fonts = Object()
		}

	Driver: false
	SetDriver(.Driver) { }

	Default(@args)
		{
		method = args[0]
		if .Driver.Method?(method)
			return .Driver[method](@+1 args)
		throw 'method not found in Report: ' $ method
		}

	GetReportDefaultFont()
		{
		return .Driver is false
			? Document_Builder.GetDefaultFont()
			: .Driver.GetDefaultFont()
		}

	GetFontSize(font)
		{
		defaultFont = .GetReportDefaultFont()
		if not font.Member?('size')
			return defaultFont.size

		fontSize = font.size
		if String?(font.size) and font.size =~ '^[+-]'
			{
			if not Object?(curFont = .GetFont())
				curFont = defaultFont
			fontSize = curFont.GetDefault('size', defaultFont.size) + Number(font.size)
			}
		return fontSize
		}

	Construct(item)
		{
		return Report.Construct(item)
		}

	Destroy()
		{
		for font in .Fonts
			DeleteObject(font)
		}
	}