// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
HtmlDriver
	{
	AddImage(x, y, w, h, data)
		{
		if String?(data) and Paths.IsValid?(data)
			Format.Hotspot(x, y, w, h, [],
				access: [control: "AttachmentGoTo", file: data])
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
