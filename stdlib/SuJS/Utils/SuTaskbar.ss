// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
class
	{
	styles: '
		.su-taskbar {
			position: fixed;
			left: 0px;
			right: 0px;
			bottom: 0px;
			height: calc(1.5rem + 10px);
			display: flex;
			background-color: var(--su-color-buttonface);
			border-top: 1px solid lightgrey;
			z-index: 1000;
		}
		.su-taskbar-item {
			border-right: 1px solid lightgrey;
			outline: none;
			cursor: pointer;
			padding: 5px;
			transition: 0.3s;
			font-size: 1rem;
			line-height: 1.5rem;
			overflow: hidden;
			white-space: nowrap;
			text-overflow: ellipsis;
			font-weight: bold;
			color: grey;
		}
		.su-taskbar-item:first-child {
			border-left: 1px solid lightgrey;
		}
		.su-taskbar-item:hover {
			background-color: #777;
			color: white;
		}
		.su-taskbar-item.selected {
			background-color: white;
			color: black;
		}
		@keyframes su-flashing {
			0% {background-color: var(--su-color-buttonface);}
			50% {background-color: orange;}
			100% {background-color: var(--su-color-buttonface);}
		}
		.su-taskbar-item-flashing {
			animation-name: su-flashing;
			animation-duration: 2s;
			animation-iteration-count: infinite;
		}'
	New()
		{
		LoadCssStyles('su-taskbar', .styles)
		body = SuUI.GetCurrentDocument().body
		.taskbar = CreateElement('div', parent: body, className: 'su-taskbar')
		.windows = Object()
		.flashingWindows = Object()
		}

	taskbarHeight: '0px'
	allMinimized: false
	Update(titles)
		{
		.taskbar.innerText = ''
		oldFlashingWindows = .flashingWindows
		.windows = Object()
		.flashingWindows = Object()
		.allMinimized = true
		for title in titles
			{
			if false is window = SuRender().GetRegisteredComponent(title.id)
				{
				SuRender().Event(false, 'ProgrammerError', [
					msg: 'SuTaskbar - window not found',
					params: titles.Copy().Add(title, at: 'title'),
					caughtMsg: 'ignored'])
				continue
				}
			item = CreateElement('div', parent: .taskbar,
				className: 'su-taskbar-item' $ (title.active? ? ' selected' : ''))
			item.innerText = title.title
			item.AddEventListener('click', .factory(title.id))
			item.AddEventListener('contextmenu', .contextmenuFactory(title.id))
			.windows[title.id] = item
			.SetFlashing(title.id, oldFlashingWindows.Has?(title.id), skipShowHide?:)
			if window.State isnt WindowPlacement.minimized
				.allMinimized = false
			}
		.updateShowHide()
		}

	updateShowHide()
		{
		if .windows.Size() <= 1 and .flashingWindows.Empty?() and .allMinimized is false
			{
			.taskbarHeight = '0px'
			.taskbar.SetStyle('display', 'none')
			}
		else
			{
			.taskbarHeight = '1.5rem + 10px'
			.taskbar.SetStyle("display", '')
			}
		.windows.Members().Each()
			{
			if false isnt window = SuRender().GetRegisteredComponent(it)
				window.UpdateMaximize()
			}
		}

	factory(id)
		{
		return {
			window = SuRender().GetRegisteredComponent(id)
			if Same?(SuRender().ActiveWindow, window)
				window.MINIMIZE()
			else
				SuRender().ActivateWindow(id)
			}
		}

	contextmenuFactory(id)
		{
		return {
			|event|
			item = event.target
			rect = SuRender.GetClientRect(item)
			SuRender().RunWhenNotFrozen(
				{ SuRender().Event(false, 'SuTaskbarContextMenu',
					[id, event.clientX, event.clientY, rect], showOverlay?:) })
			event.StopPropagation()
			event.PreventDefault()
			}
		}

	GetTaskbarHeight()
		{
		return .taskbarHeight
		}

	GetDimension()
		{
		return SuRender.GetClientRect(.taskbar)
		}

	SetFlashing(id, flashing?, skipShowHide? = false)
		{
		if not .windows.Member?(id)
			return
		if flashing?
			{
			.flashingWindows.AddUnique(id)
			.windows[id].classList.Add('su-taskbar-item-flashing')
			}
		else
			{
			.flashingWindows.Remove(id)
			.windows[id].classList.Remove('su-taskbar-item-flashing')
			}
		if not skipShowHide?
			.updateShowHide()
		}
	}
