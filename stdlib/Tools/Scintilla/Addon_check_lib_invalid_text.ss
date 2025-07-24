// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
// TEMPORARY: remove when 32381 is completed
ScintillaAddon
	{
	New(@args)
		{
		super(@args)
		.subs = [
			PubSub.Subscribe(#LibraryTreeChange, .Set)
			PubSub.Subscribe(#LibraryRecordChange, .Set)
			]
		}

	alert?: false
	Set()
		{
		if .alert?
			return
		.alert? = true
		table = .Parent.Controller.Table
		name =  .Parent.Controller.RecName
		if false is rec = Query1(table, :name, group: -1)
			return
		rec.text = rec.lib_current_text
		if rec.lib_invalid_text isnt '' and CodeState.Valid?(table, rec)
			.AlertWarn(#LibView, .prompt(table, name))
		}

	prompt(table, name)
		{
		return 'Record - ' $ table $ ':' $ name $ '\r\n\r\n' $
			'This record has lib_invalid_text, despite being valid. ' $
			'Modifying the record will correct the issue.'
		}

	Destroy()
		{
		.subs.Each(#Unsubscribe)
		}
	}