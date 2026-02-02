// Copyright (C) 2019 Axon Development Corporation All rights reserved worldwide.
class
	{
	Init(wsHandler, token, key)
		{
		Suneido.SuRenderBackend = new this(wsHandler, token, key)
		JsSessionToken.Register(token, key)
		}

	CallClass(noThrow = false)
		{
		if false is session = Suneido.GetDefault(#SuRenderBackend, false)
			{
			if noThrow is false
				throw "session not found "
			return false
			}
		return session
		}

	InErrorHandler: false
	pageVisible?: false
	New(.WSHandler, .Token, .Key)
		{
		.controls = Object()
		.status = Object().Set_default(false)

		.timeoutMin = Database.Info().GetDefault(#timeoutMin, 240/*= 4 hrs*/)
		.overrideProc =  Object()
		.beforeDisconnectFn = Object()
		.LastEvent(Date())
		.WindowManager = WindowManager()
		.TimerManager = SuTimerManager()
		.actionToAck = Object()
		.mostRecentAck = false
		.warnActionSize = OptContribution(#WarnSuJsActionSize, 43_000 /*= 2/3 of 64_000*/)

		.eventRecorder = new EventRecorder()
		}

	SetTimeoutMin(.timeoutMin) { }

	AddLog(log)
		{
		.eventRecorder.Add(#log, log)
		}

	DumpStatus(msg)
		{
		try
			{
			events = .eventRecorder.Get()
			actionToAckN = .actionToAck.Size()
			s = msg $ '\r\n'
			s $= 'Token: ' $ Base64.Decode(.Token) $ '\r\n'
			s $= 'ActionToAckN: ' $ actionToAckN $ '\r\n'
			s $= 'Events (' $ events.Size() $ '):\r\n'
			for event in events
				s $= EventRecorder.Format(event) $ '\r\n'

			if false isnt path = GetContributions('LogPaths').GetDefault(#sujslog, false)
				Rlog(path, s, multi_line:)
			}
		catch (e)
			SuneidoLog('ERROR: (CAUGHT) SuRenderBackend.DumpStatus - ' $ e,
				params: [:msg], caughtMsg: 'For debugging only')
		}

	id: 1
	NextId()
		{
		return .id++
		}

	LastEvent(lastEventTime = false)
		{
		if lastEventTime isnt false
			.lastEvent = lastEventTime
		return .lastEvent
		}

	AddOverrideProc(uniqueId, proc)
		{
		.overrideProc.GetInit(uniqueId, Object()).Add(proc)
		}

	RemoveOverrideProc(uniqueId, proc)
		{
		Assert(.overrideProc[uniqueId].Last() is: proc)
		.overrideProc[uniqueId].PopLast()
		if .overrideProc[uniqueId].Empty?()
			.overrideProc.Delete(uniqueId)
		}

	Register(uniqueId, control)
		{
		.controls[uniqueId] = control
		}

	GetRegisteredControl(uniqueId)
		{
		return .controls.GetDefault(uniqueId, false)
		}

	UnRegister(uniqueId)
		{
		.controls.Delete(uniqueId)
		}

	reconnectSocket: false
	SetReconnectSocket(socket)
		{
		.reconnectSocket = socket
		}
	closeReconnectSocket()
		{
		if .reconnectSocket isnt false
			{
			.reconnectSocket.Close()
			.reconnectSocket = false
			}
		}

	RegisterBeforeDisconnectFn(fn)
		{
		.beforeDisconnectFn.Add(fn)
		}

	beforeDisconnectFn: ()
	BeforeDisconnect()
		{
		.beforeDisconnectFn.Each()
			{ |fn|
			try
				fn()
			catch (e)
				SuneidoLog('ERROR: BeforeDisconnect - ' $ e, params: [fn])
			}
		.closeReconnectSocket()
		JsSessionToken.Unregister(.Token)
		}

	Close()
		{
		if not Suneido.Member?(#Persistent)
			return

		master = false
		for w in Suneido.Persistent.Windows
			{
			if w.Master? isnt false and .controls.Member?(w.UniqueId)
				{
				master = w
				break
				}
			}
		if master isnt false
			master.DESTROY()
		}

	getter_actions()
		{
		return .actions = Object()
		}

	Getter_Actions()
		{
		actions = .actions
		.actions = Object()
		if actions.Size() > .warnActionSize
			.logObjectTooLarge(actions, warn:)
		result = SuMessageFormatter.FormatResponse(actions, arg1: .eventId)
		.actionToAck[.eventId] = [time: Date(), :result]
		return result
		}

check(ob)
	{
	if Object?(ob)
		for item in ob
			if .check(item) is true
				{
				SuneidoLog('item', params: item)
				return true
				}
	return Class?(ob) or Instance?(ob)
	}

	Status()
		{
		return .status
		}

	eventId: false
	EventHandler(@args)
		{
		.eventRecorder.Add(#event, args)
		.processAck(args[2])

		if false isnt response = .actionToAck.GetDefault(args[1]/*eventId*/, false)
			return response.result
		.eventId = args[1]
		switch (args[0])
			{
		case SuMessageFormatter.Type.UpdateStatus:
			.status[args[3/*=member*/]] = args[4/*=value*/]
		case SuMessageFormatter.Type.SuJsTimeOut:
			.TimerManager.Timeout(args[3/*=id*/])
		case SuMessageFormatter.Type.Heartbeat:
			.checkTimeout()
		case SuMessageFormatter.Type.SyncVisibility:
			.syncVisibility(args[3/*=state*/])
		default: // event
			.handleEvent(args[3/*=uniqueId*/], args[0], args[4/*=args*/])
			}
		.TimerManager.FlushDelays()
		return .Actions
		}

	processAck(ack)
		{
		if ack > .mostRecentAck
			{
			for id in .actionToAck.Members().Copy()
				if id <= ack
					.actionToAck.Erase(id)
			.mostRecentAck = ack
			}
		}

	TimedOutMsg: `INFO: Suneido.js Timeout - closing idle connection `
	checkTimeout()
		{
		if Date().MinusMinutes(.LastEvent()) > .timeoutMin
			Finally(
				{ SuneidoLog(.TimedOutMsg $ Thread.Name()) },
				{ .Terminate(reason: 'Disconnected due to lack of activity') })
		}

	syncVisibility(state)
		{
		.pageVisible? = state is 'visible'
		.WSHandler.GetSocket().SetTimeout(.pageVisible? ? 3 : 75/*=timeout*/)
		}

	handleEvent(uniqueId, event, args)
		{
		.LastEvent(Date())
		if uniqueId is false
			Global(event)(@args)
		else
			.CallProc(uniqueId, event, args)
		}

	CallProc(uniqueId, event, args = [])
		{
		for (i = .overrideProc.GetDefault(uniqueId, #()).Size() - 1; i >= 0; i--)
			if ((.overrideProc[uniqueId][i])(uniqueId, event, args) is false)
				return

		if .controls.Member?(uniqueId) and .controls[uniqueId].Method?(event)
			(.controls[uniqueId][event])(@args)
		}

	ReserveAction()
		{
		ob = Object(reserved?:, at: .actions.Size())
		.addAction(ob)
		return ob
		}

	// for Dialog
	ExtractReserved()
		{
		if false is at = .actions.FindIf({ it.Member?(#reserved?) })
			return #()
		reserved = .actions[at..]
		.actions = .actions[..at]
		return reserved
		}

	MergeReserved(reserved)
		{
		base = .actions.Size()
		for i in reserved.FindAllIf({ it.Member?(#reserved?) })
			reserved[i].at = base + i
		.actions.Append(reserved)
		}

	CancelReserve(at)
		{
		.actions = .actions[..at]
		}

	CancelAllReserved()
		{
		if false is at = .actions.FindIf({ it.Member?(#reserved?) })
			return
		.CancelAllAfter(at-1)
		}

	RecordAction(uniqueId, action, args, at = false)
		{
		ob = Object(:uniqueId, :action, :args)
		if at isnt false
			{
			Assert(.actions[at] hasMember: #reserved?)
			.actions[at] = ob
			}
		else
			.addAction(ob)
		}

	addAction(ob)
		{
		try
			.actions.Add(ob)
		catch (e/*unused*/, 'object too large')
			{
			actions = .actions
			.actions = Object()
			.logObjectTooLarge(actions)
			.Terminate(reason: 'Fatal')
			}
		}

	logObjectTooLarge(actions, warn = false)
		{
		LogErrors('SuRenderBackend.logObjectTooLarge')
			{
			temp = Object().Set_default(0)
			for action in actions
				{
				if action.Member?(#reserved?)
					continue
				temp[action.Project(#('uniqueId', 'action'))]++
				}
			stats = Object()
			for m in temp.Members()
				stats.Add(Object(count: temp[m]).Merge(m))
			params = stats.Sort!(By(#count)).Reverse!().Take(5/*=top 5*/)
			params.Each()
				{
				if false isnt c = .GetRegisteredControl(it.uniqueId)
					{
					it.control = Display(c)
					it.name = c.Name
					}
				}
			msg = 'Suneido.js action list too large'
			SuneidoLog((warn ? 'WARNING' : 'ERROR') $ ' - ' $ msg $
				' (' $ ReadableSize(actions.Size()) $ ')',
				:params, calls:)
			.DumpStatus(msg)
			}
		}

	CancelAction(uniqueId, action, block = false)
		{
		for toCancel in .actions.Filter({
			not it.Member?(#reserved?) and not it.Member?(#canceled) and
			(uniqueId is #ignore or it.uniqueId is uniqueId) and
			it.action is action and
			(block is false or block(it.args)) })
			{
			toCancel.canceled = true
			toCancel.args = false
			}
		}

	CancelAllAfter(at)
		{
		.actions = .actions[..at+1]
		}

	Overlay(msg, hide = false)
		{
		.WSHandler.Send(#BINARY, Pack(SuMessageFormatter.FormatResponse(
			SuMessageFormatter.Type.OVERLAY, arg1: msg, arg2: hide)))
		}

	Terminate(@args) /*usage: reason = false, e = false */
		{
		.WSHandler.Terminate(@args)
		}
	}
