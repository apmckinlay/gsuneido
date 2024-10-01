#ifndef suneido_webview2_environment_options
#define suneido_webview2_environment_options

#include "WebView2.h"
#include "helpers.h"

#define CORE_WEBVIEW_TARGET_PRODUCT_VERSION L"122.0.2357.0"

// Custom implementation of ICoreWebView2EnvironmentOptions
class MyCoreWebView2EnvironmentOptions : 
    public ICoreWebView2EnvironmentOptions,
    public ICoreWebView2EnvironmentOptions2,
    public ICoreWebView2EnvironmentOptions3,
    public ICoreWebView2EnvironmentOptions4,
    public ICoreWebView2EnvironmentOptions5,
    public ICoreWebView2EnvironmentOptions6 {
public:
    MyCoreWebView2EnvironmentOptions(): m_refCount(1) {
        m_TargetCompatibleBrowserVersion.Set(CORE_WEBVIEW_TARGET_PRODUCT_VERSION);
    }

    virtual ~MyCoreWebView2EnvironmentOptions() {
        ReleaseCustomSchemeRegistrations();
    }

    // IUnknown methods
    HRESULT STDMETHODCALLTYPE QueryInterface(REFIID riid, void** ppvObject) override {
        if (ppvObject == nullptr) {
            return E_POINTER;
        }

        if (riid == IID_IUnknown) {
            *ppvObject = static_cast<ICoreWebView2EnvironmentOptions*>(this);
        } else if (riid == IID_ICoreWebView2EnvironmentOptions) {
            *ppvObject = static_cast<ICoreWebView2EnvironmentOptions*>(this);
        } else if (riid == IID_ICoreWebView2EnvironmentOptions2) {
            *ppvObject = static_cast<ICoreWebView2EnvironmentOptions2*>(this);
        } else if (riid == IID_ICoreWebView2EnvironmentOptions3) {
            *ppvObject = static_cast<ICoreWebView2EnvironmentOptions3*>(this);
        } else if (riid == IID_ICoreWebView2EnvironmentOptions4) {
            *ppvObject = static_cast<ICoreWebView2EnvironmentOptions4*>(this);
        } else if (riid == IID_ICoreWebView2EnvironmentOptions5) {
            *ppvObject = static_cast<ICoreWebView2EnvironmentOptions5*>(this);
        } else if (riid == IID_ICoreWebView2EnvironmentOptions6) {
            *ppvObject = static_cast<ICoreWebView2EnvironmentOptions6*>(this);
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

    // ICoreWebView2EnvironmentOptions methods
    public: 
        HRESULT STDMETHODCALLTYPE get_AdditionalBrowserArguments(LPWSTR* value) override {
            if (!value) 
                return E_POINTER; 
            *value = m_AdditionalBrowserArguments.Copy(); 
            if ((*value == nullptr) && (m_AdditionalBrowserArguments.Get() != nullptr)) 
                return HRESULT_FROM_WIN32(GetLastError()); 
            return S_OK;
        } 
        HRESULT STDMETHODCALLTYPE put_AdditionalBrowserArguments(LPCWSTR value) override {
            LPCWSTR result = m_AdditionalBrowserArguments.Set(value); 
            if ((result == nullptr) && (value != nullptr)) 
                return HRESULT_FROM_WIN32(GetLastError()); 
            return S_OK;
        } 
    protected: 
        MyString m_AdditionalBrowserArguments;

    public: 
        HRESULT STDMETHODCALLTYPE get_Language(LPWSTR* value) override {
            if (!value) 
                return E_POINTER; 
            *value = m_Language.Copy(); 
            if ((*value == nullptr) && (m_Language.Get() != nullptr)) 
                return HRESULT_FROM_WIN32(GetLastError()); 
            return S_OK;
        } 
        HRESULT STDMETHODCALLTYPE put_Language(LPCWSTR value) override {
            LPCWSTR result = m_Language.Set(value); 
            if ((result == nullptr) && (value != nullptr)) 
                return HRESULT_FROM_WIN32(GetLastError()); 
            return S_OK;
        } 
    protected: 
        MyString m_Language;

    public: 
        HRESULT STDMETHODCALLTYPE get_TargetCompatibleBrowserVersion(LPWSTR* value) override {
            if (!value) 
                return E_POINTER; 
            *value = m_TargetCompatibleBrowserVersion.Copy(); 
            if ((*value == nullptr) && (m_TargetCompatibleBrowserVersion.Get() != nullptr)) 
                return HRESULT_FROM_WIN32(GetLastError()); 
            return S_OK;
        } 
        HRESULT STDMETHODCALLTYPE put_TargetCompatibleBrowserVersion(LPCWSTR value) override {
            LPCWSTR result = m_TargetCompatibleBrowserVersion.Set(value); 
            if ((result == nullptr) && (value != nullptr)) 
                return HRESULT_FROM_WIN32(GetLastError()); 
            return S_OK;
        } 
    protected: 
        MyString m_TargetCompatibleBrowserVersion;

    public: 
        HRESULT STDMETHODCALLTYPE get_AllowSingleSignOnUsingOSPrimaryAccount(BOOL* value) override {
            if (!value) 
                return E_POINTER; 
            *value = m_AllowSingleSignOnUsingOSPrimaryAccount; 
            return S_OK;
        } 
        HRESULT STDMETHODCALLTYPE put_AllowSingleSignOnUsingOSPrimaryAccount(BOOL value) override {
            m_AllowSingleSignOnUsingOSPrimaryAccount = value; 
            return S_OK;
        } 
    protected: 
        BOOL m_AllowSingleSignOnUsingOSPrimaryAccount = false;

    // ICoreWebView2EnvironmentOptions2
    public: 
        HRESULT STDMETHODCALLTYPE get_ExclusiveUserDataFolderAccess(BOOL* value) override {
            if (!value) 
                return E_POINTER; 
            *value = m_ExclusiveUserDataFolderAccess; 
            return S_OK;
        } 
        HRESULT STDMETHODCALLTYPE put_ExclusiveUserDataFolderAccess(BOOL value) override {
            m_ExclusiveUserDataFolderAccess = value; 
            return S_OK;
        } 
    protected: 
        BOOL m_ExclusiveUserDataFolderAccess = false;

    // ICoreWebView2EnvironmentOptions3
    public: 
        HRESULT STDMETHODCALLTYPE get_IsCustomCrashReportingEnabled(BOOL* value) override {
            if (!value) 
                return E_POINTER; 
            *value = m_IsCustomCrashReportingEnabled; 
            return S_OK;
        } 
        HRESULT STDMETHODCALLTYPE put_IsCustomCrashReportingEnabled(BOOL value) override {
            m_IsCustomCrashReportingEnabled = value; 
            return S_OK;
        } 
    protected: 
        BOOL m_IsCustomCrashReportingEnabled = false;

    // ICoreWebView2EnvironmentOptions4
     public:
        HRESULT STDMETHODCALLTYPE GetCustomSchemeRegistrations(UINT32* count,
            ICoreWebView2CustomSchemeRegistration*** schemeRegistrations) override {
            if (!count || !schemeRegistrations) {
                return E_POINTER;
            }
            *count = 0;
            if (m_customSchemeRegistrationsCount == 0) {
                *schemeRegistrations = nullptr;
                return S_OK;
            } else {
                *schemeRegistrations =reinterpret_cast<ICoreWebView2CustomSchemeRegistration**>(
                    CoTaskMemAlloc(sizeof(ICoreWebView2CustomSchemeRegistration*) *m_customSchemeRegistrationsCount));
                if (!*schemeRegistrations) {
                    return HRESULT_FROM_WIN32(GetLastError());
                }
                for (UINT32 i = 0; i < m_customSchemeRegistrationsCount; i++) {
                    (*schemeRegistrations)[i] = m_customSchemeRegistrations[i];
                    (*schemeRegistrations)[i]->AddRef();
                }
                *count = m_customSchemeRegistrationsCount;
                return S_OK;
            }
        }
        HRESULT STDMETHODCALLTYPE SetCustomSchemeRegistrations(UINT32 count,
            ICoreWebView2CustomSchemeRegistration** schemeRegistrations) override {
            ReleaseCustomSchemeRegistrations();
            m_customSchemeRegistrations = reinterpret_cast<ICoreWebView2CustomSchemeRegistration**>(CoTaskMemAlloc(
                sizeof(ICoreWebView2CustomSchemeRegistration*) * count));
            if (!m_customSchemeRegistrations) {
                return GetLastError();
            }
            for (UINT32 i = 0; i < count; i++) {
                m_customSchemeRegistrations[i] = schemeRegistrations[i];
                m_customSchemeRegistrations[i]->AddRef();
            }
            m_customSchemeRegistrationsCount = count;
            return S_OK;
        }
    
    protected:
        void ReleaseCustomSchemeRegistrations() {
            if (m_customSchemeRegistrations) {
                for (UINT32 i = 0; i < m_customSchemeRegistrationsCount; i++) {
                    m_customSchemeRegistrations[i]->Release();
                }
                CoTaskMemFree(m_customSchemeRegistrations);
                m_customSchemeRegistrations = nullptr;
                m_customSchemeRegistrationsCount = 0;
            }
        }

    // ICoreWebView2EnvironmentOptions5
    public: 
        HRESULT STDMETHODCALLTYPE get_EnableTrackingPrevention(BOOL* value) override {
            if (!value) 
                return E_POINTER; 
            *value = m_EnableTrackingPrevention; 
            return S_OK;
        } 
        HRESULT STDMETHODCALLTYPE put_EnableTrackingPrevention(BOOL value) override {
            m_EnableTrackingPrevention = value; 
            return S_OK;
        } 
    protected: 
        BOOL m_EnableTrackingPrevention = true;

    // ICoreWebView2EnvironmentOptions6
    public: 
        HRESULT STDMETHODCALLTYPE get_AreBrowserExtensionsEnabled(BOOL* value) override {
            if (!value) 
                return E_POINTER; 
            *value = m_AreBrowserExtensionsEnabled; 
            return S_OK;
        } 
        HRESULT STDMETHODCALLTYPE put_AreBrowserExtensionsEnabled(BOOL value) override {
            m_AreBrowserExtensionsEnabled = value; 
            return S_OK;
        } 
    protected: 
        BOOL m_AreBrowserExtensionsEnabled = false;

    private:
        ULONG m_refCount;
        ICoreWebView2CustomSchemeRegistration** m_customSchemeRegistrations = nullptr;
        unsigned int m_customSchemeRegistrationsCount = 0;
};

#endif