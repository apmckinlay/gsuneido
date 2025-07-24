// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
#(
	startupLogs: // Displays suneidologs that occur durring startup
		function ()
			{
			if Date?(Suneido.GetDefault('init_end', false))
				return // Contribution has already run by this point (standalone)

			start = Suneido.GetDefault('start_time', false)
			if TableExists?('suneidolog') and start isnt false
				QueryApply('suneidolog
					where sulog_timestamp >= ' $ Display(start.Minus(seconds: 1)) $
						' and sulog_timestamp < ' $ Display(Timestamp()))
					{
					Print('SuneidoLog [' $ Display(it.sulog_timestamp) $ ']: ' $
						it.sulog_message)
					}
			},
	cleanupWebView2:
		function ()
			{
			if Sys.Win32?()
				Thread(WebView2.CleanUp)
			}
)