#include "MyWebView2Base.h"
#include "helpers.h"

#include <shlwapi.h>

extern "C" buf_t suneidoAPP(char* s);

class MyCustomSchemeHandler : public ICoreWebView2WebResourceRequestedEventHandler {
public:
    // Implement IUnknown
    ULONG STDMETHODCALLTYPE AddRef() override {
        return InterlockedIncrement(&refCount);
    }

    ULONG STDMETHODCALLTYPE Release() override {
        ULONG result = InterlockedDecrement(&refCount);
        if (result == 0) {
            delete this;
        }
        return result;
    }

    HRESULT STDMETHODCALLTYPE QueryInterface(REFIID riid, void** ppvObject) override {
        if (riid == IID_IUnknown || riid == IID_ICoreWebView2WebResourceRequestedEventHandler) {
            *ppvObject = static_cast<ICoreWebView2WebResourceRequestedEventHandler*>(this);
            AddRef();
            return S_OK;
        }
        *ppvObject = nullptr;
        return E_NOINTERFACE;
    }

    // Implement ICoreWebView2WebResourceRequestedEventHandler
    HRESULT STDMETHODCALLTYPE Invoke(ICoreWebView2* sender, ICoreWebView2WebResourceRequestedEventArgs* args) override {
        ICoreWebView2WebResourceRequest* request = nullptr;
        HRESULT hr = args->get_Request(&request);
        if (FAILED(hr) || !request) return hr;
        
        LPWSTR uri = nullptr;
        hr = request->get_Uri(&uri); 
        if (FAILED(hr) || !uri) {
            request->Release();
            return hr;
        }

        // Handle the custom scheme
        if (wcslen(uri) >= suneidoPrefixSize && 
            wcsncmp(uri, suneidoPrefix, suneidoPrefixSize) == 0) {
            ICoreWebView2Environment* environment = nullptr;

            ICoreWebView2_2* webview2_2 = nullptr;
            hr = sender->QueryInterface(IID_ICoreWebView2_2, (void**)&webview2_2);
            if (FAILED(hr) || !webview2_2) {
                CoTaskMemFree(uri);
                request->Release();
                return hr;
            }

            hr = webview2_2->get_Environment(&environment);
            if (FAILED(hr) || !environment) {
                CoTaskMemFree(uri);
                request->Release();
                webview2_2->Release();
                return hr;
            }

            char *converted = ConvertWcharToChar(uri);
            buf_t buf = suneidoAPP(converted);

            if (buf.size == 0) {
                // Create a response
                ICoreWebView2WebResourceResponse* response = nullptr;
                environment->CreateWebResourceResponse(nullptr, 204, L"No Content", 
                    L"Content-Disposition: inline\r\nContent-Length: 0\r\nContent-Type: text/plain", &response);
                args->put_Response(response);
                response->Release();
            } else {
                IStream* stream = SHCreateMemStream((BYTE*)buf.data, (ULONG)(buf.size));
                if (stream) {
                    std::wstring header = buildContentTypeHeader(converted) + L"Content-Length: " + std::to_wstring(buf.size);

                    // Create a response
                    ICoreWebView2WebResourceResponse* response = nullptr;
                    environment->CreateWebResourceResponse(stream, 200, L"OK", header.c_str(), &response);
                    args->put_Response(response);
                    response->Release();
                    stream->Release();
                }
            }
            
            delete[] converted;
            delete buf.data;
            environment->Release();
            webview2_2->Release();
        }

        CoTaskMemFree(uri);
        request->Release();
        return S_OK;
    }

private:
    ~MyCustomSchemeHandler() = default;
    LONG refCount = 1;
};

HRESULT RegisterCustomScheme(ICoreWebView2* webView) {
    // Set the custom scheme handler
    MyCustomSchemeHandler* handler = new MyCustomSchemeHandler();

    // Add the custom scheme to handle all resources starting with "suneido:"
    HRESULT hr = webView->AddWebResourceRequestedFilter(L"*", COREWEBVIEW2_WEB_RESOURCE_CONTEXT_ALL);
    if (FAILED(hr)) {
        handler->Release();
        return hr;
    }

    // Register the handler
    EventRegistrationToken token;
    hr = webView->add_WebResourceRequested(handler, &token);
    handler->Release();
    return hr;
}