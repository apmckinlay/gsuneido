// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass()
		{
		return '<html>
			<head>
			<title>$title</title>
			</head>
			<style type="text/css"">
			<!-- 64em should evaluate, on a standard browser, to 1024px. -->
			#body { width: 64em; }
			</style>
			<style type="text/css" media="print">
			.noPrint { display: none; }
			</style>
			<body>' $
			.Find $
			'<h1 style="margin-top: 0;">' $
				OptContribution('WikiTemplateLogo', '') $
				'<a href="Wiki?find=$page" alt="Find links to this page"
					title="Find links to this page">$title</a>$orphaned</h1>
			<p>
			$body
			<div class="noPrint">
			<p>
			<hr>
			$action ($created $lastedit)
			<br>
			<a href="Wiki?RecentChanges">Recent Changes</a>
			&nbsp;&nbsp;|&nbsp;&nbsp;
			Return to <a href="Wiki?StartPage">Start Page</a></div>
			</body>
			</html>'
		}
	Find: '<div style="float: right; margin-top: 10;" class="noPrint">
			<form method="get" action="Wiki">
			$smallaction &nbsp;
			<input type="text" size="20" name="find" spellcheck="true">
			<input type="submit" value="Find">
			</form><br />
			</div>\n'
	}
