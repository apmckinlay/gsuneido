// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
function(html, add?/*unused*/ = false, text = false)
	{
	SuRenderBackend().RecordAction(false, #SuClipboardWriteHtml, [html, :text])
	}
