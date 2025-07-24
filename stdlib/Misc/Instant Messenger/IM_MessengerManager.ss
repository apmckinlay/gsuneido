// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
class
	{
	EnsureTables()
		{
		Database('ensure im_channels
			(imchannel_name, imchannel_num, imchannel_abbrev, imchannel_status,
			imchannel_desc)
			key(imchannel_name)
			key(imchannel_num)
			index unique(imchannel_abbrev)')
		Database('ensure im_history
			(im_num, im_from, im_to, im_message, imchannel_num)
			key(im_num)
			index (imchannel_num) in im_channels')
		Database('ensure im_last_get
			(im_to, im_last_get)
			key(im_to)')
		}

	MessengerThread()
		{
		if not IM_Available?() or Suneido.GetDefault('IM_Thread', false)
			return

		Suneido.IM_Thread = true
		SujsAdapter.CallOnRenderBackend('RegisterBeforeDisconnectFn', .StopThread)
		Thread(.startThread)
		}

	startThread()
		{
		Thread.Name('IM_MessengerManagerStartThread-thread')
		msg = false
		if false isnt .Request("/IM_ContactsUpdate", quiet?:)
		// the IM_CheckNewMessages request here is used to check new messages sent before
		// IM starts. The long poll IM_CheckNewMessages request only checks for new
		// messages sent during the long poll
			msg = .Request("/IM_CheckNewMessages?", content: "user=" $
				Url.Encode(Suneido.User) $ '\r\nmissed=true', quiet?:)

		if msg is false
			{
			Suneido.IM_Thread = false
			return
			}

		if msg.messages
			RunOnGui({ PubSub.Publish('messengerMessages') })
		Thread(.thread)
		}

	thread()
		{
		Thread.Name('IM_MessengerManagerThread-thread')

		if Suneido.IM_Thread isnt true
			return

		if false is msg = .Request("/IM_CheckNewMessages?", content: "user=" $
			Url.Encode(Suneido.User), quiet?:)
			{
			Suneido.IM_Thread = false
			return
			}

		.publishMessengerEvents(msg)

		// start a new thread instead of looping
		// so we don't keep a database connection open all the time
		Thread(.thread)
		}

	publishMessengerEvents(msg)
		{
		if msg.setRead
			RunOnGui({ PubSub.Publish('messengerMessages', setRead:, notify: false) })
		else if msg.messages
			{
			if msg.notify
				RunOnGui({ PubSub.Publish('messengerMessages') })
			else if not msg.notify
				RunOnGui({ PubSub.Publish('messengerMessages', notify: false) })
			}
		if msg.contacts and Suneido.Member?('IM_Window')
			RunOnGui({ PubSub.Publish('messengerContacts') })
		}

	// .Request can be called from the UI thread
	Request(request, content = "", quiet? = false, skipSubscribe? = false)
		{
		try
			{
			result = Json.Decode(.postResult(request, content, ServerIP()))
			if not skipSubscribe?
				.publishMessengerOnline()
			return result
			}
		catch(e)
			.logErr(e, quiet?)
		return false
		}

	publishMessengerOnline()
		{
		if Suneido.GetDefault(#IM_Status, '') isnt 'online'
			{
			Suneido.IM_Status = 'online'
			.ensureOnGui({ PubSub.Publish('messengerOnline') })
			}
		}

	// Extracted for tests purposes
	postResult(request, content, ip)
		{
		url = 'http://' $ (ip is "" ? '127.0.0.1' : ip) $ ':' $ HttpPort()
		return Http.Post(Url.Encode(url $ request), content)
		}

	logErr(e, quiet?)
		{
		if e =~ '(?i)socket|connect'
			.ensureOnGui({ .imOffline(quiet?) })
		else
			SuneidoLog("ERRATIC: " $ e)
		}

	ensureOnGui(block)
		{
		if not Sys.MainThread?()
			RunOnGui(block)
		else
			block()
		}

	imOffline(quiet? = false)
		{
		if Suneido.GetDefault(#IM_Status, '') isnt 'offline'
			{
			Suneido.IM_Status = 'offline'
			PubSub.Publish('messengerOffline')
			}
		if not quiet? and Suneido.User isnt 'default'
			AlertError('Unable to connect to Instant Messenger server.')
		if Suneido.Member?('IM_Window')
			{
			DestroyWindow(Suneido.IM_Window.Hwnd)
			Suneido.Delete('IM_Window')
			}
		}

	StopThread()
		{
		Suneido.IM_Thread = false
		}
	FindGroup(name)
		{
		groups = .GetGroupList()
		return (false isnt (pos = groups.FindIf({ name is it.im_user }))
			? groups[pos]
			: false)
		}

	contactsImpl()
		{
		return SoleContribution('IM_ContactsImplementation')
		}

	GetGroupList()
		{
		return .contactsImpl().Member?('groupList')
			? (Global(.contactsImpl().groupList))()
			: Object()
		}
	FindUser(name)
		{
		users = .GetUserList()
		return (false isnt (pos = users.FindIf({ name is it.im_user }))
			? users[pos]
			: false)
		}
	GetUserList()
		{
		return (Global(.contactsImpl().userList))()
		}

	GetChannelList()
		{
		return QueryList('im_channels', 'imchannel_name')
		}

	GetActiveChannelsList()
		{
		return QueryList("im_channels where imchannel_status is 'active'",
			'imchannel_name')
		}

	GetActiveChannelsListByNum()
		{
		return QueryList("im_channels where imchannel_status is 'active'",
			'imchannel_num')
		}

	GetChannelNameFromNum(num)
		{
		if false is rec = Query1('im_channels', imchannel_num: num)
			return '???'
		return rec.imchannel_name
		}

	FindChannel(name)
		{
		pos = .GetChannelList().FindIf({it is name})
		return pos
		}

	ParseRecipients(recipients)
		{
		recipients_list = Object()
		for x in recipients.Split('+')
			{
			if x is 'GLOBAL'
				recipients_list = .GetUserList().Map({ it.im_user })
			else if .FindUser(x) isnt false
				recipients_list.AddUnique(x)
			else if .FindGroup(x) isnt false
				.expandGroup(x, recipients_list)
			// It is currently possible to get here without successfully
			// parsing the recipient. For example if the user was deleted from
			// the users table. We don't want to stop users from deleting users
			// from the table therefore this can not be an error case.
			}
		return recipients_list
		}

	SysMsgTemplate: "~~~Channel Maintenance~~~"
	IsSystemMessage?(msg) // System messages are meant for internal use not for user
		{
		return msg.im_message.Prefix?(.SysMsgTemplate)
		}

	SendSystemMessage(entry)
		{
		.OutputMessage('GLOBAL', entry, Timestamp())
		}

	OutputMessage(recipient, entry, time_now = '', imchannel_num = '')
		{
		if time_now is ''
			time_now = Timestamp()

		if false is IM_MessengerManager.Request("/IM_SendNewMessage?",
			content: "to=" $ Url.Encode(recipient) $
				"\r\nfrom=" $ Url.Encode(Suneido.User) $
				"\r\nmsg=" $ Url.Encode(entry) $
				"\r\ndate=" $ Url.Encode(Display(time_now)) $
				'\r\nimchannel_num=' $
					(imchannel_num is '' ? '' : Url.Encode(Display(imchannel_num))))
			return false

		if not entry.Prefix?(.SysMsgTemplate)
			QueryOutput("im_history", Record(
				im_from: Suneido.User,
				im_to: recipient,
				im_message: entry,
				im_num: time_now,
				:imchannel_num))

		return true
		}

	ParseQuery(content)
		{
		query = Object().Set_default('')
		for line in content.Split('\r\n')
			{
			i = line.Find('=')
			if (i >= line.Size())
				query.Add(ConvertNumeric(line))
			else
				query[ConvertNumeric(line[.. i])] =
					ConvertNumeric(Url.Decode(line[i + 1 ..]))
			}
		return query
		}

	expandGroup(group, recipients_list)
		{
		if not .contactsImpl().Member?('parseGroup')
			return
		recipients_list.MergeUnion(Global(.contactsImpl().parseGroup)(group))
		}

	Register()
		{
		if false is result = .Request("/IM_Register?",
			content: "user=" $ Url.Encode(Suneido.User), quiet?:)
			return "Unable to connect to Instant Messenger server."
		if result.GetDefault('allowOpen', false) is false
			return "Please exit the other instance of Instant Messenger."
		return ""
		}

	UnRegister()
		{
		.Request("/IM_UnRegister?", content: "user=" $ Url.Encode(Suneido.User), quiet?:)
		}

	LogConnectionsIfChanged()
		{
		preConnections = Suneido.GetDefault('PreConnections', #()).Copy().Sort!()
		curConnections = Sys.Connections().Sort!()
		if curConnections isnt preConnections
			{
			Suneido.PreConnections = curConnections
			diff = curConnections.Size() > preConnections.Size()
				? curConnections.Difference(preConnections)
				: preConnections.Difference(curConnections)
			increase? = curConnections.Size() > preConnections.Size()
			size = 'unknown'
			try size = ReadableSize(
				Suneido.GoMetric("/memory/classes/heap/objects:bytes"))
			try BookLog('connections (' $ Display(curConnections.Size()) $
				', ' $ size $ '): ' $ (increase? ? 'new ' : 'closed ') $ Display(diff),
				systemLog:)
			}
		}
	}
