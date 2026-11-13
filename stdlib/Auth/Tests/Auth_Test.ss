// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
Test
	{
//	Test_user_password()
//		{
//		auth = .auth
//		table = .MakeTable("(user, passhash) key(user)")
//		user = .TempName()
//		password = .TempName()
//		Assert(auth(user, password, :table)) // empty table
//		QueryOutput(table,
//			[:user, passhash: Md5(user $ password)])
//		Assert(auth(user, .TempName(), :table) is: false)
//		Assert(auth(user, password, :table))
//		x = MustQuery1(table, :user)
//		Assert(x.salt isnt: "")
//		passhash = Argon2id(password, x.salt)
//		Assert(x.passhash is: passhash)
//		Assert(auth(user, password, :table))
//		Assert(auth(user, password, :table))
//		y = Query1(table, :user) // no change after the upgrade
//		Assert(x is: y)
//		}

	Test_invalid_arguments()
		{
		Assert({ Auth() } throws: #unhandled)
		Assert({ Auth(#one, #two, #three) } throws: #unhandled)
		}

	Test_token_is_single_use()
		{
		auth = .auth
		token = auth.Auth_token()
		Assert(token isnt: "")
		Assert(auth(token))
		Assert(auth(token) is: false)
		}

	Test_unknown_token()
		{
		auth = .auth
		Assert(auth(.TempName()) is: false)
		}

	auth: Auth
		{
		Getter_Permit()
			{
			return class
				{
				Default(@unused) { return false }
				}
			}

		Auth_setPermit(@unused)
			{
			}
		}
	}
