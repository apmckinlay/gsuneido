// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		date = Date().Plus(days: -2)
		log = Object(sulog_timestamp: date,
			sulog_message: '',
			sulog_locals: #(),
			sulog_params: '',
			sulog_calls: '')
		Assert(RetrieveEmailError(log) is: Object(errMsg: "", date: date.ShortDateTime()))

		.forwardSendFailed()
		.forwardSendError()
		.cantConnectForwarders()

		.dontRelyOnHyphenSeparator()
		}

	sulog_calls: "ForwardSendMsg /* stdlib function */
BookSendEmail.BookSendEmail_emailSES /* stdlib method */
BookSendEmail.BookSendEmail_sendEmailSES /* stdlib method */
BookSendEmail.CallClass /* stdlib method */
Working /* stdlib function */
BookSendEmail.CallClass /* stdlib method */
Action_SendEmailTo /* stdlib function */
EventConditionActions.EventConditionActions_excute /* stdlib method */
EventConditionActions.PerformActions /* stdlib method */
 /* function */"

	forwardSendFailed()
		{
		log = .createLog(
			sulog_message: "ERROR: ForwardSendMessage - Non-Successful HTTP response",
			sulog_params: Object(response: #(header: "HTTP/1.0 400 Bad Request
Date: Fri, 12 Feb 2016 2...ed: Fri, 12 Feb 2016 22:25:40 GMT
Server: Suneido",
				content: "invalid email address format: "),
				header: #(X_Suneido_From: "",
					Authorization: "Basic <Authorization String>",
					X_Suneido_To: "a@b.c,hi.there@dude.man"),
				),
			sulog_calls: .sulog_calls)

		.validateContext(log,
			'Non-Successful HTTP response - invalid email address format: ',
			'MIME-Version: 1.0
From: "" <email@axoneta.com>
T...nt-Transfer-Encoding: 7bit

Order 1234 Updated
'
			)
		}

	createLog(sulog_message, sulog_params, sulog_calls = "", sulog_locals = #())
		{
		sulog_timestamp = Date().Plus(days: -2)
		sulog_params.to = "a@b.c,hi.there@dude.man"
		sulog_params.from = ""
		sulog_params.message = 'MIME-Version: 1.0
From: "" <email@axoneta.com>
T...nt-Transfer-Encoding: 7bit

Order 1234 Updated
'

		return Object(
			:sulog_timestamp,
			:sulog_message,
			:sulog_locals,
			:sulog_params,
			:sulog_calls)

		}

	validateContext(log, errMsg, message)
		{
		context = RetrieveEmailError(log)
		Assert(context.from is: '')
		Assert(context.to is: "a@b.c,hi.there@dude.man")
		Assert(context.errMsg is: errMsg)
		Assert(context.message is: message)
		}

	forwardSendError()
		{
		response = Object(
			header: "HTTP/1.0 400 Bad Request
Date: Tue, 16 Feb 2016 1...ed: Tue, 16 Feb 2016 15:14:59 GMT
Server: Suneido"
			content: "invalid email address format: "
			)
		header = Object(
			X_Suneido_From: '',
			X_Suneido_To: 'a@b.c,hi.there@dude.man',
			Authorization: "Basic <Authorization String>"
			)

		localsOb = Object(
			response_code: "400",
			:response,
			:header,
			authenticationContribs: #("<object>"),
			to: "a@b.c,hi.there@dude.man",
			from: "",
			message: 'MIME-Version: 1.0
From: "" <email@axoneta.com>
T...nt-Transfer-Encoding: 7bit

Order 1235 Updated
',
			relay: "service1.axoneta.net",
			this: "ForwardSendMsg /* stdlib function */",
			authentication: #(Authorization: "Basic <Authorization String>"))

		log = .createLog(
			sulog_message: "ERRATIC: BookSendEmail - problem sending e-mail, " $
				"switch to service2.axoneta.net",
			sulog_locals: localsOb,
			sulog_params: Object(),
			sulog_calls: .sulog_calls)

		.validateContext(log,
			'problem sending e-mail, switch to service2.axoneta.net - ' $
				'invalid email address format: ',
			'MIME-Version: 1.0
From: "" <email@axoneta.com>
T...nt-Transfer-Encoding: 7bit

Order 1235 Updated
'
			)
		}

	cantConnectForwarders()
		{
		log = .createLog(
			sulog_message: "ERRATIC: BookSendEmail - unable to connect to " $
				"service1.axoneta.net, switch to service2.axoneta.net",
			sulog_params: Object(ok: "400 invalid email address format: "))
		.validateContext(log,
			'unable to connect to service1.axoneta.net, ' $
				'switch to service2.axoneta.net - ' $
				'400 invalid email address format: ',
			'MIME-Version: 1.0
From: "" <email@axoneta.com>
T...nt-Transfer-Encoding: 7bit

Order 1234 Updated
'
			)
		}

	dontRelyOnHyphenSeparator()
		{
		log = .createLog(
			sulog_message: "ERROR: ForwardSendMessage: Non-Successful HTTP response",
			sulog_params: Object(response: #(header: "HTTP/1.0 400 Bad Request
Date: Fri, 12 Feb 2016 2...ed: Fri, 12 Feb 2016 22:25:40 GMT
Server: Suneido",
				content: "invalid email address format: "),
				header: #(X_Suneido_From: "",
					Authorization: "Basic <Authorization String>",
					X_Suneido_To: "a@b.c,hi.there@dude.man"),
				),
			sulog_calls: .sulog_calls)

		.validateContext(log,
			'Non-Successful HTTP response - invalid email address format: ',
			'MIME-Version: 1.0
From: "" <email@axoneta.com>
T...nt-Transfer-Encoding: 7bit

Order 1234 Updated
'
			)
		}
	}
