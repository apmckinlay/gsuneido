// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Title: 'Warning'

	CallClass(msg = false, hwnd = 0, log_rec = #())
		{
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
		Alert(msg, title: "Warning", :hwnd, flags: MB.ICONWARNING)
		}
	New(msg = '', log_rec = #())
		{
		super(.layout(msg))
		.sulog_timestamp = log_rec.GetDefault('sulog_timestamp', false)
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
		.update_log('DESC: ' $ message)
		.Window.Result(true)
		}

	update_log(message = '')
		{
		if .sulog_timestamp is false or message is ""
			return
		QueryApply1('suneidolog', sulog_timestamp: .sulog_timestamp)
			{|x|
			x.sulog_message = x.sulog_message $ ', ' $ message
			x.Update()
			}
		}
	}
