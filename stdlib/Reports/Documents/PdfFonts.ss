// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
// NOTE: use ReportFontTester to test changes
class
	{
	FontNames: (
		(name: 'Helvetica', ref: 'F1', standard: true)
		(name: 'Helvetica-Bold', ref: 'F1B', standard: true)
		(name: 'Helvetica-Oblique', ref: 'F1I', standard: true)
		(name: 'Helvetica-BoldOblique',  ref: 'F1BI', standard: true)
		(name: 'Courier',  ref: 'F2', standard: true)
		(name: 'Courier-Bold',  ref: 'F2B', standard: true)
		(name: 'Courier-Oblique', ref: 'F2I', standard: true)
		(name: 'Courier-BoldOblique',  ref: 'F2BI', standard: true)
		(name: 'Times-Roman',  ref: 'F3', standard: true)
		(name: 'Times-Bold',  ref: 'F3B', standard: true)
		(name: 'Times-Italic', ref: 'F3I', standard: true)
		(name: 'Times-BoldItalic', ref: 'F3BI', standard: true)
		//flags is a collection of flags defining various characteristics of the font
		//see PDF Reference Section 5.7.1
		(name: 'Free 3 of 9 Regular', ref: 'F4', standard: false, file: 'FREE3OF9.TTF',
			flags: 33) // 33 = 2 ^ 6 + 1
		(name: 'Free 3 of 9 Extended Regular', ref: 'F5', standard: false,
			file: 'FRE3OF9X.TTF', flags: 32)) // 32 = 2 ^ 6

	fonts: (
		(name: "Helvetica" ref: "/F1")
		(name: "Courier" ref: "/F2")
		(name: "Times" ref: "/F3")
		(name: "Free 3 of 9 Regular" ref: '/F4')
		(name: "Free 3 of 9 Extended Regular" ref: '/F5'))

	IsStandardFont(fontRef)
		{
		fontRef = fontRef.AfterFirst('/')
		return .FontNames.FindOne({it.ref is fontRef}).standard
		}

	ConvertFont(font)
		{
		font = font is false ? Object() : font
		fontStyle = PdfFontMetrics.GetFontReference(font)
		font = .GetCompatibleFont(font)
		x = .fonts.FindOne({it.name is font.name })
		if x isnt false
			return x.ref $ fontStyle
		return "/F1" $ fontStyle
		}

	GetCompatibleFont(font)
		{
		// set Courier family to Courier, leave Helvetica as is,
		// and change all others to Helvetica
		name = font.GetDefault('name', '')
		if .fonts.HasIf?({ it.name is name })
			return font
		font = font.Copy()
		font.name = .GetCompatibleFontName(name)
		return font
		}

	GetCompatibleFontName(fontName)
		{
		return fontName.Has?('Courier')
			? 'Courier'
			: fontName.Has?('Times')
				? 'Times'
				: fontName is "Free 3 of 9 Regular" or
					fontName is "Free 3 of 9 Extended Regular"
					? fontName
					: 'Helvetica'
		}

	GetCharWidth(font, widthChar)
		{
		return PdfFontMetrics().GetCharWidth(font, widthChar)
		}

	GetFontHeight(font)
		{
		return PdfFontMetrics().GetFontHeight(font) * font.size
		}

	GetFontDescender(font)
		{
		return PdfFontMetrics().GetFontDescender(font) * font.size
		}

	BuildEmbeddedFont(font, writer)
		{
		fontInfo = .buildFontInfo(font)
		writer.AddObject(.buildEmbededFontObject(fontInfo, writer.NextObjId))
		writer.AddObject(.buildEmbededFontStream(font))
		}

	buildFontInfo(font)
		{
		range = PdfFontMetrics().GetFontCharRange(font)
		flags = .FontNames.FindOne({it.name is font.name}).flags
		fontInfo = Object(
			'BaseFont': 	PdfFontMetrics().GetFontName(font)
			'FirstChar': 	range.FirstChar
			'LastChar': 	range.LastChar
			'Widths': 		PdfFontMetrics().GetFontWidths(font)
			'Subtype':		'TrueType'
			'Type':			'Font'
			'Encoding':		'WinAnsiEncoding'
			'FontDescriptor': Object(
				'Ascent': 		PdfFontMetrics().GetFontAscender(font)
				'Descent': 		PdfFontMetrics().GetFontDescender(font)
				'CapHeight': 	PdfFontMetrics().GetFontCapHeight(font)
				'Flags':		flags
				'FontBBox':		PdfFontMetrics().GetFontBBox(font)
				'FontName':		PdfFontMetrics().GetFontName(font)
				'ItalicAngle':	PdfFontMetrics().GetFontItalicAngle(font)
				//stackoverflow.com/questions/35485179/stemv-value-of-the-truetype-font
				'StemV':		80
				'Type':			'FontDescriptor'))
		return fontInfo
		}

	buildEmbededFontStream(font)
		{
		if false is ttf = Query1Cached('imagebook', name: font.file)
			throw 'PdfFont cannot find the font file for ' $ font.file
		file = ttf.text
		compressedFile = Zlib.Compress(file)
		s = "\n<<\n/Filter /FlateDecode" $
			"\n/Length " $ compressedFile.Size() $
			"\n/Length1 " $ file.Size() $
			"\n/Type /Stream\n>>" $
			"\nstream\n" $ compressedFile $ "\nendstream\nendobj\n"
		return s
		}

	buildEmbededFontObject(fontInfo, resourceIndex)
		{
		fontInfo.FontDescriptor.Add((resourceIndex+1) $ ' 0 R', at: 'FontFile2')
		s = '\n<<\n' $ .buildRecord(fontInfo) $ '>>\nendobj\n'
		return s
		}

	buildRecord(record)
		{
		s = ''
		for m, v in record
			{
			s $= '/' $ m $ ' '
			switch Type(v)
				{
				case 'String':
					if v.Match('^\d+\s+\d+\s+R$') isnt false
						s $= v $ '\n'
					else if v.Has?(' ')
						s $= '(' $ v $ ')'
					else
						s $= '/' $ v $ '\n'
				case 'Number':
					s $= v $ '\n'
				case 'Object':
					if v.Size() is v.Members(list:).Size()
						s $= '[' $ v.Join(' ') $ ']\n'
					else
						s $= '<<\n' $ .buildRecord(v) $ '>>\n'
				}
			}
		return s
		}

	StripInvalidChars(data)
		{ // out of character range, default to space
		return data.Tr('^\x20-\xfb', ' ')
		}
	}
