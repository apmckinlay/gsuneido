// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
/*
(      client      )         (          server                )
client <=> SvcClient <=net=> SvcServer <=> SvcCore <=> database

SvcServer recieves commands from the SvcClient SocketClient and interprets them, passing
them along to SvcCore. It returns the output from SvcCore back to the client.

This class cycles through a series of 'states', which define how the class interprets the
input.

'Verify' state is used to establish an authenticated connection with the server.
'Open' state is interprets the incoming commands as function calls. These calls are
	passed to SvcCore.
'Closed' state will close the SocketClient on the next cycle.
*/
SocketServer
	{
	Name: 'SVC Server'
	Port: 2222
	Run()
		{
		requestSize = args = false
		try
			{
			.state = 'Verify'
			.Writeline('Suneido Version Control Server')
			while .state isnt 'Closed' and false isnt request = .Readline()
				if 0 is requestSize = Number(request)
					.state = 'Closed'
				else
					.stateless((args = Unpack(.Read(requestSize))))
			}
		catch(e)
			.errorHandler(e, requestSize, args)
		}

	errorHandler(e, requestSize, args)
		{
		if e.Has?('lost connection') or e.Prefix?('socket') or .state is 'Closed'
			return
		remoteUser = 'REMOTE USER FAILED'
		try
			remoteUser = .RemoteUser()
		SuneidoLog('ERROR: Svc Server - ' $ remoteUser $ ': ' $ e,
			params: [:requestSize, :args, svcuser_id: .svcuser_id])
		.write('ERR ' $ e)
		}

	key: false
	NONCE(@unused)
		{ .write(.key = .buildKey()) }

	buildKey()
		{
		return Seq(12 /*= key length*/).
			Map({|unused| ('a'.Asc() + Random(26 /*= characters*/)).Chr() }).
			Join()
		}

	svcuser_id: false
	LOGIN(.svcuser_id, svcuser_passhash)
		{
		result = SvcUsers.LoginRequest(.key, .svcuser_id, svcuser_passhash)
		.write(result)
		.state = result is 'Invalid' ? 'Closed' : 'Open'
		}

	ADDUSER(serverPassword, svcuser_id, svcuser_passhash)
		{
		result = SvcUsers.
			AddUserRequest(.key, serverPassword, svcuser_id, svcuser_passhash)
		.write(result)
		}

	CHANGEPASSWORD(serverPassword, svcuser_id, oldPasshash, newPasshash)
		{
		result = SvcUsers.
			ChangePasswordRequest(.key, serverPassword, svcuser_id, oldPasshash,
				newPasshash)
		.write(result)
		}

	DELETEUSER(serverPassword, svcuser_id)
		{
		result = SvcUsers.DeleteUserRequest(.key, serverPassword, svcuser_id)
		.write(result)
		}

	public: #(NONCE, LOGIN, ADDUSER, CHANGEPASSWORD, DELETEUSER, QUIT)
	stateless(args)
		{
		if not Object?(args)
			throw 'invalid request, arguments are not in an Object'
		args = args.Copy()
		cmd = args.Extract(0)
		if .state is 'Verify' and not .public.Has?(cmd)
			throw 'connection must be verified prior to command: ' $ cmd
		if .Member?(cmd)
			try
				this[cmd](@args)
			catch (e)
				throw cmd $ ': ' $ e
		else
			throw 'invalid command: ' $ cmd
		}

	LISTPACK(table, since)
		{ .write(SvcCore.AllMasterChanges(table, since)) }

	LISTALL(since, dir)
		{
		sendObject = Object()
		list = SvcCore.ListAllChanges(since, dir)
		for x in list.Members()
			sendObject.Add(Object(list[x], table: x))
		.write(sendObject)
		}

	LISTRANGEPACK(table, since, to)
		{ .write(SvcCore.AllMasterChanges(table, since, to)) }

	GET(table, name)
		{ .sendRec(#GET, SvcCore.Get(table, name)) }

	GETOLD(table, name, committed)
		{ .sendRec(#GETOLD, SvcCore.GetOld(table, name, committed)) }

	GETDEL(table, name)
		{ .sendRec(#GETDEL, SvcCore.GetDel(table, name)) }

	GETDELBYDATE(table, name, committed)
		{ .sendRec(#GETDELBYDATE, SvcCore.GetDelByDate(table, name, committed)) }

	PUT(table, type, id, asof)
		{ .write(SvcCore.Put(table, type, id, asof, .readRec())) }

	REMOVE(table, name, id, asof, comment)
		{ .write(SvcCore.Remove(table, name, id, asof, comment)) }

	GETBEFORE(table, name, when)
		{ .write(SvcCore.GetBefore(table, name, when)) }

	GET10BEFORE(table, name, when)
		{ .write(SvcCore.Get10Before(table, name, when)) }

	EXISTS?(table)
		{ .write(SvcCore.Exists?(table)) }

	GETCHECKSUMS(table, from, to)
		{ .write(SvcCore.GetChecksums(table, from, to)) }

	ONLYDELETEDCHANGESBETWEEN(table, localMaxLibCommitted, savedMaxLibCommitted)
		{
		results = SvcCore.
			OnlyDeletedChangesBetween?(table, localMaxLibCommitted, savedMaxLibCommitted)
		.write(results)
		}

	CHECKSTATUS()
		{
		msg = ''
		for c in Contributions('Svc_AdditionalChecks')
			if '' isnt msg = c()
				break
		.write(msg.Blank?() ? 'OK' : 'ERR ' $ msg)
		}

	SEARCHFORRENAME(table, name)
		{ .write(SvcCore.SearchForRename(table, name)) }

	SVCTIME()
		{ .write(SvcCore.SvcTime()) }

	write(result)
		{
		pack = Pack(result)
		.Writeline(pack.Size())
		.Write(pack)
		}

	QUIT()
		{
		.write('OK bye')
		.state = 'Closed'
		}

	sendRec(call, rec)
		{
		if rec is false
			.write('INFO ' $ call $ ': failed to return a value other than false')
		else
			{
			.write([
				name: 			rec.name,
				path: 			rec.path,
				id: 			rec.id,
				lib_committed: 	rec.lib_committed,
				comment: 		.sanitize(rec.comment)
				])
			.Writeline(String(rec.text.Size()))
			.Write(rec.text)
			}
		}

	readRec()
		{
		rec = .read()
		rec.comment = .unsanitize(rec.comment)
		rec.text = .Read(Number(.Readline()))
		return rec
		}

	sanitize(s)
		{ return s.Tr('\r\n', '\x03\x04') }

	unsanitize(s)
		{ return s.Tr('\x03\x04', '\r\n') }

	read()
		{ return Unpack(.Read(Number(.Readline()))) }
	}
