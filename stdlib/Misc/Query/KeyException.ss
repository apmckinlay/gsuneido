// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
// TODO: rename since it's no longer just "key" exceptions
class
	{
	Exceptions: (
		("duplicate key", duplicate)
		("key too large", toolarge)
		("blocked by foreign key", blocked)
		("update transaction longer than", timeout)
		("can't query ended transaction", timeout)
		("transaction exceeded max age", timeout) // gSuneido
		("too many reads", toomany)
		("too many writes", toomany)
		("too many overlapping update transactions", serverslow)
		// these need to be last
		("transaction conflict", conflict)
		("record had been modified", conflict)
		("block commit failed", conflict)
		("transaction.Complete failed", conflict)
		)
	TryCatch(block,
		catch_block = function (e)
			{
			KeyException.Alert(e, 'perform action')
			throw "interrupt: KeyException (CAUGHT) " $ e // will log but user won't see
			})
		{
		try
			{
			block()
			return true
			}
		catch (e)
			{
			if .Exceptions.Any?({ e.Has?(it[0]) })
				return catch_block(e)
			else
				throw e
			}
		}
	Alert(e, action)
		{
		if e is msg = .Translate(e, action)
			throw "KeyException unhandled: " $ e // should never happen
		.alert(msg)
		}
	CallClass(e, action = 'perform action')
		{
		if Suneido.User isnt 'default'
			SuneidoLog("warning: (CAUGHT) " $ e, calls:, caughtMsg: 'user alerted: ' $ e)
		.Alert(e, action)
		}

	Translate(e, action) // returns original e if not translated
		{
		// suppress unused errors
		.duplicate
		.toolarge
		.blocked
		.timeout
		.toomany
		.conflict
		.serverslow

		for x in .Exceptions
			if e.Has?(x[0])
				return this['KeyException_' $ x[1]](:e, :action)
		return e
		}
	duplicate(e)
		{
		keys = e.AfterFirst('duplicate key:').Trim().Extract('([^ ]+?)( |$)')
		keys = keys.Split(',').Map(SelectPrompt).Map(TranslateLanguage).Join('+')
		return not keys.Tr('^a-zA-Z').Blank?() // not an empty key
			? 'Duplicate value in field ' $ keys $ .extractKeyValue(e, keys)
			: 'Duplicate Entry'
		}
	toolarge(e /*unused*/)
		{
		// currently the error contains no specific information about which field it is
		return 'Indexed field contains too much data'
		}
	blocked(e)
		{
		return 'Record cannot be updated or deleted because it is used' $
			.foreignTableInfo(e)
		}
	timeout(action)
		{
		return 'Unable to ' $ action $ '\n\nAction timed out'
		}
	toomany(action)
		{
		return 'Unable to ' $ action $ '\n\nToo many lines to process'
		}
	conflict(e, action)
		{
		user = .build_user_msg(e)
		SuneidoLog("INFO: transaction conflict caught by KeyException", calls:,
			params: Object(:action, :e, :user))
		return 'Unable to ' $ action $ '\n\nAnother user' $ user $ ' has made changes'
		}
	serverslow(action)
		{
		return 'Unable to ' $ action $ '\n\nServer too slow responding'
		}

	extractKeyValue(x, keys)
		{
		// composite keys not handled yet, assumes key values enclosed in []
		// timestamps not returned as these key values are usually system
		// generated and will not likely mean anything to the user
		if keys.Has?('+') or not x.Has?('[') or not x.Has?(']')
			return ''

		val = x.AfterFirst('[').BeforeFirst(']')
		val = val.Has?('"')
			? val.BeforeLast('"') $ '"'
			: val.BeforeFirst(',')
		val = val.Trim()
		return not val.Prefix?('#') ? ': ' $ val : ''
		}

	foreignTableInfo(x)
		{
		// Replaces are used to remove everything within () and [] which should
		// leave us with the basic foriegn key message with the table at the end.
		// See the associated test class for examples of exception strings.
		msgSplit = x.Replace('\(.*?\)', '').Replace('\[.*?\]', '').
			Trim().AfterFirst('foreign key').Split(' ').Map!(#Trim).Remove('cascade')
		if msgSplit.Empty?() or msgSplit.Last() is ''
			return ''
		table = msgSplit.Last()
		if table isnt tableDesc = GetTableName(table)
			return ' (' $ tableDesc $ ')'
		return ''
		}

	build_user_msg(x)
		{
		conflict_user = x.Extract('conflict with (.*?)@')
		return conflict_user is false or Suneido.User is conflict_user
			? "" : ' (' $ conflict_user $ ')'
		}

	alert(msg)
		{
		AlertDelayed(msg)
		}
	}
