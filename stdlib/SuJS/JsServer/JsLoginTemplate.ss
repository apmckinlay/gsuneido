// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
`
<html>
<head lang="en">
	<meta charset="UTF-8">
	<link rel="icon" type="image/x-icon" href="@.favIcon">
	@.touchIcon
	<link rel="manifest" href="@.manifest">
	<style type="text/css">
		html,
		body {
			margin: 0;
			height: 100%;
		}
	</style>
	<style id="loginStyles" type="text/css">
		header {
			width: 100%;
			height: 100px;
			background: #29a1d8;
			background: linear-gradient(90deg, #29a1d8 0, #1a59a4);
			display: flex;
			align-items: center;
			justify-content: flex-start;
		}
		header a {
			padding-left: 15px;
			display: flex;
			align-items: center;
			justify-content: center;
		}
		header img {
			width: 200px;
			height: 80px;
		}
		#bg {
			display: flex;
			justify-content: center;
			align-items: center;
			background-image : url(/Res?name=login-axon-soft-bg-1.png);
			height: calc(100% - 100px);
		}
		.loginBody {
			display: flex;
			align-items: center;
			justify-content: center;
		}
		#container {
			background-image: url("data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHhtbG5zOnhsaW5rPSJodHRwOi8vd3d3LnczLm9yZy8xOTk5L3hsaW5rIiB3aWR0aD0iNTAwIiBoZWlnaHQ9IjI4MCIgdmlld0JveD0iMCAwIDUwMCAyODAiPg0KICA8ZGVmcz4NCiAgICA8bGluZWFyR3JhZGllbnQgaWQ9ImxpbmVhci1ncmFkaWVudCIgeDE9IjAuMTEiIHkxPSIxLjE5NiIgeDI9IjAuODkiIHkyPSItMC4xOTYiIGdyYWRpZW50VW5pdHM9Im9iamVjdEJvdW5kaW5nQm94Ij4NCiAgICAgIDxzdG9wIG9mZnNldD0iMCIgc3RvcC1jb2xvcj0iIzI5YTFkOCIvPg0KICAgICAgPHN0b3Agb2Zmc2V0PSIxIiBzdG9wLWNvbG9yPSIjMDA0ODk1Ii8+DQogICAgPC9saW5lYXJHcmFkaWVudD4NCiAgPC9kZWZzPg0KICA8cmVjdCBpZD0iYmFja2dyb3VuZC1ncmFkaWVudCIgd2lkdGg9IjUwMCIgaGVpZ2h0PSIyODAiIGZpbGw9InVybCgjbGluZWFyLWdyYWRpZW50KSIvPg0KPC9zdmc+DQo=");
			background-size: cover;
			color: white;
			font-family: Arial;
			padding: 10px;
			display: flex;
			flex-direction: column;
			justify-content: center;
			align-items: center;
			width: 38em;
			height: 380px;
			position: relative;
		}
		#container form { width: 100%; height: 200px;}
		#container h1 { margin: 10px 0;}

		#logo {
			display: inline-block;
			height: 5em;
		}

		h1 {
			margin-top: 0.4em;
			margin-bottom: 1em;
			font-size: .9em;
			font-weight: bold;
			align-self: center;
		}

		input {
			border: 1px solid white;
			border-radius: 5px;
			font-size: 1.2em;
			padding: 2px;
			margin-bottom: 0.2em;
		}

		#code:focus::placeholder {
			color: transparent;
		}

		label {
			font-size: .8em;
			color: lightgray;
			margin-bottom: .6em;
			margin-top: .2em;
		}

		a {
			color: white;
		}

		.su-submit {
			color: white;
			border: 0;
			background-color: #024795;
			padding: .4em 2em;
			border-radius: 4px;
			font-size: .9em;
			margin: 10px 0;
			cursor: pointer;
		}

		.msg {
			width: 100%;
			height: auto;
			color: black;
			font-weight: bold;
			display: none;
			position: static;
			bottom: 30px;
			left: 10px;
		}

		#copyright {
			color: #024795;
			font-size: 0.7em;
			margin-right: 10px;
		}

		#overlay {
			border: 1px solid grey;
			pointer-events: none;
		}
		#overlay::backdrop {
			background-color: transparent;
			cursor: wait;
		}

		#forgot {
			background: none;
			color: rgba(255, 255, 255, .5);
			text-decoration: underline;
			border: none;
			cursor: pointer;
			padding: 0;
			margin: 0;
			margin-left: 10px;
			font-size: 0.7em;
			font-weight: bold;
		}

		#vert {
			display: flex;
			flex-direction: column;
			flex-grow: 1;
			flex-shrink: 0;
			align-items: center;
			width: 100%;
		}
		footer{
			display: flex;
			width: 100%;
			justify-content: space-between;
			position: absolute;
			bottom: 10px;
		}
		footer.noforgetpassword {
			justify-content: flex-end;
		}
		footer.noforgetpassword input{
			display: none;
		}
		.footer {
			text-align: right;
		}
		.checkboxlabel-container {
			display: flex;
			align-items: center;
		}
		.checkboxlabel-container label {
			margin: 0;
			margin-right: 10px;
			color: lightgray;
			font-size: .7em;
		}
		.checkboxlabel-container input[type="checkbox"] {
			margin-left: 10px;
		}
	</style>
</head>

<body>
	<header id="header">
		<a href="https://axonsoftware.com/" target="_blank">
			<img src="@.logo" alt="Axon Software Logo" title="Axon Software">
		</a>
	</header>
	<div id="bg">
		<div id="container" class="loginBody">
			<img id="logo" src="@.logo" alt="Axon Software Logo" title="Axon Software">
			<h1>Transportation Management System</h1>
			<form id="step1" enctype="text/plain" action="javascript:loginSubmit()">
				<div id="vert">
					<div class="checkboxlabel-container" style="margin-bottom: 10px;">
						<input type="checkbox" id="newuser" name="newreset"
								onclick="togglePasswordEnabled()">
						<label>New User or Password Reset</label>
					</div>
					<input type="text" id="user" name="user" size="20"
						autocomplete="su-do-not-autofill">
					<label for="user">User Name</label>
					<input type="password" id="password" name="password" size="20"
						autocomplete="su-do-not-autofill"
						readonly onfocus="passwordFocus();">
					<div class="checkboxlabel-container">
						<label for="password" style="font-size: .8em;">Password</label>
						<input type="checkbox" id="showPassword"
							onclick="togglePasswordVisibility()">
						<label for="showPassword">Show</label>
					</div>
					<input class="su-submit" type="submit" value="Login">
				</div>
			</form><!-- #step1 -->
			<form id="step1_1" enctype="text/plain" action="javascript:loginSubmitWithEmail()" style="display:none">
				<div id="vert">
					<input type="email" id="email" size="20" placeholder="@.domain" pattern=".+@.domain">
					<label for="email">Email to receive Login Code</label>
					<input class="su-submit" type="submit" value="Submit">
				</div>
				<!-- <div style="display: flex; flex-direction: column; 	flex-grow: 1; justify-content: flex-end; height: 100%">
					<span class="footer" id="copyright">&copy; Copyright 2000-@.toYear Axon Development Corporation</span>
				</div> -->
			</form><!-- #step1_1 -->
			<form id="step2" enctype="text/plain" action="javascript:twoFASubmit()" style="display:none">
				<div id="vert">
					<input type="text" required="" placeholder="000000" spellcheck="false" autocomplete="false" autocapitalize="none" id="code" inputmode="numeric">
					<label id="codeLabel" for="code"></label>
					<input class="su-submit" type="submit" value="Submit">
				</div>
				<!-- <div style="display: flex; flex-direction: column; 	flex-grow: 1; justify-content: flex-end; height: 100%">
					<span class="footer" id="copyright">&copy; Copyright 2000-@.toYear Axon Development Corporation</span>
				</div> -->
			</form><!-- #step2 -->
			<div id="msg1" class="msg"></div>
			<div id="msg2" class="msg"></div>
			<div id="msg3" class="msg"></div>
			<footer>
				<input id="forgot" type="submit" value="Forgot Password"
					onclick="handleForgotPassword()">
				<input type="checkbox" id="forgotPassword" name="forgotPassword"
					style="visibility: hidden; display: none">
				<span class="footer" id="copyright">&copy; Copyright 2000-@.toYear Axon Development Corporation</span>
			</footer>
		</div><!-- #container -->
		<dialog id="overlay">
			<p>Working...</p>
		</dialog>
	</div><!-- #bg -->

	<script>
		function isTextMetricsSupported() {
			if (typeof CanvasRenderingContext2D === 'undefined') {
				return false; // Canvas not supported
			}
			var canvas = document.createElement('canvas');
			var ctx = canvas.getContext('2d');
			if (typeof ctx.measureText === 'undefined' || typeof ctx.measureText('test').actualBoundingBoxAscent === 'undefined') {
				return false; // measureText or actualBoundingBoxAscent not supported
			}
			var textMetrics = ctx.measureText('Hg');
			return typeof textMetrics.fontBoundingBoxAscent !== 'undefined';
		}
		function isDialogSupported() {
			return 'HTMLDialogElement' in window;
		}
		function isDecompressionSupported() {
			return 'DecompressionStream' in window;
		}
		var checkBrowserCompatibility = function() {
			if (!isTextMetricsSupported() || !isDialogSupported() || !isDecompressionSupported()) {
				return 'Your browser is not supported. Please update your browser to ensure the best user experience.'
			}
			return '';
		};
		var checkSecureConnection = function() {
			if (location.protocol === 'http:' &&
				(location.hostname !== 'localhost' && location.hostname !== '127.0.0.1')) {
				return 'Your connection is not secure. Please use a HTTPS connection.';
			}
			return '';
		};
		var suLoginSession;
		var suLoginFocus;
		var displayMsg = function(msg, form) {
			var msgElm = document.getElementById("msg" + form);
			msgElm.innerText = msg;
			msgElm.style.display = 'block';
		};
		var clearAfterLogin = function() {
			document.getElementById("header") && document.getElementById("header").remove();
			document.getElementById("bg") && document.getElementById("bg").remove();
			document.body.classList.remove("loginBody")
			document.getElementById("loginStyles") && document.getElementById("loginStyles").remove();
		};
		var loadAfterLogin = function(res) {
			var item;
			var el;
			var attr
			var scriptToLoad = [];
			var scriptLoaded = 0;
			/* ensure the script load sequence */
			var onload = function() {
				if (scriptToLoad.length === 0) {
					eval(res.onload);
					return;
				}
				item = scriptToLoad.shift();
				makeNode(item, onload);
			};
			var makeNode = function(item, onload) {
				el = document.createElement(item.tag);
				if (onload) {
					el.onload = onload;
				}
				for (attr in item) {
					if (attr === "tag") {
						continue;
					}
					el[attr] = item[attr];
				}
				document.head.appendChild(el);
			};
			for (item of res.sources) {
				if (item.tag === "script") {
					scriptToLoad.push(item);
				} else {
					makeNode(item)
				}
			}
			onload();
		}

		function passwordFocus() {
			var passwordInput = document.getElementById("password")
			var newuser = document.getElementById("newuser")
			if (!newuser.checked)
				{
				passwordInput.removeAttribute('readonly');
				}
			}

		function focusElement(el) {
			suLoginFocus = el;
			setTimeout(function () {
				if (suLoginFocus) {
					suLoginFocus.focus();
					suLoginFocus.select();
					suLoginFocus = null;
				}
			}, 10);
		}

		function show(form) {
			function disableForm(form, disable) {
				var elements = form.elements;
				for (let i = 0; i < elements.length; i++) {
					elements[i].disabled = disable;
				}
			}
			var forms = ["step1", "step1_1", "step2"];
			var i, cur;
			var footer = document.getElementsByTagName('footer')[0];
			for (var i = 0; i < forms.length; i++) {
				if (forms[i] == form) {
					cur = document.getElementById(forms[i]);
					cur.style.display = '';
					disableForm(cur, false);
				} else {
					cur = document.getElementById(forms[i]);
					cur.style.display = 'none';
					disableForm(cur, true);
				}
				if(form === 'step1_1' || form === 'step2') {
					footer.classList.add('noforgetpassword');
				}
			}
		}

		var sendingRequest = false;
		var afterLogin = false
		function sendRequest(path) {
			if (sendingRequest || afterLogin) {
				return;
			}
			sendingRequest = true;
			var dialog = document.getElementById("overlay");
			dialog.showModal();
			var formData = JSON.stringify(suLoginSession);
			var request = new XMLHttpRequest();
			request.onreadystatechange = function () {
				if (request.readyState === XMLHttpRequest.DONE) {
					if (request.status === 200) {
						var res = JSON.parse(request.responseText);
						var form = "1";
						if (res.form) {
							form = res.form;
						}
						if (res.err) {
							displayMsg(res.err, form);
							if (res.back) {
								document.getElementById("code").value = '';
								show("step1");
							}
							if (res.focus) {
								focusElement(document.getElementById(res.focus))
							}
						} else if (res.step2) {
							document.getElementById("codeLabel").innerHTML = res.msg;
							delete res.msg;
							show("step2");
							suLoginSession = Object.assign(suLoginSession, res);
							focusElement(document.getElementById("code"));
						} else {
							doLogin(res);
						}
					} else {
						displayMsg("There was a problem with the request", "1");
					}
				dialog.close();
				sendingRequest = false;
				}
			};
			request.open("POST", path);
			request.send(formData);
		}

		function doLogin(info) {
			afterLogin = true;
			clearAfterLogin();
			loadAfterLogin(info);
		}

		function loginSubmit() {
			displayMsg("", "1");
			var digest=function(r){function n(r,n){return r<<n|r>>>32-n}function t(r,n){var t,o,e,u,f;return e=2147483648&r,u=2147483648&n,f=(1073741823&r)+(1073741823&n),(t=1073741824&r)&(o=1073741824&n)?2147483648^f^e^u:t|o?1073741824&f?3221225472^f^e^u:1073741824^f^e^u:f^e^u}function o(r,o,e,u,f,i,a){return t(n(r=t(r,t(t(function(r,n,t){return r&n|~r&t}(o,e,u),f),a)),i),o)}function e(r,o,e,u,f,i,a){return t(n(r=t(r,t(t(function(r,n,t){return r&t|n&~t}(o,e,u),f),a)),i),o)}function u(r,o,e,u,f,i,a){return t(n(r=t(r,t(t(function(r,n,t){return r^n^t}(o,e,u),f),a)),i),o)}function f(r,o,e,u,f,i,a){return t(n(r=t(r,t(t(function(r,n,t){return n^(r|~t)}(o,e,u),f),a)),i),o)}function i(r){var n,t="",o="";for(n=0;n<=3;n++)t+=(o="0"+(r>>>8*n&255).toString(16)).substr(o.length-2,2);return t}var a,c,C,g,h,d,v,S,m,l=Array();for(l=function(r){for(var n,t=r.length,o=t+8,e=16*((o-o%64)/64+1),u=Array(e-1),f=0,i=0;i<t;)f=i%4*8,u[n=(i-i%4)/4]=u[n]|r.charCodeAt(i)<<f,i++;return f=i%4*8,u[n=(i-i%4)/4]=u[n]|128<<f,u[e-2]=t<<3,u[e-1]=t>>>29,u}(r=function(r){r=r.replace(/\r\n/g,"\n");for(var n="",t=0;t<r.length;t++){var o=r.charCodeAt(t);o<128?n+=String.fromCharCode(o):o>127&&o<2048?(n+=String.fromCharCode(o>>6|192),n+=String.fromCharCode(63&o|128)):(n+=String.fromCharCode(o>>12|224),n+=String.fromCharCode(o>>6&63|128),n+=String.fromCharCode(63&o|128))}return n}(r)),d=1732584193,v=4023233417,S=2562383102,m=271733878,a=0;a<l.length;a+=16)c=d,C=v,g=S,h=m,v=f(v=f(v=f(v=f(v=u(v=u(v=u(v=u(v=e(v=e(v=e(v=e(v=o(v=o(v=o(v=o(v,S=o(S,m=o(m,d=o(d,v,S,m,l[a+0],7,3614090360),v,S,l[a+1],12,3905402710),d,v,l[a+2],17,606105819),m,d,l[a+3],22,3250441966),S=o(S,m=o(m,d=o(d,v,S,m,l[a+4],7,4118548399),v,S,l[a+5],12,1200080426),d,v,l[a+6],17,2821735955),m,d,l[a+7],22,4249261313),S=o(S,m=o(m,d=o(d,v,S,m,l[a+8],7,1770035416),v,S,l[a+9],12,2336552879),d,v,l[a+10],17,4294925233),m,d,l[a+11],22,2304563134),S=o(S,m=o(m,d=o(d,v,S,m,l[a+12],7,1804603682),v,S,l[a+13],12,4254626195),d,v,l[a+14],17,2792965006),m,d,l[a+15],22,1236535329),S=e(S,m=e(m,d=e(d,v,S,m,l[a+1],5,4129170786),v,S,l[a+6],9,3225465664),d,v,l[a+11],14,643717713),m,d,l[a+0],20,3921069994),S=e(S,m=e(m,d=e(d,v,S,m,l[a+5],5,3593408605),v,S,l[a+10],9,38016083),d,v,l[a+15],14,3634488961),m,d,l[a+4],20,3889429448),S=e(S,m=e(m,d=e(d,v,S,m,l[a+9],5,568446438),v,S,l[a+14],9,3275163606),d,v,l[a+3],14,4107603335),m,d,l[a+8],20,1163531501),S=e(S,m=e(m,d=e(d,v,S,m,l[a+13],5,2850285829),v,S,l[a+2],9,4243563512),d,v,l[a+7],14,1735328473),m,d,l[a+12],20,2368359562),S=u(S,m=u(m,d=u(d,v,S,m,l[a+5],4,4294588738),v,S,l[a+8],11,2272392833),d,v,l[a+11],16,1839030562),m,d,l[a+14],23,4259657740),S=u(S,m=u(m,d=u(d,v,S,m,l[a+1],4,2763975236),v,S,l[a+4],11,1272893353),d,v,l[a+7],16,4139469664),m,d,l[a+10],23,3200236656),S=u(S,m=u(m,d=u(d,v,S,m,l[a+13],4,681279174),v,S,l[a+0],11,3936430074),d,v,l[a+3],16,3572445317),m,d,l[a+6],23,76029189),S=u(S,m=u(m,d=u(d,v,S,m,l[a+9],4,3654602809),v,S,l[a+12],11,3873151461),d,v,l[a+15],16,530742520),m,d,l[a+2],23,3299628645),S=f(S,m=f(m,d=f(d,v,S,m,l[a+0],6,4096336452),v,S,l[a+7],10,1126891415),d,v,l[a+14],15,2878612391),m,d,l[a+5],21,4237533241),S=f(S,m=f(m,d=f(d,v,S,m,l[a+12],6,1700485571),v,S,l[a+3],10,2399980690),d,v,l[a+10],15,4293915773),m,d,l[a+1],21,2240044497),S=f(S,m=f(m,d=f(d,v,S,m,l[a+8],6,1873313359),v,S,l[a+15],10,4264355552),d,v,l[a+6],15,2734768916),m,d,l[a+13],21,1309151649),S=f(S,m=f(m,d=f(d,v,S,m,l[a+4],6,4149444226),v,S,l[a+11],10,3174756917),d,v,l[a+2],15,718787259),m,d,l[a+9],21,3951481745),d=t(d,c),v=t(v,C),S=t(S,g),m=t(m,h);return(i(d)+i(v)+i(S)+i(m)).toLowerCase()};
			var user = document.getElementById("user");
			var password = document.getElementById("password");
			var newresetCheckbox = document.getElementById("newuser");
			var newReset = newresetCheckbox.checked;
			var forgotPasswordCheckbox = document.getElementById("forgotPassword")
			var forgotPassword = forgotPasswordCheckbox.checked;

			if ((!user || !user.value) && (!password || !password.value) && !newReset &&
				!forgotPassword) {
				displayMsg("Please enter your user name and password", "1");
				focusElement(user);
				return;
			}

			if (user.value && (!password || !password.value) && !newReset &&
				!forgotPassword) {
				displayMsg("Please enter your password", "1");
				return;
			}

			if (!user || !user.value) {
				displayMsg("Please enter your user name", "1");
				focusElement(user);
				return;
			}

			suLoginSession = {
				user: user.value,
				password: digest(user.value + password.value),
				newuserreset: newReset,
				forgotPassword: forgotPassword,
				user_agent: navigator.userAgent};

			newresetCheckbox.checked = false;
			togglePasswordEnabled();
			forgotPasswordCheckbox.checked = false;
			document.getElementById("showPassword").checked = false
			togglePasswordVisibility();

			if ((user.value === "axon" || user.value === "default") && (!forgotPassword) && (!newReset)) {
				show("step1_1");
				focusElement(document.getElementById("email"));
				return;
			}

			sendRequest("login_submit");
		}

		function loginSubmitWithEmail() {
			var email = document.getElementById("email");
			if (!email || !email.value) {
				displayMsg("Email cannot be empty", "2");
				focusElement(email);
				return;
			}
			if (!email.value.endsWith("@.domain")) {
				displayMsg("You must use an email of axonsoft.com", "2");
				focusElement(email);
				return;
			}

			suLoginSession.email = email.value;
			sendRequest("login_submit");
		}

		function twoFASubmit() {
			displayMsg("", "3");
			if (!suLoginSession) {
				displayMsg("Invalid session", "3");
			}
			var code = document.getElementById("code");
			if (!code || !code.value) {
				displayMsg("Code cannot be empty", "3");
				focusElement(code);
				return;
			}
			suLoginSession.code = code.value;
			sendRequest("twoFA_submit");
		}

		function preauthOpenBook() {
			// Handle user clicking "Open Another Book" from another tab
			const urlParams = new URLSearchParams(window.location.search);
			const preauth = urlParams.get('preauth') || false;
			const user = urlParams.get('user') || '';
			const book = urlParams.get('book') || '';
			const token = urlParams.get('token') || '';
			if (preauth === 'true') {
				suLoginSession = {
					user: user,
					preauth: preauth,
					user_agent: navigator.userAgent,
					book: book
					}
				sendRequest("login_submit?token=" + token);
				}
			}
		window.onload = function (event) {
			var msg = checkBrowserCompatibility() || checkSecureConnection();
			if (msg) {
				for (var el of document.getElementsByTagName('input')) {
					el.setAttribute('disabled', true);
				}
				alert(msg);
			} else {
				preauthOpenBook();
				var userInput = document.getElementById("user");
				userInput && userInput.focus();
				@.extraOnLoad
			}
			window.onload = null;
		};
		function togglePasswordEnabled() {
			var passwordInput = document.getElementById("password")
			var newuser = document.getElementById("newuser")
			if (newuser.checked)
				{
				passwordInput.readOnly = true;
				passwordInput.style.backgroundColor = 'gray';
				}
			else
				{
				passwordInput.readOnly = null;
				passwordInput.style.backgroundColor = 'white';
				}
			}
		function handleForgotPassword() {
			document.getElementById("forgotPassword").checked = true;
			loginSubmit();
			}
		function togglePasswordVisibility() {
			var passwordInput = document.getElementById("password");
			var showPasswordCheckbox = document.getElementById("showPassword");
			if (showPasswordCheckbox.checked) {
				passwordInput.type = "text";
			} else {
				passwordInput.type = "password";
			}
		}
	</script>
</body>

</html>`
