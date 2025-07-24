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
		.channTab = .FindControl("selection_tabs")
		.contacts = .FindControl('Contacts')
		.channels = .FindControl("Channels")
		.warning = .FindControl('im_warnUser')
		.warning.SetColor(CLR.RED)

		.sub = PubSub.Subscribe('messengerMessages', .checkMessengerServer)

		BookLog('IM:open')
		Delay(100, .startUp) /*= delayed 100 milliseconds in case
			there is an dialog open right before this window open */
		}

	Controls()
		{
		contactsImpl = SoleContribution('IM_ContactsImplementation')
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
				constructAll:,
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
					if false isnt channels = .FindControl('Channels')
						{
						channels.ReloadAndUpdateVirtualList()
						if sysMsg.type is 'channel_changed'
							channelChanged? = true
						}
					}
				}
			}
		if channelChanged?
			.conversation_tabs.UpdateChangedTabs()
		return messages
		}

	handleMessages(messages, _fromStartup = false)
		{
		originalFocus = GetFocus()
		.setLastGet()
		.UpdateContactHistory()
		flashWindow? = .processMessages(messages)
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
		atleastOneMessage = false
		flashWindow? = false
		curUser = .conversation_tabs.TabGetSelectedName()
		isTyping? = .conversation_tabs.IsTyping?()
		lastUser = false

		for message in messages
			{
			.convertMessage(message)
			atleastOneMessage = true
			isChannel? = message.imchannel_num isnt ''
			user = .userFromMessage(message, isChannel?)

			if not isChannel? or .isChannelSubscribed?(message.imchannel_num)
				{
				.conversation_tabs.AddMessage(user, message,
					message.GetDefault(#imchannel_num, ''))
				flashWindow? = true
				lastUser = user
				}
			else
				flashWindow? = .updateChannelNotificationButton(message) or flashWindow?
			}

		.sendReadAllSystemMessage(atleastOneMessage)

		if isTyping? is true
			.select(curUser)
		else if lastUser isnt false
			.select(lastUser)

		return flashWindow?
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

	updateChannelNotificationButton(message)
		{
		if .hasUserMention(message.im_message)
			{
			.channels.SetNotificationButtonAsMentioned(message.imchannel_num)
			.channTab.Select(1)
			return true
			}
		else
			{
			.channels.SetNotificationButtonAsUnread(message.imchannel_num)
			return false
			}
		}

	hasUserMention(msg)
		{
		return msg.Has?('@everyone') or msg.Has?('@' $ Suneido.User)
		}

	sendReadAllSystemMessage(atleastOneMessage)
		{
		if atleastOneMessage
			IM_MessengerManager.OutputMessage(Suneido.User,
				IM_SendNewMessage.AllReadSystemMessage)
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
		subStatus = paramObj.subStatus
		subscribedChannels = UserSettings.Get('IM_SubscribedChannels', def: Object())

		if subStatus is 'subscribe'
			subscribedChannels.AddUnique(paramObj.rec.imchannel_num)
		else if subStatus is 'unsubscribe'
			subscribedChannels.Remove(paramObj.rec.imchannel_num)

		UserSettings.Put('IM_SubscribedChannels', subscribedChannels)
		.channels.ReloadAndUpdateVirtualList()
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
