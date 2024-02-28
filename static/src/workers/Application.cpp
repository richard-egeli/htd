#include <zephyr/kernel.h>
#include <zephyr/drivers/gpio.h>
#include <stdlib.h>
#include <stdio.h>
#include <sys/time.h>
#include "compile_time.h"

#include "Application.hpp"
#include "..\drivers\power.hpp"
#include "AliveLed.hpp"
#include "Sdi12.hpp"

#if 0
#include "datetime.hpp"


#include "Modem.hpp"
//#include "semaphore.hpp"
//#include "task.hpp"
#include "..\drivers\button.hpp"
#include "..\drivers\eeprom.hpp"

#include "..\drivers\i420ma.hpp"

using namespace Utilities::DateTime;

#endif

using namespace Hardware;
namespace Application {

const int kApplicationSleep = 10000;
const int kThreadSleep = 100;
const int kStackSize = 500;
/*
cpp_semaphore semApp;

void Wait() { semApp.wait(); }
*/
/*
#define LED1_NODE DT_NODELABEL(red_led)
gpio_dt_spec redLed = GPIO_DT_SPEC_GET_OR(LED1_NODE, gpios, {0});
uint16_t gpio_pin = 0;
*/

#if 0 // Enable Shell
#include <zephyr/shell/shell.h>

static int cmd_get(const struct shell *shell, size_t argc, char *argv[])
{
	shell_fprintf(shell, SHELL_NORMAL, "Current GPIO state: %d\n", gpio_pin);
	return 0;
}

static int32_t validateInput(char* cmdstr, int32_t lowLimit, int32_t high_limit)
{
    /* Received value is a string so do some test to convert and validate it */

    if ((strlen(cmdstr) == 1) && (cmdstr[0] == '0')) {
        return 0;
    }
    else {
        int32_t value = strtol(cmdstr, NULL, 10);
        if (value == 0) {
            //There was no number at the beginning of the string
            return -1;
        }
		/* Reject invalid values */
		if ((value < lowLimit) || (value > high_limit)) {
			return -1;
		}
		return value;
	}
} 

static int cmd_set(const struct shell *shell, size_t argc, char *argv[])
{
	const int32_t lowLimit = 0;
	const int32_t highLimit = 31;

    shell_fprintf(shell, SHELL_NORMAL, "DeviceTree GPIO pin number is %d\n", redLed.pin);

	int32_t value = validateInput(argv[1], lowLimit, highLimit);
	if (value<0) {
		shell_fprintf(shell, SHELL_ERROR, "Invalid value: %s; expected [%d..%d]\n", argv[1], lowLimit, highLimit);
	}
    else {/* Otherwise set and report to the user with a shell message */
	    gpio_pin = (uint16_t)value;
 		redLed.pin = gpio_pin;
		int ret = gpio_pin_configure_dt(&redLed, GPIO_OUTPUT);
		if (ret != 0) {
			printk("Error %d: failed to configure pin %d \n", ret, redLed.pin);
			return 0;
		}
       shell_fprintf(shell, SHELL_NORMAL, "GPIO pin number is %d\n", gpio_pin);
		gpio_pin_toggle(redLed.port, gpio_pin);
    }
 
    return 0;
}

static int cmd_clr(const struct shell *shell, size_t argc, char *argv[])
{
	const int32_t lowLimit = 0;
	const int32_t highLimit = 31;

	int32_t value = validateInput(argv[1], lowLimit, highLimit);
	if (value<0) {
		shell_fprintf(shell, SHELL_ERROR, "Invalid value: %s; expected [%d..%d]\n", argv[1], lowLimit, highLimit);
	}
    else {/* Otherwise set and report to the user with a shell message */
	    gpio_pin = (uint16_t)value;
		gpio_pin_set(redLed.port, gpio_pin, 0);
        shell_fprintf(shell, SHELL_NORMAL, "GPIO pin number is %d\n", gpio_pin);
    }
 
    return 0;
}

SHELL_STATIC_SUBCMD_SET_CREATE(
	shell_cmds,
	SHELL_CMD_ARG(set, NULL,
		"set the GPIO\n"
		"usage: $ gpio set <bit>\n",
		cmd_set, 2, 0),
	SHELL_CMD_ARG(clr, NULL,
		"reset the GPIO\n"
		"usage: $ gpio clr <bit>\n",
		cmd_clr, 2, 0),
	SHELL_CMD_ARG(get, NULL,
		"get the GPIO state\n"
		"usage: $ gpio get 2",
		cmd_get, 2, 0),
	SHELL_SUBCMD_SET_END
	);

struct timeval localTime;

static void KeyPressHandler(Button::ButtonEvent evt)
{
	printk("Button event: Key %s!\n", evt == Button::ButtonEvent::keyPressed ? "Pressed" : "Released");
//	if(evt == Button::ButtonEvent::keyPressed)
//		Eeprom::Test(0);
	gpio_pin_set_dt(&Application::redLed, (int)evt == Button::ButtonEvent::keyPressed);
}
#endif

} // namespace Application

int Application::Main()
{
    printk("Hello World from %s!!!\n", CONFIG_BOARD);
	printk("Application: Initializing...\n");                           // Starting main application
    if(!Power::Initialize())                                            // Initialize and Switch OFF all power domains
    {
        printk("ERROR: Failed to initialize power sytem!\n");
    }
	Power::Smps_On();
   	LedAlive::Start();                                                  // Start Alive LED
	//SDI12::Start();													// Commented by GRA
	//int ret = gpio_pin_configure_dt(&redLed, GPIO_OUTPUT);
	//if(I420mA::Initialize() != I420mA::resultError::kOK) return 2;

#if 0    
    if(!Modem::Initialize()) return 1;
    Modem::Test();

    //Utilities::DateTime::InitializeRTC(UNIX_TIMESTAMP); 				// Workaround for RTC initialization
   	//LedAlive::Start();                                                  // Start Alive LED
    //LedAlive::Initialize();                                             // Initialize LED pin

    //TaskWorker::Start();                                                // Start dependent worker Thread
//	k_timer timer;
//	k_timer_init(&timer, NULL, NULL);

//    SHELL_CMD_REGISTER(gpio, &shell_cmds, "GPIO set,clear and read.", NULL);
/*
	int ret = gpio_pin_configure_dt(&redLed, GPIO_OUTPUT);
	if (ret != 0) {
		printk("Error %d: failed to configure pin %d (LED '%s')\n", ret, redLed.pin, "red led");
		return 1;
	}
*/
	#if 0
	if(! Button::Initialize(KeyPressHandler))
	{
		printk("Error: Failed to configure button!\n");
		return 2;
	}

	if(! Eeprom::Initialize())
	{
		printk("Error: Failed to configure EEProm!\n");
		return 3;
	}
	#endif
/*
	if(! Modem::Initialize())
	{
		printk("Error: Failed to configure Modem!\n");
		return 4;
	}
*/
#endif
	printk("Application: Initializing done!\n");                           
    Power::DeepSleep();												//un-commented by GRA
    //LedAlive::PauseResume();
	//printk("Application: Resuming!\n");                           // Running main application
	while (1) 
    {
		//I420mA::Test();

		//printk("Application: Running...\n");                           
		// Utilities::DateTime::printTime();
    	// k_msleep(kApplicationSleep);
		// printk(" Done!\n");                                             // Finishing processing
		// printk("Application: Task execution...\n");
		// TaskWorker::Enable();                                           // Allow task to run
		// printk("Application: Waiting task completion...\n");
		// Wait();                                                         // Wait for task to end
		k_msleep(1000);
	}
}

void Application::Continue() { 
	//semApp.give(); 
}

