// Copyright (C) 2022 Axon Development Corporation All rights reserved worldwide.
class
	{
	CallClass(env)
		{
		if not env.queryvalues.Member?(0)
			return JsInvalidPage(title: 'Attachment',
				reason: 'Your request is invalid')
		if JsSessionToken.Validate(env) is false
			return JsInvalidPage(title: 'Attachment',
				reason: 'Your session is invalid or expired')
		filename = Base64.Decode(env.queryvalues[0]).Xor(EncryptControlKey())
		token = env.queryvalues.token
		if .notAllowPreview?(filename)
			{
			if false isnt result = .presignedUrl(filename)
				return Object('302', [Location: result], '')
			return JsDownload.Download(env, filename, preview?: false)
			}

		if false is url = .presignedUrl(filename)
			url = Url.Encode('download', Object(env.queryvalues[0], :token, preview:))
		return .preview(filename, url)
		}

	previewSuffixes: (pdf, ts, json, js, java, go, py, config)
	notAllowPreview?(filename)
		{
		suffix = filename.AfterLast('.').Lower()
		if .previewSuffixes.Has?(suffix)
			return false
		if false is type = MimeTypes.GetDefault(suffix, false)
			return true
		return not type.Prefix?('image') and not type.Prefix?('text') and
			not type.Prefix?('video')
		}

	presignedUrl(filename)
		{
		OptContribution("Attachment_PresignedUrl", {|unused| false })(filename)
		}

	preview(filename, url)
		{
		suffix = filename.AfterLast('.').Lower()
		type = MimeTypes.GetDefault(suffix, '')

		return '<!DOCTYPE html>' $ Razor(
			type.Prefix?('image') ? .template : .template,
			Object(title: Paths.Basename(filename), src: url, :type))
		}

	template: `
<html>
	<head>
		<title>@.title</title>
		<style>
			html, body, div, iframe
				{
				 margin: 0px;
				 padding: 0px;
				 height: 100%;
				 border: none;
				}
			iframe
				{
				 display: block;
				 width: 100%;
				 border: none;
				 overflow-y: auto;
				 overflow-x: hidden;
				}
			.wrapper
				{
				display:flex;
				justify-content:center;
				background-color: rgba(0,0,0,0.8)
				}
			img
				{
				margin: auto;
				max-width:100%;
				max-height:100%;
				}
		</style>
	</head>
	<body>
		<iframe id="pdfFrame"
			title='@.title'
			frameborder="0"
			marginheight="0"
			marginwidth="0"
			width="100%"
			height="100%"
			scrolling="auto" style="display: none" >Loading</iframe>
		<div class="wrapper" width="100%"
			height="100%" style="display: none" id="imgDiv">
			<img id="img"/>
		</div>
		<script>
		var iframe = document.getElementById("pdfFrame");
		var imgDiv = document.getElementById("imgDiv");

		if ('@.type'.indexOf('image/') >= 0) {
			imgDiv.style.display = 'flex';
			document.getElementById("img").src = '@HtmlString(.src)';
		} else {
			iframe.src = '@HtmlString(.src)';
			imgDiv.style.display = 'none';
			iframe.style.display = 'block';
		}
		</script>
		</body>
</html>`
	}