// IMPORTANT: If you change this record, also need to update the date arg in Calendar
// for the CodeRes call that loads this js code. This argument does not need to be
// anything specific, just needs to be changed to invalidate browser caches

Date.prototype.toDateUrlString = function(){ return this.getFullYear() + '-' + (this.getMonth() + 1)  + '-' + this.getDate(); };
Date.prototype.plusDays = function (plusDate){ return new Date(this.getFullYear(), this.getMonth(), this.getDate() + plusDate); };
Date.prototype.minusDate = function (date) { return Math.round((this - date) / 86400000);};
Date.prototype.firstDateOfWeek = function(){ return this.plusDays( -1 * this.getDay());};
String.prototype.toDate = function (){	var ymd = this.split('-'); return new Date(ymd[0], ymd[1] - 1, ymd[2]);};

var subtypeTableSelector = 'table.subtype_table';
var classTypes = ['event_type', 'event_name_invert', 'event_name_invert_none'];
var expandSelector = 'tbody tr td div.event_type, tbody tr td div.event_name_invert, tbody tr td div.event_name_invert_none';

/*
 * "plugin" variable was originally from a jQuery JSON Plugin that adds to the jQuery functionality
 * now it is a variable only local to Axon calendar logics
 * version: 2.1 (2009-08-14)
 *
 * This document is licensed as free software under the terms of the
 * MIT License: http://www.opensource.org/licenses/mit-license.php
 *
 * Brantley Harris wrote this plugin. It is based somewhat on the JSON.org
 * website's http://www.json.org/json2.js, which proclaims:
 * "NO WARRANTY EXPRESSED OR IMPLIED. USE AT YOUR OWN RISK.", a sentiment that
 * I uphold.
 *
 * It is also influenced heavily by MochiKit's serializeJSON, which is
 * copyrighted 2005 by Bob Ippolito.
 */
var plugin = (function() {
	var _escapeable = /["\\\x00-\x1f\x7f-\x9f]/g;

	var _meta = {
		'\b': '\\b',
		'\t': '\\t',
		'\n': '\\n',
		'\f': '\\f',
		'\r': '\\r',
		'"' : '\\"',
		'\\': '\\\\'
	};

	return {
		toJSON: function(o)
		{
			if (typeof(JSON) == 'object' && JSON.stringify) {
				return JSON.stringify(o);
			}

			var type = typeof(o);

			if (o === null)
				return "null";

			if (type == "undefined")
				return undefined;

			if (type == "number" || type == "boolean")
				return o + "";

			if (type == "string")
				return this.quoteString(o);

			if (type == 'object')
			{
				if (typeof o.toJSON == "function")
					return this.toJSON( o.toJSON() );

				if (o.constructor === Date)
				{
					var month = o.getUTCMonth() + 1;
					if (month < 10) month = '0' + month;

					var day = o.getUTCDate();
					if (day < 10) day = '0' + day;

					var year = o.getUTCFullYear();

					var hours = o.getUTCHours();
					if (hours < 10) hours = '0' + hours;

					var minutes = o.getUTCMinutes();
					if (minutes < 10) minutes = '0' + minutes;

					var seconds = o.getUTCSeconds();
					if (seconds < 10) seconds = '0' + seconds;

					var milli = o.getUTCMilliseconds();
					if (milli < 100) milli = '0' + milli;
					if (milli < 10) milli = '0' + milli;

					return '"' + year + '-' + month + '-' + day + 'T' +
								 hours + ':' + minutes + ':' + seconds +
								 '.' + milli + 'Z"';
				}

				if (o.constructor === Array)
				{
					var ret = [];
					for (var i = 0; i < o.length; i++)
						ret.push( this.toJSON(o[i]) || "null" );

					return "[" + ret.join(",") + "]";
				}

				var pairs = [];
				for (var k in o) {
					var name;
					var type = typeof k;

					if (type == "number")
						name = '"' + k + '"';
					else if (type == "string")
						name = this.quoteString(k);
					else
						continue;  //skip non-string or number keys

					if (typeof o[k] == "function")
						continue;  //skip pairs where the value is a function.

					var val = this.toJSON(o[k]);

					pairs.push(name + ":" + val);
				}

				return "{" + pairs.join(", ") + "}";
			}
		},

		evalJSON: function(src)
		{
			if (typeof(JSON) == 'object' && JSON.parse)
				return JSON.parse(src);
			return eval("(" + src + ")");
		},

		secureEvalJSON: function(src)
		{
			if (typeof(JSON) == 'object' && JSON.parse)
				return JSON.parse(src);

			var filtered = src;
			filtered = filtered.replace(/\\["\\\/bfnrtu]/g, '@');
			filtered = filtered.replace(/"[^"\\\n\r]*"|true|false|null|-?\d+(?:\.\d*)?(?:[eE][+\-]?\d+)?/g, ']');
			filtered = filtered.replace(/(?:^|:|,)(?:\s*\[)+/g, '');

			if (/^[\],:{}\s]*$/.test(filtered))
				return eval("(" + src + ")");
			else
				throw new SyntaxError("Error parsing JSON, source is not valid.");
		},

		quoteString: function(string)
		{
			if (string.match(_escapeable))
			{
				return '"' + string.replace(_escapeable, function (a)
				{
					var c = _meta[a];
					if (typeof c === 'string') return c;
					c = a.charCodeAt();
					return '\\u00' + Math.floor(c / 16).toString(16) + (c % 16).toString(16);
				}) + '"';
			}
			return '"' + string + '"';
		}
	}
})();

var deepExtend = function(out) {
	if (!out) {
		return {};
	}

	var args = Array.prototype.slice.call(arguments, 1);
	for (var i = 0; i < args.length; i++) {
		if (!args[i]) {
			continue;
		}

		for (var key in args[i]) {
			if (!args[i].hasOwnProperty(key)) {
				continue;
			}

			switch (Object.prototype.toString.call(args[i][key])) {
				case '[object Object]':
				  out[key] = out[key] || {};
				  deepExtend(out[key], args[i][key]);
				  break;
				case '[object Array]':
				  out[key] = args[i][key];
				  break;
				default:
				  out[key] = args[i][key];
			}
		}
	}

  return out;
};

var CAL = function(){

	return {
	  all_events 					: [],
	  not_need_download 	: [],

	  other_prefs : Object(),
	  prefs 			: Object(),
	  event_types : Object(),    // all type info
	  zoom 				: 6,
	  search_text : '',
		types_count : 0,

		read_week_events: function (data){
		var events = data.events;
		var week_date = data.from.toDate();
		var ind = 0;
		while(ind < events.length)
			ind = this.set_week_events(week_date, events, ind);
		//this.set_downloaded(data.from, data.src);
		return events.length < 0 ? null : week_date;
	  },

		set_week_events: function (week_date, events, ind){
			var key = week_date.toDateUrlString();
			var week = this.all_events[key] ? this.all_events[key] : this.init_week_events(key);
			var nLen = events.length;
			while(ind < nLen) {
			  var i = events[ind].start_date.toDate().plusDays(events[ind].completed ? events[ind].completed : 0).minusDate(week_date);
			  if(i >= 7 || i < 0) break;
			  this.set_date_events(week, i, events[ind]);
			  ind++;
			}
			//sort by length
			for(var i = 0; i < 7; i++)
			  week[i].events.sort(function (a,b){ return b.multidays - a.multidays;});
			this.reset_week_events(week);
			return ind;
		},

	  init_week_events: function(key){
		delete this.all_events[key];
		return this.all_events[key] = [{events:[]}, {events:[]}, {events:[]},
		  {events:[]}, {events:[]}, {events:[]},{events:[]}];
	  },

	  set_date_events: function(week, i, e){
		e.index = e.type === 'ERROR' ? 0 : this.event_types[e.type].index;
		week[i].events.push(e);
	  },

	  get_week_events: function(week_date){
		var key = week_date.toDateUrlString();
		if(!this.all_events[key])
		  return this.init_week_events(key);

			var week = this.all_events[key];
			this.reset_week_events(week);
			return week;
	  },

	  reset_week_events: function(week){
		for(var i = 0; i < 7; i++) {
		  var wi = week[i];
		  wi.num = wi.showed= wi.co_showed = wi.co_num = 0;
		}
		for(var i = 0; i < 7; i++) {
		  var wi = week[i];
		  for(var j = 0, nLen = wi.events.length; j < nLen; j++) {
			var e = wi.events[j];
					if(!this.is_hidden(e)) {
						wi.num++;
						for(var k = 1; k < e.span && i + k < 7; k++)
						  week[i + k].co_num++;
			}
		  }
		}
	  },

	  get_date_events: function(d){
		var date_events = [];
		var key = d.plusDays(-1 * d.getDay()).toDateUrlString();
		if(this.all_events[key]) {
		  var w_events = this.all_events[key];
		  for(var i = 0; i < 7; i++)
			for(var j = 0; j < w_events[i].events.length; j++) {
			  var e = w_events[i].events[j];
			  var diff = -1 * e.start_date.toDate().plusDays(e.completed ? e.completed : 0).minusDate(d);
			  if(!this.is_hidden(e) && diff < e.span && diff >= 0)
				date_events[date_events.length] = e;
			}
		}
		return date_events;
	  },

	  query_types: function(callback){
		var xhr = new XMLHttpRequest();
		xhr.open('GET', page + '?types' + '&session=' + session);
		xhr.onload = function() {
			if (xhr.status >= 200 && xhr.status < 300) {
				callback(JSON.parse(xhr.responseText));
			}
		};
		xhr.send();
	  },

		query_events: function(week_date, callback){
			var req = this.get_request(week_date);
			function request_events(url, callback, retry) {
				var xhr = new XMLHttpRequest();
				xhr.open('GET', url);
				xhr.onload = function () {
					if (xhr.status >= 200 && xhr.status < 300) {
						callback(xhr.responseText);
					}
				};
				xhr.onerror = function () {
					if (retry) {
						var dateFrom = week_date.toDateUrlString()
						var dateTo = week_date.plusDays(6).toDateUrlString()
						callback('(' + JSON.stringify({
							from: dateFrom, to: dateTo,
							events: [{completed: 0, desc: '', end: true,
								end_date: dateTo, multidays: 7, span: 7, start_date: dateFrom, subtype: '',
								title: 'Connection Issue: Unable to load, please refresh (F5)',
								type: 'ERROR'}]}) + ')');
					} else {
						request_events(url, callback, true);
					}
				}
				xhr.send();
			}

			if (req == '') {
				callback(week_date);
			} else {
				var url = page + '?from=' + week_date.toDateUrlString() + '&weeks=1&req=' + encodeURIComponent(req) + '&session=' + session;
				request_events(url, callback);
			}
		},

		clear_cache: function(){
			document.location.reload(true);
		},

		need_download: function(d, src){
		var k = d.toDateUrlString();
		if(!this.not_need_download[k])
			this.not_need_download[k] = [];
		if(this.not_need_download[k][src])
			return false;
		else {
			this.not_need_download[k][src] = true;
			return true;
		}
	  },

		get_request: function(week_date){
			var arr = [];
			var types = this.event_types;
			for(var t in types) {
		if(!types[t].hasOwnProperty('SubTypes') && types[t].Display != 'none' && this.need_download(week_date, t))
					arr.push(t);
				else {
					var subs = types[t].SubTypes;
					for(var st in subs){
						var st_str = t + '__' + st;
						if(subs[st].Display != 'none' && this.need_download(week_date, st_str))
							arr.push(st_str);
					}
				}
			}
			return arr.join('$*$');
		},

		save_settings: function(){
			var left = deepExtend({}, this.event_types);
			CAL.extract(left, this.prefs);
			for(var t in this.other_prefs) {
				if(!left[t]) {
					left[t] = {};
				}
				deepExtend(left[t], this.other_prefs[t]);
			}

			var cookie_ob = Object();
			if(this.zoom != 6)
				cookie_ob.zoom = this.zoom;
			cookie_ob.event_types = left;
			this.setCookie('CAL', plugin.toJSON(cookie_ob));
		},

		get_settings: function(data){
			var n = 0;
			for(var t in data) { // default value for Type.Expanded
				n++;
				data[t].Expanded = 'false';
			}
			this.types_count = n;

			var cookies = this.getCookie('CAL');
			cookies = (cookies == '' ? {} : plugin.evalJSON(cookies));
			this.zoom = cookies.zoom ? cookies.zoom : 6;

			cookies = cookies.event_types;
			// save and delete event type not in data list
			this.filter_cookies(cookies, data);
			// copy data to prefs which is used for default settings
			deepExtend(this.prefs, data);
			// read old and delete cookie data
			for(var t in data)
				this.read_old_cookie(data[t], t);
			// extend to event_types
			deepExtend(this.event_types, data, cookies);
			// save settings for old cookies
			this.save_settings();
		},

		filter_cookies: function(cookies, data){
			var cookies_copy = deepExtend({}, cookies);
			for(var t in cookies_copy) {
				if(!data[t]) {
					this.other_prefs[t] = deepExtend({}, cookies[t]);
					delete cookies[t];
				} else if(data[t].SubTypes && cookies[t].SubTypes) {
					for(var subt in cookies[t].SubTypes)
						if(!data[t].SubTypes[subt]) {
							if(!this.other_prefs[t])
								this.other_prefs[t] = {SubTypes: {}};
							this.other_prefs[t].SubTypes[subt] = deepExtend({}, cookies[t].SubTypes[subt]);
							delete cookies[t].SubTypes[subt]
						}
				}
			}
		},

		read_old_cookie: function(type_ob, t){
			this.get_old_cookie(type_ob, 'Display', t + '.display');
			this.get_old_cookie(type_ob, 'Color', t + '.color');
			this.get_old_cookie(type_ob, 'Inverted', t + '.invert');
			var subtypes = type_ob.SubTypes;
			if(subtypes)
				for(var subt in subtypes) {
					this.get_old_cookie(subtypes[subt], 'Display', t + '__' + subt + '.display');
					this.get_old_cookie(subtypes[subt], 'Color', t + '__' + subt + '.color');
				}
		},

		get_old_cookie: function(type_ob, t, old_cookie) {
			var c = this.getCookie(old_cookie);
			if(type_ob[t] && c != '') {
				type_ob[t] = c;
				this.deleteCookie(old_cookie);
			}
		},

		is_hidden: function(e){
			if(this.search_text != '' && e.title.toLowerCase().indexOf(this.search_text) < 0)
				return true;
			if(e.type === 'ERROR')
				return false;
			var subtypes = this.event_types[e.type].SubTypes;
			var type_ob = (e.subtype && e.subtype != '' && subtypes && subtypes[e.subtype] ? subtypes[e.subtype] : this.event_types[e.type]);
			return type_ob.Display == 'none';
		},

		moveCursor: function(events, cur) {
			while(cur < events.length && CAL.is_hidden(events[cur]))
				cur++;
			return cur;
		},

		sort_event_types: function(){
			var types_copy = [];
			for(var type in this.event_types){
				var ob = Object();
				ob.type = type;
				ob.Index = this.event_types[type].Index;
				types_copy.push(ob);
			}
			return types_copy.sort(function(a,b){return a.Index - b.Index;});
		},

		extract: function(target){
			var len = arguments.length, options;
			for (var i = 1; i < len; i++ )
				if ( (options = arguments[ i ]) != null )
					for ( var name in options ) {
						var src = target[ name ], copy = options[ name ];
						// Recurse if we're merging object values
						if ( copy && typeof copy === "object" && !copy.nodeType ) {
							this.extract(src, copy );
							if( this.isEmpty(src) )
								delete target[ name ];
						}
						else if ( copy === target[ name ])
							delete target[ name ];
				  }
		},

		isEmpty: function(ob){
			for(var i in ob){ if(ob.hasOwnProperty(i)){return false;}}
			return true;
		},

		getCookie: function(name){
		  var cookies = document.cookie;
		  name += '=';
		  if (cookies.indexOf(name) != -1){
			var startpos = cookies.indexOf(name) + name.length;
			var endpos = cookies.indexOf(';', startpos);
			if (endpos == -1)
				endpos = cookies.length;
			return unescape(cookies.substring(startpos, endpos));
		  }
		  else
			return ''; // the cookie couldn't be found!
		},

		setCookie: function(name, value){
			// no expiration date specified? 1 year to expire
			if(arguments[2])
			var expires = arguments[2];
			else{
				var expires = new Date();
				expires.setFullYear(expires.getFullYear() + 1);
			}
		  this.deleteCookie(name);
		  document.cookie = (name + '=' + escape(value) + ';expires=' + expires.toGMTString() + ';path=/;');
		},

		deleteCookie: function(name) {
		  var value = this.getCookie(name);
			if(value != '')
					document.cookie = (name + '=' + value + ';expires=Thu, 01-Jan-1970 00:00:01 GMT;path=/;');
		}
	}
}();

var CAL_BODY = function(CAL){
	var today_color = '#FFFF88';
	var MAX_EXPAND 	= 6;

	var b_cal_released= true;
	var b_first_load  = true;
	var cal_rez_last	= null;
	var cal_rez_timer = null;

	return {
		first_day : new Date(),
		cur_month : 0,  	//current month
		cur_year	: 0,   	//currnt year
		cal_lines : 6,

		originY : 0,
		pressY 	: 0,

		cell_min_width 		: 54,
		cell_min_height 	: 54,
		cal_height 				: 0, 	//without border size
		cal_width 				: 0,	//can be set to style.
		cell_height 			: 0,
		cell_width 				: 0,
		event_height			: 24,
		cell_date_height 	: 18,
		cell_hold 				: 0,

		cal_frame : 	null,
		cal_box 	: 	null,
		cal_grid	: 	null,
		cal_week	: 	null,
		watermarks:		[],
		cal_expand: 	null,
		cal_warning:  null,

		month_abbrv : ['Jan','Feb','Mar','Apr','May','Jun', 'Jul','Aug','Sep','Oct','Nov','Dec'],
		month_name 	: ['JANUARY', 'FEBRUARY', 'MARCH', 'APRIL', 'MAY', 'JUNE', 'JULY', 'AUGUST', 'SEPTEMBER', 'OCTOBER', 'NOVEMBER', 'DECEMBER'],

		/// Interface
		onDateChanged	:   function(){},
		onTypesLoaded	: 	function(){},
		onFirstLoaded	:		function(){},

		/// Init
		initCalendar: function(div_calendar_frame){
			this.cal_frame = document.getElementById(div_calendar_frame);
			var strHtml = "<div id='div_calendar_box' class='div_calendar_box'>";
			strHtml += "<div id='div_calendar_week' class='div_calendar_week'></div>";
			strHtml += "<div id='div_calendar' class='div_calendar'>";
			strHtml += "<table id='calendar' cellpadding='0' cellspacing='0'></table>";
			strHtml += "</div>";
			strHtml += "<div id='div_watermark' class='div_watermark'>";
			strHtml += "<div id='watermark1' class=watermark></div><div id='watermark2' class=watermark></div>";
			strHtml += "<div id='watermark3' class=watermark></div><div id='watermark4' class=watermark></div>";
			strHtml += "</div><div id='div_expand'></div>";
			strHtml += "<div id='tooltip' onselectstart='return false;'><div id='tooltip_head'>title</div><div id='tooltip_body'>desc</div></div>";
			//<img src='/Res?name=loading.gif'/>
			strHtml += "</div>";
			this.cal_frame.innerHTML = strHtml;
			this.hookEvent(document, 'selectstart', this.cancelEvent);
			this.hookEvent(document, 'dragstart', this.cancelEvent);

			this.cal_box = document.getElementById('div_calendar_box');
			this.cal_grid = document.getElementById('calendar');
		  this.cal_week = document.getElementById('div_calendar_week');
		  this.watermarks[0] = document.getElementById('watermark1');
		  this.watermarks[1] = document.getElementById('watermark2');
		  this.watermarks[2] = document.getElementById('watermark3');
		  this.watermarks[3] = document.getElementById('watermark4');
		  this.cal_expand = document.getElementById('div_expand');
		  this.cal_warning = document.getElementById('div_warning');

		  this.hookEvent(this.cal_box, 'mousewheel', this.mouseScrolled);
		  this.hookEvent(this.cal_box, 'mousedown', this.calPressed);
		  this.hookEvent(this.cal_expand, 'mousedown', function(e) { b_cal_released = true; return CAL_BODY.cancelEvent(e);});
		  this.hookEvent(document, 'keydown', this.calKeyDown);
		  this.hookEvent(document, 'keypress', this.calKeyDown);
		  this.hookEvent(document, 'click', this.clear);
		  this.hookEvent(window, 'dragstart', this.cancelEvent);
		  this.hookEvent(window, 'resize', function() {
			cal_rez_last = new Date();
			if(!cal_rez_timer)
				cal_rez_timer = setInterval(CAL_BODY.calResized, 50);
		  });

		  if(arguments[1])
				for(var arg in arguments[1])
					this[arg] = arguments[1][arg];

			this.computeSize();
			this.buildCalendar();

			this.showWarning('Loading...');
			setTimeout(function() {CAL.query_types(CAL_BODY.typesDownloaded)}, 0);
		},

		/// Event Callbacks
		typesDownloaded: function(data){
			CAL.get_settings(data);
			CAL_BODY.getAllEvents();
			CAL_BODY.onTypesLoaded();
		},

		mouseScrolled: function(e){
		  e = e ? e : window.event;
		  var raw = e.detail ? -1 * e.detail : e.wheelDelta;
		  if(raw == 0)
			return false;
		CAL_BODY.cal_box.style.top = (raw > 0 ? 1 : -1) * 0.8 * CAL_BODY.cell_height + parseInt(CAL_BODY.cal_box.style.top) + 'px';
		  CAL_BODY.topChanged();
		},

		calKeyDown: function(e){
			e = e ? e :window.event;
			switch(e.keyCode)
			{
			case 27:
				CAL_BODY.clear();
				break;
			case 33: // page up
				CAL_BODY.navigateDate(-1);
				break;
			case 34: // page down
				CAL_BODY.navigateDate(1);
				break;
			case 35: // end
				CAL_BODY.gotoToday();
				break;
			case 36: // home
				CAL_BODY.gotoToday();
				break;
			case 37: // left
				CAL_BODY.navigateDate(-1);
				break;
			case 38: // arrow up
				CAL_BODY.cal_box.style.top = CAL_BODY.cell_height / 3 + parseInt(CAL_BODY.cal_box.style.top) + 'px';
				CAL_BODY.topChanged();
				break;
			case 39: // right
				CAL_BODY.navigateDate(1);
				break;
			case 40: // arrow down
				CAL_BODY.cal_box.style.top = CAL_BODY.cell_height / -3 + parseInt(CAL_BODY.cal_box.style.top) + 'px';
				CAL_BODY.topChanged();
				break;
			}
		},

		calPressed: function(e){
			b_cal_released = false;
			e = e ? e : window.event;
			var src = (e && e.target) ? e.target : window.event.srcElement;
		  if(CAL_BODY.isInExpandBox(e))
			return false;
		  document.onmousemove = CAL_BODY.calMoving;
			document.onmouseup = CAL_BODY.calRelease;
		CAL_BODY.pressY = e.pageY ? e.pageY : e.clientY;
		  CAL_BODY.originY = parseInt(CAL_BODY.cal_box.style.top);
		},

		calMoving: function(e){
			if(b_cal_released) return;
			e = e ? e : window.event;
		  if(e.button == 1 || e.which == 1) {
			document.onmousemove = null;
			setTimeout("document.onmousemove = CAL_BODY.calMoving;", 24);
			var y = e.pageY	? e.pageY : e.clientY;
			CAL_BODY.cal_box.style.top = y - CAL_BODY.pressY + CAL_BODY.originY + 'px';
		  }
		  else if(e.button == 0 || e.which == 0)
			CAL_BODY.calRelease(e);
		},

		calRelease: function(e){
			b_cal_released = true;
			document.onmousemove = null;
			document.onmouseup = null;
			CAL_BODY.topChanged();
		},

		calResized: function(){
			if(new Date() - cal_rez_last < 200 || document.body.offsetWidth < 400) return;
			clearInterval(cal_rez_timer);
			cal_rez_timer = null;
			cal_rez_last = 0;

			CAL_BODY.clear();
			CAL_BODY.computeSize();

			var divs = document.querySelectorAll("div.div_cell_background");
			for (var i = 0; i < divs.length; i++) {
				divs[i].style.height = CAL_BODY.cell_height + 'px';
				divs[i].style.width = CAL_BODY.cell_width + 'px';
			}
			for(var i = 0, nLen = 3 * CAL_BODY.cal_lines; i < nLen; i++)
				CAL_BODY.buildWeekTable(CAL_BODY.first_day.plusDays(7 * i));
			CAL_BODY.buildWatermark();
			CAL_BODY.cal_box.style.top = -1 * CAL_BODY.cal_height + 'px';
		},

		weekEventsDownloaded: function(data){
			if(typeof(data) != 'string')
				var week_date = data;
			else
				try{
					var week_date = CAL.read_week_events(eval(data));
				}catch(err){CAL_BODY.showWarning('Oops! Catch error data.'); return;}
		  CAL_BODY.buildWeekTable(week_date);
		},

		weekEventsDownloadedFirstOnLoad: function(data){
			CAL_BODY.weekEventsDownloaded(data);
			CAL_BODY.hideWarning();
			setTimeout(CAL_BODY.getPageEvents_f(CAL_BODY.cal_lines, CAL_BODY.weekEventsDownloadedTop), 500);
			CAL_BODY.onFirstLoaded();
		},

		weekEventsDownloadedLast: function(data){
			CAL_BODY.weekEventsDownloaded(data);
			CAL_BODY.hideWarning();
		},

		weekEventsDownloadedMiddle: function(data){
			CAL_BODY.weekEventsDownloaded(data);
			CAL_BODY.hideWarning();
			CAL_BODY.getPageEvents(CAL_BODY.cal_lines, CAL_BODY.weekEventsDownloadedTop);
		},

		weekEventsDownloadedTop: function(data){
			CAL_BODY.weekEventsDownloaded(data);
			CAL_BODY.getPageEvents(0, CAL_BODY.weekEventsDownloadedBottom);
		},

		weekEventsDownloadedBottom: function(data){
			CAL_BODY.weekEventsDownloaded(data);
			CAL_BODY.getPageEvents(2 * CAL_BODY.cal_lines);
		},

		/// Operations
		getAllEvents: function(){
			if(b_first_load) {
				b_first_load = false;
				this.getPageEvents(this.cal_lines, this.weekEventsDownloadedFirstOnLoad);
			}
			else {
				CAL_BODY.showWarning('Loading...');
				this.getPageEvents(this.cal_lines, this.weekEventsDownloadedMiddle);
			}
	  },

		getPageEvents: function(base){
			var last_call = arguments[1] ? arguments[1] : CAL_BODY.weekEventsDownloaded;
			this.getWeeksEvents(base, this.cal_lines - 1, CAL_BODY.weekEventsDownloaded);
			this.getWeeksEvents(base + this.cal_lines - 1, 1, last_call);
		},

		getPageEvents_f: function(base){
			var last_call = arguments[1] ? arguments[1] : CAL_BODY.weekEventsDownloaded;
			return function(){
				CAL_BODY.getWeeksEvents(base, CAL_BODY.cal_lines - 1, CAL_BODY.weekEventsDownloaded);
				CAL_BODY.getWeeksEvents(base + CAL_BODY.cal_lines - 1, 1, last_call);
			};
		},

		buildCalendar: function(){
			var str_cal = '<table id=calendar cellpadding=0 cellspacing=0 class=cal_table><tbody>';
			var str_cal_week = '';
			for(var i = 0, nLen = this.cal_lines * 3; i < nLen; i++) {
			  var week_date = this.first_day.plusDays(i * 7);
			  str_cal += this.buildCalRow(week_date);
			  str_cal_week += '<div id=' + week_date.toDateUrlString() +  ' class=div_week><table id=calendar_week cellpadding=0 cellspacing=0><tbody>';
			  str_cal_week += this.buildWeekHead(week_date);
			  str_cal_week += '<tr style=height:' + (this.cell_height - this.cell_date_height + 2) + 'px;><td>&nbsp;</td></tr></tbody></table></div>';
			}
			str_cal += '</tbody></table>';
			document.getElementById('div_calendar').innerHTML = str_cal;
			this.cal_week.innerHTML = str_cal_week;
			this.cal_grid = document.getElementById('calendar');
			this.cal_week = document.getElementById('div_calendar_week');
			this.cal_box.style.top = -1 * this.cal_height + 'px';
			this.buildWatermark();
		},

		buildCalRow: function(week_date){
			var str = '<tr>';
			for(var i = 0; i < 7; i++){
			  var d = week_date.plusDays(i);
			  var year = d.getFullYear();
			  var month = d.getMonth();
			  var date = d.getDate();
			  str += '<td class=calendar_cell cell_date=' + d.toDateUrlString() + '>';
				str += '<div class=div_cell_background style=height:' + this.cell_height + 'px;width:' + this.cell_width + 'px;';

			  var now = new Date();
			  if(year == now.getFullYear() && month == now.getMonth() && date == now.getDate())
				str += 'background-color:' + today_color;
			  else if((year != this.cur_year || month != this.cur_month) && this.cal_lines > 3)
				str += 'background-color:#E1E1E1';
			  str += '>&nbsp;</div></td>';
			}
			str += '</tr>';
			return str;
		},

		buildWeekTable: function(week_date)	{
			if(week_date == null) return;
			var str = '<table cellpadding=0 cellspacing=0 class=cal_table><tbody>';
			str += this.buildWeekHead(week_date)
			str += this.buildWeekBody(week_date);
			str += '</tbody></table>';
			var dv_week = document.getElementById(week_date.toDateUrlString());
			if(dv_week)
				dv_week.innerHTML = str;
		},

		buildWeekHead: function(week_date){
			var str = '<tr>';
			for(var i = 0; i < 7; i++) {
			  var d = week_date.plusDays(i);
			  var month = d.getMonth();
			  var day = d.getDate();
			  str += '<td class=week_head cell_date=' + d.toDateUrlString() + '><div style=width:' + (this.cell_width + 2) + 'px; >';

			  var newD = d.plusDays(1);
			  if(day == 1 || newD.getMonth() != month || this.cal_lines <= 3)	//first/last day of month
				str += this.month_abbrv[month] + ' ' + day;
			  else
				str += day;

			  str += '</div></td>';
			}
			str += '</tr>';
			return str;
		},

		buildWeekBody: function(week_date) {
			var str = '';
			var curs = [0,0,0,0,0,0,0];
			var week_events = CAL.get_week_events(week_date);

			for(var i = 0, nLen = this.cell_hold - 1; i < nLen; i++)	{
			  var str_row = '<tr>';
			  var col = 0;
			  var b_event_found = false;
			  while(col < 7) {
				curs[col] = CAL.moveCursor(week_events[col].events, curs[col]);
				if(curs[col] >= week_events[col].events.length) {
				  str_row += '<td class=event_cell>&nbsp;</td>';
				  col++;
				}
				else {
					b_event_found = true;
				  var evt = week_events[col].events[curs[col]];
				  str_row += this.buildEventCell(evt, 'event_cell', -1);
				  this.increaseShowed(week_events, col, evt.span);
				  curs[col]++;
				  col = col + evt.span;
				}
			  }
			  str_row += '</tr>';
			  if(b_event_found == true)
				str += str_row;
			  else
				break;
			}

			//last line
			var last_height = this.cell_height + 2 - i * this.event_height - this.cell_date_height;
			var col = 0;
			str += '<tr>';
			while(col < 7) {
			  var w_e = week_events[col];
			  curs[col] = CAL.moveCursor(w_e.events, curs[col]);
			  if(curs[col] < w_e.events.length) {
				var evt = week_events[col].events[curs[col]];
				if(this.hasMore(week_events, col, evt.span))
				  str += this.buildCellMore(week_events, col, evt.span, week_date, last_height);
				else {
				  str += this.buildEventCell(evt, 'last_cell', last_height);
				  this.increaseShowed(week_events, col, evt.span);
				}
				col = col + evt.span;
			  }
			  else if(w_e.num + w_e.co_num - w_e.showed - w_e.co_showed > 0) {
				str += this.buildCellMore(week_events, col, 1, week_date, last_height);
				col++;
			  }
			  else {
				str += '<td class=last_cell style=height:' + last_height + 'px;>&nbsp;</td>';
				col++;
			  }
			}
			str += '</tr>';
			return str;
		},

		buildEventCell: function(evt, className, height){
			var str = '<td class=' + className;
			if(height > 0)
			  str += ' style=height:' + height + 'px; ';

			if(evt.multidays == 1)
			  str  += '>' + this.buildEvent(evt, false, false);
			else
			  str += ' colSpan=' + evt.span + '>' + this.buildEvent(evt, evt.completed > 0, !evt.end);

			str += '</td>';
			  return str;
		},

		buildCellMore: function(week_events, col, span, week_date, last_height)	{
			var str = ''
			for(var j = 0; j < span && col + j < 7; j++) {
			  var d = week_date.plusDays(col + j);
			  str += '<td class=cell_more style=height:' + last_height + 'px>';
			  var wi = week_events[col+j];
			  var more = wi.num + wi.co_num - wi.showed - wi.co_showed;
			  if(more > 0)
				str += '<a href=# onclick=CAL_BODY.expandCell(' + d.getFullYear() + ',' + d.getMonth() + ',' + d.getDate() +  ',event)>+' + more + ' more</a>';
			  str += '</td>';
			}
			return str;
		},

		increaseShowed: function(week_events, col, span) {
			week_events[col].showed++;
		  for(var i = 1; i < span && col + i < 7; i++)
			week_events[col + i].co_showed++;
		},

		hasMore: function(week_events, col, span)	{
		  for(var j = 0; j < span && col + j < 7; j++) {
			var wi = week_events[col + j];
			if(wi.num + wi.co_num - wi.showed - wi.co_showed > 1)
			  return true;
		  }
		  return false;
		},

		buildEvent: function(evt, leftArrow, rightArrow) {
		  var str = '<div event_type="' + evt.type + '" sub_type="' + evt.subtype;
			if (evt.type !== "ERROR")
				str += '" onmouseover="CAL_BODY.tooltip.show(this, event);" onmouseout="CAL_BODY.tooltip.hide(this);" ';
			str += ' tt_head="' + evt.title.replace(/"/g, '&quot;') + '"';

			var desc = '';
			if(evt.desc && evt.desc != "")
				desc += evt.desc.replace(/\\r\\n/g, '<br/>').replace(/"/g, '&quot;') + '<br/>';
			if(evt.multidays > 1)	{
				if(evt.start_date != '1700-01-01')
				desc += 'Start Date: ' + evt.start_date + '<br/>';
				if(evt.end_date != '3000-01-01')
					desc += 'End Date: ' + evt.end_date + '<br/>';
			}

			str += ' tt_body="' + desc + '" ';
			if(evt.id)
				str += ' tt_id="' + evt.id + '" ';
			var cont_style = "";
			if(CAL.event_types.hasOwnProperty(evt.type))
				cont_style = CAL.event_types[evt.type].Inverted == 'true' ? 'event_container_invert' : 'event_container';
			var width = evt.span * (this.cell_width + 2) - 8;
			str += ' class=' + cont_style + ' style=width:' + width + 'px;>'
			str += this.buildEventBody(evt, leftArrow, rightArrow);
			str += '</div>';
			return str;
		},

		buildEventBody: function(evt, leftArrow, rightArrow) {
		  var color = evt.type === 'ERROR' ? 'red' : CAL.event_types[evt.type].Color;
		  var extraStyle = evt.type === 'ERROR' ? ";font-size:1.2em;font-weight:bold;" : "";
		  var strDot = (evt.subtype ? ('<div class=div_dot ' +
				'style=color:' + CAL.event_types[evt.type].SubTypes[evt.subtype].Color + ';' +
				(leftArrow ? 'left:8px;' : '') + '>&#x25CF;</div>') : '');
			var txtDis = (leftArrow ? '&nbsp;&nbsp;' : '') + (evt.subtype ? '&nbsp;&nbsp;' : '') + '&nbsp;';
			var str = '';
			if(evt.type === "ERROR" || CAL.event_types[evt.type].Inverted != 'true'){
				str += '<div class=event_bar>';
				str += '<div class=corner style=background-color:' + color + ';></div>';
				str += '<div class=event_type style=background-color:' + color + extraStyle + ';>';
				str += (leftArrow ? '<div class="event_left_arrow event_arrow">&nbsp;</div>' : '');
				str += strDot + txtDis + evt.title;
				str += '</div><div class=corner style=background-color:' + color + ';></div>';
				str += (rightArrow ? '<div class="event_right_arrow event_arrow">&nbsp;</div>' : '');
				str += '</div>';
				}
			else{
				str += '<div class=event_bar_invert><div class=back_invert>';
				str += '<div class=corner_invert style=background-color:' + color + '; ></div>';
				str += '<div class=event_type_invert style=background-color:' + color + '; >&nbsp;</div>';
				str += '<div class=corner_invert style=background-color:' + color + '; ></div>';
				str += (rightArrow ? '<div class="event_right_arrow event_arrow event_right_arrow_invert">&nbsp;</div>' : '');
				str += '</div><div class=event_name_invert>';
				str += (leftArrow ? '<div class="event_left_arrow event_arrow event_left_arrow_invert">&nbsp;</div>' : '');
				str += strDot + txtDis + evt.title;
				str += '</div></div>';
			}
			return str;
		},

		topChanged: function()	{
			var top = parseInt(this.cal_box.style.top), whole_cell = this.cell_height + 2;
			var half_cell = whole_cell / 2, direction = 0;
			if(top > -1 * this.cal_height + half_cell) //See Before, direction < 0
			  direction = -1 * Math.floor((top + half_cell + this.cal_height)/ whole_cell);
			else if(top < -1 * this.cal_height - half_cell) //See Future, direction > 0
			  direction = -1 * Math.floor((top + this.cal_height + half_cell) / whole_cell);
			else
				return;
			this.moveCalendar(direction);
			CAL_BODY.onDateChanged();
		},

		getWeeksEvents: function(base, weeks){
			var last_call = arguments[2] ? arguments[2] : CAL_BODY.weekEventsDownloaded;
			for(var i = 0, nLen = weeks - 1; i < nLen; i++)
				CAL.query_events(this.first_day.plusDays((i + base) * 7), CAL_BODY.weekEventsDownloaded);
			CAL.query_events(this.first_day.plusDays((i + base) * 7), last_call);
		},

		moveCalendar: function(direction)	{
			if(direction == 0) return;

			CAL_BODY.showWarning('Loading...');
			this.changeFirstDay(7 * direction);
			if(direction < 0) {//before
			  for(var i = 0, nLen = -1 * direction; i < nLen; i++) {
				var week_date = this.first_day.plusDays(i * 7);

				this.cal_grid.deleteRow(-1);
				this.cal_week.removeChild(this.cal_week.lastChild);

				this.createCalRow(week_date, i);
				this.cal_week.insertBefore(this.createCalWeekRow(week_date), this.cal_week.childNodes[i]);
			  }
			  for(var i = -1 * direction, nLen = 3 * this.cal_lines; i < nLen; i++)
				for(var j = 0; j < 7; j++)
				  this.paintCell(this.cal_grid.rows[i].cells[j]);
			  this.getWeeksEvents(0, -1 * direction, CAL_BODY.weekEventsDownloadedLast);
			}
			else {//future
			  for(var i = 0; i < direction; i++) {
				var week_date = this.first_day.plusDays((3 * this.cal_lines - direction + i) * 7);

				this.cal_grid.deleteRow(0);
				this.cal_week.removeChild(this.cal_week.firstChild);

				this.createCalRow(week_date, -1);
				this.cal_week.appendChild(this.createCalWeekRow(week_date));
			  }
			  for(var i = 0, nLen = 3 * this.cal_lines - direction; i < nLen; i++)
				for(var j = 0; j < 7; j++)
				  this.paintCell(this.cal_grid.rows[i].cells[j]);
			  this.getWeeksEvents(this.cal_lines * 3 - direction, direction, CAL_BODY.weekEventsDownloadedLast);
			}
			this.buildWatermark();
			if(this.cal_expand.style.display != 'none' && this.cal_expand.style.top != '')
				this.cal_expand.style.top = parseInt(this.cal_expand.style.top) - direction * (this.cell_height + 2) + 'px';
			this.cal_box.style.top = parseInt(this.cal_box.style.top) + direction * (this.cell_height + 2) + 'px';
		},

		createCalWeekRow: function(week_date)	{
			var cal_week_row = document.createElement('div');
			var strDate = week_date.toDateUrlString()
			cal_week_row.id = strDate;
			cal_week_row.className = 'div_week';
			var str = '<div week_date=';
			str += strDate;
			str +=  ' class=div_week><table id=calendar_week cellpadding=0 cellspacing=0><tbody>';
			str += this.buildWeekHead(week_date);
			str += '<tr style=height:' + (this.cell_height - this.cell_date_height + 2) + 'px;><td>&nbsp;</td></tr></tbody></table></div>';
			cal_week_row.innerHTML = str;
			return cal_week_row;
		},

		createCalRow: function(week_date, ins) {
			var cal_row = this.cal_grid.insertRow(ins);
		  for(var j = 0; j < 7; j++)
		  {
			var cell = cal_row.insertCell(-1);
			cell.className = 'calendar_cell';
			cell.setAttribute('cell_date', week_date.plusDays(j).toDateUrlString());
			cell.innerHTML = '<div class=div_cell_background style=height:' + this.cell_height + 'px;width:' + this.cell_width + 'px;>&nbsp;</div>';
			this.paintCell(cell);
		  }
		},

		paintCell: function(cell)	{
		  var div_bg = cell.childNodes[0];
		  var cell_date = cell.getAttribute('cell_date').toDate();
			var y = cell_date.getFullYear();
			var m = cell_date.getMonth();
			var d = cell_date.getDate();

			var now = new Date();
			if(y == now.getFullYear() && m == now.getMonth() && d == now.getDate())
			  div_bg.style.backgroundColor = today_color;
			else if((y != this.cur_year || m != this.cur_month) && this.cal_lines > 3)
			  div_bg.style.backgroundColor = '#E1E1E1';
			else
			  div_bg.style.backgroundColor = '';
		},

		computeSize : function ()	{
		  var div_cal_head = document.getElementById('div_cal_head');
		  //compute the width of cells
			this.cal_width = document.body.offsetWidth - 150 - 6;
			if(this.cal_width < this.cell_min_width * 7)
			  this.cal_width = this.cell_min_width * 7;
			this.cal_frame.style.width = this.cal_width + 'px';
			div_cal_head.style.width = this.cal_width + 6 + 'px';
			this.cell_width = Math.floor(this.cal_width/7) - 2;
			this.cal_width = (this.cell_width + 2)* 7;
			//compute the height of cells
			this.cal_height = document.body.offsetHeight - 15 - 8;
			if(this.cal_height < 15 + 6 + this.cell_min_height * this.cal_lines)
			  this.cal_height = 15 + 6 + this.cell_min_height * this.cal_lines;
			this.cal_frame.style.height = this.cal_height + 'px';
			this.cell_height = Math.floor(this.cal_height / this.cal_lines);
			//make sure the this.cell_height should be an even
			this.cell_height = this.cell_height - this.cell_height % 2 - 2;
			this.cal_height = (this.cell_height + 2) * this.cal_lines;
			this.cell_hold = Math.floor((this.cell_height - this.cell_date_height)/ this.event_height);
			var div_cal_events = document.getElementById('div_cal_events');
			var cal_events_height = this.cal_height - 140;
			div_cal_events.style.height = (cal_events_height > 0 ? cal_events_height : '0') + 'px';
		},

		/// watermakr operation
		buildWatermark: function() {
		  var start = this.first_day.getDate() == 1 ? 0 : 1;
		  for(var i = 0; i < 4; i++)
		  {
			if(this.cal_lines <= 3){this.watermarks[i].style.display='none';continue;}

			  var d = new Date(this.first_day.getFullYear(), this.first_day.getMonth() + i + start, 1);
			  var m = d.getMonth();
			  d.setDate(d.getDate() - d.getDay());
			  var row = Math.round((d - this.first_day) / (86400000 * 7));
			  var month_width = this.cal_width * 14 / 15;
			  this.watermarks[i].style.fontSize = this.cal_width / 7 + 'px';
			  this.watermarks[i].style.width = month_width + 'px';
			  this.watermarks[i].style.top = (row + 2) * (this.cell_height + 2)+ 'px';
			  this.watermarks[i].style.left = (this.cal_width - month_width) / 2 + 'px';
			  this.watermarks[i].innerHTML = this.month_name[m].toUpperCase();
			  this.watermarks[i].style.display = 'block';
		  }
		},

		gotoToday: function(){
			var now = new Date();
			if(this.cal_lines == 6)
				this.setFirstDay(now.getFullYear(), now.getMonth());
			else
				this.first_day = now.plusDays(-1 * this.halfCalWeeks() - now.getDay());
			this.onDateChanged();
			this.buildCalendar();
			this.getAllEvents();
		},

		refreshCache: function(){
			CAL.clear_cache();
			//window.location.href = window.location.href;
			//location.reload(true);
		},

		navigateDate: function(plus){
		  if(this.cal_lines == 6)
			  this.setFirstDay(this.cur_year, this.cur_month + plus);
		  else
			this.changeFirstDay(plus * this.cal_lines * 7);
		  this.onDateChanged();
		  this.buildCalendar();
		  this.getAllEvents();
		},

		setFirstDay : function (cur) {
			var d = (typeof cur == 'number') ? new Date(cur, arguments[1]) : new Date(cur.getFullYear(), cur.getMonth(), 1);
		  if(d.getDay() == 0)
			d.setDate(d.getDate() - 7);
		  d.setDate(d.getDate() - d.getDay() - this.cal_lines * 7);
		  this.first_day = d;
		  this.updateCurYearMonth();
		},

		changeFirstDay : function (num) {
			this.first_day.setDate(this.first_day.getDate() + num);
		  this.updateCurYearMonth();
		},

		updateCurYearMonth: function() {
			var centerDay = this.first_day.plusDays(this.halfCalWeeks() + 3);
		  this.cur_month = centerDay.getMonth();
		  this.cur_year = centerDay.getFullYear();
		},

		expandCell: function(y, m ,d, e) {
		  var date = new Date(y, m, d);
		  var row = document.getElementById(date.firstDateOfWeek().toDateUrlString());
		  var cell = document.querySelectorAll('td.calendar_cell[cell_date="' + date.toDateUrlString() + '"]')[0];
		  var date_events = CAL.get_date_events(date);
		  var strClose = '<div title=close class=div_expand_close>&times;</div>';
		  var strHTML = '<div class="week_head div_expand_head">' + this.month_abbrv[m] + ' ' + d + strClose + '</div>';
		  this.cal_expand.style.display = 'block';

		  var color = cell.childNodes[0].style.backgroundColor;
		  this.cal_expand.style.backgroundColor = (color != '' ? color : 'white');

		  var nWidth = this.cell_width * 2.2 + 15;
		  this.cal_expand.style.width = nWidth + 'px';
		  this.cal_expand.style.top = row.offsetTop + 'px';
		  if(cell.offsetLeft + 2 + nWidth < CAL_BODY.cal_box.offsetWidth) {
			this.cal_expand.style.right = 'auto';
			this.cal_expand.style.left = cell.offsetLeft + 'px';
		  }
		  else {
			this.cal_expand.style.right = '6px';
			this.cal_expand.style.left = 'auto';
		  }

		  var nHeightNum = (this.cal_lines == 1) ? this.cell_hold : (this.cell_hold + MAX_EXPAND);
		  var bScroll = false;
		  if(date_events.length > nHeightNum) {
			bScroll = true;
			strHTML += '<div style=height:' + (this.event_height * nHeightNum + this.cell_date_height) + 'px;' +
				'overflow-x:hidden;overflow-y:scroll;position:relative;>';
		  }
		  else {
			bScroll = false;
			strHTML += '<div>';
		  }
		  for(var i = 0, nLen = date_events.length; i < nLen; i++) {
			var ev = date_events[i];
			if(CAL.is_hidden(ev)) continue;
			if(ev.multidays == 1) {
			  var leftArrow = false;
			  var rightArrow = false;
			}
			else {
				var e_left = -1 * ev.start_date.toDate().minusDate(date) + ev.completed;
			  var leftArrow = (e_left !== 0);
			  var rightArrow = ((e_left + 1) < ev.multidays);
			}
			strHTML += '<div style=height:' + this.event_height + 'px;width:100%;>' + this.buildEvent(ev, leftArrow, rightArrow) + '</div>';
		  }
		  strHTML += '<div style=height:3px></div></div>';

		  this.cal_expand.innerHTML = strHTML;
		  for(var i = 0, nLen = date_events.length; i < nLen; i++) {
			var node = this.cal_expand.childNodes[1].childNodes[i].childNodes[0];
			if(bScroll == true)
				node.style.width = nWidth - 21 + 'px';
			else
				node.style.width = nWidth - 6 + 'px';
			node.style.marginLeft = '3px';
		  }
		  return this.cancelEvent(e);
		},

		halfCalWeeks: function(){
			return 7 * (this.cal_lines + Math.floor( (this.cal_lines - 1) / 2));
		},

		isInExpandBox: function(e){
			if(this.cal_expand.style.display == 'none') return false;
			var obj = (e && e.target) ? e.target : window.event.srcElement;
			do {
				if(obj === this.cal_frame)
					return false;
				if(obj === this.cal_expand)
					return true;
			}while (obj = obj.offsetParent);
			return false;
		},

		offset: function(){
			var off = this.getRelPos(this.cal_frame);
			off.w = this.cal_frame.offsetWidth;
			off.h = this.cal_frame.offsetHeight;
			return off;
		},

		clear: function(){
			CAL_BODY.cal_expand.style.display = 'none';
		CAL_BODY.tooltip.hide();
		},

		changeSearchText: function(v){
			v = v.replace(/^\s+|\s+$/g, '');
			if(v != CAL.search_text) {
				CAL.search_text = v.toLowerCase();
				CAL_BODY.getAllEvents();
			}
		},

		updateColor: function(type){
			const selector = 'div.event_container[event_type="' + type + '"] div.corner, ' +
				'div.event_container[event_type="' + type + '"] div.event_type, ' +
				'div.event_container_invert[event_type="' + type + '"] div.corner_invert, ' +
				'div.event_container_invert[event_type="' + type + '"] div.event_type_invert';
			var divs = this.cal_frame.querySelectorAll(selector);
			for (var i = 0; i < divs.length; i++) {
				divs[i].style['background-color'] = CAL.event_types[type].Color;
			}
		},

		updateSubColor: function(type, subtype){
			const selector = 'div.event_container[event_type="' + type + '"][sub_type="' + subtype + '"] div.div_dot, ' +
				'div.event_container_invert[event_type="' + type + '"][sub_type="' + subtype + '"] div.div_dot';
			var divs = this.cal_frame.querySelectorAll(selector);
			for (var i = 0; i < divs.length; i++) {
				divs[i].style.color = CAL.event_types[type].SubTypes[subtype].Color;
			}
		},

		showWarning: function(str){
			CAL_BODY.cal_warning.innerHTML = str;
			CAL_BODY.cal_warning.style.display = 'block';
		},

		hideWarning: function(){
			CAL_BODY.cal_warning.style.display = 'none';
		},

		getRelPos: function(elm){
			var curtop = 0, curleft = 0, pos = {}, p = arguments[1];
			try{
			  var obj = elm;
			  do {
				if(obj === p) break;
				curleft += obj.offsetLeft;
				curtop += obj.offsetTop;
			  } while (obj = obj.offsetParent);

			  pos.x = curleft;
			  pos.y = curtop;
			  return pos;
			}
		  catch(e){return {x: 200, y: 200};}
		},

		getMousePos: function (e){
		  return (e.pageX)
			? {x:e.pageX, y:e.pageY}
			: {x:e.clientX + document.body.scrollLeft - document.body.clientLeft,
				 y:e.clientY + document.body.scrollTop - document.body.clientTop};
		},

		hookEvent: function(element, eventName, callback){
		  if(typeof(element) == 'string')
			element = document.getElementById(element);
		  if(element == null) return;
		  if(element.addEventListener){
			  if(eventName == 'mousewheel')
				element.addEventListener('DOMMouseScroll', callback, false);
			  element.addEventListener(eventName, callback, false);
		  }
		  else if(element.attachEvent)
			element.attachEvent('on' + eventName, callback);
		},

		unhookEvent: function(element, eventName, callback){
		  if(typeof(element) == 'string')
			element = document.getElementById(element);
		  if(element == null) return;
		  if(element.removeEventListener){
			  if(eventName == 'mousewheel')
				element.removeEventListener('DOMMouseScroll', callback, false);
			  element.removeEventListener(eventName, callback, false);
		  }
		  else if(element.detachEvent)
			element.detachEvent('on' + eventName, callback);
		},

		cancelEvent: function(e){
		  e = e ? e : window.event;
		  if(e.stopPropagation)
			e.stopPropagation();
		  if(e.preventDefault)
			e.preventDefault();
		  e.cancelBubble = true;
		  e.cancel = true;
		  e.returnValue = false;
		  return false;
		},

		tooltip: function(){
			var speed = 20, timer = 20, end_alpha = 90, start_alpha = 0, alpha = 0, arr_down = [],
				last_elm = null, last_x = null, in_expand = false, show_timer = null, fade_timer = null, b_init = false;

			var init = function(){
				if(!b_init) {
					for(var tp in CAL.event_types)
						if(CAL.event_types[tp].Tooltip && CAL.event_types[tp].Tooltip == 'true')
								arr_down[tp] = [];
					b_init = true;
				}
			}

			var desc = function(elm){
				var type = elm.getAttribute('event_type');
				var subt = elm.getAttribute('sub_type');
				if(type === "ERROR")
					return "";
				return (subt != '' ?	'Sub Type: <font  color=' + CAL.event_types[type].SubTypes[subt].Color + ' size=4>&#x25CF;</font>&nbsp;' + subt + '<br/>' : '') +
								last_elm.getAttribute('tt_body');
			};

			var tt_downloaded = function(data){
			if(last_elm) {
				var last_type = last_elm.getAttribute('event_type');
				arr_down[last_type][last_elm.getAttribute('tt_id')] = data;
				document.getElementById('tooltip_body').innerHTML = desc(last_elm) + data;
			}
			last_elm = null;
			};

			return{
				show: function(elm, e){
					if(!b_cal_released) return false;
					document.getElementById('tooltip_head').innerHTML = elm.getAttribute('tt_head');
					var type = elm.getAttribute('event_type');
					var type_ob = CAL.event_types[type];
					var tootip_body = document.getElementById('tooltip_body');
					if(type_ob.Tooltip && type_ob.Tooltip == 'true') {
						var id = elm.getAttribute('tt_id');
						if(!arr_down[type]) init();
						if(!arr_down[type][id]) {
							arr_down[type][id] = '_need_download_';
							last_elm = elm;
							tootip_body.innerHTML = desc(last_elm) + '<div bgColor=#888888><img src=/Res?name=tooltiploading.gif></div>';
						}
						else if(arr_down[type][id] == '_need_download_')
							return;
						else {
							last_elm = elm;
							tootip_body.innerHTML = desc(last_elm) + arr_down[type][id];
						}
					}
					else {
						last_elm = elm;
						tootip_body.innerHTML = desc(last_elm);
					}
					last_x = e.pageX ? e.pageX : e.clientX;
					in_expand = CAL_BODY.isInExpandBox(e);
					if(show_timer) clearTimeout(show_timer);
					if(fade_timer) clearInterval(fade_timer);
					show_timer = setTimeout(this.delayedShow, 300);
				},

				delayedShow: function(){
					var baseUrl = window.location.pathname; // for CalendarJS
					if(!b_cal_released || last_elm == null) return false;
					var last_elm_pos = CAL_BODY.getRelPos(last_elm, CAL_BODY.cal_box);
					last_elm.style.cursor = 'pointer';
					var tooltip = document.getElementById('tooltip');
					tooltip.style.display = 'block';
					tooltip.style.opacity = start_alpha * .01;
					tooltip.style.filter = 'alpha(opacity=' + start_alpha + ')';
					tooltip.style.zoom = 1;
					alpha = start_alpha;
					if(last_x + 200 > document.body.offsetWidth) {
						tooltip.style.left = '';
						tooltip.style.right = '10px';
					} else {
						tooltip.style.left = last_x - CAL_BODY.offset().x + 10 + 'px';
						tooltip.style.right = '';
					}
					var scroll = in_expand ? CAL_BODY.cal_expand.childNodes[1].scrollTop : 0;
					tooltip.style.top = last_elm_pos.y - scroll + CAL_BODY.event_height + 'px';
					fade_timer = setInterval(CAL_BODY.tooltip.fade, 50);
					var last_type = last_elm.getAttribute('event_type');
					var last_type_ob = CAL.event_types[last_type];
					if(last_type_ob.Tooltip && last_type_ob.Tooltip == 'true' && arr_down[last_type][last_elm.getAttribute('tt_id')] == '_need_download_') {
						var xhr = new XMLHttpRequest();
						xhr.open('GET', '/CalTooltip?page=' + page.substring(1) + '&event_type=' + last_type+	'&id=' + last_elm.getAttribute('tt_id') + session);
						xhr.onload = function() {
							if (xhr.status >= 200 && xhr.status < 300) {
								tt_downloaded(xhr.responseText);
							}
						};
						xhr.send();
					}
				},

				fade: function(){
					var a = alpha;
					if(a != end_alpha) {
						var i = speed;
						if(end_alpha - a < speed)
							i = end_alpha - a;
						alpha = a + i;
						tooltip.style.opacity = alpha * .01;
						tooltip.style.filter = 'alpha(opacity=' + alpha + ')';
						tooltip.style.zoom = 1;
					}
					else
						if(show_timer) clearInterval(show_timer);
				},

				hide: function(elm){
					if(!b_cal_released) return false;
					if(fade_timer) clearInterval(fade_timer);
					if(show_timer) clearTimeout(show_timer);
					if(last_elm) {
						var last_type = last_elm.getAttribute('event_type');
						if(arr_down[last_type] && arr_down[last_type][last_elm.getAttribute('tt_id')] == '_need_download_')
							arr_down[last_type][last_elm.getAttribute('tt_id')] = undefined;
						var tooltip = document.getElementById('tooltip');
						if (tooltip) {
							tooltip.style.display = 'none';
						}
					}
					if(elm)
						elm.style.cursor = '';
				}
			};
		}()
	}
}(CAL);

// Left Control
var CAL_LEFT = (function(CAL, CAL_BODY){

	var div_left = 'div_left';
	var sel_colors = ['Navy','DarkGreen','SlateBlue','DarkSlateGray', 'Indigo','CadetBlue','Teal','CornflowerBlue','Olive',	'SeaGreen','Purple','Sienna','Maroon'];
	var sel_sub_colors = ['Red', 'Orange', 'Yellow', 'Green', 'Blue', 'Indigo', 'Violet', 'Black', 'Pink', 'Aqua', 'Palegreen', 'Salmon', 'Magenta'];
	var gray_color = '#e1e1e1';

	var TypeSelected = '';
	var SubTypeSelected = '';

	var dragCtlBar = null;
	var origCtlTop = 0;
	var origCtlLeft = 0;
	var bCtlMoved = false;
	var bReleased = true;
	var bMoveInit = false;
	var pressX = 0;
	var pressY = 0;

	var leftCtl = 0;

	var barHeight = 24;
	var off_y = 0;

	return {
		initLeftDiv: function(div){
			div_left = div;
			document.getElementById(div_left).addEventListener('mousedown', this.leftDivPressed)
			CAL_BODY.hookEvent(document, 'keydown', this.leftKeyDown);
			CAL_BODY.hookEvent(document, 'keypress', this.leftKeyDown);
			CAL_BODY.hookEvent(document, 'click', this.clear);
		},

		buildTypes: function(){
			var copys = CAL.sort_event_types();
			var str = '<table cellspacing=0 cellpadding=0 id=leftCtl><tbody>';
			for(var i = 0, nLen = copys.length; i < nLen; i++)
				str += '<tr><td class=left_cell vAlign=top>' + CAL_LEFT.createEventControl(copys[i].type) + '</td></tr>';
			str += '</tbody></table>';
			str += "<div id='div_color_menu'></div>";
			document.getElementById(div_left).innerHTML = str;
			leftCtl = document.getElementById('leftCtl');
		  off_y = CAL_BODY.getRelPos(leftCtl).y;
		},

		leftKeyDown: function(e){
			e = e ? e : window.event;
			switch(e.keyCode)
			{
			case 27:
				CAL_LEFT.clear();
				break;
			}
		},

		createEventControl: function (type){
			var type_ob = CAL.event_types[type];
			var color = type_ob.Display === 'none' ? gray_color : type_ob.Color;
			var expand = type_ob.Expanded == 'true' ? '--' : '+';
			var str =  '<table class=left_cell_table cellpadding=0 cellspacing=0 type="' + type + '"><tbody><tr><td>'
			if(type_ob['SubTypes'] != undefined)
				str += this.createBar(expand, color, 18, 0, '', "CAL_LEFT.expandType('" + type + "',this);", '', type_ob.Inverted) +	'</td><td>' +
					this.createBar(type, color, 100, 2, '', "CAL_LEFT.switchEvent('" + type + "',this);", '', type_ob.Inverted);
			else
				str += this.createBar(type, color, 120, 0, '', "CAL_LEFT.switchEvent('" + type + "',this);", '', type_ob.Inverted);
			str += '</td><td>' + this.createBar('&gt;', color, 20, 2, '', "CAL_LEFT.popColorMenu('" + type + "',-1,0,this,event);", '', type_ob.Inverted) +
				'</td></tr></tbody></table>';

			if(expand == '--')
				str += this.buildSubType(type);
			return str;
		},

		createBar: function (txt, color, width, left, style, click, title, invert){
			var text_css = (color === gray_color) ? 'event_name_invert_none' : 'event_name_invert';
			var disabled_attr = (color !== gray_color) ? '' : ' gray="true" ';
			return (invert == 'false'
				? '<div class=event_ctlbar style=width:' + width + 'px;padding-left:' + left + 'px;' + style +
						' title="' + title + '" onclick="' + click + '" >' +
						'<div class=corner ' + disabled_attr + 'style=background-color:' + color + ';></div>' +
						'<div class=event_type ' + disabled_attr + 'style=background-color:' + color + ';>' + txt + '</div>' +
						'<div class=corner ' + disabled_attr + 'style=background-color:' + color + ';></div>' +
					'</div>'
				:	'<div class=event_ctlbar_invert style=width:' + width + 'px;padding-left:' + left + 'px;' + style +
						' onclick="' + click + '" title=' + title + '>' +
						'<div class=back_invert style=width:' + width + 'px;>' +
							'<div class=corner_invert ' + disabled_attr + 'style=background-color:' + color + '; ></div>' +
							'<div class=event_type_invert ' + disabled_attr + 'style=background-color:' + color + '; >&nbsp;</div>' +
							'<div class=corner_invert ' + disabled_attr + 'style=background-color:' + color + '; ></div>' +
						'</div>' +
						'<div class=' + text_css + '>' + txt + '</div>' +
					'</div>');
		},

		expandType: function (type, ob) {
			var left_cell = document.querySelector('td.left_cell table[type="' + type + '"]');
			if (left_cell.parentNode.childNodes.length > 1) {
				var nodes = left_cell.parentNode.querySelectorAll(subtypeTableSelector);
				var parent = nodes[0].parentNode;
				if (parent) {
					for (var i = 0; i < nodes.length; i++) {
						parent.removeChild(nodes[i]);
					}
				}
				CAL_LEFT.setExpand(left_cell, type, 'false');
			} else {
				var htmlDoc = document.createElement('div');
				htmlDoc.innerHTML = this.buildSubType(type);
				var tableNode = htmlDoc.querySelector('table');
				left_cell.parentNode.appendChild(tableNode);
				CAL_LEFT.setExpand(left_cell, type, 'true');
			}
			CAL.save_settings();
		},

		setExpand: function (left_cell, type, expand) {
			var nodes = left_cell.parentNode.childNodes[0].querySelectorAll(expandSelector);
			for (var i = 0; i < nodes.length; i++) {
				if (classTypes.indexOf(nodes[i].className) !== -1 && ['+', '--'].indexOf(nodes[i].innerHTML) !== -1) {
					nodes[i].innerHTML = expand === 'false' ? '+' : '--';
				}
			}
			CAL.event_types[type].Expanded = expand;
		},

		buildSubType: function(type){
			var evt_type = CAL.event_types[type];
			var sub_types = evt_type.SubTypes;
			var str = '<table class=subtype_table cellpadding=0 cellspacing=0><tbody>';
			for(sub in sub_types) {
				var sub_type = sub_types[sub];
				var bar_color = (sub_type.Display == 'none' ? gray_color : evt_type.Color);
				str += '<tr class=subtype_row><td style=padding-left:10px; sub_type="' + sub + '">' +
					this.createBar('<font color=' + sub_type.Color + '>&#x25CF;&nbsp;</font>' + sub, bar_color, 106, 3,
						'', "CAL_LEFT.switchSubEvent('" + type + "','" + sub + "',this);", '', evt_type.Inverted) +
					'</td><td sub_type="' + sub + '">' +
					this.createBar('>', bar_color, 20, 3,
						'', "CAL_LEFT.popColorMenu('" + type + "','" + sub + "',1,this,event);", 'select dot color', evt_type.Inverted) +
					'</td></tr>';
			}
			str += '</tbody></table>';
			return str;
		},

		popColorMenu: function(type, sub_type, level, ob, e){
			if((level == 0 && CAL.event_types[type].Display == 'none') || (level == 1 && CAL.event_types[type].SubTypes[sub_type].Display == 'none')) return;
			var ctl = ob.parentNode;
			var pos = CAL_BODY.getRelPos(ctl);
			var colorMenu = document.getElementById('div_color_menu');
			if (colorMenu) {
				colorMenu.innerHTML = this.buildColorMenu(level);
				colorMenu.style.left = pos.x + ctl.offsetWidth + 'px';
				colorMenu.style.top = pos.y + ctl.offsetHeight - 10 + 'px';
				colorMenu.style.display = 'block';
			}
			TypeSelected = type;
			SubTypeSelected = sub_type;
			return CAL_BODY.cancelEvent(e);
		},

		buildColorMenu: function(level){
			var colors = (level == 0 ? sel_colors : sel_sub_colors);
		  str = '<table style=width:120px; cellspacing=0 cellpadding=0><tbody>';
		  for(var i =0, nLen = colors.length; i < nLen; i++)
			str += '<tr><td style=background-color:' + colors[i] + '><div sel_color=' + i +
				' onclick=' + (level == 0 ? 'CAL_LEFT.selectColor' : 'CAL_LEFT.selectSubColor') +
				'(' + i + ');>' + colors[i] + '</div></td></tr>';

		  if(level == 0)
			str += '<tr><td><input type=button value="Invert Color" style=width:120px; onclick="CAL_LEFT.invertColor();"/></td></tr>';
		  str += '</tbody></table>';
		  return str;
		},

		selectColor: function(ind){
			CAL.event_types[TypeSelected].Color = sel_colors[ind];
			CAL.save_settings();
			const parent = document.querySelector('td.left_cell table[type="' + TypeSelected + '"]').parentNode;
			var nodes = parent.querySelectorAll('.corner:not([gray="true"]), .event_type:not([gray="true"]), .corner_invert:not([gray="true"]), .event_type_invert:not([gray="true"])');
			for (var i = 0; i < nodes.length; i++) {
				nodes[i].style['background-color'] = CAL.event_types[TypeSelected].Color;
			}
			setTimeout(function(t) { return function() {CAL_BODY.updateColor(t)};}(TypeSelected), 500);
			var colorMenu = document.getElementById('div_color_menu');
			if (colorMenu) {
				colorMenu.style.display = 'none';
			}
			TypeSelected = SubTypeSelected= '';
		},

		selectSubColor: function(ind){
			CAL.event_types[TypeSelected].SubTypes[SubTypeSelected].Color = sel_sub_colors[ind];
			CAL.save_settings();
			const el = document.querySelector('td.left_cell table[type="' + TypeSelected + '"]');
			const next = this.next(el, 'td.left_cell table[type="' + TypeSelected + '"]');
			var sub = next.querySelector('td[sub_type="' + SubTypeSelected + '"] font');
			if (sub) {
				sub.style.color = sel_sub_colors[ind]
			}
			setTimeout(function(t, s) { return function() {CAL_BODY.updateSubColor(t, s);}}(TypeSelected, SubTypeSelected), 500);
			var colorMenu = document.getElementById('div_color_menu');
			if (colorMenu) {
				colorMenu.style.display = 'none';
			}
			TypeSelected = SubTypeSelected= '';
		},

		next: function (el) {
			const nextEl = el.nextElementSibling;
			return nextEl ? nextEl : null
		},

		switchEvent: function(type, ob){
			if(bCtlMoved) return;

			var type_ob = CAL.event_types[type];
			var display = (type_ob.Display == 'none' ? 'block' : 'none');
			type_ob.Display = display;
			if(type_ob.SubTypes != undefined)
				for(subtype in type_ob.SubTypes)
					type_ob.SubTypes[subtype].Display = display;
			CAL.save_settings();

			var bar_color = (CAL.event_types[type].Display == 'none' ? gray_color : CAL.event_types[type].Color);
			var p = document.querySelector('td.left_cell table[type="' + type + '"]').parentNode;
			this.switchColor(p, bar_color);

			setTimeout(this.refreshCalendar, 500);
		},

		switchSubEvent: function(type, subtype, ob){
			if(bCtlMoved) return;

			var sub_type_ob = CAL.event_types[type].SubTypes[subtype];
			sub_type_ob.Display = (sub_type_ob.Display == 'none' ? 'block' : 'none');
			CAL.save_settings();

			const next = this.next(document.querySelector('td.left_cell table[type="' + type + '"]'));
			const bar_color = (sub_type_ob.Display == 'none' ? gray_color : CAL.event_types[type].Color);
			var p = next.querySelectorAll('td[sub_type="' + subtype + '"]');
			for (var i = 0; i < p.length; i++) {
				this.switchColor(p[i], bar_color);
			}
			setTimeout(this.refreshCalendar, 500);
		},

		switchColor: function(p, bar_color){
			var eventTypes = p.querySelectorAll('.corner, .event_type, .corner_invert, .event_type_invert');
			for (var i = 0; i < eventTypes.length; i++) {
				eventTypes[i].style['background-color'] = bar_color;
			}
			var text_orgin_css = (bar_color === gray_color) ? 'event_name_invert' : 'event_name_invert_none';
			var text_new_css = (bar_color === gray_color) ? 'event_name_invert_none' : 'event_name_invert';
			var elements = p.querySelectorAll('.' + text_orgin_css);
			for (var l = 0; l < elements.length; l++) {
				this.removeClass(elements[l], text_orgin_css);
				this.addClass(elements[l], text_new_css);
			}
		},

		addClass: function (el, clName) {
			if (el.className && el.className.split(' ').indexOf(clName) === -1) {
				el.className += ' ' + clName;
			} else {
				el.className = clName;
			}
		},

		removeClass: function (el, clName) {
			var classes = el.className.split(' ');
			var index = classes.indexOf(clName);
			if (index !== -1) {
				classes.splice(index, 1);
				el.className = classes.join(' ');
			}
		},

		invertColor: function(){
			var type_ob = CAL.event_types[TypeSelected];
			type_ob.Inverted = (type_ob.Inverted == 'false' ? 'true' : 'false');
			var parent = document.querySelector('td.left_cell table[type="' + TypeSelected + '"]').parentNode;
			if (parent) {
				parent.innerHTML = this.createEventControl(TypeSelected);
			}
			CAL.save_settings();
			var colorMenu = document.getElementById('div_color_menu');
			if (colorMenu) {
				colorMenu.style.display = 'none';
			}
			TypeSelected = SubTypeSelected = '';
			setTimeout(this.refreshCalendar, 500);
		},

		updateEventBar: function(left_cell, type){
			if(left_cell.childNodes.length == 1)
				left_cell.innerHTML = this.createEventControl(type);
			else
			{
				left_cell.innerHTML = this.createEventControl(type) + buildSubType(type);
				left_cell.childNodes[0].rows[0].cells[0].childNodes[0].childNodes[1].innerHTML = '--';//'&ndash;';
			}
		},

		//Drag and drop
		leftDivPressed: function(e){
			e = e ? e : window.event;
			bCtlMoved = false;
			bReleased = false;
			if(CAL_LEFT.findControlBar(e)) {
				bMoveInit = false;
			document.onmousemove = CAL_LEFT.dragCtlBarMoving;
			document.onmouseup = CAL_LEFT.dragCtlBarRelease;
			var posCtl = CAL_BODY.getRelPos(dragCtlBar);
				origCtlLeft = posCtl.x;
				origCtlTop = posCtl.y;
				pressX = e.pageX ? e.pageX : e.clientX;
				pressY = e.pageY ? e.pageY : e.clientY;
		}
		},

		findControlBar: function(e){
		  var src = (e && e.target) ? e.target : window.event.srcElement;
		  while(src.parentNode) {
			if(src.className == 'subtype_table' || (src.className.indexOf('event_ctlbar') >= 0 && parseInt(src.style.width) < 22) ){
				dragCtlBar = null;
				return false;
			}
			if(src.className && src.className == 'left_cell') {
			  dragCtlBar = src.childNodes[0];
			  return true;
			}
			src = src.parentNode;
		  }
		  dragCtlBar = null;
		  return false;
		},

		dragCtlBarMoving: function(e){
			e = e ? e : window.event;
		  if(!bReleased && dragCtlBar && (e.button == 1 || e.which == 1)) {
			bCtlMoved = true;
			CAL_LEFT.initForMoving();
			document.onmousemove = null;
			setTimeout("document.onmousemove = CAL_LEFT.dragCtlBarMoving;", 24);
			var y = e.pageY ? e.pageY : e.clientY;
			var t = origCtlTop + y - pressY;
			dragCtlBar.style.top = t + 'px';

			var dropIndex = parseInt((t - off_y) / barHeight);
			if(dropIndex > -1 && dropIndex < CAL.types_count) {
				var dropCell = leftCtl.rows[dropIndex].cells[0];
				dragCtlBar.parentNode.appendChild(dropCell.firstChild);
				dropCell.appendChild(dragCtlBar);
			}
		  }
		  else if(dragCtlBar && e.button == 0 || e.which == 0)
			CAL_LEFT.dragCtlBarRelease(e);
		},

		initForMoving: function()	{
			if(!bMoveInit && dragCtlBar) {
				bMoveInit = true;
				dragCtlBar.style.zIndex = 100;
				dragCtlBar.style.position = 'absolute';
				for(var i = 0, nLen = leftCtl.rows.length; i < nLen; i++) {
					var cell = leftCtl.rows[i].cells[0];
					cell.style.height = barHeight + 'px';
					if(cell.childNodes.length > 1) {
						cell.removeChild(cell.childNodes[1]);
						cell.childNodes[0].rows[0].cells[0].childNodes[0].childNodes[1].innerHTML = '+';
					}
					cell.style.paddingBottom = '0px';
				}
			}
		},

		dragCtlBarRelease: function(e){
			bReleased = true;
			if(bCtlMoved && dragCtlBar) {
			  document.onmousemove = null;
			  document.onmouseup = null;
			  dragCtlBar.style.position = '';
			  dragCtlBar.style.left = '';
			  dragCtlBar.style.top = '';
			  dragCtlBar.style.zIndex = '';
			  CAL_LEFT.saveEventTypesOrder();
			}
		},

		saveEventTypesOrder: function(){
		  var leftCtl = document.getElementById('leftCtl');
		  for(var i = 0, nLen = CAL.types_count; i < nLen; i++) {
			var type = leftCtl.rows[i].cells[0].firstChild.getAttribute('type');
			CAL.event_types[type].Index = i;
		  }
		  CAL.save_settings();
		},

		/// Functions
		clear: function(){
			var colorMenu = document.getElementById('div_color_menu');
			if (colorMenu) {
				colorMenu.style.display = 'none';
			}
		},

		refreshCalendar: function(){
			CAL_BODY.getAllEvents();
		}
	}
})(CAL, CAL_BODY);


var CAL_NAVI = function(CAL, CAL_BODY){

	return {
		initNavigator: function(div_navi){
			var str = "<table cellspacing=0 cellpadding=0 class=navi_table><tr>";
			str += "<td align=center style=width:65px;><select id='sel_year' style=margin-left:3px; onchange='CAL_NAVI.dateChanged();'></select></td>";
			str += "<td align=center style=width:65px;><select id='sel_month' style=width:60px; onchange='CAL_NAVI.dateChanged();'>";
			str += "<option value='1'>Jan</option><option value='2'>Feb</option><option value='3'>Mar</option><option value='4'>Apr</option>";
			str += "<option value='5'>May</option><option value='6'>Jun</option><option value='7'>Jul</option><option value='8'>Aug</option>";
			str += "<option value='9'>Sep</option><option value='10'>Oct</option><option value='11'>Nov</option><option value='12'>Dec</option>";
			str += "</select></td></tr><tr><td colspan=2 align=center><div id=cur_range></div></td></tr>";
			str += "<tr><td align='left' width=50% style=padding-left:4px;><a onclick='CAL_BODY.navigateDate(-1);' class='a_navi' style=font-size:12px href='#'>&lt;&lt;Prev</a></td>";
			str += "<td align='right' width=50% style=padding-right:4px;><a onclick='CAL_BODY.navigateDate(1);' class='a_navi' style=font-size:12px href='#'>Next&gt;&gt;</a></td>";
			str += "</tr><tr><td align='center' width=50% style=padding-top:5px;><a onclick='CAL_BODY.gotoToday();' class='a_navi' href='#'>Today</a></td>";
			str += "<td align='center' width=50% style=padding-top:5px;><a onclick='CAL_BODY.refreshCache();' class='a_navi' href='#'>Refresh</a></td></tr></table>";
			str += "<div style=margin-left:5px;margin-top:4px;>search:<br/>";
			str += "<input type='text' id='search' style='width:110px;' oninput='CAL_NAVI.searchChanged();'/>";
			str += "</div>";
			document.getElementById(div_navi).innerHTML = str;
			this.updateCurDate();
		},

		updateCurDate: function(){
			CAL_BODY.updateCurYearMonth();

		  var sel_year = document.getElementById('sel_year');
		  var sel_month = document.getElementById('sel_month');
		  sel_year.options.length = 0;
		  for(var i = -10; i < 10; i++) {
			var strYears = '' + (CAL_BODY.cur_year + i); 	    //add year
			var opt = document.createElement('OPTION');
			opt.text = strYears;
			opt.value = strYears;
			try{ sel_year.add(opt); }catch(e){ sel_year.add(opt, null); }
		  }
		  sel_year.options[10].selected = true;
		  sel_month.options[CAL_BODY.cur_month].selected = true;

			var cur_range = document.getElementById('cur_range');
			if(CAL_BODY.cal_lines == 6)
			cur_range.innerHTML = CAL_BODY.month_abbrv[CAL_BODY.cur_month] + ' ' + CAL_BODY.cur_year;
		  else {
			var cur_start = CAL_BODY.first_day.plusDays(CAL_BODY.cal_lines * 7);
			var cur_end = CAL_BODY.first_day.plusDays(CAL_BODY.cal_lines * 2 * 7 - 1);
			cur_range.innerHTML = CAL_BODY.month_abbrv[cur_start.getMonth()] + ' ' + cur_start.getDate() + ' ~ ' + CAL_BODY.month_abbrv[cur_end.getMonth()] + ' ' + cur_end.getDate();
		  }
		},

		dateChanged: function(){
		  var year = parseInt(document.getElementById('sel_year').options[document.getElementById('sel_year').selectedIndex].value);
		  var month = document.getElementById('sel_month').selectedIndex;
			CAL_BODY.setFirstDay(year, month);
			CAL_BODY.buildCalendar();
		  setTimeout(this.refreshCalendar, 500);
		},

		searchChanged: function(suffix) {
			suffix = suffix || '';
			CAL_BODY.changeSearchText(document.querySelectorAll('#search' + suffix)[0].value);
		},

		refreshCalendar: function(){
			CAL_BODY.getAllEvents();
		}
	}
}(CAL, CAL_BODY);


/// ZOOM Bar
var CAL_ZOOM = function(CAL, CAL_BODY, CAL_NAVI){
	var originZoomTop = 0;
	var zoom_bar = null;

	return {

		initZoom: function(map_bar, cur){
			var mapbar = document.getElementById(map_bar);
			if (mapbar) {
				mapbar.innerHTML = "<div id=zoom_bar title='drag' onselectstart='return false;'>" +
					"&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</div>" +
					"<img src='/Res?name=bar.png' usemap='#mapZoom' style=border-width:0;/>" +
						"<map name=mapZoom><area onclick='CAL_ZOOM.zoom(-1);' title='zoom in' shape='rect' coords='0, 0, 20, 16' href=#>" +
							"<area onclick='CAL_ZOOM.setZoom(1);' shape='rect' title=1 coords='0, 16, 20, 34' href=#><area onclick='CAL_ZOOM.setZoom(2);' shape='rect' title=2 coords='0, 34, 20, 52' href=#>" +
							"<area onclick='CAL_ZOOM.setZoom(3);' shape='rect' title=3 coords='0, 52, 20, 70' href=#><area onclick='CAL_ZOOM.setZoom(4);' shape='rect' title=4 coords='0, 70, 20, 88' href=#>" +
							"<area onclick='CAL_ZOOM.setZoom(5);' shape='rect' title=5 coords='0, 88, 20, 106' href=#><area onclick='CAL_ZOOM.setZoom(6);' shape='rect' title=6 coords='0, 106, 20, 124' href=#>" +
							"<area onclick='CAL_ZOOM.zoom(1);' title='zoom out' shape='rect' coords='0, 124, 16, 140' href=#><area id=fit_button onclick='CAL_ZOOM.fit();' shape='rect' coords='0, 142, 20, 156' href=#>" +
						"</map>";
			}

			var z = this.readZoomCookie();
			if(z > 0 && z < 6) {
				CAL_BODY.cal_lines = parseInt(z);
				var now = (typeof cur == 'number') ? new Date() : cur;
				var center_day = new Date(now.getFullYear(), now.getMonth(), now.getDate());
				center_day.setDate(center_day.getDate() - center_day.getDay());
				CAL_BODY.first_day = center_day.plusDays( -1 * CAL_BODY.halfCalWeeks());
			}
			else
				CAL_BODY.setFirstDay(cur);
			document.getElementById('fit_button').title = (z == 'fit' ? 'return to month view' : 'fit to screen');
			this.paintZoom();
			zoom_bar = document.getElementById('zoom_bar');
			zoom_bar.onmousedown = CAL_ZOOM.zoomPressed;
		},

		readZoomCookie: function(){
			var z_old = CAL.getCookie('zoom');
			if(z_old != '') {
				CAL.zoom = z_old;
				CAL.save_settings();
				CAL.deleteCookie('zoom');
				return z_old;
			}

			var cookies = CAL.getCookie('CAL');
			cookies = (cookies == '' ? {} : plugin.evalJSON(cookies));
			var z = cookies.zoom ? cookies.zoom : 6;
			return z;
		},

		zoomPressed: function(e){
			document.onmousemove = CAL_ZOOM.zoomMoving;
		document.onmouseup = CAL_ZOOM.zoomRelease;
		return false;
		},

		zoomMoving: function(e){
			e = e ? e : window.event;
			if(e.button <= 1) {
				var y = e.pageY ? e.pageY : e.clientY;
				if(y > 16 && y < 120)
				zoom_bar.style.top = y + 'px';
			}
			return false;
		},

		zoomRelease: function(e){
			document.onmousemove = null;
		document.onmouseup = null;
		  var top = parseInt(zoom_bar.style.top)
		  var z = Math.floor((top - 16) / 18) + 1;
			setTimeout('CAL_ZOOM.setZoom(' + z + ')', 20);
			return false;
		},

		fit: function(){
			if(CAL.zoom == 'fit') {
				CAL.zoom = 6;
				if(CAL_BODY.cal_lines < 6) this.setZoom(6);
				CAL.save_settings();
				document.getElementById('fit_button').title = 'fit to screen';
			}
			else {
				this.fitToWindow();
				CAL.zoom = 'fit';
			}
		},

		fitToWindow: function(){
			var max_nums = [];
			for(var i = 0, nLen = CAL_BODY.cal_lines; i < nLen; i++){
				var week_date = CAL_BODY.first_day.plusDays((CAL_BODY.cal_lines + i) * 7);
				var week_events = CAL.get_week_events(week_date);
				var max_num = 0;
				for(var j = 0; j < 7; j++) {
					var wi = week_events[j];
					if(wi.num + wi.co_num > max_num)
						max_num = wi.num + wi.co_num;
				}
				max_nums[i] = max_num;
			}

			var z = CAL_BODY.cal_lines;
			var nStart = 0;
			var nEnd = CAL_BODY.cal_lines;
			do {
				var try_can_hold = Math.floor((CAL_BODY.cal_height - z * CAL_BODY.cell_date_height) / (z * CAL_BODY.event_height));
				var can_hold_ok = true;
				for(var i = nStart; i < nEnd; i++)
				{
					if(try_can_hold < max_nums[i])
						{can_hold_ok = false; break;}
				}
				if(can_hold_ok == true)	break;
				if(z % 2 == 0)	nEnd--; else nStart++;
				z--;
			}while(z > 1)

			if(z < CAL_BODY.cal_lines)
				this.setZoom(z);
			CAL.zoom = 'fit';
			CAL.save_settings();
			document.getElementById('fit_button').title = 'return to month view';
		},

		fitOnLoad: function(){
			if(CAL.zoom == 'fit')
				CAL_ZOOM.fitToWindow();
		},

		zoom: function(distance){
			this.setZoom(CAL_BODY.cal_lines + distance);
		},

		setZoom: function(lines){
			if (CAL_BODY.cal_lines == lines || lines < 1 || lines > 6) return;

			var center_day = CAL_BODY.first_day.plusDays(CAL_BODY.halfCalWeeks());
			CAL_BODY.cal_lines = lines;
			CAL_BODY.first_day = center_day.plusDays(-1 * CAL_BODY.halfCalWeeks());
			CAL_BODY.computeSize();
			CAL.zoom = CAL_BODY.cal_lines;
			CAL.save_settings();

			this.paintZoom();
			CAL_NAVI.updateCurDate();
			CAL_BODY.buildCalendar();
			CAL_BODY.cal_box.style.top = -1 * CAL_BODY.cal_height + 'px';
			setTimeout(this.refreshCalendar, 500);
		},

		refreshCalendar: function(){
			CAL_BODY.getAllEvents();
		},

		paintZoom: function(){
			document.getElementById('zoom_bar').style.top = (CAL_BODY.cal_lines - 1) * 18 + 22 + 'px';
		}
	}
}(CAL, CAL_BODY, CAL_NAVI);

function debug(str){ document.getElementById('div_cal_events').innerHTML += str + '<br/>'; }

function load()
{
	CAL_LEFT.initLeftDiv('div_left');
	CAL_ZOOM.initZoom('map_bar', today);
	CAL_NAVI.initNavigator('div_navi', today);
	CAL_BODY.initCalendar('div_calendar_frame', {
		onDateChanged: CAL_NAVI.updateCurDate,
		onTypesLoaded: CAL_LEFT.buildTypes,
		onFirstLoaded: CAL_ZOOM.fitOnLoad
	});
}