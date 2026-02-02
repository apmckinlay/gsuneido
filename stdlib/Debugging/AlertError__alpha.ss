// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Title: 'Warning'

	CallClass(msg = false, hwnd = 0, log_rec = #(), logOnly? = false)
		{
		if logOnly?
			{
			.output_log(:log_rec, setOutput?: false)
			return
			}
		if msg is false
			{
			msg = "An unexpected problem has occurred."
			if not log_rec.Empty?()
				{
				ToolDialog(hwnd, Object(this, msg, log_rec), closeButton?: false,
					keep_size: false)
				return
				}
			}
		.output_log(:log_rec, setOutput?: false)
		Alert(msg, title: "Warning", :hwnd, flags: MB.ICONWARNING)
		}
	New(msg = '', .log_rec = #())
		{
		super(.layout(msg))
		}
	layout(msg)
		{
		return Object('Vert',
			Object('Static', msg, textStyle: 'main'),
			'Skip',
			Object('StaticWrap',
				'If you have time to enter a short description of what you were ' $
				'doing when this problem occurred, it will help us improve the ' $
				'software. Thank you.',
				xstretch: 1)
			'Skip'
			#(Editor width: 40, height: 5, name: 'message'),
			'Skip',
			#(Horz Fill OkButton)
			)
		}
	On_OK()
		{
		message = .Vert.message.Get().Trim()
		if message is ''
			{
			.Window.Result(false)
			return
			}

		// append user's message to the error log
		.output_log(message: .descPrefix $ message, log_rec: .log_rec)
		.Window.Result(true)
		}

	outputlog?: true
	descPrefix: 'DESC: '
	output_log(message = '', log_rec = #(), setOutput? = true)
		{
		message = message is .descPrefix ? '' : message
		SuneidoLog(log_rec.sulog_message[0..] $ Opt(', ', message),
			calls: log_rec.calls)
		if setOutput?
			.outputlog? = false
		}

	Destroy()
		{
		if .outputlog?
			.output_log(message: Opt(.descPrefix, .Vert.message.Get().Trim()),
				log_rec: .log_rec)
		super.Destroy()
		}
	}