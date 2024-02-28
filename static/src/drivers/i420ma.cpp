#include <zephyr/drivers/gpio.h>
#include "i420ma.hpp"
#include "power.hpp"

namespace Hardware::I420mA {



const struct device* const gpio_dev = DEVICE_DT_GET(DT_NODELABEL(gpio0));
const struct device* const adc_dev  = DEVICE_DT_GET(DT_NODELABEL(adc));

static uint16_t adc_buffer;

static adc_channel_cfg adc_config = {
    .gain             = ADC_GAIN_1_4,
    .reference        = ADC_REF_INTERNAL,
    .acquisition_time = ADC_ACQ_TIME_DEFAULT,
    .channel_id       = ADC_CHANNEL_1,
    .input_positive   = ADC_CHANNEL_1,
    .input_negative   = ADC_CHANNEL_DISABLED,
};

static adc_sequence sequence = {
    .channels     = BIT(ADC_CHANNEL_1),
    .buffer       = &adc_buffer,
    .buffer_size  = sizeof(adc_buffer),
    .resolution   = ADC_RESOLUTION_12BIT,
    .oversampling = ADC_OVERSAMPLE_256X,
    .calibrate    = true,
};


resultError mConfigureChannel(adc_channel channel)
{
    adc_config.channel_id = channel;
    adc_config.input_positive = channel;
    sequence.channels = BIT(channel);

    int error = adc_channel_setup(adc_dev, &adc_config);

    if (error != 0) {
        printk("Failed to setup channel: %d\n", error);
        Power::I420_Off();
        return resultError::kChSetupError;
    };

    Power::I420_On();
    k_sleep(K_MSEC(100));

    return resultError::kOK;
}

}//namespace Hardware::I420mA

using namespace Hardware;

I420mA::resultError I420mA::Initialize()
{
    if (!device_is_ready(gpio_dev) || !device_is_ready(adc_dev)) {
        printk("Device not ready\n");
        return resultError::kDeviceError;
    }

    Power::I420_On();
    k_sleep(K_MSEC(100));
    return mConfigureChannel(adc_channel::ADC_CHANNEL_1); // Default to channel 1
}

I420mA::resultError I420mA::Read()
{
    int error = adc_read(adc_dev, &sequence);
    if (error != 0) {
        printk("Failed to read ADC: %d\n", error);
        Power::I420_Off();
        return resultError::kRdError;
    }
    Power::I420_Off();

    return resultError::kOK;
}
void I420mA::Test()
{
    adc_channel ch = adc_channel::ADC_CHANNEL_2;
    mConfigureChannel(ch);
    I420mA::Read();
    //uint16_t result = (adc_buffer[1] << 8) | adc_buffer[0];
    int32_t result = adc_buffer;
    adc_raw_to_millivolts(adc_ref_internal(adc_dev),
					      adc_config.gain,
					      sequence.resolution,
					      &result);
    //printk("ADC result>> %d >>> %dmV\n", adc_buffer, result);
    printk("4-20mA Channel %d = %.2f mA\n", ch-3, adc_buffer/136.0F);
}