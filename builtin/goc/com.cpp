// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

#define WIN32_LEAN_AND_MEAN
#include <windows.h>
#include <objbase.h>
#include <malloc.h>
#include "cside.h"

extern "C" void release(uintptr iunk) {
	if (iunk)
		((IUnknown*) iunk)->Release();
}

extern "C" uintptr queryIDispatch(uintptr iunk) {
	if (iunk == 0) {
		return 0;
	}
	IDispatch* idisp = nullptr;
	HRESULT hr =
		((IUnknown*) iunk)->QueryInterface(IID_IDispatch, (void**) &idisp);
	if (FAILED(hr) || idisp == 0) {
		return 0;
	}
	release(iunk);
	return (uintptr) idisp;
}

static uintptr invoke2(IDispatch* idisp, char* name, WORD flags,
	DISPPARAMS* args, VARIANT* result) {
	if (idisp == 0)
		return -1;
		
	// get id from name
	int n = MultiByteToWideChar(CP_ACP, 0, name, -1, NULL, 0);
	LPWSTR wname = (LPWSTR) _alloca(n * 2);
	MultiByteToWideChar(CP_ACP, 0, name, -1, wname, n);
	DISPID dispid;
	HRESULT hr = idisp->GetIDsOfNames(
		IID_NULL, &wname, 1, LOCALE_SYSTEM_DEFAULT, &dispid);
	if (FAILED(hr))
		return hr;

	// convert BSTR args
	for (int i = 0; i < args->cArgs; i++) {
		VARIANT* v = &args->rgvarg[i];
		if (V_VT(v) == VT_BSTR) {
			char* s = (char*) V_BSTR(v);
			int n = MultiByteToWideChar(CP_ACP, 0, s, -1, NULL, 0);
			BSTR bs = ::SysAllocStringLen(NULL, n);
			MultiByteToWideChar(CP_ACP, 0, s, -1, bs, n);
			V_BSTR(v) = bs;
		}
	}

	hr = idisp->Invoke(dispid, IID_NULL, LOCALE_SYSTEM_DEFAULT, flags, args,
		result, NULL, NULL);

	// free BSTR args
	for (int i = 0; i < args->cArgs; i++) {
		VARIANT* v = &args->rgvarg[i];
		if (V_VT(v) == VT_BSTR && V_BSTR(v) != 0)
			SysFreeString(V_BSTR(v));
	}
	return hr;
}

extern "C" uintptr invoke(
	uintptr idisp, uintptr name, uintptr flags, uintptr args, uintptr result) {
	return invoke2((IDispatch*) idisp, (char*) name, (WORD) flags,
		(DISPPARAMS*) args, (VARIANT*) result);
}
