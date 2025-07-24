// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
`
<html>
<head lang="en">
	<meta charset="UTF-8">

	<style type="text/css">
		html,
		body {
			margin: 0;
			height: 100%;
		}
	</style>
	@for(stylesheet in .stylesheets)
		{
		<link href="@stylesheet" rel="stylesheet"></link>
		}
	@for(script in .scripts)
		{
		<script src="@script"></script>
		}
	@for(style in .styles)
		{
		<style>@style</style>
		}
<body>
	<script>
		window.onload = function (event) {
			@.onload
			window.onload = null;
		};
	</script>
</body>
</html>`
