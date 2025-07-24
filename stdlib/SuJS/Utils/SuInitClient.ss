// Copyright (C) 2018 Axon Development Corporation All rights reserved worldwide.
class
	{
	CallClass(@args)
		{
		try
			{
			.localeSettings()
			args.isHttps? = IsHttps?()
			args.userAgent = SuUI.GetCurrentWindow().navigator.userAgent
			SuRender.Init(args.token)
			SuRender().Event(false, 'SuInit', args)
			}
		catch (e)
			{
			SuUI.GetCurrentWindow().Alert('Error occurred when initializing: ' $ e)
			SuRender.Reload()
			}
		}
	localeSettings()
		{
		Settings.Set('ShortDateFormat', "yyyy-MM-dd")
		Settings.Set('SystemShortDateFormat', "yyyy-MM-dd")
		Settings.Set('LongDateFormat', "dddd, MMMM dd, yyyy")
		Settings.Set('TimeFormat', "h:mm tt")
		Settings.Set('ThousandsSeparator', ",")
		}
	}
