// Copyright (C) 2015 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_FormatMessage()
		{
		format = BookNotification.FormatMessage

		columns = #(name, num, phone)
		linkOb = #(name: (Access, name))
		rec = [num: #20000101, name: "Test Guy"]

		msg = format(columns, linkOb, rec)
		Assert(msg
			is: `Name: <a href="suneido:/eval?AccessGoTo%28%22Access%22%2C` $
				`%22name%22%2C%22Test%20Guy%22%2Cwindow%3A%22Dialog%22%29">` $
				`Test Guy</a>&emsp;Date/Time Created: 2000-01-01`)

		rec.name = "#Test Guy"
		msg = format(columns, linkOb, rec)
		Assert(msg
			is: `Name: <a href="suneido:/eval?AccessGoTo%28%22Access%22%2C` $
				`%22name%22%2C%22%23Test%20Guy%22%2Cwindow%3A%22Dialog%22%29">` $
				`#Test Guy</a>&emsp;Date/Time Created: 2000-01-01`)

		linkOb = #(name: ("#(Biz_Employees, 'Trucking')", name))
		msg = format(columns, linkOb, rec)
		Assert(msg.Has?("%23%28Biz_Employees%2C%20%27Trucking%27%29"))
		Assert(.linkArgs(msg)
			is: ["#(Biz_Employees, 'Trucking')", #name,
				"#Test Guy", window: #Dialog])
		}

	Test_FormatMessage_linkValues()
		{
		linkOb = #(name: (Test_Control, test_value))
		for value in #("#1982", "1982", "A+B %23 ? # & 'quoted'", 1982)
			{
			msg = BookNotification.FormatMessage(#(name), linkOb, [name: value])
			url = msg.AfterFirst('href="').BeforeFirst('"')
			Assert(not url.Has?('#'))
			Assert(.linkArgs(msg)
				is: [#Test_Control, #test_value, value,
					window: #Dialog])
			}
		}

	linkArgs(msg)
		{
		url = msg.AfterFirst('href="').BeforeFirst('"')
		url = SuneidoAPP.SuneidoAPP_decode(url)
		code = Url.Decode(url.RemovePrefix("suneido:/eval?"))
		return code.Replace("^AccessGoTo", #Object).Eval()
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

			BookNotification_formatMessage(x) { return String(x.i) }
			}

		msgOb = cl.GetNewEntries(false)
		Assert(msgOb isSize: 50)
		Assert(msgOb[0] is: "<li>0</li>")
		Assert(msgOb.Last() is: "<li>49</li>")

		msgOb = cl.GetNewEntries(false, 4, 12)
		Assert(msgOb isSize: 8)
		Assert(msgOb[0] is: "<li>4</li>")
		Assert(msgOb.Last() is: "<li>11</li>")
		}
	}
