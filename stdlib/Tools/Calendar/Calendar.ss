// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
class
	{
	Dev: false
	CallClass(allPage, date, session = '')
		{
		return	.Head(allPage, date, session) $ .Body()
		}
	Head(allPage, date, session)
		{
		return '<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.1//EN"
			"http://www.w3.org/TR/xhtml11/DTD/xhtml11.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
<head>
	<meta http-equiv="content-type" content="text/html; charset=utf-8" />
	<title>Calendar</title>' $
(.Dev
	? '<link rel="stylesheet" href="/CodeRes?calendar1_release.css" type="text/css"/>'
	: '<link rel="stylesheet" href="/CodeRes?name=calendar1_7.css" type="text/css" />') $
.Script_Main(allPage, date, session) $
'</head>'
		}
	Script_Main(allPage, date, session = '')
		{
		return '<script type="text/javascript">
		var today = new Date(' $ String(date.Year()) $
			',' $ String(date.Month()-1) $ ',' $ String(date.Day()) $ ');
		var page = "' $ allPage $ '";
		var session = "' $ session $ '";
	</script>'
	}
	Body()
		{
		return `<body><table cellpadding='0' cellspacing='0'>
<tr><td ROWSPAN=2 width=150px style=vertical-align:top;>
	<div style=width:150px>
		<table cellpadding=0 cellspacing=0><tr>
			<td valign=top style=width:130px;>
				<div style=height:15px;>&nbsp;</div>

				<div id='div_navi'></div>

			</td>
			<td align=center width=20px style=vertical-align:top;>
				<div style=height:3px;>
					<div class=c1 style=margin-left:3px;></div>
					<div class=c1 style=margin-left:2px;></div>
					<div class=c1 ></div>
				</div>

				<div id=map_bar></div>

			</td>
		</tr></table>

		<div id=div_cal_events style=overflow:hidden;>
			<div style=height:3px>
				<div class=c1 style=margin-left:3px;></div>
				<div class=c1 style=margin-left:2px;></div>
				<div class=c1 ></div>
			</div>
			<div class=left_head>Calendar Events</div>

			<div style=background-color:#C3D9FF; id='div_left'></div>

			<div class=left_foot></div>
			<div style=height:3px>
				<div class=c1></div>
				<div class=c1 style=margin-left:2px;></div>
				<div class=c1 style=margin-left:3px;></div>
			</div>
		</div>
	</div>

</td>
<td valign='top'>

	<div id='div_cal_head'>
		<table id='cal_head' class='calendar_head'	cellpadding='0' cellspacing='0'>
			<tr>
				<td class='head_cell'>Sun</td><td class='head_cell'>Mon</td>
				<td class='head_cell'>Tue</td><td class='head_cell'>Wed</td>
				<td class='head_cell'>Thu</td><td class='head_cell'>Fri</td>
				<td class='head_cell'>Sat</td>
			</tr>
		</table>
		<div id='div_warning'>&nbsp;</div>
	</div>
	<div id='div_calendar_frame'></div>
<noscript>
	<br/>
	<strong>
		&nbsp;Calendar Requires Javascript enabled on your Internet Explorer.
	</strong><br/>
	&nbsp;&nbsp;&nbsp;You can find setting at:<br/>
	&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;Internet Explorer >> Tools >> Internet Options >> <br/>
	&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;` $
	`Security >> Custom level >> Scripting >> Active scripting(need to be Enable)<br/>
	<br/>
</noscript>

</td>
</tr>
</table>
	<script type=text/javascript src="/CodeRes?name=calendar1_9_9.js&20250316"></script>
	<script type="text/javascript">
		document.addEventListener("DOMContentLoaded", function() {
		load();
	});
	</script>
</body>
</html>`
		}
	}
