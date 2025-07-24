// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
// abstract class for CheckMarkFormat and CheckBoxForamt
ImageFormat
	{
	New(width = 4)
		{
		super(top:, scale: .01)
		.tf = _report.Construct(['Text', :width])
		}

	GetDefaultWidth()
		{
		return .tf.GetDefaultWidth()
		}

	ImageText(name = 'checkmark', ext = false)
		{
		ext = ext is false ? _report.GetAcceptedImageExtension() : ext
		return Query1Cached('imagebook', name: name $ ext).text
		}

	GetSize(data /*unused*/ = "")
		{
		img = .ImageText(ext: '.jpg')
		size = super.GetSize(img)
		size.w = Max(size.w, .tf.GetSize().w)
		return size
		}

	marginLeft: .4
	marginTop: .1
	Print(x, y, w, h)
		{
		_report.DrawWithinClip(x, y, w, h)
			{
			img = .ImageText(ext: .getImageExt())
			.doWithSize(x, y, w)
				{ |left, top, width, height|
				super.Print(left, top, width, height, img)
				}
			}
		}

	PrintInvalidData(x, y, w, h, data)
		{
		if not String?(data)
			data = Display(data)
		.DoWithFont(false)
			{|font|
			_report.DrawWithinClip(x, y, w, h)
				{
				_report.AddText(data, x, y, w, h, font)
				}
			}
		}

	PrintWithBox(x, y, w, h /*unused*/, data = '')
		{
		.doWithSize(x, y, w)
			{ |left, top, width, height|
			if data is true
				{
				img = .ImageText(ext: .getImageExt())
				super.Print(left + 1, top + 1, width - 2, height - 2, img)
				}
			.drawbox(left, top, left + width - 1, top + height - 1)
			}
		}

	drawbox(left, top, right, bottom)
		{
		_report.AddLine(left, top, right, top, 1)
		_report.AddLine(right, top, right, bottom, 1)
		_report.AddLine(right, bottom, left, bottom, 1)
		_report.AddLine(left, bottom, left, top, 1)
		}

	doWithSize(x, y, w, block)
		{
		h = .tf.GetSize().h
		left = x + w * .marginLeft
		top = y + h * .marginTop
		height = width = h * (1 - .marginTop * 2)
		block(left, top, width, height)
		}

	getImageExt()
		{
		// gif is for printing on gdi driver, jpg is for pdf
		return _report.GetAcceptedImageExtension() is '.gif' ? '.emf' : '.jpg'
		}

	ExportCSV(data = "")
		{
		return .CSVExportString(data is true ? 'Yes' : 'No')
		}
	}