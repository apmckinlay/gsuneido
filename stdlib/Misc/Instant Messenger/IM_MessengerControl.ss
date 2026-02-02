// Copyright (C) 2022 Axon Development Corporation All rights reserved worldwide.
Controller
	{
	Xmin: 600
	Ymin: 500
	Title: "Instant Messenger"

	CallClass()
		{
		if not .Opened?()
			{
			if '' isnt errMsg = IM_MessengerManager.Register()
				{
				.AlertWarn('Instant Messenger', errMsg)
				return
				}
			else
				Suneido.IM_Window = Window(this, w: 500, h: 300, keep_placement:)
			}
		else
			WindowActivate(Suneido.IM_Window.Hwnd)
		}

	Opened?()
		{
		return Suneido.Member?('IM_Window')
		}

	New()
		{
		IM_MessengerManager.EnsureTables()

		.conversation_tabs = .FindControl('conversation_tabs')
		.selection_tabs = .FindControl('selection_tabs')
		.contacts = .FindControl('Contacts')
		.warning = .FindControl('im_warnUser')
		.warning.SetColor(CLR.RED)

		.sub = PubSub.Subscribe('messengerMessages', .checkMessengerServer)

		BookLog('IM:open')
		Delay(100, .startUp) /*= delayed 100 milliseconds in case
			there is an dialog open right before this window open */
		}

	Controls()
		{
		contactsImpl = LastContribution('IM_ContactsImplementation')
		return Object('HorzSplit'
			Object('Vert'
				Object('IM_MessengerTabs'
					Object('IM_MessengerHome', Tab: 'Home'),
					name: 'conversation_tabs', close_button:,
					orientation: 'bottom',
					staticTabs: #(Home))
				#(Heading, '', name: 'im_warnUser')
				xmin: 500)
			Object('Tabs'
				#(IM_Contacts, name: "Contacts", Tab: 'Contacts'),
				#(IM_Channels, name: "Channels", Tab: 'Channels'),
				extraControl: contactsImpl.GetDefault('extraTop', false),
				name: 'selection_tabs'
				)
			)
		}


	timer: false
	startUp()
		{
		if .Destroyed?()
			return

		_fromStartup = true
		.getNewMessages()

		if .Destroyed?()
			return
		.conversation_tabs.FocusEditor()
		}

	checkMessengerServer()
		{
		if Suneido.GetDefault('GettingNewMessage', false) is false
			.getNewMessages()
		return 0
		}

	getNewMessages()
		{
		Suneido.GettingNewMessage = true
		pendingMessages = Object()
		if false isnt messages = IM_MessengerManager.Request(
			"/IM_GrabNewMessage?", content: "user=" $ Url.Encode(Suneido.User))
			{
			if messages.Member?('error')
				{
				.showTooManyMessagesPopup()
				messages.Delete('error')
				}
			pendingMessages.MergeUnion(messages)
			}

		if not pendingMessages.Empty?()
			{
			messages = .processSystemMessages(pendingMessages)
			.handleMessages(messages)
			}
		Suneido.GettingNewMessage = false
		}

	showTooManyMessagesPopup()
		{
		errMsg = 'You have too many messages!' $
				'Only the first ' $ Display(IM_GrabNewMessage.MaxMessageFetchLimit) $
	' will be loaded, this means you might miss notifications for some of these messages'$
	' All your messages wil still be there but tabs might not automatically open up for' $
	'you. So be sure to manually click on important conversations and load history' $
	' to NOT miss any messages'
		.AlertWarn('Instant Messenger', errMsg)
		}

	processSystemMessages(pendingMessages)
		{
		channelChanged? = false
		messages = Object()

		for (i = 0; i < pendingMessages.Size(); i++)
			{
			message = pendingMessages[i].im_message
			if not message.Prefix?(IM_MessengerManager.SysMsgTemplate)
				messages.Add(pendingMessages[i])
			else
				{
				sysMsg = message.RemovePrefix(IM_MessengerManager.SysMsgTemplate).
					SafeEval()
				if sysMsg.token is 'a0cb03a6f32a2448'
					{
					.reloadChannels()
					if sysMsg.type is 'channel_changed'
						channelChanged? = true
					}
				}
			}
		if channelChanged?
			.conversation_tabs.UpdateChangedTabs()
		return messages
		}

	reloadChannels(channelNums = false)
		{
		if .channels isnt false
			.channels.Reload(channelNums)
		}

	getter_channels()
		{
		if false is channels = .FindControl('Channels')
			return false
		return .channels = channels
		}

	handleMessages(messages, _fromStartup = false)
		{
		originalFocus = GetFocus()
		.setLastGet()
		.UpdateContactHistory()
		selectUser, reloadChannelNums, flashWindow? = .processMessages(messages)
		if false isnt selectUser
			.select(selectUser)
		if reloadChannelNums.NotEmpty?()
			.reloadChannels(reloadChannelNums)
		SetFocus(originalFocus)
		if flashWindow? and not super.HasFocus?() and not fromStartup
			FlashWindowEx(Object(cbSize: FLASHWINFO.Size(), hwnd: .Window.Hwnd,
				dwFlags: FLASHW.ALL, uCount: 5))
		}

	setLastGet()
		{
		if QueryEmpty?('im_last_get', im_to: Suneido.User)
			QueryOutput('im_last_get', [im_to: Suneido.User, im_last_get: Timestamp()])
		else
			QueryDo('update im_last_get
				where im_to is '$ Display(Suneido.User) $
				' set im_last_get = ' $ Display(Timestamp()))
		}

	processMessages(messages)
		{
		flashWindow? = lastUser = selectUser = false
		reloadChannelNums = Object()
		for message in messages
			{
			.convertMessage(message)
			isChannel? = message.imchannel_num isnt ''
			user = .userFromMessage(message, isChannel?)

			if not isChannel? or .isChannelSubscribed?(message.imchannel_num)
				{
				.conversation_tabs.AddMessage(user, message,
					message.GetDefault(#imchannel_num, ''))
				flashWindow? = true
				lastUser = user
				}
			else if .updateChannelNotify(message)
				{
				reloadChannelNums.Add(message.imchannel_num)
				flashWindow? = true
				}
			}

		if messages.NotEmpty?()
			IM_MessengerManager.
				OutputMessage(Suneido.User, IM_SendNewMessage.AllReadSystemMessage)

		selectUser = .conversation_tabs.IsTyping?()
			? .conversation_tabs.TabGetSelectedName()
			: lastUser
		return selectUser, reloadChannelNums, flashWindow?
		}

	convertMessage(message)
		{
		if false isnt Date(message.GetDefault(#imchannel_num, ''))
			message.imchannel_num = Date(message.imchannel_num)
		else
			message.imchannel_num = ''
		message.im_num = Date(message.im_num)
		}

	RecipientTagPat: "\[([\w\d_ +]+)\]\s+"
	userFromMessage(message, isChannel?)
		{
		if isChannel?
			return IM_MessengerManager.GetChannelNameFromNum(message.imchannel_num)

		// extract conversation recipients from the message body
		groupConvo? = message.im_message.Prefix?('[')
		if not groupConvo?
			return message.im_from

		recipients = message.im_message.Extract(.RecipientTagPat)
		return recipients is false ? message.im_from : recipients
		}

	isChannelSubscribed?(num)
		{
		if IM_ShowMembersControl.GetMembersFromChannelNum(num).Has?(Suneido.User)
			return true

		return .conversation_tabs.ChannelOpen?(num) // Not a subscribed channel but open
		}

	updateChannelNotify(message)
		{
		if mentioned? = .hasUserMention(message.im_message)
			.selection_tabs.Select(1)
		notify? = IM_ChannelsControl.SetNotify(message.imchannel_num, :mentioned?)
		return notify? or mentioned?
		}

	hasUserMention(msg)
		{
		return msg.Has?('@everyone') or msg.Has?('@' $ Suneido.User)
		}

	select(user)
		{
		.conversation_tabs.SelectByName(user)
		}

	OpenConversation(user)
		{
		.conversation_tabs.OpenTab(user)
		}

	OpenChannel(data)
		{
		.conversation_tabs.OpenTab(data.imchannel_name, data.imchannel_num)
		}

	TabsControl_SelectTab(source)
		{
		if Same?(source, .conversation_tabs)
			.ClearWarning()
		}

	ShowWarning(warning)
		{
		.warning.Set(warning)
		}

	ClearWarning()
		{
		.warning.Set("")
		}

	UpdateContactHistory()
		{
		.contacts.UpdateContactHistory()
		}

	Context_subUnsub(paramObj)
		{
		cachedSubscriptions = UserSettingsCached().Get('IM_SubscribedChannels')
		subscribedChannels = Object?(cachedSubscriptions)
			? cachedSubscriptions.Copy() // Cached Objects are static
			: Object()

		subStatus = paramObj.subStatus
		if subStatus is 'subscribe'
			subscribedChannels.AddUnique(paramObj.rec.imchannel_num)
		else if subStatus is 'unsubscribe'
			subscribedChannels.Remove(paramObj.rec.imchannel_num)

		UserSettingsCached().Put('IM_SubscribedChannels', subscribedChannels,
			resetServer?:)
		.reloadChannels()
		}

	Destroy()
		{
		Suneido.Delete('IM_Window')
		IM_MessengerManager.UnRegister()
		.sub.Unsubscribe()
		BookLog('IM:close')
		super.Destroy()
		}
	}
