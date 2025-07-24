// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
/*
(      client      )         (          server                )
client <=> SvcClient <=net=> SvcServer <=> SvcCore <=> database

If SvcClient is being used, then standalone mode is not enabled and a server exists.-

Each function writes the neccessary requests to the SocketClient, which is passed to
SvcCore through the SvcSocketClient to SvcServer. It then reads the output it gets passed
back through SvcSocketClient and returns it to Svc.
*/

class
	{
	AllMasterChanges(table, since, to = "")
		{
		args = to is ''
			? [#LISTPACK, :table, :since]
			: [#LISTRANGEPACK, :table, :since, :to]
		return .svccl.Send(args, result: Object())
		}

	getter_svccl()
		{
		return .svccl = SvcSocketClient()
		}

	ListAllChanges(since, dir = false)
		{ return .svccl.Send([#LISTALL, :since, :dir], result: Object()) }

	Get(table, name)
		{
		result = false
		.svccl.Run()
			{ |sc|
			.svccl.Write(sc, [#GET, :table, :name])
			result = .readRec(sc, name)
			}
		return result
		}

	readRec(sc, name)
		{
		// SvcSocketClient handles the initial Read to handle potential errors
		x = .svccl.Read(sc, name)
		if not Object?(x)
			return false
		x.comment = .unsanitize(x.comment)
		x.text = sc.Read(Number(sc.Readline()))
		return x.name is name ? x : false
		}

	GetDel(table, name)
		{
		result = false
		.svccl.Run()
			{ |sc|
			.svccl.Write(sc, [#GETDEL, :table, :name])
			result = .readRec(sc, name)
			}
		return result
		}

	GetOld(table, name, committed)
		{
		result = false
		.svccl.Run()
			{ |sc|
			.svccl.Write(sc, [#GETOLD, :table, :name, :committed])
			result = .readRec(sc, name)
			}
		if result isnt false and committed isnt result.lib_committed and
			committed isnt Date.End()
			Alert('NAME:   ' $ name $ '\r\n' $
				'    LOCAL:  ' $ Display(committed) $ '\r\n' $
				'    MASTER: ' $ Display(result.lib_committed) $
				'\r\n\r\nPlease verify record state via the Compare button',
				'SvcClient - GetOld', flags: MB.ICONERROR)
		return result
		}

	GetDelByDate(table, name, committed)
		{
		result = false
		.svccl.Run()
			{ |sc|
			.svccl.Write(sc, [#GETDELBYDATE, :table, :name, :committed])
			result = .readRec(sc, name)
			}
		return result
		}

	Put(table, type, id, asof, rec)
		{
		result = false
		.svccl.Run()
			{ |sc|
			.svccl.Write(sc, [#PUT, :table, :type, :id, :asof])
			result = .sendRec(sc, rec)
			}
		return result
		}

	sendRec(sc, rec)
		{
		// Pack / Send bulk of record
		.svccl.Write(sc, [
			name: 				rec.name,
			path: 				rec.path,
			comment: 			.sanitize(rec.comment),
			lib_before_hash: 	rec.lib_before_hash
			])
		// Send text seperate (not packed)
		sc.Writeline(String(rec.text.Size()))
		sc.Write(rec.text)
		return .svccl.Read(sc)
		}

	Remove(table, name, id, asof, comment)
		{
		return .svccl.Send([
			#REMOVE,
			:table,
			:name,
			:id,
			:asof,
			comment: .sanitize(comment)])
		}

	sanitize(s)
		{
		return s.Tr("\r\n", "\x03\x04")
		}

	unsanitize(s)
		{
		return s.Tr("\x03\x04", "\r\n")
		}

	GetBefore(table, name, when)
		{ return .svccl.Send([#GETBEFORE, :table, :name, :when]) }

	Get10Before(table, name, when)
		{ return .svccl.Send([#GET10BEFORE, :table, :name, :when], result: Object()) }

	Exists?(table)
		{ return .svccl.Send([#EXISTS?, :table]) }

	GetChecksums(table, from, to)
		{ return .svccl.Send([#GETCHECKSUMS, :table, :from, :to], result: Object()) }

	OnlyDeleteChangesBetween(table, localMaxLibCommitted, savedMaxLibCommitted)
		{
		return .svccl.Send([
			#ONLYDELETEDCHANGESBETWEEN,
			:table,
			:localMaxLibCommitted,
			:savedMaxLibCommitted])
		}

	CheckSvcStatus()
		{
		result = .svccl.Send([#CHECKSTATUS],
			result: .svccl.ErrorPrefix $ 'SocketClient is not connected')
		return result is #OK ? '' : result.AfterFirst(.svccl.ErrorPrefix)
		}

	SearchForRename(table, name)
		{ return .svccl.Send([#SEARCHFORRENAME, :table, :name], result: []) }

	SvcTime()
		{ return .svccl.Send([#SVCTIME]) }
	}
