// Copyright (C) 2018 Axon Development Corporation All rights reserved worldwide.
Test
	{
	Setup()
		{
		.cur = Suneido.User
		}
	Test_recentContactList()
		{
		contact = IM_ContactsControl
			{
			IM_ContactsControl_groupListFunc: false
			ListContact: class
				{
				Get()
					{
					return #([im_user: "admin", im_status: "Logged In"],
						[im_user: "p1", im_status: ""],
						[im_user: "p2", im_status: ""],
						[im_user: "lq", im_status: "Logged In"],
						[im_user: "sal1", im_status: ""])
					}
				}
			IM_ContactsControl_query(block)
				{
				recentConv = #(
					[im_to: "p1+p2+lq", im_message: "[p1+p2+lq] Grp", im_from: "p1"],
					[im_to: "p1", im_message: "hi", im_from: "lq"],
					[im_to: "p2", im_message: "hi", im_from: "p1"],
					[im_to: "gzj", im_message: "hi", im_from: "p1"])
				recentConv.Each(block)
				}
			IM_ContactsControl_contactsImpl: #()
			}
		Suneido.User = "p1"
		recentContact = #([im_to: "p1+p2+lq", im_status: ""],
			[im_to: "lq", im_status: "Logged In"],
			[im_to: "p2", im_status: ""])
		Assert(contact.IM_ContactsControl_recentContactList() is: recentContact)

		Suneido.User = "p2"
		recentContact = #([im_to: "p1+p2+lq", im_status: ""],
			[im_to: "p1", im_status: ""])
		Assert(contact.IM_ContactsControl_recentContactList() is: recentContact)

		Suneido.User = "sal1"
		recentContact = #()
		Assert(contact.IM_ContactsControl_recentContactList() is: recentContact)
		}

	Test_group?()
		{
		cl = IM_ContactsControl
			{
			IM_ContactsControl_contactsImpl: #(
				parseGroup: 'IM_ContactsControl_Test.ParseGroup')
			}
		cl.IM_ContactsControl_contactsImpl
		fn = cl.IM_ContactsControl_group?

		Suneido.User = 'TestUser'
		result = fn([im_to: 'Grp'], Object(), Object('Grp'))
		Assert(result)
		result = fn([im_to:'Prgm'], Object(), Object('Grp'))
		Assert(result is: false)
		result = fn([im_to:'Prgm'], Object(), Object('Prgm'))
		Assert(result is: false)
		result = fn([im_to:'Prgm'], Object(), Object('Group0', 'Group1'))
		Assert(result is: false)
		result = fn([im_to:'Grp'], Object(), Object('Group0', 'Grp'))
		Assert(result)

		result = fn([im_to:'Grp'], Object(Object(im_to:'p1')), Object('Grp'))
		Assert(result)
		result = fn([im_to:'Grp'], Object(Object(im_to:'Grp')), Object('Grp'))
		Assert(result is: false)

		result = fn([im_to:'Grp'], Object(Object(im_to:'Grp+p2')), Object('Grp'))
		Assert(result)
		result = fn([im_to:'Grp+p2'], Object(Object(im_to:'Grp+p2')), Object('Grp'))
		Assert(result is: false)

		result = fn([im_to:'Grp+p2'], Object(), Object('Grp'))
		Assert(result)
		result = fn([im_to:'Grp+p2'], Object(Object(im_to: 'Grp')), Object('Grp'))
		Assert(result)
		}

	groups: (Prgm: (A B C D E), Grp: (TestUser))
	ParseGroup(groupName) { return .groups.GetDefault(groupName, Object()) }

	Teardown()
		{
		Suneido.User = .cur
		super.Teardown()
		}
	}