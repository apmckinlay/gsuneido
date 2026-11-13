// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
// BuiltDate > 30000101
class
	{
	CallClass(@args)
		{
		table = args.GetDefault(#table, #users)
		switch args.Size(list:)
			{
		case 1: // token
			if false is ob = .UseToken(args[0], includePersistent?:)
				return false
			if String?(ob.permit) and ob.permit is #all
				.Permit.All()
			else
				.Permit.Set(ob.permit)
		case 2: // user name and password
			if not .checkUserPass(table, args[0], args[1])
				return false
			.setPermit(args[0])
			}
		return true
		}

	checkUserPass(table, user, password)
		{
		// backwards compatible, consider alternatives
		.Permit.Table(#tables, #read)
		.Permit.Table(table, #read)
		if not TableExists?(table) or QueryEmpty?(table)
			return true

		if false is x = Query1(table, :user)
			return false

//		if x.salt isnt ""
//			return Argon2id(password, x.salt) is x.passhash

		if PassHash(user, password) isnt x.passhash
			return false

		// upgrade hash
//		.Permit.Table(table, #write)
//		.Permit.Schema(#create)
//		Database("ensure " $ table $ " (salt)")
//		QueryApply1(table, :user)
//			{|x|
//			if x.salt is "" // recheck in case it's changed
//				{
//				x.salt = Random.Bytes(16) /*= salt size */
//				x.passhash = Argon2id(password, x.salt)
//				x.Update()
//				}
//			}
		return true
		}

	// verifies a token and removes it from the list
	UseToken(token, includePersistent? = false)
		{
		if false isnt ob = .authTokens().Extract(token, false)
			return ob

		if includePersistent? is true
			{
			.Permit.Table(#tables, #read)
			if TableExists?(#a25vef90b2tlbnm)
				{
				.Permit.Table(#a25vef90b2tlbnm, #write)
				QueryApply1(#a25vef90b2tlbnm, a25vef90b2tlbg: token)
					{
					if it.a25vef9lehbpcnk isnt Date.End()
						it.Delete()
					return it.a25vef9lehbpcnk >= Date()
						? Object(permit: #all)
						: false
					}
				}
			}

		return false
		}

	setPermit(user)
		{
		// Reset temp grants for checking
		.Permit.None()

		type = user is #default ? #ide : #customer

		serverEvals = type is #ide
			? ServerEvalWhitelist()
			: ServerEvalWhitelist.GetCustomerWhitelist()
		for permission in serverEvals
			.Permit.ServerEval(permission)

		.Permit.Table(#tables, #read)
		QueryApply("tables extend columns = TableColumns(table)")
			{
			action = HasLibraryColumns?(it.columns) and type is #customer
				? #read // still need read permission for Plugins
				: #write
			.Permit.Table(it.table, action)
			}

		// FIX ME !!!
		.Permit.Schema(#create)
		.Permit.Schema(#update)
		}

	// Persistent Token exists after program is closed
	PersistentToken(expiry, scriptRegen? = false)
		{
		// This is to ensure the session has the required permissions
		Assert(not Sys.Client?(),
			msg: "Safegarurd: tokens can only be output by the server")
		// This is to ensure programmers consider the state of customer launch scripts
		if expiry is Date.End()
			Assert(scriptRegen?,
				msg: "Safegarurd: scriptRegen must be true to output a non-expiring key")

		.cleanupPersistentToken(expiry)
		Assert(Date?(expiry) and expiry > Date(),
			msg: "Safegaurd: expiry must be a valid future date")

		token = .GenerateToken()
		while false isnt Query1(#a25vef90b2tlbnm, a25vef90b2tlbg: token)
			token = .GenerateToken()
		QueryOutput(#a25vef90b2tlbnm, [a25vef90b2tlbg: token, a25vef9lehbpcnk: expiry])

		return token
		}

	cleanupPersistentToken(expiry)
		{
		Table_a25vef90b2tlbnm().Ensure()
		where = expiry isnt Date.End()
			? "a25vef9lehbpcnk < " $ Display(Date())
			: "a25vef9lehbpcnk is " $ Display(Date.End())
		QueryDo("delete a25vef90b2tlbnm where " $ where)
		}

	// creates a new token and saves it
	Token()
		{
		if Sys.Client?()
			return ServerEval("Auth.Token")

		return .token()
		}

	token()
		{
		tok = .GenerateToken()
		.authTokens()[tok] = Object(date: Date(), permit: .Permit.Get())
		return tok
		}

	tokenSize: 32
	GenerateToken()
		{
		return Base64.EncodeUrl(Random.Bytes(.tokenSize))
		}

	// returns the token list after removing expired
	authTokens()
		{
		tokens = Suneido.GetInit(#AuthTokens, Object)
		expire = Date().Minus(seconds: 120) /*= 2 minutes */
		// RemoveIf is thread safe because it's named members and only add/delete
		tokens.RemoveIf({ it.date < expire })
		return tokens
		}

	Getter_Permit() // overridden by test
		{
		return NewAuth?()
			? Permit
			: class
				{
				Default(@unused) { return false }
				}
		}
	}
