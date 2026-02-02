// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CSDevCmd: `SystemNoWait('csdev.bat');Exit();`
	CallClass(standalone = false, defaultPort = false)
		{
		SetCurrentDirectory(ExeDir())
		if standalone
			{
			cmd = 'gsuneido.exe'
			DeleteFile('suneido.args')
			}
		else
			{
			cmd = .CreateFiles(defaultPort)
			}
		if false isnt hwnd = ServerSuneido.Get('DevHwnd', false)
			PostMessage(hwnd, WM.DESTROY, 0, 0)
		PutFile('gdev_switchmode.bat', 'timeout /T 2 /NOBREAK > NUL
start "" ' $ cmd $ '
start "" cmd /c del gdev_switchmode.bat
exit')
		SystemNoWait('gdev_switchmode.bat')
		Shutdown(alsoServer:)
		}

	CreateFiles(defaultPort = false)
		{
		cmd = 'csdev.bat'
		port = Number?(defaultPort)
			? defaultPort
			: defaultPort
				? 3147 /*= default port */
				: 4000 + Random(6000) /*= from 4000 to 10000, to avoid default 3147*/
		for c in Contributions('IDESwitchMode_ExtraSetups')
			c(:port)
		PutFile('suneido.args', .CSDevCmd)
		PutFile('csdev.bat', 'title Start Suneido Server
gsport.exe -s -p ' $ port $ ' CSDevServerWindow() 1>> golang.log 2>&1
exit')
		if DirExists?(ApplicationDir())
			PutFile(Paths.Combine(ApplicationDir(), 'ports.config'),
				'http:' $ String(port + 100)) /*= http offset*/
		return cmd
		}
	}