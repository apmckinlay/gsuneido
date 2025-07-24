// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
class
	{
	Create(currency, largedata? = false, config = #(), isDev = false, dept? = true)
		{
		DoWithAlertToSuneidoLog()
			{
			fn = 2
			Plugins().ForeachContribution('DemoData', 'add_data', showErrors:)
				{ |x| (x[fn]) (currency, largedata?, :config, :isDev, :dept?) }
			}
		}
	}