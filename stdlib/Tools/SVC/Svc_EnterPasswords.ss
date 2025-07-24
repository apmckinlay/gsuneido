// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	New()
		{
		super(.layout())
		}
	layout()
		{
		Database("ensure svc_passwords (svc_password) key (svc_password)")
		return #(Browse 'svc_passwords rename svc_password to svc_password_new'
			columns: (password confirm_password),
			validField: 'svc_password_valid', statusBar:)
		}
	}