#ifndef suneido_webview2_custom_scheme_registration
#define suneido_webview2_custom_scheme_registration

#include "WebView2.h"
#include "helpers.h"

// Custom implementation of ICoreWebView2CustomSchemeRegistration
class MyCoreWebView2CustomSchemeRegistration :
    public ICoreWebView2CustomSchemeRegistration {
public:
    MyCoreWebView2CustomSchemeRegistration(LPCWSTR schemeName): m_refCount(1) {
        m_schemeName.Set(schemeName);
    }
    ~MyCoreWebView2CustomSchemeRegistration() { 
        ReleaseAllowedOrigins(); 
    }

    // IUnknown methods
    HRESULT STDMETHODCALLTYPE QueryInterface(REFIID riid, void** ppvObject) override {
        if (ppvObject == nullptr) {
            return E_POINTER;
        }

        if (riid == IID_IUnknown) {
            *ppvObject = static_cast<ICoreWebView2CustomSchemeRegistration*>(this);
        } else if (riid == IID_ICoreWebView2CustomSchemeRegistration) {
            *ppvObject = static_cast<ICoreWebView2CustomSchemeRegistration*>(this);
        } else {
            *ppvObject = nullptr;
            return E_NOINTERFACE;
        }

        AddRef();
        return S_OK;
    }

    ULONG STDMETHODCALLTYPE AddRef() override {
        return InterlockedIncrement(&m_refCount);
    }

    ULONG STDMETHODCALLTYPE Release() override {
        ULONG count = InterlockedDecrement(&m_refCount);
        if (count == 0) {
            delete this;
        }
        return count;
    }

    HRESULT STDMETHODCALLTYPE get_SchemeName(LPWSTR* schemeName) override {
        if (!schemeName)
            return E_POINTER;
        *schemeName = m_schemeName.Copy();
        if ((*schemeName == nullptr) && (m_schemeName.Get() != nullptr))
            return HRESULT_FROM_WIN32(GetLastError());
        return S_OK;
    }

    HRESULT STDMETHODCALLTYPE GetAllowedOrigins(UINT32* allowedOriginsCount, 
        LPWSTR** allowedOrigins) override {
        if (!allowedOrigins || !allowedOriginsCount) {
            return E_POINTER;
        }
        *allowedOriginsCount = 0;
        if (m_allowedOriginsCount == 0) {
            *allowedOrigins = nullptr;
            return S_OK;
        } else {
            *allowedOrigins = reinterpret_cast<LPWSTR*>(
                CoTaskMemAlloc(m_allowedOriginsCount * sizeof(LPWSTR)));
            if (!(*allowedOrigins)) {
                return HRESULT_FROM_WIN32(GetLastError());
            }
            ZeroMemory(*allowedOrigins, m_allowedOriginsCount * sizeof(LPWSTR));
            for (UINT32 i = 0; i < m_allowedOriginsCount; i++) {
                (*allowedOrigins)[i] = m_allowedOrigins[i].Copy();
                if (!(*allowedOrigins)[i]) {
                    HRESULT hr = HRESULT_FROM_WIN32(GetLastError());
                    for (UINT32 j = 0; j < i; j++) {
                        CoTaskMemFree((*allowedOrigins)[j]);
                    }
                    CoTaskMemFree(*allowedOrigins);
                    return hr;
                }
            }
            *allowedOriginsCount = m_allowedOriginsCount;
            return S_OK;
        }
    }
    HRESULT STDMETHODCALLTYPE SetAllowedOrigins(UINT32 allowedOriginsCount, 
        LPCWSTR* allowedOrigins) override {
        ReleaseAllowedOrigins();
        if (allowedOriginsCount == 0) {
            return S_OK;
        } else {
            m_allowedOrigins = new MyString[allowedOriginsCount];
            if (!m_allowedOrigins) {
                return HRESULT_FROM_WIN32(GetLastError());
            }
            for (UINT32 i = 0; i < allowedOriginsCount; i++) {
                m_allowedOrigins[i].Set(allowedOrigins[i]);
                if (!m_allowedOrigins[i].Get()) {
                    HRESULT hr = HRESULT_FROM_WIN32(GetLastError());
                    for (UINT32 j = 0; j < i; j++) {
                        m_allowedOrigins[j].Release();
                    }
                    m_allowedOriginsCount = 0;
                    delete[] (m_allowedOrigins);
                    return hr;
                }
            }
            m_allowedOriginsCount = allowedOriginsCount;
            return S_OK;
        }
    }

public: 
    HRESULT STDMETHODCALLTYPE get_TreatAsSecure(BOOL* value) override {
        if (!value) 
            return E_POINTER; 
        *value = m_TreatAsSecure; 
        return S_OK;
    } 
    HRESULT STDMETHODCALLTYPE put_TreatAsSecure(BOOL value) override {
        m_TreatAsSecure = value; 
        return S_OK;
    } 
protected: 
    BOOL m_TreatAsSecure = false;

public: 
    HRESULT STDMETHODCALLTYPE get_HasAuthorityComponent(BOOL* value) override {
        if (!value) 
            return E_POINTER; 
        *value = m_HasAuthorityComponent; 
        return S_OK;
    } 
    HRESULT STDMETHODCALLTYPE put_HasAuthorityComponent(BOOL value) override {
        m_HasAuthorityComponent = value; 
        return S_OK;
    } 
protected: 
    BOOL m_HasAuthorityComponent = false;

private:
    ULONG m_refCount;
    MyString m_schemeName;
    MyString* m_allowedOrigins = nullptr;
    unsigned int m_allowedOriginsCount = 0;

    void ReleaseAllowedOrigins() {
        if (m_allowedOrigins) {
            delete[] (m_allowedOrigins);
            m_allowedOrigins = nullptr;
        }
        m_allowedOriginsCount = 0;
    }
};

#endif