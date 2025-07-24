// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_BuildXRef()
		{
		Assert(PdfDriver.BuildXRef(1, #(0), 1, 100)
			is: #("xref"
				"0 1"
				"0000000000 65535 f "
				"trailer <</Size 1/Root 1 0 R>>"
				"startxref"
				"100"
				"%%EOF").Join("\n"))
		Assert(PdfDriver.BuildXRef(1, #(1), 1, 100)
			is: #("xref"
				"0 1"
				"0000000001 00000 n "
				"trailer <</Size 1/Root 1 0 R>>"
				"startxref"
				"100"
				"%%EOF").Join("\n"))
		}

	testDriver: PdfDriver
		{
		PdfDriver_jpeg(unused)
			{
			return FakeObject(
				GetColorSpace: 'gray',
				GetWidth: 5,
				GetHeight: 5)
			}
		PdfDriver_createWriter(filename)
			{
			return new (.testWriter)(filename)
			}
		testWriter: PdfDriverFileWriter
			{
			PdfDriverFileWriter_file(unused)
				{
				return FakeFile('')
				}
			}
		}

	Test_buildXObjectRef()
		{
		p = new (.testDriver)('testFile.pdf', compress: false)
		p.PdfDriver_xObjects = #("stream
			image1
			endstream")

		p.PdfDriver_buildXObjectRef()
		Assert(p.PdfDriver_imageReferences isSize: 1)

		p.PdfDriver_xObjects = #("stream
			image1
			endstream")
		// duplicate image on second page so no reference added
		p.PdfDriver_buildXObjectRef()
		Assert(p.PdfDriver_imageReferences isSize: 1)

		p.PdfDriver_xObjects = #("stream
			image2
			endstream")
		// unique image added on third page so reference added
		p.PdfDriver_buildXObjectRef()
		Assert(p.PdfDriver_imageReferences isSize: 2)
		}

	Test_BuildImageObjectHead()
		{
		head = PdfDriver.BuildImageObjectHead(1900, 2400, 'DeviceRGB', 9999)
		Assert(head
			is: "<</Type /XObject\n" $
			"/Subtype /Image\n" $
			"/Width 1900\n" $
			"/Height 2400\n" $
			"/ColorSpace /DeviceRGB\n" $
			"/BitsPerComponent 8\n" $
			"/Length 9999\n" $
			"/Filter /DCTDecode\n" $
			">>")

		head = PdfDriver.BuildImageObjectHead(1900, 2400, 'DeviceRGB', 9999,
			additionalEntries: #(`SMask 48 0 R`))
		Assert(head
			is: "<</Type /XObject\n" $
			"/Subtype /Image\n" $
			"/Width 1900\n" $
			"/Height 2400\n" $
			"/ColorSpace /DeviceRGB\n" $
			"/BitsPerComponent 8\n" $
			"/Length 9999\n" $
			"/Filter /DCTDecode\n" $
			"/SMask 48 0 R\n" $
			">>")
		}

	Test_main()
		{
		// don't compress because results very between implementations
		p = new (.testDriver)('testFile.pdf', compress: false)
		.testemptyPdf(p)

		p = new (.testDriver)('testFile.pdf', compress: false)

		p.AddPage(#(width: 5, height: 12))
		Assert(p.PdfDriver_pages is: 1)
		Assert(p.PdfDriver_pageRefs is: #())
		Assert(p.PdfDriver_stream is: '')
		Assert(p.PdfDriver_xObjects is: #())
		Assert(p.PdfDriver_embeddedFonts is: .fonts)
		Assert(p.PdfDriver_docHeight is: 864)
		Assert(p.PdfDriver_docWidth is: 360)
		// 3 basics + 12 fonts + 1 Producer (Suneido PDF Generator)
		Assert(p.PdfDriver_writer.NextObjId is: 16)
		Assert(p.PdfDriver_writer.TotalLength is: 1325)

		.testAddText(p)
		p.EndPage()
		Assert(p.PdfDriver_pageRefs is: #("18 0 R"))
		// 2 from embed font + 3 extra from .buildPage
		Assert(p.PdfDriver_writer.NextObjId is: 21)
		Assert(p.PdfDriver_embeddedFonts
			is: .fonts.Copy().Add([id: 16, ref: #F5], at: 'Free 3 of 9 Extended Regular'))

		p.AddPage(#(width: 5, height: 12))
		Assert(p.PdfDriver_pages is: 2)
		Assert(p.PdfDriver_pageRefs is: #("18 0 R"))
		Assert(p.PdfDriver_stream is: '')
		Assert(p.PdfDriver_xObjects is: #())
		Assert(p.PdfDriver_embeddedFonts
			is: .fonts.Copy().Add([id: 16, ref: #F5], at: 'Free 3 of 9 Extended Regular'))
		Assert(p.PdfDriver_docHeight is: 864)
		Assert(p.PdfDriver_docWidth is: 360)

		.testAddArc(p)
		.testAddImage(p)

		p.EndPage()
		Assert(p.PdfDriver_pageRefs is: #("18 0 R", "22 0 R"))
		// 1 xObject for image + 3 extra from .buildPage
		Assert(p.PdfDriver_writer.NextObjId is: 25)
		Assert(p.PdfDriver_imageReferences isSize: 1)

		.testFinalBuild(p)
		}

	fonts: (
		"Helvetica": [id: 3, ref: "F1"],
		"Helvetica-Bold": [id: 4, ref: "F1B"],
		"Helvetica-Oblique": [id: 5, ref: "F1I"],
		"Helvetica-BoldOblique": [id: 6, ref: "F1BI"],
		"Courier": [id: 7, ref: "F2"],
		"Courier-Bold": [id: 8, ref: "F2B"],
		"Courier-Oblique": [id: 9, ref: "F2I"],
		"Courier-BoldOblique": [id: 10, ref: "F2BI"],
		"Times-Roman": [id: 11, ref: "F3"],
		"Times-Bold": [id: 12, ref: "F3B"],
		"Times-Italic": [id: 13, ref: "F3I"],
		"Times-BoldItalic": [id: 14, ref: "F3BI"])
	testemptyPdf(p)
		{
		Assert(p.PdfDriver_pages is: 0)
		Assert(p.PdfDriver_pageRefs is: Object())
		Assert(p.PdfDriver_embeddedFonts is: .fonts)
		Assert(p.PdfDriver_docHeight is: 792)
		Assert(p.PdfDriver_docWidth is: 612)

		// Builds pdf, but circumvents outputting it to file
		p.Finish(ReportStatus.SUCCESS)

		Assert(p.PdfDriver_pages is: 0)
		Assert(p.PdfDriver_pageRefs is: Object())
		Assert(p.PdfDriver_embeddedFonts is: .fonts)
		Assert(p.PdfDriver_docHeight is: 792)
		Assert(p.PdfDriver_docWidth is: 612)
		Assert(p.PdfDriver_writer.NextObjId is: 16)
		Assert(p.PdfDriver_writer.TotalLength is: 1757)
		Assert(p.PdfDriver_writer.PdfDriverFileWriter_f.Get()[1375..].
			Prefix?('xref'))
		}

	testAddText(p)
		{
		data = 'Bold Text Being Added'
		x = 0
		y = 0
		w = 0
		h = 'unused'
		// regular weight <= 550, bold weight > 550
		// NOTE: Calling code (PdfFontMetrics) does not check for "bold" -
		//		 only weight numerical values
		weight = 600
		font = Object(size: 12, :weight, italic: false)
		justify = 'left'
		ellipsis? = false
		color = false

		p.AddText(data, x, y, w, h, font, justify, ellipsis?, color)

		stream = "BT /F1B 12 Tf 0 852.456 Td (Bold Text Being Added) Tj ET\n"
		Assert(p.PdfDriver_stream is: stream)
		Assert(p.PdfDriver_xObjects is: #())
		Assert(p.PdfDriver_embeddedFonts is: .fonts)

		data = 'Barcode Being Added'
		font.name = 'Free 3 of 9 Extended Regular'
		font.weight = 550
		x = 1000
		y = 1000
		p.AddText(data, x, y, w, h, font, justify, ellipsis?, color)

		stream $= "BT /F5 12 Tf 50 804.976 Td (Barcode Being Added) Tj ET\n"
		Assert(p.PdfDriver_stream is: stream)
		Assert(p.PdfDriver_xObjects is: #())
		Assert(p.PdfDriver_embeddedFonts
			is: .fonts.Copy().Add(false, at: 'Free 3 of 9 Extended Regular'))

		data = 'Italic Red Text Being Added'
		y = 250
		x = 5000
// NOTE: justify: right, sets the texts left most character to the left of the left margin
//		 at this time I am not whole heartidly convinced that this is accurate behaviour
		font.italic = true
		font.Delete(#name)
		justify = 'right'
		color = CLR.RED

		p.AddText(data, x, y, w, h, font, justify, ellipsis?, color)

		stream $= "1 0 0 rg\nBT /F1I 12 Tf 101.26 840.328 Td " $
			"(Italic Red Text Being Added) Tj ET\nf\n0 0 0 rg\n"
		Assert(p.PdfDriver_stream is: stream)
		Assert(p.PdfDriver_xObjects is: #())
		Assert(p.PdfDriver_embeddedFonts
			is: .fonts.Copy().Add(false, at: 'Free 3 of 9 Extended Regular'))
		return stream
		}

	testAddArc(p)
		{
		left = 100
		top = 2000
		right = 2000
		bottom = 10
		thick = 5
		lineColor = CLR.GREEN
		xStartArc = 'unused'
		yStartArc = 'unused'
		xEndArc = 'unused'
		yEndArc = 'unused'
		p.AddArc(:left, :top, :right, :bottom, :xStartArc, :yStartArc, :xEndArc, :yEndArc,
			:thick, :lineColor)

		stream = 'q\n5 764 m\n100 863.5 l\n100 764 l\nh\nW\nn\n' $
		'0 1 0 RG\n.33333 w\n5 813.75 m\n' $
		'5 841.2261663040823 26.266474383037 863.5 52.5 863.5 c\n' $
		'78.733525616963 863.5 100 841.2261663040823 100 813.75 c\n' $
		'100 786.2738336959177 78.733525616963 764 52.5 764 c\n' $
		'26.266474383037 764 5 786.2738336959177 5 813.75 c\nS\n0 0 0 RG\nQ\n'
		Assert(p.PdfDriver_stream is: stream)

		left = 2000
		top = 2500
		right = 100
		bottom = 2000
		lineColor = false
		p.AddArc(:left, :top, :right, :bottom, :xStartArc, :yStartArc, :xEndArc, :yEndArc,
			:thick, :lineColor)

		stream $= 'q\n100 739 m\n5 764 l\n100 764 l\nh\nW\nn\n.33333 w\n' $
			'5 751.5 m\n' $
			'5 758.403559372885 26.266474383037 764 52.5 764 c\n' $
			'78.733525616963 764 100 758.403559372885 100 751.5 c\n' $
			'100 744.596440627115 78.733525616963 739 52.5 739 c\n' $
			'26.266474383037 739 5 744.596440627115 5 751.5 c\nS\nQ\n'
		Assert(p.PdfDriver_stream is: stream)
		}

	testAddImage(p)
		{
		stream = p.PdfDriver_stream
		x = y = w = h = 100
		p.AddImage(x, y, w, h, '\x00~IMAGE')

		stream $= 'q\n' $
			'5 0 0 5 5 854 cm\n' $
			'/Im0 Do\n' $
			'Q\n'
		Assert(p.PdfDriver_stream is: stream)
		Assert(p.PdfDriver_xObjects is: #("<</Type /XObject\n/Subtype /Image\n" $
			"/Width 5\n/Height 5\n/ColorSpace /gray\n/BitsPerComponent 8\n/Length 7\n" $
			"/Filter /DCTDecode\n>>\nstream\n\x00~IMAGE\nendstream"))
		}

	testFinalBuild(p)
		{
		p.Finish(ReportStatus.SUCCESS)

		Assert(p.PdfDriver_pages is: 2)
		Assert(p.PdfDriver_pageRefs is: Object("18 0 R", "22 0 R"))
		Assert(p.PdfDriver_embeddedFonts
			is: .fonts.Copy().Add([id: 16, ref: #F5], at: 'Free 3 of 9 Extended Regular'))
		Assert(p.PdfDriver_docHeight is: 864)
		Assert(p.PdfDriver_docWidth is: 360)

		Assert(p.PdfDriver_writer.NextObjId is: 25)
		}

	Test_AddImage()
		{
		m = Mock(PdfDriver)
		m.PdfDriver_docHeight = 1000
		m.PdfDriver_xObjects = Object()
		m.PdfDriver_stream = ''
		m.When.jpeg([anyArgs:]).Return(FakeObject(GetColorSpace: 'gray'))
		m.When.GetImageSize([anyArgs:]).Return(#(width: 100, height: 100))
		m.When.AddImage([anyArgs:]).CallThrough()
		m.AddImage(.00000000425, .00000000426, 1000, 1000, '\xff\xd8')
		Assert(m.PdfDriver_stream has: '50 0 0 50 0 950 cm')
		}

	Test_GetImageSize()
		{
		size = PdfDriver.GetImageSize('', FakeObject(GetWidth: 100, GetHeight: 10))
		Assert(size.width.Round(2) is: 1574.80)
		Assert(size.height.Round(2) is: 157.48)

		.MakeLibraryRecord([name: "ImageMagick", text: `class
			{
			GetWidthHeight(s)
				{
				size = s.BeforeFirst('_')
				return Object(w: Number(size), h: Number(size))
				}
			}`])
		size = PdfDriver.GetImageSize('200_jpg_content'.Repeat(100))
		Assert(size.width.Round(2) is: 3149.61)
		Assert(size.height.Round(2) is: 3149.61)

		filename = .MakeFile('300_jpg_file')
		filedata = GetFile(filename)
		size = PdfDriver.GetImageSize(filedata)
		Assert(size.width.Round(2) is: 4724.41)
		Assert(size.height.Round(2) is: 4724.41)
		}
	}