// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
// simplest possible window using only win32 (no stdlib)
function ()
	{
	wc = Object(instance:  Instance(),
		className: "myclass",
		wndProc: GetDefWindowProc())
	if 0 is RegisterClass(wc)
		Print("RegisterClass failed")

	hwnd = CreateWindowEx(
		0, "myclass",
		"mytitle",
		WS.OVERLAPPEDWINDOW,
		500, 500, 500, 500,
		0, 0, Instance(), NULL)
	if hwnd is 0
		throw "CreateWindow failed"

	ShowWindow(hwnd, SW.SHOW)
	}