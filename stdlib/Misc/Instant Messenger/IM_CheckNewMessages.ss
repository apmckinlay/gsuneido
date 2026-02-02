// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(env)
		{
		IM_MessengerManager.LogConnectionsIfChanged()
		args = IM_MessengerManager.ParseQuery(env.body)
		if args.missed is 'true'
			{
			return Json.Encode(
				Object(messages:
					.checkMessagesForSubscription(args.user,Suneido.GetDefault(
						#IM_Messages, #()))))
			}
		return LongPoll(args, .check, .hasUpdate?, .defaultResult)
		}

	checkMessagesForSubscription(userName,messages)
		{
		subs = UserSettingsCached().Get('IM_SubscribedChannels', #(), userName)
		userMessages = messages.GetDefault(userName, #()).Copy()
		for message in userMessages
			{ // Checking for contact messages
			if not Date?(Date(message.GetDefault(#imchannel_num, ''))) and
				not message.im_message.Prefix?(IM_MessengerManager.SysMsgTemplate)
				{
				return true
				}
			if subs.Has?(Date(message.GetDefault(#imchannel_num, false)))
				return true
			}
		return false
		}

	check(args, checkingTime)
		{
		// after restarted, re-initialize IM_CheckMessages
		if not Suneido.Member?(#IM_CheckMessages)
			{
			Suneido.IM_CheckMessages = Object()
			.loadUnreadMsgs()
			}

		cstr =  Object(
			messages: Suneido.IM_CheckMessages.GetDefault(
				args.user, Date.Begin()) > checkingTime
			contacts: Suneido.GetDefault(#IM_Contacts, Date.Begin()) > checkingTime
			setRead: Suneido.GetDefault(#IM_MessagesRead, Object()).GetDefault(
				args.user, Date.Begin()) > checkingTime
			)

		notify = true
		hasSub = false
		sys_msg = false

		if .checkMessagesForSubscription(args.user,Suneido.GetDefault(#IM_Messages, #()))
			hasSub = true

		if Suneido.IM_CheckMessages.GetDefault(
			args.user $ '_sys_msg', Date.Begin()) > checkingTime
			{
			sys_msg = true
			}

		if sys_msg or not hasSub
			notify = false

		cstr.notify = notify

		return cstr
		}

	hasUpdate?(result)
		{
		return result.messages or result.contacts or result.setRead
		}

	defaultResult()
		{
		return #(messages: false, contacts: false, setRead: false)
		}

	loadUnreadMsgs()
		{
		IM_MessengerManager.EnsureTables()
		Suneido.IM_Messages = Object()

		lastGetRecs = QueryAll('im_last_get sort im_last_get')
		if lastGetRecs.Empty?()
			return

		daysInWeek = 7
		lastGet_farthest = Max(lastGetRecs[0].im_last_get, Date().Minus(days: daysInWeek))

		userLastGet = Object().Set_default(Date().Minus(days: daysInWeek))
		lastGetRecs.Each({ userLastGet[it.im_to] = it.im_last_get })

		QueryApply('im_history where im_num > ' $ Display(lastGet_farthest))
			{
			msgOb = Object(im_to: "", im_message: String(it.im_message),
				im_from: it.im_from, im_num: it.im_num, imchannel_num: it.imchannel_num)
		//  ^^^ [WARNING!] This should be consistent with the format in IM_SendNewMessage!
			if #() isnt recipientsList = IM_MessengerManager.ParseRecipients(it.im_to)
				{
				for recipient in recipientsList.Remove(it.im_from)
					if userLastGet[recipient] < it.im_num
						{
						msgOb.im_to = recipient
						IM_SendNewMessage.AddMessage(msgOb)
						}
				}
			}
		}
	}
