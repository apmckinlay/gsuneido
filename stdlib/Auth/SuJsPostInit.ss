// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
function (port)
	{
	Suneido.RunAsStandalone = true
	try Query1('postinit').text.Eval();
	RunSuJSHttpServer(port)
	}