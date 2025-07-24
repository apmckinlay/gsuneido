// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Title: 'Disconnect Users'
	New(conns)
		{
		super(.layout(conns))
		.list = .Vert.List
		}
	layout(conns)
		{
		data = Object()
		for conn in conns.Map(.SessionIdWithoutToken).UniqueValues()
			data.Add(Object(loggedinusers_col: conn))
		return Object('Vert'
			Object('List', #(loggedinusers_col), data, defWidth: 200)
			'Skip'
			#('Horz', 'Fill', #(Button 'Disconnect')))
		}
	SessionIdWithoutToken(session)
		{
		return session.Suffix?('(jsS)') ? session.BeforeLast('<') $ '(jsS)' : session
		}
	On_Disconnect()
		{
		rows = .list.GetSelection()
		if rows.Empty?()
			{
			.AlertInfo('Disconnect Users', 'You must select a user to disconnect')
			return
			}
		data = .list.Get()
		copy = data.Copy()
		for row in rows
			{
			if not data.Member?(row)
				return
			.disconnect(data[row].loggedinusers_col)
			copy.Remove(data[row].Delete('listrow_flags'))
			}
		.list.Clear()
		.list.AddRows(copy)
		}
	disconnect(user)
		{
		BookLog("Kill " $ user)
		Sys.Kill(user)
		}
	List_DoubleClick(@unused)
		{
		return false
		}
	}