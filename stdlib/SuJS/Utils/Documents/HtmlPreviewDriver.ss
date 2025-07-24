// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
HtmlDriver
	{
	AddImage(x, y, w, h, data, origData = false)
		{
		if String?(data) and Paths.IsValid?(data)
			{
			if origData is false
				origData = data
			Format.Hotspot(x, y, w, h, [],
				access: [control: "AttachmentGoTo", file: origData])
			}
		super.AddImage(x, y, w, h, data)
		}

	AddPage(dimens)
		{
		super.AddPage(dimens)
		return .page = Object(:dimens)
		}

	EndPage()
		{

		}

	Getter_Page()
		{
		return .page
		}
	}
