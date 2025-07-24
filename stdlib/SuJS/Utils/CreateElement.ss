// Copyright (C) 2018 Axon Development Corporation All rights reserved worldwide.
function (tag, parent = false, className = false, _document = false, _at = false,
	namespace = false)
	{
	if document is false
		document = SuUI.GetCurrentDocument()
	el = namespace is false
		? document.CreateElement(tag)
		: document.CreateElementNS(namespace, tag)

	if className isnt false
		el.className = className

	AttachElement(el, parent, at)
	return el
	}
