// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
Singleton
	{
	printStyles: `
		body {
			margin: 0;
		}
		body>div:not(.su-printer),
		dialog {
			display: none;
		}
		.su-printer { /* to remove the gap between its children */
			display: flex;
			flex-direction: column;
		}`
	screenStyles: '
		.su-printer {
			display: none;
		}'
	New()
		{
		LoadCssStyles('su-printer-print', .printStyles, media: 'print')
		LoadCssStyles('su-printer-screen', .screenStyles, media: 'screen')
		.tasks = Object().Set_default(Object())
		}

	printEl: false
	AddPage(page, id)
		{
		.tasks[id].Add(page)
		}

	Cancel(id)
		{
		.tasks.Delete(id)
		}

	Print(id)
		{
		if .tasks[id].Empty?()
			return

		pages = .tasks[id]
		.printEl = CreateElement('div', SuUI.GetCurrentDocument().body,
			className: 'su-printer')

		for page in pages
			{
			pageEl = CreateElement('svg', .printEl, namespace: SvgDriver.Namespace)
			pageEl.style.width = SvgDriver.ConvertToPixcel(page.dimens.width, 1)
			pageEl.style.height = SvgDriver.ConvertToPixcel(page.dimens.height, 1)
			driver = SvgDriver(pageEl, page.dimens)
			for cmd in page.Values(list:)
				{
				if not driver.Member?(cmd[0])
					continue
				(driver[cmd[0]])(@+1cmd)
				}
			}

		width = pages[0].dimens.width
		height = pages[0].dimens.height
		pageCss = Object(
			size: width $ 'in ' $ height $ 'in'
			margin: 0)
		LoadCssStyles('su-printer-page',
			'@page {' $ pageCss.Map2({ |m, v| m $ ': ' $ v $ ';' }).Join('\n') $ '}',
			override?:, media: 'both')
		SuDelayed(0, { .print(id) })
		}

	print(id)
		{
		SuUI.GetCurrentWindow().Print()
		.printEl.Remove()
		.printEl = false
		.tasks.Delete(id)
		}
	}
