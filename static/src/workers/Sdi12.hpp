#pragma once
namespace Hardware::SDI12 {

    enum SDI_SENSORS { SDI_SENSOR_NONE, SDI_SENSOR1, SDI_SENSOR2, SDI_SENSOR_ALL };
    enum SDI_CMD {
        SDI_CMD_ACTIVE, 
        SDI_CMD_ID, 
        SDI_CMD_CHG_ADDR, 
        SDI_CMD_START, 
        SDI_CMD_CONCURRENT, 
        SDI_CMD_GET_DATA, 
        SDI_CMD_CONTINUOUS, 
        SDI_CMD_START_OTHER, 
        SDI_CMD_CONTINUOUS_OTHER, 
        SDI_CMD_VERIFY,
        SDI_CMD_COUNT
    };
    bool Initialize();                          // Initialize the SDI-12 Interface
    void Start();                               // Start SDI-12 Handler thread
    /*
    void Enable(SDI_SENSORS ch);                // Enable SDI-12 Power and sensor in specified channel (one-by-one mode)
    void PowerOff();                            // Switch off all SDI-12 interface power
    void SendCmd(char adr, SDI_CMD cmd, char param = '?'); // Send a command to selected address (? is a dummy parameter)
    char* Answer();                            // get Bus answer (empty string "" on timeout)
    bool IsAnswerReady();
    bool IsSensorActive();                      // Check if sensor is present in the BUS
*/
} // namespace Hardware::SDI12
