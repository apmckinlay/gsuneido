// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
#(
	('Get', '/Display', function (env) { Xml('html', Xml('body', Display(env))) })
	('Get', '/TestResponse', function () { return 'TestResponse' })
	('Get', '/TestRackResponse', function () { return 'Test Rack Response' })
	('Get', '/TestException', function () { throw 'TestException' })
)