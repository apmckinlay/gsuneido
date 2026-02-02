// Copyright (C) 2020 Axon Development Corporation All rights reserved worldwide.
_CheckMarkBaseFormat
	{
	New(.width = 4)
		{
		super(width)
		.tf = _report.Construct(['Text', 'M', :width])
		}

	GetDefaultWidth()
		{
		return .isHtmlDriver?() ? .width : super.GetDefaultWidth()
		}

	GetSize(data = "")
		{
		return .isHtmlDriver?() ? .tf.GetSize() : super.GetSize(data)
		}

	Print(x, y, w, h)
		{
		if .isHtmlDriver?()
			.DoWithFont(false)
				{|font|
				_report.AddText('&#10004;', x, y, w, h, font, justify: 'center', html:,
					extra: #('font-size': '150%', 'font-style': 'normal',
						'font-weight': 'normal'))
				}
		else
			super.Print(x, y, w, h)
		}

	PrintWithBox(x, y, w, h, data)
		{
		if .isHtmlDriver?()
			.DoWithFont(false)
				{|font|
				_report.AddText(data is true ? '&#9745;' : '&#9744;', x, y, w, h, font
					justify: 'center', html:, extra: #('font-size': '150%',
						'font-style': 'normal', 'font-weight': 'normal'))
				}
		else
			super.PrintWithBox(x, y, w, h, data)
		}

	isHtmlDriver?()
		{
		return _report.Driver.Base?(HtmlDriver) or _report.Driver.Base?(SuJsHtmlDriver)
		}
	}