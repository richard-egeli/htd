#include "AliveLed.hpp"
#include <zephyr/kernel.h>
#include <zephyr/drivers/gpio.h>


namespace Hardware::LedAlive {

#define LED0_NODE DT_NODELABEL(green_led)

gpio_dt_spec spec = GPIO_DT_SPEC_GET_OR(LED0_NODE, gpios, {0});
const char *led_name = DT_PROP_OR(LED0_NODE, label, "");
bool mLedEnabled=false;
void ledAliveWorker()
{
	const int32_t sleep_ms = 100;
	 if (!Initialize()) return;			// Exit if initialization failed

	uint32_t cnt=0;
	while (true) {
		if(mLedEnabled)
			gpio_pin_set(spec.port, spec.pin, !(cnt++ % 10));
		k_msleep(sleep_ms);
	}
}

const int kStackSize = 500;
k_thread ledAlive_thread;
K_THREAD_STACK_DEFINE(ledAlive_stack, kStackSize);

} //namespace Hardware::LedAlive

using namespace Hardware;

void LedAlive::Start()
{
	mLedEnabled=true;
	k_thread_create(&ledAlive_thread, ledAlive_stack, kStackSize, 
					(k_thread_entry_t) ledAliveWorker, NULL, NULL, NULL, 
					K_PRIO_COOP(7), 0, K_NO_WAIT);
	k_thread_name_set(&ledAlive_thread, "ledAlive");
}

void LedAlive::PauseResume()
{
	mLedEnabled = !mLedEnabled;
	if(!mLedEnabled)
		gpio_pin_set(spec.port, spec.pin, 0);
}

bool LedAlive::Initialize()
{
	int ret = gpio_pin_configure_dt(&spec, GPIO_OUTPUT);
	if (ret != 0) {
		printk("Error %d: failed to configure pin %d (LED '%s')\n", ret, spec.pin, led_name);
		return false;
	}
	return true;
}