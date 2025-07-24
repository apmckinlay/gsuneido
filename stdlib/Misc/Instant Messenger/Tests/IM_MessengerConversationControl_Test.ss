// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_formatMessage()
		{
		cl = IM_MessengerConversationControl
			{
			IM_MessengerConversationControl_getMessageHeaderDateTime(unused)
				{
				return 'August 7, 2024 12:01 PM'
				}
			IM_MessengerConversationControl_getFullName(unused) { return '' }
			}
		fn = cl.IM_MessengerConversationControl_formatMessage
		msg = [im_message: '', im_from: 'user1', im_num: #20240807.1201]
		Assert(fn(msg) is: '<div class="msg-container msg-in" id="#20240807.1201" ' $
			'align="left"><div class="msg-header">August 7, 2024 12:01 PM</div>' $
			'<div class="msg-body"></div></div>')

		msg.im_message = 'hello world'
		Assert(fn(msg) is: '<div class="msg-container msg-in" id="#20240807.1201" ' $
			'align="left"><div class="msg-header">August 7, 2024 12:01 PM</div>' $
			'<div class="msg-body">hello world</div></div>')

		msg.im_message = 'rec\[test\] = #()'
		Assert(fn(msg) is: '<div class="msg-container msg-in" id="#20240807.1201" ' $
			'align="left"><div class="msg-header">August 7, 2024 12:01 PM</div>' $
			'<div class="msg-body">rec[test] = #()</div></div>')

		msg.im_message = `File is in \\share\work\myfile.txt`
		Assert(fn(msg) is: '<div class="msg-container msg-in" id="#20240807.1201" ' $
			'align="left"><div class="msg-header">August 7, 2024 12:01 PM</div>' $
			`<div class="msg-body">File is in \\\\share\\work\\myfile.txt</div></div>`)

		msg.im_message = '[user1+user2+user3] hello world'
		Assert(fn(msg) is: '<div class="msg-container msg-in" id="#20240807.1201" ' $
			'align="left"><div class="msg-header">August 7, 2024 12:01 PM</div>' $
			'<div class="msg-body">hello world</div></div>')

		msg.im_message = '[user1+user2+user3] hello\[test\] world'
		Assert(fn(msg) is: '<div class="msg-container msg-in" id="#20240807.1201" ' $
			'align="left"><div class="msg-header">August 7, 2024 12:01 PM</div>' $
			'<div class="msg-body">hello[test] world</div></div>')

		msg.im_message = '[GLOBAL] hello\[test\] world'
		Assert(fn(msg) is: '<div class="msg-container msg-in" id="#20240807.1201" ' $
			'align="left"><div class="msg-header">August 7, 2024 12:01 PM</div>' $
			'<div class="msg-body">hello[test] world</div></div>')

		msg.im_message = '[channel name]     some\[content\] is provided'
		Assert(fn(msg) is: '<div class="msg-container msg-in" id="#20240807.1201" ' $
			'align="left"><div class="msg-header">August 7, 2024 12:01 PM</div>' $
			'<div class="msg-body">some[content] is provided</div></div>')
		}

	Test_makeLinksClickable()
		{
		fn = IM_MessengerConversationControl.
			IM_MessengerConversationControl_makeLinksClickable

		Assert(fn('') is: '')
		Assert(fn('This is a test message.') is: 'This is a test message.')
		Assert(fn('Check out https://www.example.com/')
			is: 'Check out <a href="suneido:/eval?ShellExecute(0,&quot;open&quot;,' $
				'&quot;https%3A%2F%2Fwww.example.com%2F&quot;)">' $
				'https://www.example.com/</a>')
		Assert(fn(XmlEntityEncode(`Visit "https://www.example.com/" for more info. ` $
			`Also, check out www.google.com/`))
			is: 'Visit &quot;<a href="suneido:/eval?ShellExecute(0,&quot;open&quot;,' $
				'&quot;https%3A%2F%2Fwww.example.com%2F&quot;)">' $
				'https://www.example.com/</a>&quot; for more info. Also, check out ' $
				'<a href="suneido:/eval?ShellExecute(0,&quot;open&quot;,' $
				'&quot;www.google.com%2F&quot;)">www.google.com/</a>')
		Assert(
			fn('Click &lt;a href=&quot;https://www.example.com/&quot;&gt;here&lt;/a&gt;')
			like: 'Click &lt;a href=&quot;' $
				'<a href="suneido:/eval?ShellExecute(0,&quot;open&quot;,' $
				'&quot;https%3A%2F%2Fwww.example.com%2F&quot;)">' $
				'https://www.example.com/</a>&quot;&gt;here&lt;/a&gt;')
		Assert(fn(
			XmlEntityEncode('This message contains a URL with special characters: ' $
				'http://www.example.com/page?query=hello%20world'))
			is: 'This message contains a URL with special characters: ' $
				'<a href="suneido:/eval?ShellExecute(0,&quot;open&quot;,' $
				'&quot;http%3A%2F%2Fwww.example.com%2F' $
				'page%3Fquery%3Dhello%2520world&quot;)">' $
				'http://www.example.com/page?query=hello%20world</a>')
		}

	Test_messageOrder()
		{
		// send > load
		mock = .setupMock()
		_testImOb = Object()
		mock.On_Send()
		Assert(mock.IM_MessengerConversationControl_oldestImNum isnt: false)
		Assert(mock.IM_MessengerConversationControl_newestImNum isnt: false)
		Assert(mock.IM_MessengerConversationControl_newestImNum is:
			mock.IM_MessengerConversationControl_oldestImNum)
		mock.IM_MessengerConversationControl_display.Verify.Times(1).
			InsertAdjacentHTML([anyArgs:])
		Assert(_testImOb isSize: 1)
		Assert(XmlParser(_testImOb[0]).Attributes().class is: "msg-container msg-out")
		mock.Verify.Times(1).scrollToMessage([anyArgs:])
		mock.Verify.scrollToMessage([greaterThan: #20010101.154521929])
		mock.On_History(true)
		Assert(mock.IM_MessengerConversationControl_oldestImNum is: #20010101.154521927)
		Assert(mock.IM_MessengerConversationControl_newestImNum
			greaterThan: #20010101.154521929)
		mock.IM_MessengerConversationControl_display.Verify.Times(4).
			InsertAdjacentHTML([anyArgs:])
		Assert(_testImOb isSize: 4)
		Assert(XmlParser(_testImOb[0]).Attributes().id is: "#20010101.154521927")
		Assert(XmlParser(_testImOb[1]).Attributes().id is: "#20010101.154521928")
		Assert(XmlParser(_testImOb[2]).Attributes().id is: "#20010101.154521929")
		Assert(XmlParser(_testImOb[3]).Attributes().class is: "msg-container msg-out")
		mock.Verify.Times(2).scrollToMessage([anyArgs:])
		mock.Verify.scrollToMessage(#20010101.154521929)

		// load > send
		mock = .setupMock()
		_testImOb = Object()
		mock.On_History(true)
		Assert(mock.IM_MessengerConversationControl_oldestImNum isnt: false)
		Assert(mock.IM_MessengerConversationControl_newestImNum isnt: false)
		Assert(mock.IM_MessengerConversationControl_oldestImNum is: #20010101.154521927)
		Assert(mock.IM_MessengerConversationControl_newestImNum is: #20010101.154521929)
		mock.IM_MessengerConversationControl_display.Verify.Times(3).
			InsertAdjacentHTML([anyArgs:])
		Assert(_testImOb isSize: 3)
		Assert(XmlParser(_testImOb[0]).Attributes().id is: "#20010101.154521927")
		Assert(XmlParser(_testImOb[1]).Attributes().id is: "#20010101.154521928")
		Assert(XmlParser(_testImOb[2]).Attributes().id is: "#20010101.154521929")
		mock.Verify.Times(1).scrollToMessage([anyArgs:])
		mock.Verify.scrollToMessage(#20010101.154521929)
		mock.On_Send()
		Assert(mock.IM_MessengerConversationControl_newestImNum
			greaterThan: #20010101.154521929)
		mock.IM_MessengerConversationControl_display.Verify.Times(4).
			InsertAdjacentHTML([anyArgs:])
		Assert(_testImOb isSize: 4)
		Assert(XmlParser(_testImOb[0]).Attributes().id is: "#20010101.154521927")
		Assert(XmlParser(_testImOb[1]).Attributes().id is: "#20010101.154521928")
		Assert(XmlParser(_testImOb[2]).Attributes().id is: "#20010101.154521929")
		Assert(XmlParser(_testImOb[3]).Attributes().class is: "msg-container msg-out")
		mock.Verify.Times(2).scrollToMessage([anyArgs:])
		mock.Verify.scrollToMessage([greaterThan: #20010101.154521929])

		// receive > load
		mock = .setupMock()
		_testImOb = Object()
		mock.AppendMessages(.append)
		Assert(mock.IM_MessengerConversationControl_oldestImNum is: #20010101.154521940)
		Assert(mock.IM_MessengerConversationControl_newestImNum is: #20010101.154521942)
		mock.IM_MessengerConversationControl_display.Verify.Times(3).
			InsertAdjacentHTML([anyArgs:])
		Assert(_testImOb isSize: 3)
		Assert(XmlParser(_testImOb[0]).Attributes().id is: "#20010101.154521940")
		Assert(XmlParser(_testImOb[1]).Attributes().id is: "#20010101.154521941")
		Assert(XmlParser(_testImOb[2]).Attributes().id is: "#20010101.154521942")
		mock.Verify.Times(1).scrollToMessage([anyArgs:])
		mock.Verify.scrollToMessage(#20010101.154521942)
		mock.On_History(true)
		Assert(mock.IM_MessengerConversationControl_oldestImNum is: #20010101.154521927)
		Assert(mock.IM_MessengerConversationControl_newestImNum is: #20010101.154521942)
		mock.IM_MessengerConversationControl_display.Verify.Times(6).
			InsertAdjacentHTML([anyArgs:])
		Assert(_testImOb isSize: 6)
		Assert(XmlParser(_testImOb[0]).Attributes().id is: "#20010101.154521927")
		Assert(XmlParser(_testImOb[1]).Attributes().id is: "#20010101.154521928")
		Assert(XmlParser(_testImOb[2]).Attributes().id is: "#20010101.154521929")
		Assert(XmlParser(_testImOb[3]).Attributes().id is: "#20010101.154521940")
		Assert(XmlParser(_testImOb[4]).Attributes().id is: "#20010101.154521941")
		Assert(XmlParser(_testImOb[5]).Attributes().id is: "#20010101.154521942")
		mock.Verify.Times(2).scrollToMessage([anyArgs:])
		mock.Verify.scrollToMessage(#20010101.154521929)

		// load > receive
		mock = .setupMock()
		_testImOb = Object()
		mock.On_History(true)
		Assert(mock.IM_MessengerConversationControl_oldestImNum is: #20010101.154521927)
		Assert(mock.IM_MessengerConversationControl_newestImNum is: #20010101.154521929)
		mock.IM_MessengerConversationControl_display.Verify.Times(3).
			InsertAdjacentHTML([anyArgs:])
		Assert(_testImOb isSize: 3)
		Assert(XmlParser(_testImOb[0]).Attributes().id is: "#20010101.154521927")
		Assert(XmlParser(_testImOb[1]).Attributes().id is: "#20010101.154521928")
		Assert(XmlParser(_testImOb[2]).Attributes().id is: "#20010101.154521929")
		mock.Verify.Times(1).scrollToMessage([anyArgs:])
		mock.Verify.scrollToMessage(#20010101.154521929)
		mock.AppendMessages(.append)
		Assert(mock.IM_MessengerConversationControl_oldestImNum is: #20010101.154521927)
		Assert(mock.IM_MessengerConversationControl_newestImNum is: #20010101.154521942)
		mock.IM_MessengerConversationControl_display.Verify.Times(6).
			InsertAdjacentHTML([anyArgs:])
		Assert(_testImOb isSize: 6)
		Assert(XmlParser(_testImOb[0]).Attributes().id is: "#20010101.154521927")
		Assert(XmlParser(_testImOb[1]).Attributes().id is: "#20010101.154521928")
		Assert(XmlParser(_testImOb[2]).Attributes().id is: "#20010101.154521929")
		Assert(XmlParser(_testImOb[3]).Attributes().id is: "#20010101.154521940")
		Assert(XmlParser(_testImOb[4]).Attributes().id is: "#20010101.154521941")
		Assert(XmlParser(_testImOb[5]).Attributes().id is: "#20010101.154521942")
		mock.Verify.Times(2).scrollToMessage([anyArgs:])
		mock.Verify.scrollToMessage(#20010101.154521942)

		mock.On_Clear()
		Assert(mock.IM_MessengerConversationControl_oldestImNum is: false)
		Assert(mock.IM_MessengerConversationControl_newestImNum is: false)
		mock.On_History()
		Assert(mock.IM_MessengerConversationControl_oldestImNum is: #20010101.154521927)
		Assert(mock.IM_MessengerConversationControl_newestImNum is: #20010101.154521929)
		mock.Verify.Times(3).scrollToMessage([anyArgs:])
		mock.Verify.Times(2).scrollToMessage(#20010101.154521929)
		}

	setupMock()
		{
		mock = Mock(IM_MessengerConversationControl)
		mock.IM_MessengerConversationControl_tabName = 'mock tab'
		mock.When.On_Send().CallThrough()
		mock.When.On_History([anyArgs:]).CallThrough()
		mock.When.AppendMessages([anyArgs:]).CallThrough()
		mock.When.On_Clear([anyArgs:]).CallThrough()
		mock.When.IM_MessengerConversationControl_getHistory([anyArgs:]).Return(.history)
		mock.When.IM_MessengerConversationControl_getRecipientName([anyArgs:]).Return('')
		mock.When.IM_MessengerConversationControl_getFullName([anyArgs:]).Return('')
		mock.When.IM_MessengerConversationControl_getHtmlEmpty([anyArgs:]).Return('')
		mock.When.IM_MessengerConversationControl_outputMessage([anyArgs:]).Return(true)

		mock.IM_MessengerConversationControl_display = Mock()
		mock.IM_MessengerConversationControl_display.When.InsertAdjacentHTML([anyArgs:]).
			Do()
			{ |call|
			pos = call[2]
			if pos is 'beforeend'
				_testImOb.Add(call[3])
			else if pos is 'afterbegin'
				_testImOb.Add(call[3] at: 0)
			}
		mock.IM_MessengerConversationControl_editor = Mock()
		mock.IM_MessengerConversationControl_editor.When.Get().Return('A test message')

		Assert(mock.IM_MessengerConversationControl_oldestImNum is: false)
		Assert(mock.IM_MessengerConversationControl_newestImNum is: false)
		return mock
		}
	history: (
		(im_from: 'send', im_to: 'rec', im_message: 'msg3', im_num: #20010101.154521929)
		(im_from: 'send', im_to: 'rec', im_message: 'msg2', im_num: #20010101.154521928)
		(im_from: 'send', im_to: 'rec', im_message: 'msg1', im_num: #20010101.154521927)
		)
	append: (
		(im_to: 'send', im_from: 'rec', im_message: 'msg4', im_num: #20010101.154521940)
		(im_to: 'send', im_from: 'rec', im_message: 'msg5', im_num: #20010101.154521941)
		(im_to: 'send', im_from: 'rec', im_message: 'msg6', im_num: #20010101.154521942)
		)
	}