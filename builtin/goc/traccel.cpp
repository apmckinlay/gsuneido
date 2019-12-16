// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

#undef UNICODE
#undef _UNICODE
#define WIN32_LEAN_AND_MEAN
#include <objbase.h>
#include <cstdio>

typedef unsigned int uint32;
typedef unsigned long long uintptr;

// Used to forward keyboard messages to a browser control
// to get Tab etc. to work.
extern "C"
long traccel(uintptr ob, uintptr msg) {
	if (!ob)
		return S_FALSE;
	void* p = reinterpret_cast<void*>(ob);
	auto iunk = static_cast<IUnknown*>(p);
	IOleInPlaceActiveObject* pi;
	HRESULT hr = iunk->QueryInterface(
		IID_IOleInPlaceActiveObject, reinterpret_cast<void**>(&pi));
	if (!SUCCEEDED(hr) || !pi)
		return S_FALSE;
	hr = pi->TranslateAcceleratorA(reinterpret_cast<MSG*>(msg));
	pi->Release();
	return hr;
}
