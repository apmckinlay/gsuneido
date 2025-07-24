// Copyright (C) 2023 Axon Development Corporation All rights reserved worldwide.
function (@unused)
	{
	return ['200 OK', [Expires: Date().Plus(days: 20).InternetFormat()],
		'User-agent: *
Disallow: /']
	}