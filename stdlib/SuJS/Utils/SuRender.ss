// Copyright (C) 2019 Axon Development Corporation All rights reserved worldwide.
class
	{
	Init(token)
		{
		Suneido.SuRender = new this(token)
		}
	CallClass()
		{
		return Suneido.SuRender
		}
	connectid: 0
	socketStatus: 0
	mousemoveCB: false
	mouseupCB: false
	delayedTask: false
	delayedTaskId: 0
	nextAck: false
	initialized: false
	New(.token)
		{
		.socketStatus = SuSocketStatus.init
		.components = Object()
		.eventToAck = Object()
		.status = Object()
		.Timers = Object()
		.windows = Object()
		.dropFilesList = Object()
		.keydownListeners = Object()
		.iframes = Object()
		.recorder = EventRecorder()

		.connect([:token, connectid: .connectid])

		SuUI.GetCurrentDocument().AddEventListener('keydown', .keydown, useCapture:)
		SuUI.GetCurrentDocument().AddEventListener('keydown', .keydown2)
		SuUI.GetCurrentDocument().AddEventListener('contextmenu', .disableDefaultMenu)
		SuUI.GetCurrentWindow().AddEventListener('beforeunload', .beforeUnload)
		SuUI.GetCurrentWindow().AddEventListener('error', .onError)

		SuUI.GetCurrentDocument().AddEventListener('contextmenu',
			.cancelMouseTracking, useCapture:)
		SuUI.GetCurrentWindow().AddEventListener("dragstart", .cancelMouseTracking)
		SuUI.GetCurrentDocument().AddEventListener("visibilitychange",
			.onvisibilitychange)

		.resizeObserver = SuUI.MakeWebObject('ResizeObserver', .onResize)
		.resizeObserver.Observe(SuUI.GetCurrentDocument().body)

		SuUI.GetCurrentDocument().AddEventListener('mousemove', .mousemove)
		SuUI.GetCurrentDocument().AddEventListener('mouseup', .mouseup)

		Global('SuKillTimer') // preload the record

		// this is to disable file drag and drop on the page
		// ImageComponent handles file drag and drop by its own
		.NotAllowDragAndDrop(SuUI.GetCurrentWindow())

		.overlay = SuOverlay()
		.Taskbar = SuTaskbar()
		.Notification = SuNotification(.Taskbar)
		.SnapManager = SnapManager()
		.heartbeat = SuHeartbeat()

		.init()
		.initialized = true
		}

	connect(query)
		{
		host = (SuUI.GetCurrentWindow())['location']['hostname']
		port = (SuUI.GetCurrentWindow())['location']['port']
		.ws = SuUI.WebSocketClient((IsHttps?() ? `wss://` : `ws://`) $
			host $ ':' $ port $ `/connect` $ Url.BuildQuery(query))
		.ws.AddEventListener(#open, .open)
		.ws.AddEventListener(#message, .message)
		.ws.AddEventListener(#close, .close)
		.ws.AddEventListener(#error, .error)
		.logWS('create')
		}

	logWS(event)
		{
		event = event $ ' (' $ .connectid $ ')'
		.recorder.Add(#ws, event)
		Print(event)
		}

	Engine: 'unknown'
	init()
		{
		SuUI.GetCurrentDocument().body.style.SetProperty(
			'--su-color-buttonface', ToCssColor(CLR.ButtonFace))
		SuUI.GetCurrentDocument().body.style.SetProperty(
			'background-color', 'darkgray')
		.SetWindowHeaderColor()

		userAgent = (SuUI.GetCurrentWindow())['navigator']['userAgent']
		if userAgent.Has?('Chrome')
			.Engine = 'Blink'
		else if userAgent.Has?('AppleWebKit')
			.Engine = 'WebKit'
		}

	SetWindowHeaderColor(color = false)
		{
		SuUI.GetCurrentDocument().body.style.SetProperty(
			'--su-color-windowheader', ToCssColor(color is false ? 'lightgray': color))
		}

	disableDefaultMenu(event)
		{
		event.PreventDefault()
		}

	cancelMouseTracking(event)
		{
		if .mouseupCB isnt false
			(.mouseupCB)(event)
		}

	onvisibilitychange(@unused)
		{
		.syncVisibility()
		}

	syncVisibility()
		{
		state = SuUI.GetCurrentDocument().visibilityState
		.sendWithAck(SuMessageFormatter.Type.SyncVisibility, eventId: .eventId++,
			arg1: state)
		}

	onResize(@unused)
		{
		.syncWindowSize()
		}

	syncWindowSize()
		{
		width = SuUI.GetCurrentWindow().innerWidth
		height = SuUI.GetCurrentWindow().innerHeight
		.sendWithAck(SuMessageFormatter.Type.UpdateStatus, eventId: .eventId++,
			arg1: #dimension, arg2: Object(:width, :height))
		}

	Heartbeat()
		{
		eventId = .eventId++
		.sendWithAck(SuMessageFormatter.Type.Heartbeat, :eventId)
		return eventId
		}

	close(event)
		{
		.logWS('close ' $ event.code $ Opt(' ', event.reason))
		Print(close: event)
		++.connectid
		.pause()
		reason = event.reason
		if event.code is 1006/*=Abnormal Closure*/ or
			reason.Has?('CloudFlare')
			{
			if .retryCount < .maxRetries
				{
				.tryReconnect()
				return
				}
			else
				reason = 'Reached maximum reconnect retries' $ Opt('(', reason, ')')
			}

		.socketStatus = SuSocketStatus.terminated
		SuUI.GetCurrentWindow().Alert(
			event.code isnt 1006/*=Abnormal Closure*/ and reason is ''
				? "You are logged out"
				: "Lost connection")

		.Reload()
		}

	Reload()
		{
		href = SuUI.GetCurrentWindow().location.href
		if not href.Has?('preauth=true')
			SuUI.GetCurrentWindow().location.Reload()
		else
			SuUI.GetCurrentWindow().location = href.BeforeFirst('/?')
		}

	pause()
		{
		.heartbeat.Stop()
		}

	retryCountDown: false
	timer: false
	retryCount: 0
	maxRetries: 10
	tryReconnect()
		{
		.socketStatus = SuSocketStatus.reconnecting

		.overlay.Show(id: #reconnect, level: 999)
		if .retryCount is 0
			{
			.retryCountDown = 0
			.doReconnect()
			return
			}
		.retryCountDown = 5
		.timer = SuSetTimer(NULL, 0, 1.SecondsInMs(), .doReconnect)
		}

	doReconnect(@unused)
		{
		if .retryCountDown > 0
			{
			.overlay.SetMsg('Network problems, trying to reconnect in ' $
				.retryCountDown $ ' sec(s)', id: #reconnect)
			.retryCountDown--
			return
			}
		if .timer isnt false
			{
			SuKillTimer(NULL, .timer)
			.timer = false
			}
		.showReconnecting()
		.retryCount++
		.connect([token: .token, reconnect:, connectid: .connectid])
		}

	showReconnecting()
		{
		// Not show the message in the first try
		if 0 isnt .retryCount
			.overlay.SetMsg('Network problems, reconnecting...', id: #reconnect)
		}

	closeReconnecting()
		{
		.overlay.Close(id: #reconnect)
		}

	error(event/*unused*/)
		{
		.socketStatus = SuSocketStatus.error
		.pause()
		.logWS('error')
		}

	NotAllowDragAndDrop(window)
		{
		window.AddEventListener("dragenter", .dragDropHandler)
		window.AddEventListener("dragover", .dragDropHandler)
		window.AddEventListener("drop", .dragDropHandler)
		}
	dragDropHandler(event)
		{
		event.PreventDefault()
		event.dataTransfer.dropEffect = "none"
		}

	RegisterIframe(el)
		{
		.iframes.Add(el)
		}

	UnregisterIframe(el)
		{
		.iframes.Remove(el)
		}

	freezeIframes()
		{
		for iframe in .iframes
			iframe.SetStyle('pointer-events', 'none')
		}

	restoreIframes()
		{
		for iframe in .iframes
			iframe.SetStyle('pointer-events', '')
		}

	lastMoveEvent: false
	mousemove(event)
		{
		.lastMoveEvent = event
		if .mousemoveCB is false
			return
		(.mousemoveCB)(event)
		}

	GetCursorPos()
		{
		return .lastMoveEvent isnt false
			? [x: .lastMoveEvent.clientX, y: .lastMoveEvent.clientY]
			: [x: 0, y: 0]
		}

	mouseup(event)
		{
		if .mouseupCB is false or event.button isnt 0
			return
		(.mouseupCB)(event)
		}

	SetMouseMoveCB(mousemoveCB)
		{
		if .mousemoveCB isnt false
			.ClearMouseMoveCB()
		.mousemoveCB = mousemoveCB
		}

	ClearMouseMoveCB()
		{
		.mousemoveCB = false
		}

	SetMouseUpCB(mouseupCB)
		{
		if .mouseupCB isnt false
			.ClearMouseUpCB()
		.freezeIframes()
		.mouseupCB = mouseupCB
		}

	ClearMouseUpCB()
		{
		.restoreIframes()
		.mouseupCB = false
		}

	Register(uniqueId, component)
		{
		id = Display(uniqueId)
		.components[id] = component
		}

	UnRegister(uniqueId)
		{
		.components.Delete(Display(uniqueId))
		}

	GetRegisteredComponent(uniqueId)
		{
		id = Display(uniqueId)
		return .components.GetDefault(id, false)
		}

	ActiveWindow: false
	keydown(event)
		{
		if .disableAccelCmd?()
			{
			event.PreventDefault()
			event.StopPropagation()
			return
			}

		if .lastMoveEvent isnt false
			.cancelMouseTracking(.lastMoveEvent)

		.keydownListeners.Each({ it(event) })
		if .ActiveWindow is false or .ActiveWindow.Destroyed?()
			return

		if event.key is 'Escape' and (.ActiveWindow.Base?(ModalWindowComponent) or
			.ActiveWindow.Base?(ListEditWindowComponent))
			{
			.ActiveWindow.On_Cancel()
			event.PreventDefault()
			event.StopPropagation()
			}

		ctrlKey = event.GetDefault(#ctrlKey, false)
		altKey = event.GetDefault(#altKey, false)
		shiftKey = event.GetDefault(#shiftKey, false)
		if false is cmd = .ActiveWindow.GetAccelCmd(ctrlKey, altKey, shiftKey, event.key)
			return
		preventDefault? = not (ctrlKey is true and event.key is 'c')
		if preventDefault?
			{
			event.PreventDefault()
			event.StopPropagation()
			}
		.ActiveWindow.EventWithFreeze(#COMMAND, cmd)
		}

	disableAccelCmd?()
		{
		return .overlay.Status isnt #Closed or .Frozen?()
		}

	keydown2(event)
		{
		if .ActiveWindow is false or .ActiveWindow.Destroyed?()
			return
		if event.key isnt 'Enter'
			return
		if event.target.tagName in ('BUTTON', 'TEXTAREA')
			return
		.ActiveWindow.CallDefaultButton()
		}

	RegisterKeydownListener(block)
		{
		.keydownListeners.Add(block)
		}

	UnregisterKeydownListener(block)
		{
		.keydownListeners.Remove(block)
		}

	eventId: 0
	overlayEvent: false
	freezeEvent: false
	Event(uniqueId, event, args, showOverlay? = false, freeze? = false)
		{
		if showOverlay? is true
			{
			.overlay.Show(id: #event)
			.overlayEvent = .eventId
			.CancelDelayedTask()
			}
		if freeze? is true
			{
			.freezeEvent = .eventId
			.CancelDelayedTask()
			}
		eventId = .eventId++
		.sendWithAck(event, :eventId, arg1: uniqueId, arg2: args)
		}

	sendWithAck(@eventOb)
		{
		.eventToAck.Add([id: eventOb.eventId, event: eventOb])
		.send(.prepareSend(eventOb))
		}

	prepareSend(eventOb)
		{
		eventOb.ack = .nextAck
		return Pack(SuMessageFormatter.FormatEvent(@eventOb))
		}

	HasOutstandingEvents?()
		{
		return .eventToAck.NotEmpty?()
		}

	sendLog?: true
	send(s, noOverlay? = false)
		{
		if .socketStatus is SuSocketStatus.connected
			{
			if .ws.readyState isnt 1
				{
				if .ws.readyState is 0 /*=CONNECTING*/ and .sendLog? is true
					{
					.sendLog? = false
					.Event(false, 'SuneidoLog', Object(
						'ERROR: (CAUGHT) sending events when ws is in CONNECTING',
						params: [state: .ws.readyState, :s],
						calls: .recorder.Get().Map(EventRecorder.Format).Join('\r\n'),
						caughtMsg: 'for debug 34395'))
					}
				return false
				}

			.ws.SendPacked(s)
			if noOverlay? is false
				.startTimer()
			.sendLog? = true
			return true
			}

		if .socketStatus is SuSocketStatus.terminated
			Print('socket terminated': s)

		return false
		}

	connectingTimer: false
	showConnecting?: false
	startTimer()
		{
		if .connectingTimer isnt false or .showConnecting? is true
			return
		.connectingTimer = SuDelayed(5000/*=5 secs*/, .showConnecting)
		}

	showConnecting()
		{
		.connectingTimer = false
		.showConnecting? = true
		.overlay.Show(id: #connecting, msg: 'Working...')
		}

	closeConnecting()
		{
		if .connectingTimer isnt false
			{
			.connectingTimer.Kill()
			.connectingTimer = false
			}
		if .showConnecting? isnt false
			{
			.overlay.Close(id: #connecting)
			.showConnecting? = false
			}
		}

	UpdateStatus(member, value)
		{
		.status[member] = value
		eventId = .eventId++
		.sendWithAck(SuMessageFormatter.Type.UpdateStatus, :eventId,
			arg1: member, arg2: value)
		}

	TimeOut(id)
		{
		eventId = .eventId++
		.sendWithAck(SuMessageFormatter.Type.SuJsTimeOut, :eventId, arg1: id)
		}

	GetToken()
		{
		return .token
		}

	Confirm(msg, onOK, onCancel)
		{
		.overlay.Show(#confirm, msg, okHandler: onOK, cancelHandler: onCancel, level: 99)
		}

	Getter_Overlay()
		{
		return .overlay
		}

	Getter_OverlayId()
		{
		return .overlayEvent
		}

	open()
		{
		.logWS('open')
		}

	connected()
		{
		Print("server connected")
		.closeReconnecting()
		.socketStatus = SuSocketStatus.connected
		.heartbeat.Start()
		.resendEventToAck()
		.retryCount = 0
		.syncVisibility()
		}

	resendEventToAck()
		{
		for ob in .eventToAck.Copy()
			.send(.prepareSend(ob.event))
		}

	ackEvent(id)
		{
		if .eventToAck.Empty?() or .eventToAck[0].id > id // processed event - ignore
			{
			return false
			}

		if .eventToAck[0].id < id // missing event response
			{
			.resendEventToAck()
			return false
			}

		.eventToAck.Delete(0)
		.nextAck = id
		return true
		}

	ErrorLog(s)
		{
		.Event(false, #ErrorLog, [s $ ' (from browser)'])
		}

	message(event)
		{
		if not .initialized
			{
			Print('Message skipped')
			return
			}

		try
			{
			response = FlatObject.Build(event.data)

			if response[0] is SuMessageFormatter.Type.CONNECTED
				{
				.logWS('connected ' $ response[1])
				if response[1] is .connectid
					.connected()
				}
			else if response[0] is SuMessageFormatter.Type.OVERLAY
				{
				if response[2/*=hide*/] is false
					.overlay.Show(#working, msg: response[1/*=msg*/], level: 50)
				else
					.overlay.Close(#working)
				}
			else
				{
				eventId = response[1]
				actions = response[0]
				.heartbeat.Event(eventId)

				if .ackEvent(eventId) is false
					return

				if .overlayEvent is eventId
					.overlay.Close(id: #event)

				.closeConnecting()
				.processActions(actions)

				if .overlayEvent is eventId
					.overlayEvent = false
				if .freezeEvent is eventId
					.freezeEvent = false
				if not .Frozen?() and .delayedTask isnt false
					{
					(.delayedTask.run)()
					.delayedTask = false
					}
				}
			}
		catch (e)
			{
			.overlayEvent = .freezeEvent = .delayedTask = false
			.Event(false, 'SuBrowserError', [message: e, stack: e.Callstack()])
			}
		}

	Frozen?()
		{
		return .freezeEvent isnt false or .overlayEvent isnt false
		}

	RunWhenNotFrozen(task)
		{
		if not .Frozen?()
			{
			task()
			return false
			}
		id = .delayedTaskId++
		.delayedTask = Object(:id, run: task)
		return id
		}

	CancelDelayedTask(id = 'any')
		{
		if id isnt 'any' and (.delayedTask is false or .delayedTask.id isnt id)
			return
		.delayedTask = false
		}

	processActions(actions)
		{
		// action is Object(uniqueId, action, args)
		for action in actions
			{
			if action.GetDefault(#canceled, false) is true
				continue

			if action.uniqueId is false
				{
				Global(action.action)(@action.args)
				}
			else
				{
				id = Display(action.uniqueId)
				if not .components.Member?(id)
					{
					Print("id (" $ id $ ") is not found - ", action)
					continue
					}
				method = false
				try
					method = .components[id][action.action]
				if method is false
					.debug36105(action, id, action.uniqueId)
				else
					method(@action.args)
				}
			}
		}

	debug36105(action, id, uniqueId)
		{
		nearbys = Object()
		for (i = -5; i <= 5/*=upper*/; i++)
			{
			if .components.Member?(Display(uniqueId + i))
				nearbys.Add(Display(.components[Display(uniqueId + i)]) $
					' - ' $ (uniqueId + i))
			}

		all = Object()
		for m in .components.Members()
			all[Number(m)] = Display(.components[m])

		try
			componentMax = .components.Members().MaxWith(Number)
		catch (e)
			componentMax = e
		params = [
			component: Display(.components[id]),
			componentId: .components[id].UniqueId,
			componentSize: .components.Size(),
			:componentMax,
			componentNearbys: nearbys]

		.Event(false, 'Debug36105', Object(action, params, all))
		}

	canvas: false
	GetTextMetrics(el, text)
		{
		if .canvas is false
			.canvas = CreateElement('canvas')

		context = .canvas.GetContext('2d')
		context.font = SuUI.GetCurrentWindow().
			GetComputedStyle(el).
			GetPropertyValue('font')
		width = height = 0
		ascent = 0
		descent = 0
		text.Lines().Each()
			{
			metrics = context.MeasureText(it)
			width = Max(width, metrics.width.Ceiling())
			height += metrics.fontBoundingBoxAscent + metrics.fontBoundingBoxDescent
			ascent = Max(ascent, metrics.fontBoundingBoxAscent)
			descent = metrics.fontBoundingBoxDescent
			}
		return Object(:width, :height, :ascent, :descent)
		}

	GetClientRect(el = false)
		{
		if el is false
			el = SuUI.GetCurrentDocument().documentElement
		boundingRect = el.GetBoundingClientRect()
		return Object(
			left: boundingRect.left,
			right: boundingRect.right,
			top: boundingRect.top,
			bottom: boundingRect.bottom,
			width: boundingRect.right - boundingRect.left,
			height: boundingRect.bottom - boundingRect.top)
		}

	scrollbarWidth: false
	GetScrollbarWidth()
		{
		if .scrollbarWidth isnt false
			return .scrollbarWidth

		outer = CreateElement('div', parent: SuUI.GetCurrentDocument().body)
		outer.SetStyle('visibility', 'hidden')
		outer.SetStyle('overflow', 'scroll')
		inner = CreateElement('div', parent: outer)
		.scrollbarWidth = outer.offsetWidth - inner.offsetWidth
		outer.Remove()
		return .scrollbarWidth
		}

	hDrop: 0
	AddDropFiles(files)
		{
		.dropFilesList[.hDrop] = files
		return .hDrop++
		}
	GetDropFiles(hDrop)
		{
		dropFiles = .dropFilesList.GetDefault(hDrop, false)
		.dropFilesList.Delete(hDrop)
		return dropFiles
		}

	beforeUnload(event)
		{
		if .socketStatus is SuSocketStatus.init or
			.socketStatus is SuSocketStatus.terminated
			return

		event.PreventDefault()
		// Chrome, Safari and Firefox don't support custom message any more
		return event.returnValue = "Are you sure you want to logout?"
		}

	ignores: #('ResizeObserver loop limit exceeded',
		'ResizeObserver loop completed with undelivered notifications.',
		// from stands AdBlocker
		'Cannot redefine property: googletag',
		)
	onError(event)
		{
		if .socketStatus is SuSocketStatus.terminated
			return false

		if .ignores.Any?({ event.message =~ it })
			return false

		.Event(false, 'SuBrowserError',
			[message: event.message, stack: event.error.stack])
		return true
		}

	ActivateWindow(id)
		{
		.Event(false, 'WindowActivate', [id], showOverlay?:)
		}

	Shutdown()
		{
		.ws.Close(1000/*=normal close*/)
		}
	}
