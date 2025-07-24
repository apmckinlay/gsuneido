// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
function (env) // get a library record, transpiled to JavaScript
	{
	name = env.query
	for lib in SuCode().Libraries
		if false isnt x = Query1(lib, group: -1, :name)
			{
			SuCode().Add(lib, name)
			return [200/*=ok*/, ['Content-Type': 'text/plain'],
				JsTranslate(x.text, name, lib)]
			}
	return Object('404 Not Found', [], 'not found')
	}