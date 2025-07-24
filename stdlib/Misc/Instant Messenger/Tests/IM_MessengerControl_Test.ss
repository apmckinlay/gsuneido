// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_userFromMessage()
		{
		mock = Mock(IM_MessengerControl)
		mock.IM_MessengerControl_userList = #(user0, user1, user2)
		mock.When.userFromMessage([anyArgs:]).CallThrough()

		msg = Object(im_from: #admin, im_message: 'this is the message')
		Assert(mock.userFromMessage(msg, false) is: #admin)

		msg.im_message = '[this is the message'
		Assert(mock.userFromMessage(msg, false) is: #admin)
		Assert(msg.im_message is: '[this is the message')

		msg.im_message = '\\[this is the message'
		Assert(mock.userFromMessage(msg, false) is: #admin)
		Assert(msg.im_message is: '\\[this is the message')

		msg.im_message = '[user0] \\[this is the message'
		Assert(mock.userFromMessage(msg, false) is: #user0)
		Assert(msg.im_message is: '[user0] \\[this is the message')

		msg.im_message = '[user0+user1] \\[this is the message\\]'
		Assert(mock.userFromMessage(msg, false) is: 'user0+user1')
		Assert(msg.im_message is: '[user0+user1] \\[this is the message\\]')
		}
	}