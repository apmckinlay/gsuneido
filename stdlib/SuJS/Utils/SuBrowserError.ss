// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
function (message, stack)
	{
	if SuRenderBackend().InErrorHandler is true
		{
		reason = 'error in Handler (err: ' $ message $ ')'
		SuRenderBackend().Terminate(e: 'SuBrowserError FATAL - ' $ reason, :reason)
		}
	if String?(stack)
		stack = stack.Lines()[1..]
	calls = stack.Map({ Object(fn: it, locals: #()) })
	SuRenderBackend().InErrorHandler = true
	SuRenderBackend().DumpStatus(message)
	Handler(message, 0, calls)
	SuRenderBackend().InErrorHandler = false
	}