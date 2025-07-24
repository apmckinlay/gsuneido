// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Title: 'Channels'
	Xmin: 400
	New()
		{
		.listChannels = .FindControl('listChannel')
		.status_ctrl = .FindControl('imchannel_notification')
		}

	Controls()
		{
		.subscriptionStatus = .getSubscriptionStatus()

		return Object('Vert'
			Object('VirtualList',
				'im_channels
				extend imchannel_joinedStatus, imchannel_lastActivity
				where imchannel_status is "active"
				and imchannel_name isnt "user_channel"',
				name: 'listChannel',
				columns:
					#(imchannel_name,
					imchannel_desc,
					imchannel_notification,
					imchannel_joinedStatus,
					imchannel_lastActivity,
					imchannel_status,
				),
				columnsSaveName: .Title,
				resetColumns:
			))
		}

	VirtualList_ExtraSetupRecordFn()
		{
		return .BeforeRecord
		}

	BeforeRecord(rec)
		{
		allUnreadChannels = UserSettings.Get('IM_UnreadChannels')
		unreadMentionsList = UserSettings.Get('IM_MentionedChannels')
		if .subscriptionStatus.value.Has?(rec.imchannel_num)
			rec.imchannel_joinedStatus = true

		if allUnreadChannels isnt false and allUnreadChannels.Has?(rec.imchannel_num)
			rec.imchannel_notification = CLR.GREEN

		if unreadMentionsList isnt false and unreadMentionsList.Has?(rec.imchannel_num)
			rec.imchannel_notification = CLR.orange // orange button takes priority

		rec.imchannel_lastActivity = .getTimeSinceLastActivity(rec.imchannel_num)
		}

	VirtualList_DoubleClick(rec, col /*unused*/, source)
		{
		if source is .listChannels
			if .listChannels.Empty?() isnt true and rec isnt false
				{
				.SetNotificationButtonAsRead(rec.imchannel_num)
				.Send('OpenChannel', rec)
				}
		return true
		}

	VirtualList_BuildContextMenu(rec)
		{
		if rec is false
			return #('Manage Channels','Mark All as Read')
		return #('Status' #('Set Inactive' #('Set Inactive')),
			'Notifications' #('Subscribe','Unsubscribe'),
			'Show Members',
			'Mark All as Read',
			'Manage Channels')
		}
	maintMsg: "~~~Channel Maintenance~~~ " $
		"#(type: 'channel_changed', token: 'a0cb03a6f32a2448')"
	On_Context_Status_Set_Inactive_Set_Inactive(@args)
		{
		if false is args.GetDefault('rec', false)
			return
		else
			{
			num_param = String(args.rec.imchannel_num)
			update_query = "update im_channels where imchannel_num is " $ num_param $
				" set imchannel_status = 'inactive'"
			QueryDo(update_query)
			.listChannels.Refresh()
			.Send("CloseTab",args.rec)
			IM_MessengerManager.SendSystemMessage(.maintMsg)
			}
		}

	On_Context_Notifications_Subscribe(@args)
		{ // Constructing new param obj for easy parsing in parent class
		paramObj = Object(rec : '', subStatus : '')
		paramObj.rec = args.rec
		paramObj.subStatus = 'subscribe'
		.Send("Context_subUnsub",paramObj)
		}

	On_Context_Notifications_Unsubscribe(@args)
		{
		paramObj = Object(rec : '', subStatus : '')
		paramObj.rec = args.rec
		paramObj.subStatus = 'unsubscribe'
		.Send("Context_subUnsub",paramObj)
		}

	On_Context_Show_Members(@args)
		{
		channel_num = args.rec.imchannel_num
		ToolDialog(.WindowHwnd(), Object('IM_ShowMembers', channel_num))
		}

	On_Context_Mark_All_as_Read()
		{
		UserSettings.Put('IM_UnreadChannels', #())
		.listChannels.Refresh()
		}

	On_Context_Manage_Channels()
		{
		ToolDialog(.WindowHwnd(), 'IM_ChannelManager', addNew:)
		.listChannels.Refresh()
		IM_MessengerManager.SendSystemMessage(.maintMsg)
		}

	ReloadAndUpdateVirtualList()
		{
		.subscriptionStatus = .getSubscriptionStatus()
		.listChannels.Refresh()
		}

	getSubscriptionStatus()
		{
		res = Query1('user_settings where user is ' $
			Display(Suneido.User) $ ' and key is "IM_SubscribedChannels"')
		return res isnt false ? res : [value: #()]
		}

	getTimeSinceLastActivity(channel_num) // Returns formatted string of last activity
		{
		now = Timestamp()
		last = .getLastActivity(channel_num)
		diff = last isnt false ?
				Display(now.MinusDays(last.im_num).Round(0)) $ ' day(s) ago' :
				'No Activity'

		return diff
		}

	getLastActivity(channel_num)
		{
		return QueryLast('im_history
			where imchannel_num is ' $ Display(channel_num) $
			' sort im_num')
		}

	colorOrange: 42495
	SetNotificationButtonAsUnread(imchannel_num)
		{
		if false is unreadChannelsList = UserSettings.Get('IM_UnreadChannels')
			unreadChannelsList = Object()

		unreadChannelsList.AddUnique(imchannel_num)
		UserSettings.Put('IM_UnreadChannels',unreadChannelsList)
		if #() isnt data = .listChannels.GetLoadedData()
			{
			if false isnt rec = data.FindOne({ it.imchannel_num is imchannel_num })
				{
				if rec.imchannel_notification isnt .colorOrange
					{
					rec.imchannel_notification = CLR.GREEN
					.listChannels.GetGrid().RepaintRecord(rec)
					}
				}
			}
		}

	SetNotificationButtonAsRead(imchannel_num)
		{
		unreadChannelsList = Object()
		unreadMentionsList = Object()

		if false isnt unreadChannelsList = UserSettings.Get('IM_UnreadChannels')
			unreadChannelsList.Remove(imchannel_num)
		if false isnt unreadMentionsList = UserSettings.Get('IM_MentionedChannels')
			unreadMentionsList.Remove(imchannel_num)

		UserSettings.Put('IM_UnreadChannels',unreadChannelsList)
		UserSettings.Put('IM_MentionedChannels',unreadMentionsList)
		data = .listChannels.GetLoadedData()
		if false isnt rec = data.FindOne({ it.imchannel_num is imchannel_num })
			{
			rec.imchannel_notification = CLR.LIGHTGRAY
			.listChannels.GetGrid().RepaintRecord(rec)
			}
		}

	SetNotificationButtonAsMentioned(imchannel_num)
		{
		if false is unreadChannelsList = UserSettings.Get('IM_MentionedChannels')
			unreadChannelsList = Object()

		unreadChannelsList.AddUnique(imchannel_num)
		UserSettings.Put('IM_MentionedChannels',unreadChannelsList)

		if #() isnt data = .listChannels.GetLoadedData()
			{
			if false isnt rec = data.FindOne({ it.imchannel_num is imchannel_num })
				{
				rec.imchannel_notification = CLR.orange
				.listChannels.GetGrid().RepaintRecord(rec)
				}
			}
		}
	}
