// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.

// This is based on the Microsoft WebView2 API
// https://learn.microsoft.com/en-us/microsoft-edge/webview2/webview2-api-reference?tabs=win32cpp
class
	{
	webBrowser: false
	New(.parent)
		{
		if String?(result = .createWebView2())
			{
			.Release()
			throw 'WebView: ' $ result
			}

		.afterReady = Object()
		.afterLoad = Object()
		.afterNavComplete = Object()
		.webBrowser = result
		}

	callback: false
	createWebView2()
		{
		dllPath = .GetLoaderDll()
		result = dllPath is false
			? -1 /*=error code*/
			: (Global('WebBrowser2'))(.parent.Hwnd, dllPath, .getUserDataFolder(),
				.callback = .cb)
		return .checkResult(result)
		}

	ready?: false
	onReady(stage, code)
		{
		if .webBrowser is false
			return 0

		if stage isnt 0 /*= failed*/
			{
			.Release()
			.retry(stage, code)
			return 0
			}

		if .retries.NotEmpty?()
			{
			ErrorLog('INFO: WebView2 initialization succeeds after retries: ' $
				Display(.retries))
			}

		.ready? = true
		for task in .afterReady
			(task.block)()
		.afterReady = Object()

		if .parent.HasFocus?()
			.SetFocus()

		return 0
		}

	delay: false
	retry(stage, code)
		{
		.retries.Add(Object(:stage, :code))

		if .retries.Size() > 3/*=max retries*/
			{
			.onInitFailure('Failed to initialize WebView2 (max retry reached)',
				params: .retries)
			return
			}

		.delay = Delay((1 << .retries.Size()).SecondsInMs())
			{
			Suneido.Delete(#WebView2_UDF) // use a new User Data Folder
			if String?(result = .createWebView2())
				{
				ProgrammerError('Failed to initialize WebView2 - ' $ result,
					params: .retries)
				}
			else
				.webBrowser = result
			}
		}

	onInitFailure(msg, params)
		{
		// Cannot use an icon because icon is not available before login
		Alert("Can't create WebView.

You may need to restart your computer.", "Error")
		ProgrammerError(msg, params)
		}

	getter_retries()
		{
		return .retries = Object()
		}

	cbTypes: (ON_READY: 0, ON_LOADED: 1, ON_ACCEL_KEY_PRESSED: 2,
		ON_CONTEXT_MENU_REQUESTED: 3, ON_NAVCOMPLETED: 4)
	cb(type, arg1, arg2, arg3/*unused*/)
		{
		switch (type)
			{
		case .cbTypes.ON_READY:
			return .onReady(arg1, arg2)
		case .cbTypes.ON_LOADED:
			return .onLoaded()
		case .cbTypes.ON_ACCEL_KEY_PRESSED:
			return .onAccel(arg1)
		case .cbTypes.ON_CONTEXT_MENU_REQUESTED:
			return .onContextMenuRequested()
		case .cbTypes.ON_NAVCOMPLETED:
			return .onNavCompleted()
			}
		}

	loaded?: false
	onLoaded()
		{
		if .webBrowser is false
			return 0

		.loaded? = true
		for task in .afterLoad
			(task.block)()
		.afterLoad = Object()
		return 0
		}

	OnNavComplete(block)
		{
		.waitFor(false, .afterNavComplete, block, false)
		}

	onNavCompleted()
		{
		if .webBrowser is false
			return 0
		for task in .afterNavComplete
			(task.block)()
		.afterNavComplete = Object()
		return 0
		}

	skipAccel?(key)
		{
		return KeyPressed?(VK.CONTROL) and key in (VK.C, VK.F)
		}
	onAccel(key)
		{
		if .skipAccel?(key)
			return false

		if false isnt cmd = .parent.Window.QueryAccel(key,
			ctrl: KeyPressed?(VK.CONTROL), alt: KeyPressed?(VK.MENU),
			shift: KeyPressed?(VK.SHIFT))
			{
			PostMessage(.parent.Window.Hwnd, WM.COMMAND, cmd, NULL)
			return true
			}
		return false
		}

	// return false to disable the context menu
	onContextMenuRequested()
		{
		return true
		}

	doAfterReady(block, id = false)
		{
		.waitFor(.ready?, .afterReady, block, id)
		}

	doAfterLoaded(block, id = false)
		{
		.waitFor(.loaded?, .afterLoad, block, id)
		}

	waitFor(flag, tasks, block, id)
		{
		if .webBrowser isnt false and flag
			{
			block()
			return
			}

		if id isnt false
			tasks.RemoveIf({ it.id is id })
		tasks.Add(Object(:id, :block))
		}

	Load(what)
		{
		.doAfterReady(id: #load)
			{
			if what isnt false
				{
				if what.Prefix?('MSHTML:')
					.webBrowser.NavigateToString(what.RemovePrefix('MSHTML:'))
				else
					.webBrowser.Navigate(.ensureURLEncoded(what))
				.loaded? = false
				}
			}
		}

	ensureURLEncoded(url)
		{
		if url isnt Url.Decode(url)
			return url
		return Url.Encode(url)
		}

	Ready?()
		{
		return .ready?
		}

	Resize(w, h)
		{
		.doAfterReady({ .webBrowser.Resize(w, h) }, id: #resize)
		}

	SetFocus()
		{
		.doAfterReady({ .webBrowser.SetFocus() }, id: #focus)
		}

	SetCssStyle(style)
		{
		.executeScript(`
var style = document.createElement("style");
style.textContent = "` $ style.Tr('\r\n') $ `";
document.body.appendChild(style);`)
		}

	executeScript(script)
		{
		.doAfterLoaded({ .webBrowser.ExecuteScript(script) })
		}

	Getter_LocationURL()
		{
		if .webBrowser is false
			return ''
		return .webBrowser.GetSource()
		}

	TriggerKeyDown(key)
		{
		.executeScript("document.dispatchEvent(
			new KeyboardEvent('keydown', { keyCode: " $ key $ " }))")
		}

	DoFind()
		{
		}

	DoGoBack()
		{
		.executeScript("window.history.back()")
		}

	DoGoForward()
		{
		.executeScript("window.history.forward()")
		}

	DoCopy()
		{
		}

	DoPaste()
		{
		}

	DoRefresh()
		{
		.executeScript('location.reload')
		}

	DoPrint()
		{
		.doAfterLoaded({ .webBrowser.Print() })
		}

	DoPrintPreview()
		{
		}

	DoPageSetup()
		{
		}

	InsertAdjacentHTML(id, position, text)
		{
		script = `parent = document.getElementById(` $ Display(id) $ `);
parent.insertAdjacentHTML(` $ Display(position) $ `, ` $ Display(text) $ `);`
		.executeScript(script)
		}

	ScrollIntoView(id, alignToTop)
		{
		script = `el = document.getElementById(` $ Display(id) $`);
el.scrollIntoView(` $ Display(alignToTop) $ `)`
		.executeScript(script)
		}

	results: (
		1: 'WebView2 not ready'
		2: 'Invalid operation'
		3: 'DLL not found'
		4: 'Create WebView2 proc not found'
		-1: 'DLL not found in Registry')
	checkResult(result)
		{
		if result is 3 /*=DLL not found*/
			{
			.Release()
			.retry(-1, 3 /*=DLL not found*/)
			return false
			}

		if Number?(result)
			return .results.GetDefault(result, 'Error Code (' $ result $ ')')
		return result
		}

	Release()
		{
		if .webBrowser isnt false
			{
			.webBrowser.Release()
			.webBrowser = false
			}
		if .delay isnt false
			{
			.delay.Kill()
			.delay = false
			}
		if .callback isnt false
			{
			ClearCallback(.callback)
			.callback = false
			}
		}

	folderPrefix: 'SuWebView2_'
	getUserDataFolder()
		{
		return Suneido.GetInit(#WebView2_UDF,
			{ Paths.ToLocal(Paths.Combine(GetTempPath(),
				.folderPrefix $ Display(Date())[1..].Tr('.', '_') $ '_' $ Random())) })
		}

	CleanUp()
		{
		toDelete = Date().NoTime().Minus(days: 3)
		Dir(Paths.Combine(GetTempPath(), .folderPrefix $ '*'))
			{
			date = Date(it.RemovePrefix(.folderPrefix).BeforeFirst('_'),
				format: 'yyyyMMdd')
			if Date?(date) and date <= toDelete
				{
				if true isnt result = DeleteDir(Paths.Combine(GetTempPath(), it))
					SuneidoLog('INFO: WebView2.CleanUp - ' $ result)
				}
			}
		}

	browserKey: `SOFTWARE\WOW6432Node\Microsoft\EdgeUpdate\ClientState\`

	Available?()
		{
		return .GetLoaderDll() isnt false
		}

	cu_key: `HKCU:\SOFTWARE\Microsoft\EdgeUpdate\ClientState\`
	lm_key: `HKLM:\SOFTWARE\WOW6432Node\Microsoft\EdgeUpdate\ClientState\`
	stableReleaseGuid: "{F3017226-FE2A-4295-8BDF-00C3A9A7E4C5}"
	dllPath: `EBWebView\x64\EmbeddedBrowserWebView.dll`
	GetLoaderDll()
		{
		if not Sys.Windows?()
			return false

		if Suneido.Member?(#WebView2Path)
			return Suneido.WebView2Path

		if false isnt path = .queryRegistry(.lm_key)
			return Suneido.WebView2Path = path

		if false isnt path = .queryRegistry(.cu_key)
			return Suneido.WebView2Path = path

		return false
		}

	queryRegistry(key)
		{
		res = RunPipedOutput.WithExitValue(
			PowerShell() $ " Get-ItemPropertyValue -Path " $
			"'" $ key $ .stableReleaseGuid $ "' -Name EBWebView")

		return res.exitValue is 0
			? Paths.ToWindows(Paths.Combine(res.output.Trim(), .dllPath))
			: false
		}

	Install(ensureAdmin? = false)
		{
		PutFile('install_webview2.ps1',
			(ensureAdmin? ? .ensureAdminScript : '') $ .installScript)
		Finally(
			{
			result = RunPipedOutput.WithExitValue(
				PowerShell() $ ' -file install_webview2.ps1')
			},
			{
			if true isnt res = DeleteFile('install_webview2.ps1')
				SuneidoLog('ERROR: (CAUGHT) WebView2 Install - ' $ res)
			})

		if result.exitValue is 10/*=success*/
			return .GetLoaderDll()

		return result
		}

	ensureAdminScript: '
# Function to check if the script is running as an administrator
function Ensure-RunAsAdministrator {
	$currentUser = [Security.Principal.WindowsIdentity]::GetCurrent()
	$principal = New-Object Security.Principal.WindowsPrincipal($currentUser)
	if (-Not $principal.IsInRole(' $
		'[Security.Principal.WindowsBuiltInRole]::Administrator)) {
		Write-Host "The script needs to be run as Administrator. Relaunching..."
		$result = Start-Process powershell ' $
			'"-ExecutionPolicy Bypass -File `"$PSCommandPath`"" ' $
			'-Verb RunAs -Wait -PassThru
		$exitCode = $result.ExitCode
		Exit $exitCode
	}
}

# Ensure the script is running with administrative privileges
Ensure-RunAsAdministrator'

	installScript: '
# Force PowerShell to use TLS 1.2
[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12

# Define the URL for the WebView2 installer
$webview2InstallerUrl = "https://go.microsoft.com/fwlink/p/?LinkId=2124703"

# Define the path where the installer will be downloaded
$installerPath = "$env:TEMP\MicrosoftEdgeWebView2Setup.exe"

# Download the WebView2 installer
Invoke-WebRequest -Uri $webview2InstallerUrl -OutFile $installerPath

# Run the installer silently
Start-Process $installerPath -ArgumentList "/silent /install" -Wait

# Clean up by removing the installer file
Remove-Item $installerPath

Write-Host "Microsoft Edge WebView2 Runtime has been installed successfully."
Exit 10'
	}