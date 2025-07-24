// Copyright (C) 2015 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_ParseQuery()
		{
		parseQuery = IM_MessengerManager.ParseQuery
		content = ""
		query = parseQuery(content)
		Assert(query is: Object())

		content = "user=p1"
		query = parseQuery(content)
		Assert(query is: #(user: "p1"))

		content = "to=p1
from=default
msg=this+is+a+simple+message
date=#20150608.104034380"
		query = parseQuery(content)
		Assert(query is: #(msg: "this is a simple message",
			date: "#20150608.104034380", from: "default", to: "p1"))

		content = "to=p1
from=default
msg=this+is%0D%0Aa+new%0D%0Aline
date=#20150608.104034380"
		query = parseQuery(content)
		Assert(query is: #(msg: "this is
a new
line", date: "#20150608.104034380", from: "default", to: "p1"))

		content = "to=p1
from=default
msg=This+is+a+line+with+&+and+=
date=#20150608.104034380"
		query = parseQuery(content)
		Assert(query is: #(msg: "This is a line with & and =",
			date: "#20150608.104034380", from: "default", to: "p1"))
		}

	Test_logErr()
		{
		cl = IM_MessengerManager
			{
			IM_MessengerManager_postResult(request, content /*unused*/, ip /*unused*/)
				{
				return request // Pass request through to Json.Decode
				}
			}
		mock = Mock(cl)
		mock.When.Request([anyArgs:]).CallThrough()
		mock.When.IM_MessengerManager_publishMessengerOnline().Return(false)
		mock.When.IM_MessengerManager_logErr([anyArgs:]).Return(false)
		logErr = 'IM_MessengerManager_logErr'

		mock.Request("'ValidString'")
		mock.Verify.Never()[logErr]([anyArgs:])
		mock.Verify.IM_MessengerManager_publishMessengerOnline()

		mock.Request("{}", skipSubscribe?:)
		mock.Verify.Never()[logErr]([anyArgs:])
		// skipSubscribe? should prevent this from being called again
		mock.Verify.IM_MessengerManager_publishMessengerOnline()

		mock.Request("")
		mock.Verify[logErr](err = "Invalid Json format: unexpected end of string", false)

		mock.Request("{")
		mock.Verify.Times(2)[logErr](err, false)

		mock.Request("{}}")
		mock.Verify[logErr]("Invalid Json format: extra text at end", false)

		mock.Request("InvalidString")
		mock.Verify[logErr]("Invalid Json format: unexpected: InvalidString", false)
		}
	}