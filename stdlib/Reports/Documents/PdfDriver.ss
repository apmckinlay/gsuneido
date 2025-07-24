// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
// NOTE: this record is getting large, should look for ways to split it up
Document_Builder
	{
	New(.filename, .compress = true)
		{
		.writer = .createWriter(filename)
		.pages = 0
		.pageRefs = Object()
		.embeddedFonts = Object()
		.docHeight = 792 // by default letter size
		.docWidth = 612
		.imageReferences = Object()

		.buildBaseObjects()
		}

	createWriter(filename)
		{
		return new PdfDriverFileWriter(filename)
		}

	AddPage(dimens)
		{
		.pages++
		.stream = ''
		.xObjects = Object()
		.setPageSize(dimens.width, dimens.height)
		}

	EndPage()
		{
		.buildPage()
		.writer.Flush()
		}

	// pdf uses PostScript point, which is 1/72 inch
	pspToTwip: 20
	setPageSize(w, h)
		{
		.docWidth 	= w	* PointsPerInch
		.docHeight 	= h * PointsPerInch
		}

	AddText(data, x, y, w, h /*unused*/, font, justify = 'left', ellipsis? = false,
		color = false)
		{
		fontRef = PdfFonts.ConvertFont(font)
		font = PdfFonts.GetCompatibleFont(font)
		if PdfFonts.IsStandardFont(fontRef) is false and
			not .embeddedFonts.Member?(font.name)
			.embeddedFonts[font.name] = false // false is a placeholder value
		data = PdfFonts.StripInvalidChars(data)
		size = font is false ? 10 : font.size
		if ellipsis? is true
			data = .pdfEllipsis(font, data, w)
		textWidth = .GetTextWidth(font, data)
		xAdjust = .justifyAdjust(textWidth, w, justify)
		data = .escapeSpecialCharacters(data)

		// print the text on the top of rectangle
		fontHeight = PdfFonts.GetFontHeight(font) * .afmToTwip
		.doWithColorAndBorder(0, color, false)
			{
			.stream $= "BT " $ fontRef $ " " $ size $
				" Tf " $ (x + xAdjust) / .pspToTwip $ " " $
				(.docHeight - (y + fontHeight) / .pspToTwip) $ " Td (" $ data $
				") Tj ET\n"
			.addUnderlineAndStrike(font, textWidth, x + xAdjust, y, fontHeight, size)
			}
		}

	escapeSpecialCharacters(data)
		{
		return data.Replace('\\', `\\134`).Replace('\(', `\\050`).Replace('\)', `\\051`)
		}

	addUnderlineAndStrike(font, textWidth, x, y, fontHeight, size)
		{
		if font.GetDefault('underline', false) is true
			{
			baseY = y + fontHeight + size
			.AddLine(x, baseY, x + textWidth, baseY, 9)
			}
		if font.GetDefault('strikeout', false) is true
			{
			harfY = y + fontHeight/2 + size * 4
			.AddLine(x, harfY, x + textWidth, harfY, 9)
			}
		}

	justifyAdjust(textSize, maxWidth, justify)
		{
		if justify is 'left'
			return 0
		else if justify is 'right'
			return maxWidth - textSize
		else if justify is 'center'
			return (maxWidth / 2) - (textSize / 2)
		else
			return 0
		}

	pdfEllipsis(font, data, w)
		{
		do
			{
			textSize = .pdfTextSize(font, data)
			if textSize <= w
				return data
			if not data.Suffix?('...')
				data = data[..-2] $ '...' // remove 2 chars to make up for ...
			data = data[..-4] $ '...' // remove 1 char
			} while data.Size() > 3
		return data
		}

	pdfTextSize(font, data)
		{
		w = .GetTextWidth(font, data)
		return w
		}

	DrawWithinClip(x, y, w, h, block)
		{
		// save current state
		.stream $= 'q\n'
		// set clipping path equal to the rectangle passed through so that
		// only the area defined by the rectangle will be painted
		.stream $= x / .pspToTwip $ " " $ (.docHeight - y / .pspToTwip) $
			" " $ w / .pspToTwip $ " " $ (-h / .pspToTwip) $ " re h\nW\nn\n"
		block()
		//restore the previous state
		.stream $= 'Q\n'
		}

	AddMultiLineText(data, x, y, w, h /*unused*/, font, justify = 'left', color = false)
		{
		font = PdfFonts.GetCompatibleFont(font)
		lineHeight = .GetLineSpecs(font).height
		for line in data.Split('\n')
			{
			.AddText(line, x, y, w, lineHeight, font, justify, :color)
			y += lineHeight
			}
		}

	AddPolygon(points, thick, fillColor = false, lineColor = false)
		{
		if points.Empty?()
			return
		.doWithColorAndBorder(thick, fillColor, lineColor)
			{
			.stream $= points[0].x / .pspToTwip  $ ' ' $
				(.docHeight - points[0].y / .pspToTwip) $ ' m\n'
			for (i = 1; i < points.Size(); i++)
				.stream $= points[i].x / .pspToTwip  $ ' ' $
					(.docHeight - points[i].y / .pspToTwip) $ ' l\n'
			.stream $= points[0].x / .pspToTwip  $ ' ' $
				(.docHeight - points[0].y / .pspToTwip) $ ' l\nh\n'
			}
		}

	AddArc(left, top, right, bottom, xStartArc /*unused*/, yStartArc /*unused*/,
		xEndArc /*unused*/, yEndArc /*unused*/, thick, lineColor = false)
		{
		left = left / .pspToTwip
		top = .docHeight - top / .pspToTwip
		right = right / .pspToTwip
		bottom = .docHeight - bottom / .pspToTwip
		//save current state
		.stream $= 'q\n'
		//set clipping path equal to half of the bounding rectangle
		//so that only half of the ellipse will be drawn on the page
		//which is the arc wanted
		.stream $= left $ ' ' $ top $ ' m\n'
		.stream $= right $ ' ' $ bottom $ ' l\n'
		.stream $=
			(left < right and bottom < top) or (left > right and bottom > top)
				? (left $ ' ' $ bottom $ ' l\n')
				: (right $ ' ' $ top $ ' l\n')
		.stream $= 'h\nW\nn\n'
		.doWithColorAndBorder(thick, false, lineColor)
			{
			.ellipse((left - right).Abs(), (top - bottom).Abs(), Min(left, right),
				Min(bottom, top))
			}
		//restore the previous state
		.stream $= 'Q\n'
		}

	AddLine(x, y, x2, y2, thick, color = 0x00000000)
		{
		x = x / .pspToTwip
		y = .docHeight - y / .pspToTwip
		x2 = x2 / .pspToTwip
		y2 = .docHeight - y2 / .pspToTwip
		.doWithColorAndBorder(thick, false, color)
			{
			.stream $= '1 J\n' // rounded cap
			.stream $=  x $ " " $ y $ " m " $ x2 $ " " $ y2 $ " l h "
			}
		}

	AddRect(x, y, w, h, thick, fillColor = false, lineColor = false)
		{
		x = x / .pspToTwip
		y = .docHeight - y / .pspToTwip
		.doWithColorAndBorder(thick, fillColor, lineColor)
			{
			.stream $=  x $ " " $ y $ " " $
				w / .pspToTwip $ " " $ (-h / .pspToTwip) $ " re h "
			}
		}

	AddRoundRect(x, y, w, h, width, height, thick, fillColor = false,
		lineColor = false)
		{
		x = x / .pspToTwip
		y = .docHeight - y / .pspToTwip
		w = w / .pspToTwip
		h = -h / .pspToTwip
		width = width / .pspToTwip
		height = -height / .pspToTwip
		.doWithColorAndBorder(thick, fillColor, lineColor)
			{
			.roundRect(x, y, x + w, y + h, width, height)
			}
		}


	doWithColorAndBorder(thick, fillColor, lineColor, block)
		{
		if fillColor isnt false
			{
			b = ((fillColor & 0xFF0000) / 0x10000 	/ 0xFF).Round(5)
			g = ((fillColor & 0x00FF00) / 0x100 	/ 0xFF).Round(5)
			r = ((fillColor & 0x0000FF)				/ 0xFF).Round(5)
			.stream $= r $ ' ' $ g $ ' ' $ b $ ' rg\n'
			}
		if lineColor isnt false
			{
			b = ((lineColor & 0xFF0000) / 0x10000 	/ 0xFF).Round(5)
			g = ((lineColor & 0x00FF00) / 0x100 	/ 0xFF).Round(5)
			r = ((lineColor & 0x0000FF)				/ 0xFF).Round(5)
			.stream $= r $ ' ' $ g $ ' ' $ b $ ' RG\n'
			}
		if thick isnt 0
			.stream $= (thick / 15).Round(5) $ ' w\n'
		block()
		.endColorAndBorder(thick, fillColor, lineColor)
		}

	endColorAndBorder(thick, fillColor, lineColor)
		{
		if thick isnt 0 and fillColor isnt false
			.stream $= "B\n"
		else if thick isnt 0
			.stream $= "S\n"
		else if fillColor isnt false
			.stream $= "f\n"

		if fillColor isnt false
			.stream $= '0 0 0 rg\n'
		if lineColor isnt false
			.stream $= '0 0 0 RG\n'
		}

	AddEllipse(x, y, w, h, thick, fillColor = false, lineColor = false)
		{
		x = x / .pspToTwip
		y = .docHeight - y / .pspToTwip
		w = w / .pspToTwip
		h = -h / .pspToTwip
		.doWithColorAndBorder(thick, fillColor, lineColor)
			{
			.ellipse(w, h, x, y)
			}
		}

	AddCircle(x, y, radius, thick, fillColor = false, lineColor = false)
		{
		x = x / .pspToTwip
		y = .docHeight - y / .pspToTwip
		radius = radius / .pspToTwip
		.doWithColorAndBorder(thick, fillColor, lineColor)
			{
			.circle(radius, x, y)
			}
		}
	// using bezier to draw cirlce
	// Reference: http://spencermortensen.com/articles/bezier-circle/
	circle(radius, x, y)
		{
		distance = radius * 0.55191502449
		.stream $= (x - radius) $ " " $ y $ " m\n"
		.stream $= (x - radius) $ " " $ (y + distance) $ " " $
			(x - distance) $ " " $ (y + radius) $ " " $
			x $ " " $ (y + radius) $ " c\n"
		.stream $= (x + distance) $ " " $ (y + radius) $ " " $
			(x + radius) $ " " $ (y + distance) $ " " $
			(x + radius) $ " " $ y $ " c\n"
		.stream $= (x + radius) $ " " $ (y - distance) $ " " $
			(x + distance) $ " " $ (y - radius) $ " " $
			x $ " " $ (y - radius) $ " c\n"
		.stream $= (x - distance) $ " " $ (y - radius) $ " " $
			(x -radius) $ " " $ (y - distance) $ " " $
			(x - radius) $ " " $ y $ " c\n"
		}

	roundRect(x1, y1, x2, y2, width, height)
		{
		distanceX = width * 0.2761423749154
		distanceY = height * 0.2761423749154
		p1 = Object(x: x1, y: y1 + height / 2)
		p2 = Object(x: x1, y: y1 + height / 2 - distanceY)
		p3 = Object(x: x1 + width / 2 - distanceX, y: y1)
		p4 = Object(x: x1 + width / 2, y: y1)
		p5 = Object(x: x2 - width / 2, y: y1)
		p6 = Object(x: x2 - width / 2 + distanceX, y: y1)
		p7 = Object(x: x2, y: y1 + height / 2 - distanceY)
		p8 = Object(x: x2, y: y1 + height / 2)
		p9 = Object(x: x2, y: y2 - height / 2)
		p10 = Object(x: x2, y: y2 - height / 2 + distanceY)
		p11 = Object(x: x2 - width / 2 + distanceX, y: y2)
		p12 = Object(x: x2 - width / 2, y: y2)
		p13 = Object(x: x1 + width / 2, y: y2)
		p14 = Object(x: x1 + width / 2 - distanceX, y : y2)
		p15 = Object(x: x1, y: y2 - height / 2 + distanceY)
		p16 = Object(x: x1, y: y2 - height / 2)
		.stream $= p1.x $ " " $ p1.y $ " m\n"
		.stream $= p2.x $ " " $ p2.y $ " " $ p3.x $ " " $ p3.y $
			" " $ p4.x $ " " $ p4.y $ " c\n"
		.stream $= p5.x $ " " $ p5.y $ " l\n"
		.stream $= p6.x $ " " $ p6.y $ " " $ p7.x $ " " $ p7.y $
			" " $ p8.x $ " " $ p8.y $ " c\n"
		.stream $= p9.x $ " " $ p9.y $ " l\n"
		.stream $= p10.x $ " " $ p10.y $ " " $ p11.x $ " " $ p11.y $
			" " $ p12.x $ " " $ p12.y $ " c\n"
		.stream $= p13.x $ " " $ p13.y $ " l\n"
		.stream $= p14.x $ " " $ p14.y $ " " $ p15.x $ " " $ p15.y $
			" " $ p16.x $ " " $ p16.y $ " c\n"
		.stream $= p1.x $ " " $ p1.y $ " l\nh\n"
		}

	ellipse(w, h, x, y)
		{
		distanceX = w * 0.2761423749154
		distanceY = h * 0.2761423749154
		centreX = x + w / 2
		centreY = y + h / 2
		.stream $= x $ " " $ centreY $ " m\n"
		.stream $= x $ " " $ (centreY + distanceY) $ " " $
			(centreX - distanceX) $ " " $ (y + h) $ " " $ centreX $ " " $
			(y+h) $ " c\n"
		.stream $= (centreX + distanceX) $ " " $ (y + h) $ " " $
			(x + w) $ " " $ (centreY + distanceY) $ " " $ (x + w) $ " " $
			centreY $ " c\n"
		.stream $= (x + w) $ " " $ (centreY - distanceY) $ " " $
			(centreX + distanceX) $ " " $ y $ " " $ centreX $ " " $
			y $ " c\n"
		.stream $= (centreX - distanceX) $ " " $ y $ " " $
			x $ " " $ (centreY - distanceY) $ " " $ x $ " " $
			centreY $ " c\n"
		}

	/*	buildPageRefString
	*	Builds a string of page references that are used by the Pages object
	*	(object 2) to tell the document where each Page object is located.
	*	Returns: A string of formatted page references ("5 0 R 8 0 R 11 0 R")
	*/
	buildPageRefString()
		{
		references = ""
		for x in .pageRefs
			{
			references = references $ x $ " "
			}
		return references
		}
	/*	buildPage
	*	Responsible for creating each set of page, stream, and stream count
	*	objects for each physical page in the PDF. Applies the stream associated
	*	with this page to ensure that all text and line data is placed on the
	*	appropriate page.
	*	Params:
	*		i: the current page to be created
	*	Returns: An object containing all PDF object strings to be written to
	*	the main PDF document string.
	*/
	buildPage()
		{
		fontReferences = ''
		for font in .embeddedFonts.Members()
			{
			if .embeddedFonts[font] is false
				{
				f = PdfFonts.FontNames.FindOne({ it.name is font })
				Assert(f isnt false)
				.embeddedFonts[f.name] = [ref: f.ref, id: .writer.NextObjId]
				PdfFonts.BuildEmbeddedFont(f, .writer)
				}
			ob = .embeddedFonts[font]
			fontReferences $= '/' $ ob.ref $ " " $ ob.id $ " 0 R\n"
			}

		xObjRefs = .buildXObjectRef()
		.pageRefs.Add(.writer.NextObjId $ " 0 R")
		resources = "<</Type /Page /Parent 2 0 R /MediaBox [0 0 " $
			.docWidth $ " " $ .docHeight $"] " $
			"/Contents ["$ (.writer.NextObjId + 1) $" 0 R] /Resources<<\n" $
			"/Font <<\n" $ fontReferences $ ">>\n" $
			xObjRefs $
			">>>>\n" $
			"endobj\n"
		.writer.AddObject(resources)
		stream = .compress ? Zlib.Compress(.stream) : .stream
		.writer.AddObject("\n<</Length "$ (.writer.NextObjId + 1) $ " 0 R " $
				( .compress ? "/Filter /FlateDecode" : "" ) $ ">>\n" $
			"stream\n" $ stream $ "\n" $
			"endstream\n" $
			"endobj\n")
		.writer.AddObject("\n" $ stream.Size() $"\nendobj\n")
		}
	buildXObjectRef()
		{
		xObjRefs = ""
		if not .xObjects.Empty?() //If there are XObjects add them
			{
			xObjRefs = "/XObject<<\n"
			for (c = 0; c < .xObjects.Size(); c++)
				{
				contents = .xObjects[c] $ "\n" $
					"endobj\n"

				hash = Adler32()
				hash.Update(contents)
				imgCksum = hash.Value().Hex()
				if not .imageReferences.Member?(imgCksum)
					.imageReferences[imgCksum] = .writer.AddObject("\n" $ contents)
				xObjRefs $= .imgObName(c) $ " " $ .imageReferences[imgCksum] $ " 0 R\n"
				}
			xObjRefs $= ">>"
			}
		return xObjRefs
		}
	/*	buildBaseObjects
	*	Builds the required objects for the PDF document that remain (mostly)
	*	static from one document to the next. This includes the version,
	*	the catalog, the "pages" reference object, and fonts
	* 	It then proceeds to build each page
	*	in the document.
	*	Returns: An object containing all the PDF object strings for use with an
	*	xref builder.
	*/
	buildBaseObjects()
		{
		bin = .compress ?  "%\xe9\xe9\xe9\xe9" : ""

		.writer.AddObject("%PDF-1.3\n" $ bin $ "\n")
		.writer.AddObject("<</Type /Catalog /Pages 2 0 R>>\nendobj\n")
		.writer.Reserve() // for pages

		encoding = " /Encoding /WinAnsiEncoding >>\nendobj\n"
		fontDef = "<</Type /Font /Subtype /Type1 /BaseFont /"
		for f in PdfFonts.FontNames
			{
			if f.standard is true
				{
				id = .writer.AddObject(fontDef $ f.name $ encoding)
				.embeddedFonts[f.name] = [ref: f.ref, :id]
				}
			}
		.writer.AddObject("<< /Producer (Suneido PDF Generator) >>\nendobj\n")
		}

	BuildXRef(size, locs, root, totalLength) // also used by PdfMerger
		{
		xref = "xref\n" $ "0 " $ size $"\n"
		for x in locs
			{
			locString = "0".Repeat(10 - Display(x).Size())  /*= pdf xref format */
			locString = locString $ x
			if x is 0
				xref $= locString $ " 65535 f \n"
			else
				xref $= locString $ " 00000 n \n"
			}
		xref $= "trailer <</Size " $ size $"/Root " $ root $ " 0 R>>\n" $
			"startxref\n" $ Display(totalLength) $ "\n" $
			"%%EOF"
		return xref
		}

	EnsureFont(font, oldfont)
		{
		return EnsurePDFFont(font, oldfont)
		}

	GetTextWidth(font, text)
		{
		font = PdfFonts.GetCompatibleFont(font)
		totalWidth = 0
		for c in text
			totalWidth += PdfFonts.GetCharWidth(font, c)
		return totalWidth * font.size * .afmToTwip
		}

	afmToTwip: .02 // afm Spec uses measurements in 1 1000th PSP. This number
				   // Converts from that measurement to Twips.

	GetLineSpecs(font)
		{
		// calculating line height
		descender = Abs(PdfFonts.GetFontDescender(font)) * .afmToTwip
		ascender = PdfFonts.GetFontHeight(font) * .afmToTwip
		return Object(height: descender + ascender, descent: descender)
		}
	GetDefaultFont()
		{
		return #(angle: 0, size: 10, weight:  400, name: "Helvetica", italic: false)
		}
	GetCharWidth(width, font, widthChar)
		{
		if width is false
			return 0
		font = PdfFonts.GetCompatibleFont(font)
		return width * PdfFonts.GetCharWidth(font, widthChar) *
			font.size * .afmToTwip
		}

	GetTextHeight(data, lineHeight)
		{
		numrows = data.Lines().Size()
		return lineHeight * numrows
		}

	/*	AddImage
	*	Used to add an image to the current page. Uses the image string, and
	*	should attempt to shrink the image down to fit at the coordinates given.
	*	Draws the image using the coordinates as the bottom-left point for the
	* 	image.
	*	NOTE: ONLY SUPPORTS .JPEG FILES AT THIS TIME.
	*	ANY OTHER IMAGES WILL RESULT IN ERRORS.
	*	Params:
	*		imageString: A string of binary data containing information about
	*			the image.
	*		xStart: The starting x coordinate for the image
	*		yStart: The starting y coordinate for the image
	*		sWidth: The preferred scaled width of the image
	*		sHeight: The preferred scaled height of the image
	*		scale: The preferred scale to multiply the height/width values
	*			by for scaling purposes. Normal size is a scale value of 1.
	*/
	AddImage(x, y, w, h, data)
		{
		if Paths.IsValid?(data)
			{
			if not Jpeg.ValidExtension?(data)
				throw Jpeg.InvalidExtension
			if false is data = ImageHandler.GenerateThumbnail(data)
				throw "Image: Couldn't generate thumbnail"
			}
		jpeg = .jpeg(data)
		//Height and width must be integers to work in Microsoft Edge
		imgSize = .GetImageSize(data, jpeg)
		width = (imgSize.width / .pspToTwip * 1.27).Round(0)
		height = (imgSize.height / .pspToTwip * 1.27).Round(0)
		head = .BuildImageObjectHead(width, height, jpeg.GetColorSpace(), data.Size())
		.xObjects.Add(
			head $ "\n" $
			"stream\n" $
			data $ "\n" $
			"endstream")
		x = (x / .pspToTwip).Round(5)
		y = (.docHeight - y / .pspToTwip - h / .pspToTwip).Round(5)
		w = (w / .pspToTwip).Round(5)
		h = (h / .pspToTwip).Round(5)
		.stream $= "q\n" $
			w $ " 0 0 " $ h $ " " $ x $ " " $ y $ " cm\n" $
			.imgObName(.xObjects.Size() -1) $ " Do\n" $
			"Q\n"
		}

	jpeg(data)
		{
		return new Jpeg(data)
		}

	BuildImageObjectHead(width, height, colorSpace, length, additionalEntries = #())
		{
		if not additionalEntries.Empty?()
			additionalEntries = additionalEntries.Copy().Map({ `/` $ it })
		return "<</Type /XObject\n" $
			"/Subtype /Image\n" $
			"/Width " $ width $ "\n" $
			"/Height " $ height $ "\n" $
			"/ColorSpace /" $ colorSpace $ "\n" $
			"/BitsPerComponent 8\n" $
			"/Length " $ length $ "\n" $
			"/Filter /DCTDecode\n" $
			Opt(additionalEntries.Join("\n"), "\n") $
			">>"
		}

	imgObName(idx)
		{
		return "/Im" $ idx
		}

	GetImageSizeAdjustment()
		{
		return 1.04
		}

	GetImageSize(data, jpeg = false)
		{
		w = h = 0
		if jpeg is false
			{
			size = ImageHandler.GetWidthHeight(data)
			w = size.w
			h = size.h
			}
		else
			{
			w = jpeg.GetWidth()
			h = jpeg.GetHeight()
			}
		return Object(height: h * .pspToTwip / 1.27, width: w * .pspToTwip / 1.27)
		}

	GetAcceptedImageExtension()
		{
		return '.jpg'
		}

	Finish(status)
		{
		if status is ReportStatus.SUCCESS
			{
			.writer.AddObject("<</Type /Pages /Kids [" $
				.buildPageRefString() $"] /Count " $
				.pages $ ">>\n" $
				"endobj\n", id: 2)
			locations = .writer.Locations
			xref = .BuildXRef(locations.Size(), locations, '1', .writer.TotalLength)
			.writer.Write(xref)
			.writer.Finish()
			}
		else
			.writer.Abort()
		return status
		}
	}