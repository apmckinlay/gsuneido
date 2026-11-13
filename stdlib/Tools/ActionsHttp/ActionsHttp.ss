// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
// HTTP endpoint for executing action on server for which customers don't have permission.
//
// The request must supply an "action" argument plus whatever else that
// action needs. Only actions contributed by a library are allowed to be
// processed, so this endpoint can never be used to run arbitrary code.
class
	{
	ActionName:  false
	ContribName: false
	EndPoint:    false
	ErrPrefix:   false

	// ---------------- client side ----------------
	Run(@args)
		{
		Assert(.ActionName isnt false and .ErrPrefix isnt false)
		if not args.Member?(0)
			return throw "missing action"

		action = args.Extract(0)
		args[.ActionName] = action
		result = .runOnHttp?()
			? .runOnHttp(args)
			: .process(@args)
		if Object?(result) and result.Member?(#ActionsHttpErr)
			return throw .ErrPrefix $ ": " $ result.ActionsHttpErr
		return result
		}

	runOnHttp?()
		{
		return ServerSuneido.Get(#RunningHttpServer) and not TestRunner.RunningTests?()
		}

	runOnHttp(args)
		{
		Assert(.EndPoint isnt false)
		try
			return RunOnHttp(HttpPort(), .EndPoint, args, auth?:)
		catch (e)
			return Object(ActionsHttpErr: e)
		}

	// ---------------- server side ----------------
	HandleRequest(env)
		{
		return RunWithPack(.process, env, auth?:)
		}

	process(@args)
		{
		action = args.Extract(.ActionName, "")

		actions = .getActions()
		if not actions.Member?(action)
			return .log("action not allowed: " $ action)

		fn = actions[action]
		try
			return (String?(fn) ? Global(fn) : fn)(@args)
		catch (err)
			return .log(action $ ": " $ err, args)
		}

	getActions()
		{
		Assert(.ContribName isnt false)
		return GetContributions(.ContribName)
		}

	log(err, args = #())
		{
		SuneidoLog("ERROR: " $ .ErrPrefix $ ": " $ err, params: args)
		return Object(ActionsHttpErr: err)
		}
	}
