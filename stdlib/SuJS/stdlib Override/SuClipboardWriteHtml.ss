// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
function(html, text = false)
	{
	blob = SuUI.MakeWebObject(#Blob, [html], [type: "text/html"])
	item = ["text/html": blob]
	if text isnt false
		item["text/plain"] = text
	clipboardItem = SuUI.MakeWebObject(#ClipboardItem, item)
	SuUI.GetCurrentWindow().navigator.clipboard.Write([clipboardItem])
	}
