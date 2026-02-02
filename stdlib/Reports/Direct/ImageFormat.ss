// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Format
	{
	New(image = "", .scale = 1, .width = false, .height = false,
		.stretch = false, .background = false, .top = false, .maxWidth = false)
		{
		.data = image
		}

	aspectRatio: 1
	GetSize(data = "")
		{
		if Object?(.background)
			return #(w: 0, h: 0, d: 0)

		if .width isnt false and .height isnt false
			{
			width = .width
			height = .height
			}
		else
			{
			dataStr, textOnly? = .dataString(data)
			if textOnly?
				return .getTextSize(dataStr)

			if false is imageSizeOb = .imageSizeOb(dataStr)
				return .getFallbackSize()

			heightAndWidth = .heightAndWidthFromImageSizeOb(imageSizeOb)
			height = heightAndWidth.height
			width = heightAndWidth.width
			}
		d = .getDescent(height)
		adjust = _report.GetImageSizeAdjustment()
		return Object(w: width / adjust, h: height / adjust, :d)
		}

	getTextSize(dataStr)
		{
		w = .width is false ? 1 : .width
		dataStr = .printPath(dataStr)
		sz = WrapFormat(dataStr, :w).GetSize(dataStr)
		return Object(:w, h: sz.h, d: sz.d)
		}

	validImageSize?: ""
	imageSizeOb(data)
		{
		imageSizeOb = false
		if .data isnt data or
			.validImageSize? isnt false // Do not check image size again if it is invalid
			Image.RunWithErrorLog(
				{
				imageSizeOb = _report.GetImageSize(data)
				.aspectRatio = imageSizeOb.width / imageSizeOb.height
				})
		if .data is data // Don't set this flag unless it is for the original set .data
			.validImageSize? = imageSizeOb isnt false
		return imageSizeOb
		}

	getFallbackSize()
		{
		w = .width is false ? 1 : .width
		h = .height is false ? 1 : .height
		return Object(:w, :h, d: .getDescent(h))
		}

	heightAndWidthFromImageSizeOb(imageSizeOb)
		{
		height = imageSizeOb.height * .scale
		width = imageSizeOb.width * .scale
		if .height isnt false
			{
			width = ((width / height) * .height).Round(0)
			height = .height
			if .maxWidth isnt false and width > .maxWidth
				{
				height *=.maxWidth / width
				width = .maxWidth
				}
			}
		else if .width isnt false
			{
			height = ((height / width) * .width).Round(0)
			width = .width
			}
		return Object(:height, :width)
		}

	getDescent(height)
		{
		if .top // top-align images
			{
			font = _report.GetFont()
			lineSpecs = _report.GetLineSpecs(font)
			return height - (lineSpecs.height - lineSpecs.descent)
			}
		return 0
		}

	Print(x, y, w, h, data = "")
		{
		if Object?(.background)
			{
			x = .background.x.InchesInTwips()
			y = .background.y.InchesInTwips()
			w = .width.InchesInTwips()
			h = .height.InchesInTwips()
			}
		if .Xstretch isnt .Ystretch
			.stretch = true

		data, textOnly? = .dataString(data)
		if data is '' // no image to print
			return

		if textOnly? or not Image.RunWithErrorLog({ .print(x, y, w, h, data) })
			{
			path? = Paths.IsValid?(data)
			wrap = WrapFormat(path? ? .printPath(data) : '?')
			if path?
				Format.Hotspot(x, y, w, h, [],
					access: [control: "AttachmentGoTo", file: data])
			wrap.Print(x, y, w, h)
			}
		}

	dataString(data)
		{
		if not String?(data) or data is ""
			data = .data
		if data is ''
			return data, true
		if not Paths.IsValid?(data)
			return data, false
		fullPath = OpenImageWithLabelsControl.SplitFullPath(data)
		return fullPath, true
		}

	printPath(file)
		{
		return OpenImageWithLabelsControl.CopyAndLinkPath?(file)
			? Paths.Basename(file)
			: file
		}

	print(x, y, w, h, data)
		{
		if .stretch isnt true
			{
			if .aspectRatio < w / h
				w = .aspectRatio * h
			else
				h = w / .aspectRatio
			}
		_report.AddImage(x, y, w, h, data)
		}

	ExportCSV(data = '')
		{
		if .data isnt ""
			data = .data
		return .CSVExportString(data)
		}

	Variable?()
		{
		return true
		}
	}
