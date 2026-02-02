// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
Component
	{
	Name: "PreviewPage"
	New(.scale, dimens)
		{
		.CreateElement('div')
		.El.SetStyle('background-color', 'white')
		.svg = CreateElement('svg', .El, namespace: SvgDriver.Namespace)
		.pages = Object()
		.pg = false; // the page currently displayed
		if dimens isnt false
			.updateSize(dimens.width, dimens.height)
		}

	scale: 1
	UpdateScale(scale)
		{
		if scale is .scale
			return
		.scale = scale
		page = .pages[.pg]
		.updateSize(page.dimens.width, page.dimens.height)
		}

	DisplayPage(pg)
		{
		if not .pages.Member?(pg)
			{
			SuRender().Event(false, 'SuneidoLog', Object(
				'ERROR: (CAUGHT) PreviewPageComponent DisplayPage with invalid pg',
				params: [:pg, pages: .pages.Members(), npages: .pages.Size()],
				caughtMsg: 'for debug'))
			.Event('InvalidPG', pg)
			return
			}
		Assert(.pages hasMember: pg, msg: 'PreviewPageComponent.DisplayPage')
		page = .pages[.pg = pg]
		driver = SvgDriver(.svg, page.dimens)
		.updateSize(page.dimens.width, page.dimens.height)
		for cmd in page.Values(list:)
			{
			if not driver.Member?(cmd[0])
				continue
			(driver[cmd[0]])(@+1cmd)
			}
		for idx in page.accessPoints.Members()
			{
			point = page.accessPoints[idx]
			driver.AddAccessPoint(point.left, point.top,
				point.right - point.left, point.bottom - point.top, idx,
				callFn: .accessPointClicked)
			}
		}

	accessPointClicked(event)
		{
		target = event.target
		idx = Number(target.GetAttribute('data-idx'))
		.EventWithOverlay('AccessPointClicked', idx)
		}

	SetPage(pg, page)
		{
		Assert(.pages hasntMember: pg)
		.pages[pg] = page
		}

	updateSize(width, height)
		{
		.Xmin = SvgDriver.ConvertToPixcel(width, .scale)
		.Ymin = SvgDriver.ConvertToPixcel(height, .scale)
		.svg.style.width = .Xmin
		.svg.style.height = .Ymin
		.SetMinSize()
		}
	}
