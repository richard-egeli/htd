#include "button.hpp"
#include <zephyr/drivers/gpio.h>

namespace Hardware::Button {

#define SW_NODE   DT_NODELABEL(button0)
static gpio_dt_spec sw1 = GPIO_DT_SPEC_GET_OR(SW_NODE, gpios, {0});
static struct gpio_callback button_cb_data;

static KeyPressEvent userHandler;

static void onEndOfDebouncingTime(struct k_work *work) // Called after debouncing time (20ms)
{
    ARG_UNUSED(work);

    int val = gpio_pin_get_dt(&sw1);
    ButtonEvent evt = (val ? ButtonEvent::keyPressed : ButtonEvent::keyReleased);
    if (userHandler) 
        userHandler(evt); 
}

static K_WORK_DELAYABLE_DEFINE(debouncingWorker, onEndOfDebouncingTime);

void onKeyPressed(const device *dev, gpio_callback *cb, uint32_t pins)
{
   k_work_reschedule(&debouncingWorker, K_MSEC(kDebouncingMS));
}


} // namespace Hardware::Button

using namespace Hardware;

bool Button::Initialize(KeyPressEvent handler)
{
    if (!handler) {
		printk("Init Error: Invalid handler!\n");
        return false;
    }
    userHandler = handler;
	if (!device_is_ready(sw1.port))
	{
		printk("Init Error: GPIO not ready!\n");
		return false;
	}
	
	gpio_pin_configure_dt(&sw1, GPIO_INPUT);                            // configure the button pin as input
	gpio_pin_interrupt_configure_dt(&sw1, GPIO_INT_EDGE_BOTH);	        // configure the interrupt on button press
	gpio_init_callback(&button_cb_data, onKeyPressed, BIT(sw1.pin));      // setup the button press callback
	if(gpio_add_callback(sw1.port, &button_cb_data)<0)
    {
		printk("Init Error: Fail to add button callback!\n");
		return false;
    }

    return true;
}