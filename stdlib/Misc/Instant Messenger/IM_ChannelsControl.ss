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
		.beforeRecord(rec, UserSettingsCached().Get('IM_SubscribedChannels'))
		}

	beforeRecord(rec, cachedSubscriptions)
		{
		if Object?(cachedSubscriptions)
			rec.imchannel_joinedStatus = cachedSubscriptions.Has?(rec.imchannel_num)

		if .getUserSettingOb('IM_MentionedChannels').Has?(rec.imchannel_num)
			rec.imchannel_notification = CLR.orange // orange button takes priority
		else if .getUserSettingOb('IM_UnreadChannels').Has?(rec.imchannel_num)
			rec.imchannel_notification = CLR.GREEN

		last = .getLastActivity(rec.imchannel_num)
		rec.imchannel_lastActivity = last isnt false
			? Display(Timestamp().MinusDays(last.im_num).Round(0)) $ ' day(s) ago'
			: 'No Activity'
		}

	getUserSettingOb(key)
		{
		return Object?(setting = UserSettings.Get(key))
			? setting
			: Object()
		}

	getLastActivity(channel_num)
		{
		return QueryLast('im_history
			where imchannel_num is ' $ Display(channel_num) $
			' sort im_num')
		}

	VirtualList_DoubleClick(rec, col /*unused*/, source)
		{
		if source is .listChannels
			if .listChannels.NotEmpty?() and rec isnt false
				{
				.SetNotify(rec.imchannel_num, read?:)
				.listChannels.ReloadRecord(rec)
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
		num_param = String(args.rec.imchannel_num)
		update_query = "update im_channels where imchannel_num is " $ num_param $
			" set imchannel_status = 'inactive'"
		QueryDo(update_query)
		.listChannels.Refresh()
		.Send("CloseTab",args.rec)
		IM_MessengerManager.SendSystemMessage(.maintMsg)
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

	Reload(channelNums = false)
		{
		if channelNums is false
			.listChannels.Refresh()
		else
			for channel in .listChannels.GetLoadedData()
				if channelNums.Has?(channel.imchannel_num)
					.listChannels.ReloadRecord(channel)
		}

	SetNotify(imchannel_num, read? = false, mentioned? = false)
		{
		unreadChannelsList = .getUserSettingOb('IM_UnreadChannels')
		unreadMentionsList = .getUserSettingOb('IM_MentionedChannels')
		unreadChannels = unreadChannelsList.Size()
		unreadMentions = unreadMentionsList.Size()
		if read?
			{
			unreadChannelsList.Remove(imchannel_num)
			unreadMentionsList.Remove(imchannel_num)
			}
		else if mentioned?
			unreadMentionsList.AddUnique(imchannel_num)
		else
			unreadChannelsList.AddUnique(imchannel_num)
		UserSettings.Put('IM_UnreadChannels', unreadChannelsList)
		UserSettings.Put('IM_MentionedChannels', unreadMentionsList)
		return unreadChannels < unreadChannelsList.Size() or
			unreadMentions < unreadMentionsList.Size()
		}
	}
