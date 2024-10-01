#include "MyWebView2Base.h"

class MyContextMenuRequestedEventHandler : public ICoreWebView2ContextMenuRequestedEventHandler {
public:
    // Constructor
    MyContextMenuRequestedEventHandler(MyBrowserObject* pBrowserObject) : m_refCount(1), pBrowserObject(pBrowserObject) {
        pBrowserObject->AddRef();
    }

    ~MyContextMenuRequestedEventHandler() {
        pBrowserObject->Release();
    }

    // IUnknown methods
    STDMETHODIMP QueryInterface(REFIID riid, void** ppvObject) override {
        if (riid == IID_IUnknown || riid == IID_ICoreWebView2ContextMenuRequestedEventHandler) {
            *ppvObject = static_cast<ICoreWebView2ContextMenuRequestedEventHandler*>(this);
            AddRef();
            return S_OK;
        }
        *ppvObject = nullptr;
        return E_NOINTERFACE;
    }

    STDMETHODIMP_(ULONG) AddRef() override {
        return InterlockedIncrement(&m_refCount);
    }

    STDMETHODIMP_(ULONG) Release() override {
        ULONG refCount = InterlockedDecrement(&m_refCount);
        if (refCount == 0) {
            delete this;
        }
        return refCount;
    }

    // ICoreWebView2ContextMenuRequestedEventHandler method
    STDMETHODIMP Invoke(ICoreWebView2 *sender, ICoreWebView2ContextMenuRequestedEventArgs *args) override {
        if (pBrowserObject->onContextMenuRequested() == false) {
            args->put_Handled(true);
        }
        return S_OK;
    }

private:
    // Reference count for COM object
    ULONG m_refCount;
    MyBrowserObject* pBrowserObject;
};