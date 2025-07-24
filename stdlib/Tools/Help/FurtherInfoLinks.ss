// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
HelpReferenceLinks
	{
	CallClass(title, howtos = #(), crossref = #(), trouble = #())
		{
		str = ''
		str $= '<p><b>For further information' $ Opt(' about ', title) $
			', see:</b></p>\n'
		str $= .AddCrossrefLinks(crossref, side?: false)
		str $= .AddHowDoILinks(howtos, side?: false)
		str $= .AddTroubleLinks(trouble, side?: false)
		return str
		}
	}