#include "power.hpp"
#include <modem/lte_lc.h>
#include <modem/nrf_modem_lib.h>
#include <zephyr/drivers/uart.h>
namespace Hardware::Power {


#define SDION_NODE DT_NODELABEL(sdi12poweron)
#define SDIONSEN1_NODE DT_NODELABEL(sdi12sen1on)
#define SDIONSEN2_NODE DT_NODELABEL(sdi12sen2on)
#define I420ON_NODE DT_NODELABEL(i420on)
#define BLEON_NODE DT_NODELABEL(bleon)
#define MEMON_NODE DT_NODELABEL(memon)
#define I2CON_NODE DT_NODELABEL(i2con)
#define SMPSON_NODE DT_NODELABEL(smpson)

#define UART_DEVICE_NODE DT_CHOSEN(zephyr_shell_uart)

static const struct device *const uart_dev = DEVICE_DT_GET(UART_DEVICE_NODE);

gpio_dt_spec pinSdiOn = GPIO_DT_SPEC_GET_OR(SDION_NODE, gpios, {0});
gpio_dt_spec pinSdiSen1On = GPIO_DT_SPEC_GET_OR(SDIONSEN1_NODE, gpios, {0});
gpio_dt_spec pinSdiSen2On = GPIO_DT_SPEC_GET_OR(SDIONSEN2_NODE, gpios, {0});
gpio_dt_spec pinI420On = GPIO_DT_SPEC_GET_OR(I420ON_NODE, gpios, {0});
gpio_dt_spec pinBleOn = GPIO_DT_SPEC_GET_OR(BLEON_NODE, gpios, {0});
gpio_dt_spec pinMemOn = GPIO_DT_SPEC_GET_OR(MEMON_NODE, gpios, {0});
gpio_dt_spec pinI2cOn = GPIO_DT_SPEC_GET_OR(I2CON_NODE, gpios, {0});
gpio_dt_spec pinSmpsOn = GPIO_DT_SPEC_GET_OR(SMPSON_NODE, gpios, {0});

// Private implementation
static bool Sdi12EnableInit();
static bool I4_20EnableInit();
static bool BleEnableInit();
static bool MemoryEnableInit();
static bool I2CEnableInit();
static bool SmpsEnableInit();
static void setEnable(const gpio_dt_spec* pin, int value, int32_t waitTime=0);

static bool Sdi12EnableInit()
{
	if(!InitializePin(&pinSdiOn, 		"SDI12 PowerOn")) return false;
	if(!InitializePin(&pinSdiSen1On,	"SDI12 Chan1On")) return false;
	if(!InitializePin(&pinSdiSen2On,	"SDI12 Chan2On")) return false;
	return true;
}

static bool I4_20EnableInit()
{
	return InitializePin(&pinI420On, 		"4..20mA PowerOn");
}

static bool SmpsEnableInit()
{
	return InitializePin(&pinSmpsOn, 		"SMPS PowerOn");	
}

static bool BleEnableInit()
{
	return InitializePin(&pinBleOn, 		"BLE PowerOn");	
}

static bool MemoryEnableInit()
{
	return InitializePin(&pinMemOn, 		"Memory PowerOn");	
}

static bool I2CEnableInit()
{
	return InitializePin(&pinI2cOn, 		"I2C PowerOn");	
}

static void setEnable(const gpio_dt_spec* pin, int value, int32_t waitTime)
{
	gpio_pin_set_dt(pin, value); 
	if(waitTime!=0)
		k_msleep(waitTime);
}


}// namespace Hardware::Button

using namespace Hardware;

bool Power::InitializePin(const gpio_dt_spec* pin, const char* label)
{
	int ret = gpio_pin_configure_dt(pin, GPIO_OUTPUT);
	if (ret != 0) {
		printk("Power: failed to configure %s gpio's!\n", label);
		return false;
	}
	return true;	
}

bool Power::Initialize()	// Initialize all power control gpio's
{
	printk("Power: Initializing...\n");
	if(!Power::Sdi12EnableInit())	return false;
	if(!Power::I4_20EnableInit())	return false;
	if(!Power::BleEnableInit())		return false;
	if(!Power::MemoryEnableInit())	return false;
	if(!Power::I2CEnableInit()) 	return false;
	if(!Power::SmpsEnableInit()) 	return false;

	Power::Off();
	printk("Power: Initializing Done!\n");
    return true;
}

void Power::Sdi12_On()  	{ setEnable(&pinSdiOn, 1, 10); }
void Power::Sdi12_Off()
{ 
	setEnable(&pinSdiOn, 0);
	Sdi12_Ch1_Off();
	Sdi12_Ch2_Off();
}

void Power::Sdi12_Ch1_On()  { Sdi12_On();  setEnable(&pinSdiSen1On, 1, 10); }
void Power::Sdi12_Ch1_Off() { setEnable(&pinSdiSen1On, 0); }

void Power::Sdi12_Ch2_On()  { Sdi12_On();  setEnable(&pinSdiSen2On, 1, 10); }
void Power::Sdi12_Ch2_Off() { setEnable(&pinSdiSen2On, 0); }

void Power::I420_On()   	{ setEnable(&pinI420On, 1, 10); }
void Power::I420_Off()  	{ setEnable(&pinI420On, 0); }

void Power::Ble_On()     	{ setEnable(&pinBleOn, 1, 10); } 
void Power::Ble_Off()    	{ setEnable(&pinBleOn, 0); }

void Power::Memory_On()  	{ setEnable(&pinMemOn, 1, 10); }
void Power::Memory_Off() 	{ setEnable(&pinMemOn, 0); }

void Power::I2C_On()     	{ setEnable(&pinI2cOn, 1, 10); }
void Power::I2C_Off()    	{ setEnable(&pinI2cOn, 0); }

void Power::Smps_On()     	{ setEnable(&pinSmpsOn, 1, 200); }
void Power::Smps_Off()    	{ setEnable(&pinSmpsOn, 0); }

void Power::On() 	// Enable all OnBoard Power domains (For test only)
{ 
	Sdi12_Ch1_On();
	Sdi12_Ch2_On();
	I420_On();
	Ble_On();
	Memory_On(); 
	I2C_On();
	Smps_On();
} 
void Power::Off() 	// Disable all OnBoard Power domains
{
	Sdi12_Off();
	I420_Off();
	Ble_Off();
	Memory_Off(); 
	I2C_Off();
	Smps_Off();
}

void Power::LoggerDisable()
{
	uart_rx_disable(uart_dev);

	NRF_UARTE0->TASKS_STOPTX = 1;
	NRF_UARTE0->EVENTS_RXTO = 0;
	NRF_UARTE0->TASKS_STOPRX = 1;
	while(NRF_UARTE0->EVENTS_RXTO == 0){}
	NRF_UARTE0->EVENTS_RXTO = 0;
	// disable whole UART 0 interface
	NRF_UARTE0->ENABLE = 0;

	NRF_UARTE0->PSEL.TXD = 0xFFFFFFFF;
	NRF_UARTE0->PSEL.RXD = 0xFFFFFFFF;
}

void Power::LoggerEnable()
{
	//uart_rx_enable(uart_dev)
}

void Power::DeepSleep()
{
	printk("Power: Entering low power mode!\n");	// Running main application
    k_sleep(K_MSEC(1000));    
	Power::Off();
	//LoggerDisable();
	nrf_modem_lib_init();
     lte_lc_power_off();
    k_sleep(K_MSEC(1000));         					// Wait for the modem to power off
    NRF_REGULATORS->SYSTEMOFF = 1; 					// Power off the system
}