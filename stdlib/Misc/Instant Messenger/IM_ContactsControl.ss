// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Title: "Contacts"
	Xmin: 400
	New()
		{
		.subs = [
			PubSub.Subscribe('messengerContacts', .refreshContactInfo)
			PubSub.Subscribe('messengerContacts', .refreshRecentContactList)
			]
		.contactSelected = Object().Set_default(false)
		.userListFunc = Global(.contactsImpl.userList)
		.groupListFunc = .contactsImpl.Member?('groupList')
			? Global(.contactsImpl.groupList)
			: false
		.ListContact = .FindControl('ListContact')
		.ListHistory = .FindControl('ListHistory')
		.refreshContactInfo()
		.UpdateContactHistory()
		}

	UpdateContactHistory()
		{
		.ListHistory.Set(.recentContactList())
		}

	columnsContact: ("im_selected", "im_user", "im_status")
	columnsHistory: ("im_to", "im_status")
	Controls()
		{
		.contactsImpl = SoleContribution('IM_ContactsImplementation')
		contactcolumns = .columnsContact.Copy()
		if .contactsImpl.Member?('extraCols')
			contactcolumns.MergeUnion(.contactsImpl.extraCols)
		return Object('Vert',
			Object('VertSplit',
				Object('Vert',
					Object('List', name: 'ListContact', columns: contactcolumns,
						noDragDrop:, defWidth: false, xstretch: 1, ystretch: 3,
						xmin: 400, multiselect:, columnsSaveName: .Title, resetColumns:,
						checkBoxColumn: 'im_selected')
					#(HorzEqual
						Skip
						(Button, "Message with ALL Users",xstretch: 1)
						Skip
						(Button, "Message with Selected Users", xstretch: 1)
						Skip)),
				Object('ListStretch', name: 'ListHistory', columns: .columnsHistory,
					noDragDrop:, defWidth: false, xstretch: 1, ystretch: 1, xmin: 100,
					columnsSaveName: .Title $ ' - History'))
				name:'Contacts_tabs', bottom:)
		}
	userAvailableList()
		{
		avail = (.userListFunc)().RemoveIf({ it.im_user is Suneido.User })
		if false isnt .groupListFunc
			avail.MergeUnion((.groupListFunc)())
		return avail
		}
	recentContactList()
		{
		status = .ListContact.Get()
		recentContact = Object()
		groupList = Object()
		if .groupListFunc isnt false
			groupList = (.groupListFunc)().Map({ it.im_user })
		try .query()
			{
			if it.im_from is Suneido.User
				.addContact(it.im_to, recentContact, status)
			else if it.im_to.Split('+').Has?(Suneido.User)
				.addContact(
					it.im_to.Has?('+') ? it.im_to : it.im_from, recentContact, status)
			else if .group?(it, recentContact, groupList)
				.addContact(it.im_to, recentContact, status)
			}
		catch (unused, 'max_number_of_recent')
			return recentContact
		return recentContact
		}

	group?(conversation, recentContact, groupList)
		{
		if not .contactsImpl.Member?('parseGroup') or
			recentContact.HasIf?( {|x| x.im_to is conversation.im_to } )
			return false

		recipients = conversation.im_to.Split('+')
		for recipient in recipients
			if groupList.Has?(recipient) and
				Global(.contactsImpl.parseGroup)(recipient).Has?(Suneido.User)
				return true
		return false
		}

	refreshRecentContactList()
		{
		status = .ListContact.Get()
		recent = .ListHistory.Get()
		for i, item in recent
			if ((false isnt toRec = status.FindOne({ it.im_user is item.im_to })) and
				toRec.im_status isnt item.im_status)
				{
				item.im_status = toRec.im_status
				.ListHistory.RepaintRow(i)
				}
		}

	query(block)
		{
		QueryApply('im_history where imchannel_num is "" ' $
			' where im_num > ' $ Display(Date().Minus(months: 1)) $
			' sort reverse im_num')
			{
			block(it)
			}
		}

	addContact(im_to, recentContact, status)
		{
		if recentContact.HasIf?( { it.im_to is im_to } )
			return
		if im_to.Has?('+')
			recentContact.Add(Object(:im_to, im_status: ''))
		else if false isnt toRec = status.FindOne({ it.im_user is im_to })
			recentContact.Add(Object(:im_to, im_status: toRec.im_status))
		if recentContact.Size() >= 10 /*= max number of recent contacts*/
			throw 'max_number_of_recent' // will be caught to stop listing
		}

	refreshContactInfo()
		{
		data = .userAvailableList()
		data.Each({ it.im_selected = .contactSelected[it.im_user] })
		sortCol = .ListContact.GetSortCol()
		.ListContact.DoWithCurrentVScrollPos()
			{
			.ListContact.Set(data)
			for (i = 0; i < data.Size(); ++i)
				if data[i].GetDefault(#im_highlight, '') isnt ''
					.ListContact.AddHighlight(i, data[i].im_highlight)
			.ListContact.SortListData(sortCol)
			}
		}

	List_WantEditField(col /*unused*/) { return false }

	List_WantNewRow() { return false }

	List_DoubleClick(row, col, source)
		{
		if(source is .ListContact)
			.doubleClickContact(row, col)
		else
			.doubleClickRecentConversation(row)
		return 0
		}

	doubleClickContact(row, col)
		{
		if .ListContact.GetCol(col) is 'im_selected'
			return
		if .ListContact.Empty?() isnt true and row isnt false
			{
			clickedRecipient = .ListContact.Get()[row].im_user
			conversationName = .createConversationName(clickedRecipient)
			.setConversation(conversationName)
			}
		}

	createConversationName(@selectedNames)
		{
		if selectedNames.Size() is 1 and
			false is IM_MessengerManager.FindGroup(selectedNames[0])
			return selectedNames[0]
		else
			{
			conversationName = selectedNames.Sort!().Join('+')
			if not IM_MessengerManager.ParseRecipients(conversationName).Has?(
				Suneido.User)
				conversationName = selectedNames.Add(Suneido.User).Sort!().Join('+')

			return conversationName
			}
		}

	doubleClickRecentConversation(row)
		{
		if .ListHistory.Empty?() isnt true and row isnt false
			{
			conversationName = .ListHistory.Get()[row].im_to
			.setConversation(conversationName)
			}
		}

	List_AfterToggle(data)
		{
		.contactSelected[data.im_user] = data.im_selected
		}

	List_SingleClick(row, col, source)
		{
		if row is false or
			source isnt .ListContact or
			.ListContact.Empty?() or
			.ListContact.GetCol(col) isnt 'im_selected'
			return 0

		originalSelected = .ListContact.GetRow(row).im_selected
		.contactSelected[.ListContact.GetRow(row).im_user] = not originalSelected
		.ListContact.GetRow(row).im_selected = not originalSelected
		.ListContact.RepaintRow(row)
		return 0
		}

	setConversation(user)
		{
		.Send('OpenConversation', user)
		}

	On_Message_with_ALL_Users()
		{
		.Send('OpenConversation', 'GLOBAL')
		}

	On_Message_with_Selected_Users()
		{
		if .ListContact.Empty?() is true
			return

		nameSelected = .contactSelected.MembersIf( { .contactSelected[it] is true } )
		if nameSelected.Empty?()
			{
			.AlertInfo('Instant Messenger',
				'Please select at least one user to chat with.')
			return
			}

		conversationName = .createConversationName(@nameSelected)
		.setConversation(conversationName)
		.contactSelected = Object().Set_default(false)
		.refreshContactInfo()
		.refreshRecentContactList()
		}

	Destroy()
		{
		.subs.Each(#Unsubscribe)
		IM_MessengerManager.Request("/IM_ContactsUpdate")
		super.Destroy()
		}
	}
