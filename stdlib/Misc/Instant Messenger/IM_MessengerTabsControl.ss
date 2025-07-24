// Copyright (C) 2022 Axon Development Corporation All rights reserved worldwide.
TabsControl
	{
	tabPosOffset: 1 // the Home tab
	New(@args)
		{
		super(@args)
		.home = .GetControl(0)
		.SetImageList(Object(.unreadMessageImage()))
		}

	unreadMessageImage()
		{
		codeOb = IconFont().MapToCharCode(#document)
		return Sys.SuneidoJs?()
			? Object(char: codeOb.char, font: codeOb.font, color: #black)
			: Object(ImageFont(codeOb.char, codeOb.font), CLR.BLACK)
		}

	Startup()
		{
		}

	UpdateChangedTabs()
		{
		activeChannels = IM_MessengerManager.GetActiveChannelsListByNum()
		toRemove = Object()
		.forEachTab()
			{ |i, data|
			if data.imchannel_num isnt ''
				{
				newName = IM_MessengerManager.GetChannelNameFromNum(data.imchannel_num)
				if not activeChannels.Has?(data.imchannel_num)
					toRemove.Add(i)
				else if newName isnt .TabName(i)
					.SetTabData(i, #(), name: newName)
				}
			}
		toRemove.Sort!(Gt).Each(.Remove)
		}

	ChannelOpen?(num)
		{
		.forEachTab()
			{ |data|
			if data.imchannel_num is num
				return true
			}
		return false
		}

	forEachTab(block)
		{
		for (i = .tabPosOffset; i < .GetAllTabCount(); i++)
			block(ctrl: .GetControl(i), :i, data: .GetTabData(i))
		}

	AddMessage(name, message, channelNum)
		{
		i = .OpenTab(name, channelNum, noSelect:)
		.GetTabData(i).pendingMessages.Add(message)
		.SetUnread(true, i)
		}

	OpenTab(name, channelNum = '', noSelect = false)
		{
		if false isnt i = .FindTab(name)
			{
			if noSelect is false
				.Select(i)
			return i
			}

		.Insert(name, Object('IM_MessengerConversation', channelNum, Tab: name),
			data: [imchannel_num: channelNum, pendingMessages: Object()],
			:noSelect)
		return .GetAllTabCount() - 1
		}

	Tab_Close(i)
		{
		.Remove(i)
		.Select(Max(.GetSelected(), 0))
		}

	GetTabName(source)
		{
		if false is i = .findControl(source)
			return '???'
		return .TabName(i)
		}

	findControl(target)
		{
		.forEachTab()
			{ |ctrl, i|
			if Same?(ctrl, target)
				return i
			}
		return false
		}

	AutoLoadHistory?()
		{
		return .home.AutoLoadHistory?()
		}

	EnterToSend?()
		{
		return .home.EnterToSend?()
		}

	IsTyping?()
		{
		if false is ctrl = .GetControl()
			return false
		return ctrl.IsTyping?()
		}

	SetUnread(unread?, i)
		{
		.SetImage(i, unread? ? 0 : -1)
		}

	SelectTab(@args)
		{
		super.SelectTab(@args)
		.renderPendingMessages(args.GetDefault(0, { args.i }))
		}

	renderPendingMessages(i)
		{
		data = .GetTabData(i)
		messages = data.GetDefault(#pendingMessages, #())
		ctrl = .GetControl(i)
		if messages.NotEmpty?()
			{
			ctrl.AppendMessages(messages)
			messages.Delete(all:)
			}
		.SetUnread(false, i)
		}

	SelectByName(name)
		{
		if  false is i = .FindTab(name)
			return

		if .TabGetSelectedName() isnt name
			.Select(i)
		else
			.renderPendingMessages(i)
		}

	FocusEditor()
		{
		if false isnt ctrl = .GetControl()
			ctrl.FocusEditor()
		}
	}
