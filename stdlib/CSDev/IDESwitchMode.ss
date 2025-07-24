// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CSDevCmd: `SystemNoWait('csdev.bat');Exit();`
	CallClass(standalone = false, defaultPort = false)
		{
		if standalone
			{
			cmd = 'gsuneido.exe'
			DeleteFile('suneido.args')
			}
		else
			{
			cmd = 'csdev.bat'
			port = defaultPort
				? 3147 /*= default port */
				: 4000 + Random(6000) /*= from 4000 to 10000, to avoid default 3147*/
			for c in Contributions('IDESwitchMode_ExtraSetups')
				c(:port)
			PutFile('suneido.args', .CSDevCmd)
			PutFile('csdev.bat', 'title Start Suneido Server
gsport.exe -s -p ' $ port $ ' CSDevServerWindow()
exit')
			if DirExists?(ApplicationDir())
				PutFile(Paths.Combine(ApplicationDir(), 'ports.config'),
					'http:' $ String(port + 100)) /*= http offset*/
			}
		if false isnt hwnd = ServerSuneido.Get('DevHwnd', false)
			PostMessage(hwnd, WM.DESTROY, 0, 0)
		PutFile('gdev_switchmode.bat', 'sleep 2
start "" ' $ cmd $ '
start "" cmd /c del gdev_switchmode.bat
exit')
		SystemNoWait('gdev_switchmode.bat')
		Shutdown(alsoServer:)
		}
	}