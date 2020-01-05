// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

#undef UNICODE

#include <windows.h>
#include <wininet.h>
#include <malloc.h>
#include "cside.h"

extern "C" buf_t suneidoAPP(char* s);

class CSuneidoAPP : public IInternetProtocol {
public:
	CSuneidoAPP() {
	}

	virtual ~CSuneidoAPP() {
		delete str;
	}

	// IUnknown
	virtual HRESULT __stdcall QueryInterface(const IID& iid, void** ppv);
	virtual ULONG __stdcall AddRef();
	virtual ULONG __stdcall Release();

	// IInternetProtocolRoot

	HRESULT STDMETHODCALLTYPE Start(
		/* [in] */ LPCWSTR szUrl,
		/* [in] */ IInternetProtocolSink* pOIProtSink,
		/* [in] */ IInternetBindInfo* pOIBindInfo,
		/* [in] */ DWORD grfPI,
		/* [in] */ HANDLE_PTR dwReserved) override;

	HRESULT STDMETHODCALLTYPE Continue(
		/* [in] */ PROTOCOLDATA __RPC_FAR* pProtocolData) {
		return S_OK;
	}

	HRESULT STDMETHODCALLTYPE Abort(
		/* [in] */ HRESULT hrReason,
		/* [in] */ DWORD dwOptions) {
		return S_OK; // or E_NOTIMPL ???
	}

	HRESULT STDMETHODCALLTYPE Terminate(
		/* [in] */ DWORD dwOptions) {
		return S_OK;
	}

	HRESULT STDMETHODCALLTYPE Suspend(void) {
		return E_NOTIMPL;
	} // Not implemented

	HRESULT STDMETHODCALLTYPE Resume(void) {
		return E_NOTIMPL;
	} // Not implemented

	// IInternetProtocolRoot

	HRESULT STDMETHODCALLTYPE Read(
		/* [length_is][size_is][out][in] */ void __RPC_FAR* pv,
		/* [in] */ ULONG cb,
		/* [out] */ ULONG __RPC_FAR* pcbRead);

	HRESULT STDMETHODCALLTYPE Seek(
		/* [in] */ LARGE_INTEGER dlibMove,
		/* [in] */ DWORD dwOrigin,
		/* [out] */ ULARGE_INTEGER __RPC_FAR* plibNewPosition) {
		return E_NOTIMPL;
	}

	HRESULT STDMETHODCALLTYPE LockRequest(
		/* [in] */ DWORD dwOptions) {
		return S_OK;
	}

	HRESULT STDMETHODCALLTYPE UnlockRequest(void) {
		return S_OK;
	}

private:
	ULONG len = 0;
	ULONG pos = 0;
	const char* str = nullptr;
	long m_cRef = 1; // Reference count
};

// IUnknown implementation

HRESULT __stdcall CSuneidoAPP::QueryInterface(const IID& iid, void** ppv) {
	if (iid == IID_IUnknown || iid == IID_IInternetProtocol) {
		*ppv = static_cast<IInternetProtocol*>(this);
	} else {
		*ppv = NULL;
		return E_NOINTERFACE;
	}
	reinterpret_cast<IUnknown*>(*ppv)->AddRef();
	return S_OK;
}

ULONG __stdcall CSuneidoAPP::AddRef() {
	return InterlockedIncrement(&m_cRef);
}

ULONG __stdcall CSuneidoAPP::Release() {
	if (InterlockedDecrement(&m_cRef) == 0) {
		delete this;
		return 0;
	}
	return m_cRef;
}

// IInternetProtocol

#define USES_CONVERSION \
	int _convert; \
	_convert; \
	LPCWSTR _lpw; \
	_lpw;
inline LPSTR WINAPI AtlW2AHelper(LPSTR lpa, LPCWSTR lpw, int nChars) {
	lpa[0] = '\0';
	WideCharToMultiByte(CP_ACP, 0, lpw, -1, lpa, nChars, NULL, NULL);
	return lpa;
}
#define W2CA(w) \
	((LPCSTR)(((_lpw = (w)) == NULL) \
			? NULL \
			: (_convert = (lstrlenW(_lpw) + 1) * 2, \
				  AtlW2AHelper((LPSTR) _alloca(_convert), _lpw, _convert))))

HRESULT STDMETHODCALLTYPE CSuneidoAPP::Start(
	/* [in] */ LPCWSTR szUrl,
	/* [in] */ IInternetProtocolSink* pOIProtSink,
	/* [in] */ IInternetBindInfo* pOIBindInfo,
	/* [in] */ DWORD grfPI,
	/* [in] */ HANDLE_PTR dwReserved) {
	USES_CONVERSION;
	const char* url = W2CA(szUrl);
	ULONG buflen = strlen(url) + 10; // shouldn't get any bigger?
	char* buf = new char[buflen];
	InternetCanonicalizeUrl(url, buf, &buflen, ICU_DECODE | ICU_NO_ENCODE);

	buf_t result = suneidoAPP(buf);
	str = result.buf;
	len = result.size;
	pos = 0;

	pOIProtSink->ReportData(BSCF_DATAFULLYAVAILABLE | BSCF_LASTDATANOTIFICATION,
		len, len); // MUST call this
	return S_OK;
}

#define min(x, y) ((x) < (y) ? (x) : (y))

HRESULT STDMETHODCALLTYPE CSuneidoAPP::Read(
	/* [length_is][size_is][out][in] */ void __RPC_FAR* pv,
	/* [in] */ ULONG cb,
	/* [out] */ ULONG __RPC_FAR* pcbRead) {
	*pcbRead = min(cb, len - pos);
	memcpy(pv, str + pos, *pcbRead);
	pos += *pcbRead;
	return pos >= len ? S_FALSE : S_OK; // more to read ?
}

//-------------------------------------------------------------------

class CFactory : public IClassFactory {
public:
	// IUnknown
	virtual HRESULT __stdcall QueryInterface(const IID& iid, void** ppv);
	virtual ULONG __stdcall AddRef();
	virtual ULONG __stdcall Release();

	// Interface IClassFactory
	virtual HRESULT __stdcall CreateInstance(
		IUnknown* pUnknownOuter, const IID& iid, void** ppv);
	virtual HRESULT __stdcall LockServer(BOOL bLock);

	virtual ~CFactory() = default;

private:
	long m_cRef = 1;
};

HRESULT __stdcall CFactory::QueryInterface(const IID& iid, void** ppv) {
	if ((iid == IID_IUnknown) || (iid == IID_IClassFactory)) {
		*ppv = static_cast<IClassFactory*>(this);
	} else {
		*ppv = NULL;
		return E_NOINTERFACE;
	}
	reinterpret_cast<IUnknown*>(*ppv)->AddRef();
	return S_OK;
}

ULONG __stdcall CFactory::AddRef() {
	return InterlockedIncrement(&m_cRef);
}

ULONG __stdcall CFactory::Release() {
	if (InterlockedDecrement(&m_cRef) == 0) {
		delete this;
		return 0;
	}
	return m_cRef;
}

HRESULT __stdcall CFactory::CreateInstance(
	IUnknown* pUnknownOuter, const IID& iid, void** ppv) {
	if (pUnknownOuter != NULL) {
		return CLASS_E_NOAGGREGATION;
	}

	CSuneidoAPP* pA = new CSuneidoAPP;
	if (pA == NULL) {
		return E_OUTOFMEMORY;
	}

	// Get the requested interface.
	HRESULT hr = pA->QueryInterface(iid, ppv);

	// Release the IUnknown pointer.
	// (If QueryInterface failed, component will delete itself.)
	pA->Release();
	return hr;
}

HRESULT __stdcall CFactory::LockServer(BOOL bLock) {
	return S_OK;
}

//-------------------------------------------------------------------

const CLSID CLSID_SuneidoAPP = {0xBFBE2090, 0x6BBA, 0x11D4,
	{0xBC, 0x13, 0x00, 0x60, 0x6E, 0x30, 0xB2, 0x58}};
const char* clsid = "{BFBE2090-6BBA-11D4-BC13-00606E30B258}";

static CFactory factory;

extern "C"
int sunapp_register_classes() {
	HRESULT hr;
	hr = CoInitializeEx(NULL, COINIT_APARTMENTTHREADED);
	if (FAILED(hr))
		return false;

	IInternetSession* iis;
	hr = CoInternetGetSession(0, &iis, 0);
	if (FAILED(hr))
		return false;
	else {
		iis->RegisterNameSpace(&factory, CLSID_SuneidoAPP, L"suneido", 0, 0, 0);
		iis->Release();
	}
	return true;
}

extern "C"
void sunapp_revoke_classes() {
	CoUninitialize();
}
