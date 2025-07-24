// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
class
	{
	New(routesList)
		{
		.routes = routesList.
			Map({ NameArgs(it.Copy(), #(method, pat, app)) }).Instantiate()
		}
	Call(env)
		{
		if false is route = .find_route(env)
			{
// turning off the logging until the authentication is fixed (36289) because http attacks
// will kick this in and we don't want that
// SuneidoLog('can not find route', params: env)
			return ['404 page not found', #(), 'page not found']
			}

		app = route.app
		if String?(app)
			app = Global(app)
		return app(:env)
		}
	find_route(env)
		{
		return .routes.FindOne({ env.method =~ ('^' $ it.method.Upper() $ '$') and
			env.path =~ ('^' $ it.pat $ '\>') })
		}
	}