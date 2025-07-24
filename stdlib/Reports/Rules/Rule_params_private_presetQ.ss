// Copyright (C) 2024 Axon Development Corporation All rights reserved worldwide.
function ()
	{
	return .report.Has?(`~presets~`) and Object?(.report_options) and
		.report_options.GetDefault(`private?`, false)
	}