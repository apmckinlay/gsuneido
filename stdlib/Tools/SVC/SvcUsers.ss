// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
/* Helper class for handling SVC user credentials.
This class:
- Includes the Server side and Client side methods
- Purposely does NOT go through SvcSocketClient.Run or Send.
	- This is to circumvent the standard validation/error checking.
	- This will allow svc_users to be managed regardless of
		SvcSocketClient error state / credentials
	- A valid password in svc_passwords is required to manage svc_users
*/
class
	{
	// For standalone credential testing
	TestLogin(svcuser_id, password, hashed = false)
		{
		sc = .socketClient()
		key = .send(sc, [#NONCE])
		passhash = hashed
			? password
			: PassHash(svcuser_id, password)
		svcuser_passhash = Sha256(key $ passhash)
		return .send(sc, [#LOGIN, :svcuser_id, :svcuser_passhash], close?:)
		}

	socketClient()
		{
		sc = SvcSocketClient().SC()
		if sc.Readline() isnt SvcSocketClient.ValidServer
			{
			sc.Close()
			throw 'SvcUsers - Error: ' $ SvcSocketClient.InvalidServer
			}
		return sc
		}

	svcUsersTable: svc_users
	// Login does not have a client side specific request.
	// As the SocketClient needs to be validated before prior use, a (possibly temporary)
	// SocketClient must be initalized and used for this request.
	// 	- See SvcSocketClient.SvcSocketClient_verifyCredentials for Client Side request
	LoginRequest(key, svcuser_id, svcuser_passhash)
		{
		.ensure()
		if false isnt rec = Query1(.svcUsersTable, :svcuser_id)
			return .passwordMatch?(key, rec.svcuser_passhash, svcuser_passhash)
		return false
		}

	ensure()
		{
		Database('ensure ' $ .svcUsersTable $ ' (svcuser_id, svcuser_passhash)
			key (svcuser_id)')
		}

	passwordMatch?(key, serverSide, clientSide)
		{ return Sha256(key $ serverSide) is clientSide }

	// Client Call
	AddUser(serverPassword, svcuser_id, userPassword)
		{
		sc = .socketClient()
		serverPassword = .encrypt(sc, serverPassword)
		svcuser_passhash = PassHash(svcuser_id, userPassword)
		args = [#ADDUSER, :serverPassword, :svcuser_id, :svcuser_passhash]
		return .send(sc, args, close?:)
		}

	encrypt(sc, password)
		{
		passhash = PassHash('', password)
		return Sha256(.send(sc, [#NONCE]) $ passhash)
		}

	send(sc, args, close? = false)
		{
		result = false
		try
			{
			SvcSocketClient.Write(sc, args)
			result = SvcSocketClient.Read(sc)
			}
		catch (e)
			{
			try sc.Close()
			throw 'SvcUsers - Error during send: ' $ e
			}
		if close?
			sc.Close()
		return result
		}

	// Server Response
	AddUserRequest(key, serverPassword, svcuser_id, svcuser_passhash)
		{
		.ensure()
		if true isnt msg = .verifyServer(key, serverPassword)
			return msg
		rec = [:svcuser_id, :svcuser_passhash]
		return false is QueryOutputIfNew(.svcUsersTable, rec)
			? 'User already exists'
			: true
		}

	// Server requires the table svc_passwords with a password prior to setting up users
	svcPasswordTable: svc_passwords
	verifyServer(key, serverPassword)
		{
		if not TableExists?(.svcPasswordTable) or QueryCount(.svcPasswordTable) is 0
			return 'Server is not secure'
		QueryApply(.svcPasswordTable)
			{
			if .passwordMatch?(key, it.svc_password, serverPassword)
				return true
			}
		return 'Invalid Server Password'
		}

	// Client Call
	ChangePassword(serverPassword, svcuser_id, oldPassword, newPassword)
		{
		sc = .socketClient()
		serverPassword = .encrypt(sc, serverPassword)
		oldPasshash = PassHash(svcuser_id, oldPassword)
		newPasshash = PassHash(svcuser_id, newPassword)
		args = [#CHANGEPASSWORD, :serverPassword, :svcuser_id, :oldPasshash, :newPasshash]
		return .send(sc, args, close?:)
		}

	// Server Response
	ChangePasswordRequest(key, serverPassword, svcuser_id, oldPasshash, newPasshash)
		{
		.ensure()
		if true isnt msg = .verifyServer(key, serverPassword)
			return msg
		QueryApply1(.svcUsersTable, :svcuser_id)
			{
			if oldPasshash isnt it.svcuser_passhash
				return 'Old password does not match the saved password'
			it.svcuser_passhash = newPasshash
			it.Update()
			return true
			}
		return false
		}

	// Client call
	DeleteUser(serverPassword, svcuser_id)
		{
		sc = .socketClient()
		serverPassword = .encrypt(sc, serverPassword)
		return .send(sc, [#DELETEUSER, :serverPassword, :svcuser_id], close?:)
		}

	// Server Response
	DeleteUserRequest(key, serverPassword, svcuser_id)
		{
		.ensure()
		if true isnt msg = .verifyServer(key, serverPassword)
			return msg
		return QueryDelete(.svcUsersTable, [:svcuser_id]) is 1
		}
	}
