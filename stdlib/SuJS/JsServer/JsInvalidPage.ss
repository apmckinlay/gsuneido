// Copyright (C) 2022 Axon Development Corporation All rights reserved worldwide.
class
	{
	CallClass(title, reason)
		{
		return '<!DOCTYPE html>' $ Razor(.invalidTempalte, Object(:title, :reason))
		}

	invalidTempalte: `
<html>
	<head>
		<title>@.title</title>
	</head>
	<body>
		<h1>Invalid Request</h1>
		<p>Reason: @.reason</p>
	</body>
</html>
`
	}