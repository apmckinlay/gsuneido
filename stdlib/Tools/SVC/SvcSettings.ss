// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Table: svc_settings
	CallClass(openDialog = false)
		{
		// Prevents stackoverflow issue when multiple places try to access Svc
		// at the same time. ie. addons
		return Suneido.Member?('SvcDialogOpen') ? false : .getSettings(:openDialog)
		}

	getSettings(openDialog = false)
		{
		Suneido.SvcDialogOpen = true
		.Ensure()
		data = .Get()
		if openDialog
			data = OkCancel(Object(this, data), .Title)
		Suneido.Delete('SvcDialogOpen')
		return data
		}

	Ensure()
		{
		Database("ensure " $ .Table $ "
			(svc_server, svc_user, svc_local?)
			key()")
		}

	Title: "Version Control Settings",
	New(data)
		{
		.Data.Set(data)
		.server = .FindControl('svc_server')
		.passhash = .FindControl('svc_passhash')
		.user = .FindControl('svc_user')
		.Data.AddObserver(.change)
		.serverChanged(data)
		}
	Controls: #(Record (Vert
		(Pair (Static '') (CheckBox, "Standalone" name: "svc_local?"))
		(Pair (Static Server) (Field name: svc_server))
		(Pair, (Static Login), (Field name: svc_userId)),
		(Pair, (Static Password), (Field password:, name: svc_passhash))
		Skip
		(Pair, (Static 'Default User'), (Field name: svc_user)),
		))

	change(member)
		{
		data = .Data.Get()
		if member is 'svc_local?' or member is 'svc_server'
			.serverChanged(data)
		if member is 'svc_userId'
			.passhash.Set('')
		}

	serverChanged(data)
		{
		noServer? = data.svc_server is ''
		local? = data.svc_local? is true

		.server.SetReadOnly(local?)
		.user.SetReadOnly(clearUser? = not local? and noServer?)

		if clearUser?
			.user.Set('')
		}

	OK()
		{
		data = .Data.Get().Copy() // copy to remove observer
		if data.svc_local? isnt true and data.svc_server is ''
			{
			.AlertError(.Title, 'Server or Standalone required')
			return false
			}
		if TableExists?(.Table)
			.updateSettings(.Get(), data)
		return data
		}

	updateSettings(oldData, newData)
		{
		if oldData is false
			oldData = []
		dataChanged? = false
		for col in .columns()
			if oldData[col] isnt newData.GetDefault(col, '')
				dataChanged? = true
		if dataChanged?
			.update(newData)
		}

	columns()
		{
		return QueryColumns(.Table).Add('svc_userId', 'svc_passhash')
		}

	update(newData)
		{
		QueryDo('delete ' $ .Table)
		QueryOutput(.Table, newData)
		.UpdateCredentials(newData.svc_userId, newData.svc_passhash)
		PubSub.Publish('SvcSettings_ConnectionModified')
		}

	UpdateCredentials(userId, passhash)
		{
		file = .passhashFile()
		if userId is '' or passhash is ''
			{
			DeleteFile(file)
			return []
			}
		credOb = .decrypt(file)
		if userId is credOb.userId and passhash is credOb.passhash
			return credOb
		credOb = [:userId, passhash: PassHash(userId, passhash)]
		PutFile(file, .encodeDecode(credOb))
		return credOb
		}

	decrypt(file)
		{
		credentials = []
		if false is encrypted = GetFile(file)
			return credentials
		try
			credentials = .encodeDecode(encrypted).SafeEval()
		catch
			SuneidoLog('INFO: Failed to decrypt credentials. Please verify Svc Settings')
		return credentials
		}

	passhashFile()
		{
		dir = Sys.Linux?() ? 'HOME' : 'APPDATA'
		return Paths.ToLocal(Getenv(dir) $ '/svc.key')
		}

	encodeDecode(credentials)
		{
		if Object?(credentials)
			credentials = Display(credentials)
		return credentials.Xor(GetMacAddressHex())
		}

	Cancel()
		{ return .Get() }

	Get()
		{
		if not TableExists?(.Table) or false is rec = Query1(.Table) // Empty table
			return .Credentials()
		return rec.MergeNew(.Credentials())
		}

	Credentials()
		{
		rec = []
		credOb = .decrypt(.passhashFile())
		rec.svc_userId = credOb.userId
		rec.svc_passhash = credOb.passhash
		return rec
		}

	Set?()
		{
		rec = TableExists?(.Table)
			? Query1(.Table)
			: false
		return rec isnt false
			? rec.svc_server isnt '' or rec.svc_local? is true
			: false
		}
	}
