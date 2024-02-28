#pragma once

namespace Hardware::Modem {

    void Start();
    bool Initialize();
    bool ConnectNetwork();
    bool Shutdown();
    bool ConnectServer();
    void SendMessage(const char* msg);
    void getInfo();
    void getSignalStrength(void);
    void Test();
} //namespace Hardware::Modem