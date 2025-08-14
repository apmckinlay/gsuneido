// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
// ref: https://github.com/KnicKnic/WASM-ImageMagick/tree/master/apidocs
Component
	{
	styles: `
		.su-loader-container {
			padding: 20px;
			display: flex;
			border-radius: 10px;
			flex-direction: column;
			justify-content: center;
			align-items: center;
			border: none;
			border-spacing: 0px;
			overflow: hidden;
		}
		.su-loader {
			border: 16px solid #f3f3f3; /* Light grey */
			border-top: 16px solid #3498db; /* Blue */
			border-radius: 50%;
			width: 120px;
			height: 120px;
			animation: spin 2s linear infinite;
			margin-bottom: 10px;
		}
		@keyframes spin {
			0% { transform: rotate(0deg); }
			100% { transform: rotate(360deg); }
		}`
	New(filter, hDrop = false, multi = false, .fileSizeLimit = 10, .s3? = false, cdn = '')
		{
		LoadCssStyles('su-loader.css', .styles)
		.CreateElement('div', className: 'su-loader-container')
		.loader = CreateElement('div', .El, className: 'su-loader')
		SuUI.GetCurrentWindow().Eval(.initFileCount())

		if .s3?
			{
			SuUI.GetCurrentWindow().Eval(.LoadMagickScript(cdn))
			SuUI.GetCurrentWindow().LoadMagick(
				{ .uploadAll(hDrop, filter, multi) },
				{|unused|
				.magick? = false
				SuRender().Event(false, 'SuneidoLog', Object(
					'ERROR: (CAUGHT) failed to load ImageMagick.js',
					caughtMsg: 'continue uploading without rotation'))
				.uploadAll(hDrop, filter, multi)
				})
			}
		else
			.uploadAll(hDrop, filter, multi)
		}

	LoadMagickScript(cdn)
		{
		Assert(cdn isnt '')
		return `if (!window.magick) {
window.uint8Array = function (arrayBuffer) {
	if(typeof(arrayBuffer) !== 'string')
		return new Uint8Array(arrayBuffer);

	const len = arrayBuffer.length;
	const uint8 = new Uint8Array(len);
	for (let i = 0; i < len; i++) {
		uint8[i] = arrayBuffer.charCodeAt(i);
	}
	return uint8;
};

window.file = function (src, type) {
	return new File([src.blob], src.name, {
		type: type,
		lastModified: Date.now(),
		});
};

window.arrayToString = function(uint8Array) {
	const CHUNK_SIZE = 0x8000; // 32,768
	let result = '';
	for (let i = 0; i < uint8Array.length; i += CHUNK_SIZE) {
		const chunk = uint8Array.subarray(i, i + CHUNK_SIZE);
		result += String.fromCharCode(...chunk); // since argument size limit
	}
	return result;
}

window.loadMagick = function (block, catchBlock) {
	const importFunc = function (c, v) {
		import("` $ cdn $ `/magickApi.js?v="` $ `+v).then(handler).catch(c);
	}

	const handler = function (f) {
		window.magick = function (files, command) { return f.call(files, command); };
		block();
	}

	const catchWrapper = function (unused) {
		setTimeout(function () {
			console.log('Retrying to load ImageMagick');
			importFunc(catchBlock, Date.now()) },
			2000);
	}

	// 'https://cdn.jsdelivr.net/npm/wasm-imagemagick/dist/bundles/magickApi.js'
	importFunc(catchWrapper, 1);
};

window.loadPako = function (block, catchBlock) {
	// 'https://cdn.jsdelivr.net/npm/pako@latest/dist/pako.min.js'
	import("` $ cdn $ `/pako.min.js").then(
		function(f) {
			window.pakoInflate = function(str) {
				return arrayToString(pako.inflate(uint8Array(str)));
			};
			block();
		}).catch(catchBlock);
};

window.downloadFile = function(filename, content) {
	const array = uint8Array(content);
	const blob = new Blob([array], { type: 'application/pdf' });
	const link = document.createElement('a');
	link.href = URL.createObjectURL(blob);
	link.download = filename;
	document.body.appendChild(link);
	link.click();
	document.body.removeChild(link);
	URL.revokeObjectURL(link.href);
};
}`
	}

	uploadAll(hDrop, filter, multi)
		{
		.uploadTasks = Object()
		if hDrop is false
			{
			input = CreateElement('input', .El)
			input.type = 'file'
			if filter isnt ''
				input.accept = filter
			if multi
				input.multiple = 'multiple'
			// Safari doens't popup the open file window automatically on click
			if SuRender().Engine isnt 'WebKit'
				input.SetStyle('display', 'none')
			input.onchange = .onChange
			input.Click()
			}
		else if false isnt files = SuRender().GetDropFiles(hDrop)
			.upload(files)
	}

	onChange(event)
		{
		.upload(event.target.files)
		}

	upload(files)
		{
		// check for invalid files
		if not .validateFiles?(files)
			return

		.fileCount = files.length
		for i in ...fileCount
			{
			file = files.Item(i)
			msgEl = CreateElement('div', .El)
			fileNameEl = CreateElement('div', .El)
			fileNameEl.innerText = file.name $ ' (' $ .formatSize(file.size) $ ')'
			try
				.send(file, msgEl)
			catch (e)
				.updateMsg(msgEl, e, color: 'red')
			}
		}

	validateFiles?(files)
		{
		for i in ..files.length
			{
			if .El is false
				return false

			file = files.Item(i)
			if file.size > .fileSizeLimit.Mb()
				{
				.Event(#FileSizeOverLimit)
				return false
				}
			if ExecutableExtension?(file.name)
				{
				.Event(#InvalidExtenstion, file.name)
				return false
				}
			}
		return true
		}

	send(file, msgEl)
		{
		if .s3?
			.sendToS3(file, msgEl)
		else
			.sendFile(file, msgEl, "/upload" $
				Url.BuildQuery([file.name, token: SuRender().GetToken()]))
		}

	initFileCount()
		{
		return `window.getFileCount = function () {
					window.fileCount = (window.fileCount || 0) + 1;
					return window.fileCount
				}

				window.clearFileCount = function () {
					window.fileCount = 0;
				}`
		}

	sendFile(file, msgEl, url, method = 'POST')
		{
		xhr = SuUI.MakeWebObject('XMLHttpRequest')
		.uploadTasks[file.name] = Object([:xhr])
		xhr.upload.AddEventListener("progress",
			{ |event|
			if .El isnt false
				{
				percent =  (event.loaded / event.total).DecimalToPercent(0)
				.updateMsg(msgEl, 'Uploading...' $ String(percent).LeftFill(2, '0') $ '%')
				}
			})
		xhr.AddEventListener('readystatechange', { |event/*unused*/|
			if .El isnt false and xhr.readyState is 4/*=DONE*/
				{
				if xhr.status is HttpResponseCodes.OK
					{
					.updateMsg(msgEl, 'success', color: 'green')
					.getSaveName(method, url, file, xhr)
					}
				else if xhr.status is HttpResponseCodes.BadRequest
					.updateMsg(msgEl, xhr.response, color: 'red', file: file.name)
				// 0: cors preflight request failed, like permission or credential expired
				else if xhr.status is 0
					.updateMsg(msgEl, 'upload failed', color: 'red', file: file.name)

				if SuUI.GetCurrentWindow().GetFileCount() is .fileCount
					{
					SuUI.GetCurrentWindow().ClearFileCount()
					results = .uploadTasks.Values().Filter({
						it.Member?('saveName') }).Map({ it.saveName })

					if not results.Empty?()
						.Event(#UploadFinished, results)
					}
				}
			})
		xhr.Open(method, url)
		ext = file.name.AfterLast('.').Lower()
		xhr.SetRequestHeader('Content-Type',
			MimeTypes.GetDefault(ext, 'application/octet-stream'))
		xhr.Send(file)
		}

	getSaveName(method, url, file, xhr)
		{
		if method is 'PUT'
			{
			filename = SuUI.GetCurrentWindow().Eval('decodeURIComponent("' $
				Url.Split(url).basepath.AfterFirst('/').AfterFirst('/') $ '")')
			.uploadTasks[file.name].saveName = filename
			}
		else
			.uploadTasks[file.name].saveName = xhr.response
		}

	magick?: true
	sendToS3(file, msgEl)
		{
		if file.name.AfterLast('.').Lower() in ('jpg', 'jpeg') and .magick? is true
			.rotateAndSend(file, msgEl)
		else
			.signedUpload(file, msgEl)
		}

	signedUpload(file, msgEl, uploadFile = '')
		{
		suXhr = SuUI.MakeWebObject('XMLHttpRequest')
		if not uploadFile.Blank?()
			.uploadTasks[uploadFile] = Object([xhr: suXhr])
		suXhr.open('POST', "/upload" $
			Url.BuildQuery([file.name, token: SuRender().GetToken(), s3:]))
		suXhr.AddEventListener('readystatechange', { |event/*unused*/|
			if .El isnt false and suXhr.readyState is 4/*=DONE*/
				{
				if suXhr.response isnt 'invalid credential'
					.sendFile(file, msgEl, suXhr.response, 'PUT')
				else
					.updateMsg(msgEl, 'upload failed - invalid credential', color: 'red')
				}
			})
		suXhr.Send()
		}

	rotateAndSend(file, msgEl)
		{
		uploadFn = .signedUpload
		file.ArrayBuffer().Then({|arrayBuffer|
			sourceBytes = SuUI.GetCurrentWindow().Uint8Array(arrayBuffer)
			files = Object(Object( name: 'src.jpg', content: sourceBytes ))
			// \\n is required for stdout
			orientCmd = ['identify', '-format', '"%[EXIF:Orientation]\\n"', 'src.jpg']
			SuUI.GetCurrentWindow().Magick(files, orientCmd).Then(
				{|orientRes|
				stdout = orientRes.stdout[0].Tr('"')
				if stdout not in ('1','')
					{
					command = ["convert", "src.jpg", "-auto-orient", file.name]
					SuUI.GetCurrentWindow().Magick(files, command).Then(
						{ |result|
						if result.exitCode is 0
							{
							output = result.outputFiles[0]
							resultFile = SuUI.GetCurrentWindow().File(output, file.type)
							}
						else
							resultFile = file
						uploadFn(resultFile, msgEl, file.name)
						})
					}
				else
					{
					uploadFn(file, msgEl, file.name)
					}
				}).Catch(
					{|unused|
					uploadFn(file, msgEl, file.name)
					})
			})
		}

	// Not use ReadableSize because Number() is not supported in Suneido.js
	formatSize(n)
		{
		amountPerUnit = 1024
		for unit in #('', kb, mb, gb, tb)
			{
			if n < amountPerUnit
				return n.Round(1) $ Opt(' ', unit)
			n /= amountPerUnit
			}
		return n.Round(2) $ ' pb'
		}

	updateMsg(el, msg, color = 'black', file = '')
		{
		if .El is false
			return

		el.innerText = msg
		el.SetStyle('color', color)
		if color is 'red'
			{
			.loader.SetStyle('display', 'none')
			.abort(file)
			}
		}

	uploadTasks: #()
	abort(file = '')
		{
		if not file.Blank?()
			{
			.uploadTasks[file][0].xhr.Abort()
			return
			}

		for task in .uploadTasks
			if not task.Member?(#saveName)
				{
				task[0].xhr.Abort()
				}
		}

	Destroy()
		{
		.abort()
		super.Destroy()
		}
	}
