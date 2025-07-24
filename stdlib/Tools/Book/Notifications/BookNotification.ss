// Copyright (C) 2006 Suneido Software Corp. All rights reserved worldwide.
class
	{
	Title: 'New'
	//TODO rename this method - it doesn't return events
	NewEvents(user) // runs on the HTTP Server
		{
		try
			{
			Transaction(read:)
				{ |t|
				query = .MessageQuery(user) $ ' sort notify_message_num'
				lastMsg = Cursor(query).Prev(t) // using Cursor to prevent temp index
				}
			return lastMsg isnt false ? lastMsg.notify_message_num : false
			}
		catch (unused, '*nonexistent table: notifications')
			return false
		}

	MessageQuery(user = false)
		{
		if user is false
			user = Suneido.User
		return 'notifications' $
			' where notify_users_unread isnt ""' $
			' where notify_message_num > ' $ Display(Date().NoTime().Plus(months: -1)) $
			' and BookNotifyIncludes?(notify_users_unread, ' $ Display(user) $ ')'
		}

	HandleNewEvents(newMessageInfo) // runs on the client
		{
		newmessages? = false
		if Date?(newMessageInfo)
			{
			lastShown = Suneido.Member?(#ShowNewEntries)
				? Suneido.ShowNewEntries.time
				: Date.Begin()
			newmessages? = newMessageInfo > lastShown
			}
		else // if newMessageInfo isnt a date, should be false
			newmessages? = newMessageInfo

		showhide = Suneido.GetDefault('ShowAll', '')
		if showhide isnt ''
			newmessages? = showhide is true
		return newmessages?
		}

	Count(user = false)
		{
		return QueryCount(.MessageQuery(:user))
		}

	// ServerEval'ed by NotificationsHtml
	GetNewEntries(user, from = 0, to = 50)
		{
		entries = Object()
		count = 0
		states = Object(prevTitle: '', prevOption: '')
		.foreachNotification(:user)
			{
			if ++count > to
				break
			if count <= from
				continue
			it.notify_users_message = .formatMessage(it)
			.BuildOneEntry(it, states, entries)
			}
		return entries
		}
	foreachNotification(block, user = false) // broken out for testing
		{
		QueryApply(.MessageQuery(user) $
			' sort title, notify_message_heading, notify_message_num', block)
		}
	formatMessage(x) // broken out for testing
		{
		fmt = .format(x)
		return .FormatMessage(fmt.columns, fmt.linkOb, x)
		}

	BuildOneEntry(x, states, entries)
		{
		states.prevTitle = .addTitle(x.title, states.prevTitle, entries)
		states.prevOption = .addOptionHeading(x, states.prevOption, entries)
		entries.Add(Xml('li', x.notify_users_message.Replace(
			'((?i)invalid)', '<span style="font-weight:bold;color:red">\1</span>')))
		}

	addTitle(title, prevTitle, entries)
		{
		if title isnt prevTitle
			{
			entries.Add(Xml('h1', title))
			prevTitle = title
			}
		return prevTitle
		}

	addOptionHeading(x, prevOption, entries)
		{
		option = x.notify_message_heading
		if option isnt prevOption
			{
			if option isnt ''
				entries.Add(Xml('h3', option))
			prevOption = option
			}
		return prevOption
		}

	BuildReleaseNoteHelpLinks(links, book)
		{
		if links.Empty?()
			return ''

		helpLinkOb = links.Copy()
		helpLinkOb.Map!({ Xml('a',
			it.Trim('/').Replace('/Reference', '').Replace('/', ' > '),
			href: "suneido:/eval?OpenBook(" $ Display(book) $ "," $ "path:" $
				Display(it) $ ")")
			})
		helpLinksMsg = '<li>For more information see: <ul><li>' $
			helpLinkOb.Join('</li><li>') $ '</li></ul></li>'
		return helpLinksMsg
		}

	ClearNewEntries(from = 0, to = 50)
		{
		count = 0
		.foreachNotification()
			{
			if .messageNewerThanWindow?(it)
				continue
			if ++count > to
				break
			if count <= from
				continue

			RetryTransaction()
				{ |t|
				.ClearOneEntry(t, it)
				}
			}
		}

	ClearOneEntry(t, notification, user = false)
		{
		t.QueryApply1(notification.table $ ' where ' $ notification.numField $ ' is ' $
			Display(notification.notify_message_num))
			{ |x|
			users = .bookNotification_UserList(x[notification.notifyField])
			users.Remove(user is false ? Suneido.User : user)
			x[notification.notifyField] = users.Join(',')
			x.Update()
			}
		}

	bookNotification_UserList(userList)
		{
		return BookNotification_UserList(userList)
		}

	messageNewerThanWindow?(x)
		{
		return ((false isnt date = .getDate()) and (x.notify_message_num >= date))
		}

	getDate()
		{
		return Suneido.Member?(#ShowNewEntries) ? Suneido.ShowNewEntries.time : false
		}

	// used by Plugin_NotifyTables
	BuildView()
		{
		viewTables = Object()
		Plugins().ForeachContribution('NotifyUser', 'tables')
			{ |c|
			str = '(' $ c.query $
				' rename ' $ c.numField $ ' to notify_message_num, ' $
				c.notifyField $ ' to notify_users_unread' $
				' extend numField = ' $ Display(c.numField) $
				', notifyField = ' $ Display(c.notifyField) $
				', notify_message_heading = ' $ c.GetDefault('headingField', '""') $
				', title = ' $ Display(c.title) $
				', table = ' $ Display(c.table) $
				')'
			viewTables.Add(str)
			}
		return viewTables.Join('\nunion\n')
		}

	info: MemoizeSingle
		{
		Func()
			{
			info = Object()
			Plugins().ForeachContribution('NotifyUser', 'tables')
				{|c|
				info[c.title $ ' - ' $ c.table] =
					Object(columns: c.columns, linkOb: c.GetDefault('link', #()))
				}
			return info
			}
		}
	format(x)
		{
		return (.info)()[x.title $ ' - ' $ x.table]
		}

	FormatMessage(columns, linkOb, rec)
		{
		str = ''
		for col in columns
			{
			if rec[col] is ''
				continue

			prompt = PromptOrHeading(col)
			if prompt is 'Num'
				prompt = 'Created On'
			s = (prompt is '' ? '' : prompt $ ': ')
			if linkOb.Member?(col)
				{
				link = linkOb[col]
				s $= Xml('a', rec[col],
					href: "suneido:/eval?AccessGoTo(" $
						Display(.encodeHashTag(String(link[0]))) $ "," $
						Display(link[1]) $ "," $
						Display(rec[col]) $ ",window:" $ Display('Dialog') $ ")")
				}
			else
				s $= .Format_dates(rec[col])
			str $= s $ '&emsp;'
			}
		return str.BeforeLast('&emsp;')
		}

	// to handle passing "#(Biz_Employees, 'Trucking')" to AccessGoTo where "#" is
	// special in URL
	encodeHashTag(s)
		{
		return s.Replace('#', '%23')
		}

	Format_dates(data)
		{
		if not Date?(data)
			return data
		return data.NoTime?() ? data.StdShortDate() : data.StdShortDateTime()
		}
	}
