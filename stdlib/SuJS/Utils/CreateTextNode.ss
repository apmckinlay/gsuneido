// Copyright (C) 2018 Axon Development Corporation All rights reserved worldwide.
function (text, parent = false, className = false, _document = false, _at = false)
	{
	if document is false
		document = SuUI.GetCurrentDocument()
	el = document.CreateTextNode(text)

	if className isnt false
		el.className = className

	AttachElement(el, parent, at)
	return el
	}
