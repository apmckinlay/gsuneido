#include "MyWebView2Base.h"

class MyNewWindowRequestedHandler : public ICoreWebView2NewWindowRequestedEventHandler {
public:
    // IUnknown methods
    HRESULT STDMETHODCALLTYPE QueryInterface(REFIID riid, void** ppv) override {
        if (riid == IID_ICoreWebView2NewWindowRequestedEventHandler ||
            riid == IID_IUnknown) {
            *ppv = static_cast<ICoreWebView2NewWindowRequestedEventHandler*>(this);;
            AddRef();
            return S_OK;
        }
        *ppv = nullptr;
        return E_NOINTERFACE;
    }
    ULONG STDMETHODCALLTYPE AddRef(void) override {
        return InterlockedIncrement(&m_refCount);
    }
    ULONG STDMETHODCALLTYPE Release(void) override {
        ULONG count = InterlockedDecrement(&m_refCount);
        if (count == 0) {
            delete this;
        }
        return count;
    }
    MyNewWindowRequestedHandler() : m_refCount(1) {}

    // ICoreWebView2NewWindowRequestedEventHandler method
    HRESULT STDMETHODCALLTYPE Invoke(ICoreWebView2* sender, ICoreWebView2NewWindowRequestedEventArgs* args) override {
        LPWSTR uri = nullptr;
        HRESULT hr = args->get_Uri(&uri);
        if (SUCCEEDED(hr) && uri != nullptr) {
            // Check if the URI starts with "suneido:"
            if (wcslen(uri) >= suneidoPrefixSize && 
                wcsncmp(uri, suneidoPrefix, suneidoPrefixSize) == 0) {
                // Prevent creation of a new window for "suneido:" URLs
                args->put_Handled(TRUE);
            }
            CoTaskMemFree(uri);
        }
        return S_OK;
    }
private:
    LONG m_refCount;
};