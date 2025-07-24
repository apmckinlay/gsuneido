// Copyright (C) 2019 Axon Development Corporation All rights reserved worldwide.
Document_Builder
	{
	New()
		{
		.paintContainer = Object()
		}

	defaultPaint:
		#(data: '', font: false, justify: 'left', color: false, ellipsis?: false,
			html: false, type: 'text')
	Print(block)
		{
		.paintContainer = Object()
		block()
		if .paintContainer.Size(list:) <= 1
			return .paintContainer.GetDefault(0, .defaultPaint.Copy())
		.paintContainer.type = 'multi'
		return .paintContainer
		}

	SetMultiPartsRatio(ratios)
		{
		.paintContainer.ratios = ratios
		}

	NoWrapSupported?: true

	AddText(data, x/*unused*/, y/*unused*/, w/*unused*/, h/*unused*/, font = false,
		justify = 'left', ellipsis? = false, color = false, html = false, extra = false)
		{
		ob = Object(:data, type: 'text')
		if font isnt false
			ob.font = font
		if justify isnt 'left'
			ob.justify = justify
		if ellipsis? isnt false
			ob.ellipsis? = ellipsis?
		if color isnt false
			ob.color = ToCssColor(color)
		if html isnt false
			ob.html = html
		if extra isnt false
			ob.extra = extra
		return .paintContainer.Add(ob)
		}

	GetLineSpecs(@unused)
		{
		return Object(height: 1, descent: 1)
		}

	GetCharWidth(@unused)
		{
		return 1
		}

	GetTextWidth(@unused)
		{
		return 1
		}

	AddImage(x/*unused*/, y/*unused*/, w/*unused*/, h/*unused*/, data)
		{
		if Paths.IsValid?(data)
			{
			if not Jpeg.ValidExtension?(data)
				throw Jpeg.InvalidExtension
			if false is data = ImageHandler.GenerateThumbnail(data)
				throw "Image: Couldn't generate thumbnail"
			}

		new Jpeg(data) // verify that the image is a valid Jpeg
		.paintContainer.Add(Object(src: 'data:image/jpeg;base64,' $ Base64.Encode(data),
			type: 'image'))
		}

	GetImageSize(data, jpeg = false)
		{
		if jpeg is false
			{
			if Paths.IsValid?(data)
				data = GetFile(data)
			jpeg = new Jpeg(data)
			}
		w = jpeg.GetWidth()
		h = jpeg.GetHeight()
		return Object(height: h, width: w)
		}

	AddRect(x/*unused*/, y/*unused*/, w/*unused*/, h/*unused*/, thick, fillColor = false,
		lineColor = false)
		{
		if fillColor isnt false
			fillColor = ToCssColor(fillColor)
		if lineColor isnt false
			lineColor = ToCssColor(lineColor)
		return .paintContainer.Add(Object(:thick, :fillColor, :lineColor, type: 'rect'))
		}

	AddCircle(x/*unused*/, y/*unused*/, radius/*unused*/, thick, fillColor = false,
		lineColor = false)
		{
		if fillColor isnt false
			fillColor = ToCssColor(fillColor)
		if lineColor isnt false
			lineColor = ToCssColor(lineColor)
		return .paintContainer.Add(Object(:thick, :fillColor, :lineColor, type: 'circle'))
		}
	}
