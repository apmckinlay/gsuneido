// Copyright (C) 2015 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_FormatMessage()
		{
		format = BookNotification.FormatMessage

		columns = #(name, num, phone)
		linkOb = #(name: #('Access', 'name'))
		rec = [num: #20000101, name: 'Test Guy']

		msg = format(columns, linkOb, rec)
		Assert(msg
			is: `Name: <a href="suneido:/eval?AccessGoTo(&quot;Access&quot;,` $
				`&quot;name&quot;,&quot;Test Guy&quot;,window:&quot;Dialog&quot;)">` $
				`Test Guy</a>&emsp;Date/Time Created: 2000-01-01`)

		linkOb = #(name: #("#(Biz_Employees, 'Trucking')", "name"))
		msg = format(columns, linkOb, rec)
		Assert(msg
			is: `Name: <a href="suneido:/eval?AccessGoTo(` $
				`&quot;%23(Biz_Employees, &apos;Trucking&apos;)&quot;,` $
				`&quot;name&quot;,&quot;Test Guy&quot;,window:&quot;Dialog&quot;)">` $
				`Test Guy</a>&emsp;Date/Time Created: 2000-01-01`)
		}

	Test_GetNewEntries()
		{
		cl = BookNotification
			{
			BookNotification_foreachNotification(block)
				{
				for i in ..100
					{
					try
						block([:i])
					catch (ex, "block:")
						if ex is "block:break"
							break
						// else block:continue ... so continue
					}
				}
			BookNotification_formatMessage(x)
				{
				return String(x.i)
				}
			}

		msgOb = cl.GetNewEntries(false)
		Assert(msgOb isSize: 50)
		Assert(msgOb[0] is: "<li>0</li>")
		Assert(msgOb.Last() is: '<li>49</li>')

		msgOb = cl.GetNewEntries(false, 4, 12)
		Assert(msgOb isSize: 8)
		Assert(msgOb[0] is: "<li>4</li>")
		Assert(msgOb.Last() is: '<li>11</li>')
		}
	}
