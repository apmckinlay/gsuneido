html {
	overflow:hidden;
	height:100%;
}
body {
	font-size: 16px;
	font-family: Arial;
	height:100%;
	margin:0px 0px 0px 0px;
	-moz-user-select: none;
}

div {white-space:nowrap;}

table {table-layout:fixed;}


.left_cell{
	padding: 0px 0px 0px 5px;
}
.left_cell_table{
	table-layout:auto;
	height:24px;
}
.subtype_table {
	table-layout:auto;
}
.subtype_row{
	height:24px;
}

.calendar_head
{
	width: 100%;
	border-left-width: 3px;
	border-left-color: #C3D9FF;
	border-left-style: solid;
	border-right-width: 3px;
	border-right-color: #C3D9FF;
	border-right-style: solid;
	background-color: #C3D9FF;
}
.head_cell
{
	text-align: center;
	height: 15px;
	border-color: #C3D9FF;
	border-width: 1px;
	color: #112ABB;
	font-size: 12px;
	background-color: #C3D9FF;
}
#div_calendar_frame
{
	overflow: hidden;
	position: relative;
	height: 100%;
	left: 0px;
	top: 0;
	border-color: #C3D9FF;
	border-style: solid;
	border-width: 3px;
}
#div_calendar_box
{
	width: 100%;
	position: absolute;
	left: 0;
	top: 0;
	cursor:url(/Res?name=openhand.cur),move;
}
.div_watermark{
	position: absolute;
	left: 0;
	top: 0;
	z-index: 0;
}
#div_warning{
	right: 13px;
	top: 10px;
	display: none;
	background: #CC4444 none repeat scroll 0 0;
	color: white;
	font-size: 70%;
	padding: 2px 2px 2px 2px;
	position: absolute;
	z-index: 200;
}
.div_calendar{
	position: absolute;
	left: 0;
	top: 0;
	z-index: 1;
}
.div_calendar_week{
	position: absolute;
	left: 0;
	top: 0;
	z-index: 2;
}
.watermark{
	font-style: italic;
	position: absolute;
	overflow: hidden;
	text-align: center;
	vertical-align: middle;
	color: black;
	width: 10px;
	left: 0;
	opacity: 0.15;
	filter: alpha(opacity=15);
}
.div_cell_background
{
	top: 0;
	left: 0;
	height: 100%;
	margin: 0;
	padding: 0 0 0 0;
	border-width: 0;
	opacity: 0.8;
	filter: alpha(opacity=80);
}
.calendar_cell{
	border-color: Gray;
	border-width: 1px;
	border-style: solid;
	vertical-align: top;
	text-align: left;
	color: #000080;
	word-break: break-all;
}
.event_container{
	opacity: 0.9;
	filter: alpha(opacity=90);
	zoom: 1;
}
.event_ctlbar{
	cursor:pointer;
	overflow:hidden;
	position:relative;
	opacity:0.85;
	filter:alpha(opacity=85);
	zoom: 1;
}
.event_bar{
	overflow:hidden;
	position:relative;
	opacity:0.85;
	filter:alpha(opacity=85);
	zoom: 1;
}
.corner{
	overflow: hidden;
	font-size:1px;
	height:1px;
	line-height:1px;
	margin:0 1px;
}
.event_type{
	padding-right: 4px;
	padding-left: 4px;
	overflow: hidden;
	color: white;
	height: 20px;
	display: block;
}
.event_container_invert{
	position: relative;
}
.event_ctlbar_invert{
	cursor:pointer;
	overflow:hidden;
	position:relative;
}
.event_bar_invert{
	overflow:hidden;
	position:relative;
}
.back_invert{
	position:relative;
}
.corner_invert{
	position:relative;
	overflow: hidden;
	font-size:1px;
	height:1px;
	line-height:1px;
	margin:0 1px;
	opacity:0.20;
	filter:alpha(opacity=20);
	zoom: 1;
}
.event_type_invert{
	position:relative;
	padding: 0px 4px;
	height:20px;
	opacity:0.20;
	filter:alpha(opacity=20);
	zoom: 1;
}
.event_name_invert{
	top: 1px;
	position:absolute;
	color: black;
	padding-right: 4px;
	padding-left: 4px;
}
.event_name_invert_none{
	top: 1px;
	position:absolute;
	color: grey;
	padding-right: 4px;
	padding-left: 4px;
}



.event_right_arrow{
	border-color: transparent transparent transparent white;
	border-width: 6px 0 6px 6px;
	right: -3px;
}
.event_left_arrow{
	border-color: transparent white transparent transparent;
	border-width: 6px 6px 6px 0;
	left: -3px;
}
.event_right_arrow_invert{
	border-left-color: black;
}
.event_left_arrow_invert{
	border-right-color: black;
}
.event_arrow{
	top: 5px;
	position: absolute;
	border-style: solid;
	font-size-adjust: none;
	line-height: 0;
	vertical-align: middle;
	font-size: 18px;
}



.event_txt{
	overflow:hidden;
	color:white;
	font-size:15px;
}
.div_dot{
	position:absolute;
	overflow:hidden;
	font-size: 18px;
	top: -1px;
	left: 2px;
}


.tooltip_color{
	background-color:#888888;
	border-color:#888888;
}
#tooltip{
	overflow:visible;
	display:none;
	top:0;
	position:absolute;
	font-size:12px;
	z-index: 100;
	cursor:url(/Res?name=openhand.cur),move;
	background-color:#888888;
	opacity: 0.9;
	filter: alpha(opacity=90);
}
#tooltip_head{
	white-space: normal;
	overflow: visible;
	font-size: 16px;
	background-color:#888888;
	border-color:#888888;
	padding:3px 8px 3px 8px;
	margin:0;
	color: white;
}
#tooltip_body{
	background-color:#888888;
	border-color:#888888;
	overflow: visible;
	padding:0 5px 4px 6px;
	margin: 0;
	color:white;
	white-space: normal;
}

a:visited{ color:blue }
.a_navi{ font-size: 17px; }

.week_head
{
	font-size: 15px;
	height: 18px;
	text-align: left;
	font-family: Arial;
	color: #000080;
	overflow: hidden;
}
.event_cell
{
	height: 24px;
	padding-left: 4px;
	padding-right: 4px;
	vertical-align: top;
	text-align: left;
	word-break: break-all;
	overflow: hidden;
}
.last_cell{
	padding-left: 4px;
	padding-right: 4px;
	vertical-align: top;
	word-break: break-all;
	overflow: hidden;
}
.cell_more{
	text-align: center;
	vertical-align: top;
	word-break: break-all;
	overflow: hidden;
	font-size: 14px;
}
#div_expand{
	display: none;
	position: absolute;
	left: 0px;
	top: 0px;
	z-index: 60;
	border-width: 1px;
	border-style: solid;
	border-color: Gray;
	cursor:url(/Res?name=openhand.cur),move;
	background-color: white;
	overflow: hidden;
}

.div_expand_head{
	position: relative;
	background-color: #E8EEF7;
	cursor: default;
	margin-bottom: 3px;
}

.div_expand_close{
	font-size: 18px;
	position: absolute;
	top: 0px;
	right: 3px;
	cursor: pointer;
}

#zoom_bar{
	position:absolute;
	z-index:100;
	height:7px;
	cursor:pointer;
	left:3px;
	right:3px;
	top:115px;
	background-color:#669900;
	font-size:1px;
}
.c1{
	font-size: 1px;
	height:1px;
	line-height:1px;
	background-color:#C3D9FF;
	border-color:#C3D9FF;
	margin:0 0 0 1px;
}

#cur_range{
	height: 18px;
	padding-top: 10px;
	white-space: nowrap;
	color: #000
	font-size: 13px;
	font-weight: bold;
}

#map_bar{
	position:relative;
	display:block;
	background-color:#C3D9FF;
}

.left_head{
	font-size: 100%;
	background-color: #C3D9FF;
	padding: 3px 0 3px 8px;
	margin: 0;
	margin-bottom: 2px;
}

#leftCtl{
	width: 100%;
	background-color:#FFFFFF;
}

.left_foot{
	margin-top: 2px;
	background-color:#C3D9FF;
	height:10px;
}

#div_color_menu{
	position:absolute;
	display:none;
	z-index:20;
	text-align:center;
	color:white;
	cursor:pointer;
}

.navi_table{
	width:130px;
}
