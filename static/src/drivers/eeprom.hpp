#pragma once

#include <zephyr/kernel.h>
#include <zephyr/drivers/eeprom.h>

namespace Hardware::Eeprom {

const uint32_t kSystemId = 0xFADABABA;    

bool Initialize();
void PowerOn();
void PowerOff();
bool Test(uint16_t address);

}//namespace Hardware::Eeprom