// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Name: 'MarkdownEdit'
	New()
		{
		.display = .FindControl('display')
		}

	Controls()
		{
		return Object('HorzSplit'
			#('ScintillaAddonsEditor', Addon_markdown:, name: 'edit')
			Object('MshtmlControl', style: MarkdownCSS, name: 'display'))
		}

	GetDisplay()
		{
		return .display
		}
	}