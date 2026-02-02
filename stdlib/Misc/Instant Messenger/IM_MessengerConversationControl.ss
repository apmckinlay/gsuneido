// Copyright (C) 2022 Axon Development Corporation All rights reserved worldwide.
IM_MessengerTabBase
	{
	channelNum: ''
	oldestImNum: false
	newestImNum: false
	New(.channelNum)
		{
		super(.layout())
		.display = .FindControl('display')
		.display.SetReadOnly(true)
		.editor = .FindControl('editor')
		}

	Startup()
		{
		if .Send('AutoLoadHistory?') is true or .channelNum isnt ''
			.On_History(quite?:)
		}

	layout()
		{
		return [#VertSplit,
			[#Vert,
				['Mshtml', .getHtmlEmpty(), name: 'display',
					xmin: 460]
				#EtchedLine,
				[#HorzEqual
					#(Button 'History', tip: 'Load more history records'),
					#(Button 'Clear'), #Fill]],
			[#Horz,
				[#ScintillaAddonsEditor, xstretch: 60, height: 4, name: 'editor'],
				#(Skip 5),
				[#Vert, [#Button 'Send', width: 12, ystretch: 1]]
				ystretch: 0.1]]
		}

	getHtmlEmpty()
		{
		return .getHtmlHead() $ .getHtmlTail()
		}

	getHtmlHead()
		{
		return '<!DOCTYPE html>' $
			'<html lang="en">' $
			'<head>' $
			.getHtmlRefreshScript() $
			'<meta content="IE=edge" http-equiv="X-UA-Compatible">' $
			.CssStyle() $
			'</head><body onload="scrollFunc()"><div id="IM_container" class="container">'
		}

	getHtmlTail()
		{
		return BookContextMenu() $ '</div></body></html>'
		}

	getHtmlRefreshScript()
		{
		return '<script>' $
			'function scrollFunc()' $
			'{ setTimeout(function(){window.scrollTo(0, 999999999);},0);}' $
			'</script>'
		}

	FocusEditor()
		{
		.editor.SetFocus()
		}

	On_History(quite? = false)
		{
		// messages are ordered newest to oldest
		messages = .getHistory(.tabName, .channelNum, .oldestImNum)

		if messages.NotEmpty?()
			{
			.render(messages, 'afterbegin')
			.oldestImNum = messages.Last().im_num
			if .newestImNum is false
				.newestImNum = messages.First().im_num
			.scrollToMessage(messages[0].im_num)
			}
		else if quite? is false
			InfoWindowControl('No more history', titleSize: 0)
		}

	scrollToMessage(id)
		{
		.display.ScrollIntoView(String(id), false)
		}

	getter_tabName()
		{
		return .Send('GetTabName')
		}

	// expecting messages object to be ordered oldest to newest
	AppendMessages(messages)
		{
		messages = messages.Filter({ it.im_num > .newestImNum })

		.render(messages, 'beforeend')
		if messages.NotEmpty?()
			{
			if .oldestImNum is false
				.oldestImNum = messages.First().im_num
			.newestImNum = messages.Last().im_num
			}
		.scrollToMessage(.newestImNum)
		}

	render(messages, pos)
		{
		if messages.Empty?()
			return

		s = ''
		for msg in messages
			if '' isnt s = .formatMessage(msg)
				.display.InsertAdjacentHTML('IM_container', pos, s)
		}

	formatMessage(msg)
		{
		message = .removeRecipientTag(msg.im_message).
			Replace(`\\\[`, `[`).Replace(`\\]`, `]`)
		return .formatMessageHtml(message, msg.im_num, msg.im_from)
		}

	maxPerHistory: 100
	getHistory(user, channelNum, from = false)
		{
		query = .getHistoryQuery(user, channelNum)
		if from isnt false
			query $= 'where im_num < ' $ Display(from)

		messages = Object()
		QueryApply(query $ ' sort reverse im_num')
			{ |im_hist|
			if IM_MessengerManager.IsSystemMessage?(im_hist)
				continue

			messages.Add(im_hist)
			if messages.Size() > .maxPerHistory /*= max number of history*/
				break
			}
		return messages
		}

	getHistoryQuery(user, channelNum)
		{
		if channelNum isnt ''
			return 'im_history where imchannel_num is ' $ Display(channelNum)
		if IM_MessengerManager.FindUser(user) is false
			return 'im_history where imchannel_num is "" where im_to is ' $ Display(user)
		return 'im_history where imchannel_num is ""
			where (im_from is ' $ Display(Suneido.User) $
				' or im_from is ' $ Display(user) $ ')
				and (im_to is ' $ Display(Suneido.User) $
				' or im_to is ' $ Display(user) $ ')'
		}

	formatMessageHtml(msg, time, user)
		{
		msgHeader = '<div class="msg-header">' $
			Opt(.getFullName(user), " - ") $ .getMessageHeaderDateTime(time) $ '</div>'
		if msg.Has?(`\`)
			msg = msg.Replace(`[\\]`, `\\\\`)
		msg = .encodeHtml(msg)
		msg = .makeLinksClickable(msg)
		msg = '<div class="msg-body">' $ msg $ '</div>'
		id = String(time)
		if user is Suneido.User
			return '<div class="msg-container msg-out" id="' $ id $ '" align="right">' $
				msgHeader $ msg $ '</div>'
		else
			return '<div class="msg-container msg-in" id="' $ id $ '" align="left">' $
				msgHeader $ msg $ '</div>'
		}

	getFullName(user)
		{
		if user is Suneido.User
			return ''
		if false isnt rec = Query1Cached('users where user is ' $ Display(user))
			{
			fmtName = rec.bizuser_name
			if fmtName.Find(',') < fmtName.Size()
				fmtName = fmtName.Split(",").Reverse!().Join(" ").Replace(
					"[^0-9a-zA-Z:,]+"," ").LeftTrim()
			return fmtName
			}
		else
			return ''
		}

	getMessageHeaderDateTime(time)
		{
		week = 7
		timeArg = Date(time)
		timeNow = Timestamp()
		if timeArg.ShortDate() is timeNow.ShortDate()
			return timeArg.Time()
		else if timeNow.MinusDays(timeArg) < week
			return timeArg.ShortDateTime()
		else
			return timeArg.LongDateTime()
		}

	removeRecipientTag(msg)
		{
		return msg.Replace(IM_MessengerControl.RecipientTagPat, "")
		}

	encodeHtml(txt)
		{
		return XmlEntityEncode(XmlEntityDecode(txt))
		}

	makeLinksClickable(msg)
		{
		i = 0
		matches = Object()
		msg = XmlEntityDecode(msg)

		while false isnt m = Addon_url.MatchUrl(msg, i)
			{
			matches.Add(m)
			i = m[0] + m[1]
			}
		result = ''
		start = 0
		for (i = 0; i < matches.Size(); i++)
			{
			url = msg[matches[i][0]::matches[i][1]]
			replace = Sys.SuneidoJs?()
				? Xml('a', url, href: url, 'data-copy-link': url, target: '_blank')
				: Xml('a', url, href: `suneido:/eval?ShellExecute(0,"open","` $
					Url.EncodeQueryValue(url) $ '")', 'data-copy-link': url)
			result $= XmlEntityEncode(msg[start..matches[i][0]]) $ replace
			start = matches[i][0] + matches[i][1]
			}
		result $= XmlEntityEncode(msg[start..])
		return result
		}

	On_Clear()
		{
		.display.Set(.getHtmlEmpty())
		.messages = Object()
		.oldestImNum = false
		.newestImNum = false
		}

	IsTyping?()
		{
		return not .editor.Get().Blank?()
		}

	On_Send()
		{
		if '' is msg = .getMessage()
			return false

		channel? = .channelNum isnt ''
		time_now = Timestamp()
		current_recipient = channel? ? '' : .tabName

		recipients = .getRecipientName(current_recipient, .tabName)

		formattedMsg = .formatMessageHtml(msg, time_now, Suneido.User)
		msg = recipients $ msg.Replace(`\[`, `\\[`).Replace(`]`, `\\]`)
		error_message = ''
		sentStatus = .outputMessage(.channelNum, msg, time_now, current_recipient)
		if sentStatus is false
			error_message = .colourMeRedHtml("!! Error: Message not sent. !!")
		else
			.Send('UpdateContactHistory')

		.display.InsertAdjacentHTML('IM_container', 'beforeend',
			formattedMsg $ error_message)

		if .oldestImNum is false
			.oldestImNum = time_now
		.newestImNum = time_now
		.scrollToMessage(.newestImNum)

		.Send('ClearWarning')
		.editor.Set("")
		.editor.SetFocus()
		}
	outputMessage(channelNum, msg, time_now, current_recipient)
		{
		return channelNum isnt ''
			? IM_MessengerManager.OutputMessage('GLOBAL', msg, time_now, channelNum)
			: IM_MessengerManager.OutputMessage(current_recipient, msg, time_now)
		}

	messageLengthLimit: 5000
	getMessage()
		{
		msg = .editor.Get().Trim()
		msgSize = msg.Size()
		if msgSize > .messageLengthLimit
			{
			.Send('ShowWarning', "Limit message length to " $ .messageLengthLimit $
				" characters. (current size " $ msgSize $" characters)")
			msg = ''
			}
		return msg
		}

	getRecipientName(current_recipient, tabName)
		{
		return false is IM_MessengerManager.FindUser(current_recipient)
			? '[' $ tabName $ '] '
			: ''
		}

	colourMeRedHtml(err)
		{
		return '<span style="color:red">' $ err $ '</span>'
		}

	Enter_Pressed(pressed = false)
		{
		if KeyPressed?(VK.SHIFT, :pressed) or .Send('EnterToSend?') isnt true
			return 0
		.On_Send()
		return false // disable ENTER
		}
	}