// Copyright (C) 2015 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_GetCompatibleFont()
		{
		Assert(PdfFonts.GetCompatibleFont(#()) is: #(name: Helvetica))
		Assert(PdfFonts.GetCompatibleFont(#(name: Helvetica)) is: #(name: Helvetica))
		Assert(PdfFonts.GetCompatibleFont(#(name: Arial)) is: #(name: Helvetica))
		Assert(PdfFonts.GetCompatibleFont(#(name: Courier)) is: #(name: Courier))
		Assert(PdfFonts.GetCompatibleFont(#(name: 'Courier New')) is: #(name: Courier))
		Assert(PdfFonts.GetCompatibleFont(#(name: Arial, size: 12))
			is: #(name: Helvetica, size: 12))
		Assert(PdfFonts.GetCompatibleFont(#(name: 'Times')) is: #(name: Times))
		Assert(PdfFonts.GetCompatibleFont(#(name: 'Times New Roman'))
			is: #(name: 'Times'))
		}

	Test_GetCompatibleFontName()
		{
		Assert(PdfFonts.GetCompatibleFontName('') is: 'Helvetica')
		Assert(PdfFonts.GetCompatibleFontName('Courier') is: 'Courier')
		Assert(PdfFonts.GetCompatibleFontName('Courier New') is: 'Courier')
		Assert(PdfFonts.GetCompatibleFontName('Times') is: 'Times')
		Assert(PdfFonts.GetCompatibleFontName('Times New Roman') is: 'Times')
		Assert(PdfFonts.GetCompatibleFontName('Times Roman') is: 'Times')
		}

	Test_BuildEmbeddedFont()
		{
		writer = class
			{
			New()
				{
				.Ob = Object()
				}
			AddObject(content)
				{
				.Ob.Add(content)
				}
			Getter_NextObjId()
				{
				return .Ob.Size()
				}
			}()
		for f in PdfFonts.FontNames.Filter({not it.standard})
			{
			PdfFonts.BuildEmbeddedFont(f, writer)
			s = writer.Ob
			Assert(s[0] has: "FontBBox")
			Assert(s[1] has: "stream")
			}
		}
	}