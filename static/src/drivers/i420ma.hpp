#pragma once
#include <zephyr/kernel.h>
#include <zephyr/drivers/adc.h>
#include <hal/nrf_saadc.h>

namespace Hardware::I420mA {

enum class resultError {
    kOK             =  0, 
    kDeviceError    = -1, 
    kChSetupError   = -2, 
    kRdError        = -3
};

typedef enum adc_oversample {
    ADC_OVERSAMPLE_1X,
    ADC_OVERSAMPLE_2X,
    ADC_OVERSAMPLE_4X,
    ADC_OVERSAMPLE_8X,
    ADC_OVERSAMPLE_16X,
    ADC_OVERSAMPLE_32X,
    ADC_OVERSAMPLE_64X,
    ADC_OVERSAMPLE_128X,
    ADC_OVERSAMPLE_256X
} adc_oversample_t;

typedef enum adc_resolution {
    ADC_RESOLUTION_8BIT  = 8,
    ADC_RESOLUTION_10BIT = 10,
    ADC_RESOLUTION_12BIT = 12,
    ADC_RESOLUTION_14BIT = 14,
} adc_resolution_t;

typedef enum adc_channel {
    ADC_CHANNEL_DISABLED = NRF_SAADC_INPUT_DISABLED,
    ADC_CHANNEL_1        = NRF_SAADC_INPUT_AIN3,
    ADC_CHANNEL_2        = NRF_SAADC_INPUT_AIN4,
} adc_channel_t;

resultError Initialize();
resultError Read();
void Test();

} //namespace Hardware::I420mA