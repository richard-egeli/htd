#include "eeprom.hpp"
#include <zephyr/drivers/gpio.h>

namespace Hardware::Eeprom {


#define MEMON_NODE DT_NODELABEL(memon)
gpio_dt_spec pinMemOn = GPIO_DT_SPEC_GET_OR(MEMON_NODE, gpios, {0});

size_t mSize = 0;
const device* mDev = nullptr;

struct parameters {
	uint32_t systemID;
	uint32_t paramX;
};

static bool getDevice(void)
{
	mDev = DEVICE_DT_GET(DT_ALIAS(eeprom_0));

	if (!device_is_ready(mDev)) {
		printk("\nError: Device \"%s\" is not ready; \n", mDev->name);
		return false;
	}

	printk("Found EEPROM device \"%s\"\n", mDev->name);
	return true;
}


}// namespace Hardware::Button

using namespace Hardware;
bool Eeprom::Initialize()
{
    if (!getDevice())
		return false;
	int ret = gpio_pin_configure_dt(&pinMemOn, GPIO_OUTPUT);
	if (ret != 0) {
		printk("Error %d: failed to configure enable pin@%d)\n", ret, pinMemOn.pin);
		return false;
	}
    mSize = eeprom_get_size(mDev);
    printk("Using eeprom with size of: %zu.\n", mSize);
	PowerOff();
    return true;
}

bool Eeprom::Test(uint16_t offset)
{
    parameters params;
    PowerOn();
    int result = eeprom_read(mDev, offset, &params, sizeof(params));
	if (result < 0) {
		printk("Error: Couldn't read eeprom: err: %d.\n", result);
		return false;
	}
    printk("Read device ID: %zu\n", params.systemID);

	if (params.systemID != kSystemId) {
		params.systemID = kSystemId;
		params.paramX = 0x00550055;
	}

	result = eeprom_write(mDev, offset, &params, sizeof(params));
	if (result < 0) {
		printk("Error: Couldn't write eeprom: err:%d.\n", result);
		return 0;
	}

    result = eeprom_read(mDev, offset, &params, sizeof(params));
	if (result < 0) {
		printk("Error: Couldn't read eeprom: err: %d.\n", result);
		return false;
	}
    printk("2nd Read device ID: %zu\n", params.systemID);

    PowerOff();

    return true;
}

void Eeprom::PowerOn()  { gpio_pin_set_dt(&pinMemOn, 1); 	k_msleep(10); } // Turn on Flash and Eeprom
void Eeprom::PowerOff() { gpio_pin_set_dt(&pinMemOn, 0); }     			// Turn off Flash and Eeprom
