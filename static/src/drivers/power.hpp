#pragma once

#include <zephyr/kernel.h>
#include <zephyr/drivers/eeprom.h>
#include <zephyr/drivers/gpio.h>

namespace Hardware:: Power {

bool Initialize();
void On();
void Off();

void Sdi12_On();
void Sdi12_Off();
void Sdi12_Ch1_On();
void Sdi12_Ch1_Off();
void Sdi12_Ch2_On();
void Sdi12_Ch2_Off();
void I420_On();
void I420_Off();
void Ble_On();
void Ble_Off();
void Memory_On();
void Memory_Off();
void I2C_On();
void I2C_Off();
void Smps_On();
void Smps_Off();

bool InitializePin(const gpio_dt_spec* pin, const char* label);

void LoggerDisable();
void LoggerEnable();

void DeepSleep();

}//namespace Hardware::Power