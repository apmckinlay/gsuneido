// Copyright (C) 2020 Axon Development Corporation All rights reserved worldwide.
class
	{
	PreLogin: "PRELOGIN_"
	CallClass(origCmd = '')
		{
		if false is (userPass = .LoginDialog()) or
			false is Authorize(userPass.user, userPass.password)
			{
			Exit()
			return false
			}
		else
			{
			if origCmd is IDESwitchMode.CSDevCmd
				{
				IDESwitchMode.CSDevCmd.Eval()
				return true
				}
			.SetUser(userPass.user)
			.PostLoginPlugins()
			PersistentWindow.Load()
			}
		return true
		}

	LoginDialog()
		{
		return ToolDialog(0, Object(.loginCtrl),
			title: 'Login', closeButton?: false, keep_size: false)
		}

	loginCtrl: Controller
		{
		New()
			{
			.user = .FindControl('user')
			.pw = .FindControl('password')
			}
		Controls: #(Vert
			(Pair (Static Name) (Field width: 15, name: user))
			(Pair (Static Password) (Field width: 15, password:, name: password))
			Skip OkCancel)
		On_OK()
			{
			.Window.Result(Object(user: .user.Get(), password: .pw.Get()))
			}
		// prevent close without using either ok or cancel button
		ConfirmDestroy()
			{
			return false
			}
		}

	PostLoginPlugins()
		{
		for c in Contributions('PostLogin')
			c()
		}

	SetUser(user)
		{
		Suneido.User = Suneido.User_Loaded = user
		}
	}
