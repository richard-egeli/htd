#include "Sdi12.hpp"
#include <zephyr/kernel.h>
#include <zephyr/device.h>
#include <zephyr/pm/device.h>
#include <zephyr/drivers/uart.h>
#include <zephyr/drivers/gpio.h>
#include <hal/nrf_gpio.h>
#include <stdio.h>
#include "power.hpp"



//#include <zephyr/drivers/sensor.h>
//#include <zephyr/sys/ring_buffer.h>
//#include <zephyr/logging/log.h>
//LOG_MODULE_REGISTER(sdi12, CONFIG_SENSOR_LOG_LEVEL);

namespace Hardware::SDI12 {

static const size_t kBufferSize   = 82;                                      	// SDI-12 Buffer size
static const size_t kTxBufferSize = kBufferSize;
static const size_t kRxBufferSize = kBufferSize;
static const size_t kMessageSize  = kBufferSize;
static const int 	kReadTimeout  = 2000;                                     	// Read timout in SDI-12 bus (in ms)
static const int 	kStackSize 	  = 1024;

#define SDIDIR_NODE DT_NODELABEL(sdi12dir)
K_MSGQ_DEFINE(uart_msgq, kMessageSize, 10, 4);									// Queue to store up to 10 messages
K_THREAD_STACK_DEFINE(sdi12_stack, kStackSize);
const gpio_dt_spec pinSdiDir = GPIO_DT_SPEC_GET_OR(SDIDIR_NODE, gpios, {0});
k_thread sdi12_thread;
const device* const mDevUart = DEVICE_DT_GET(DT_NODELABEL(uart1));   			// SDI-12 UART device

const char cmdAcknowledge[] = "\r\n";
static char rxBuf[kRxBufferSize];			// Receiver buffer
static size_t rxBufPos=0;					// Read Pointer

//==================================
// Private implementation Prototypes
//==================================
static bool Sdi12DirectionInit();
static void mSdi12_XmitEnable();
static void mSdi12_XmitDisable();
static bool mUartConfig();
static void mSendBreak();
static void mSendCommand(uint8_t address, const char* command);

#if defined(CONFIG_UART_ASYNC_API)

static void uart_cb(const struct device *dev, struct uart_event *evt, void *user_data)
{
	switch (evt->type) {
	
	case UART_TX_DONE:
		// do something
		break;

	case UART_TX_ABORTED:
		// do something
		break;
		
	case UART_RX_RDY:
		// do something
		break;

	case UART_RX_BUF_REQUEST:
		// do something
		break;

	case UART_RX_BUF_RELEASED:
		// do something
		break;
		
	case UART_RX_DISABLED:
		// do something
		break;

	case UART_RX_STOPPED:
		// do something
		break;
		
	default:
		break;
	}
}
#elif defined(CONFIG_UART_INTERRUPT_DRIVEN)
static void serial_cb(const struct device *dev, void *user_data)
{
	if (!uart_irq_update(mDevUart)) 	return;
	if (!uart_irq_rx_ready(mDevUart))	return;

	
	uint8_t c;
	while (uart_fifo_read(mDevUart, &c, 1) == 1) 			// read until FIFO empty
	{
		if (rxBufPos < (sizeof(rxBuf) - 1))
		{
			rxBuf[rxBufPos++] = c;
		}
		else
		{
			rxBufPos=0;
		}

		if (((c == '\n') || (c == '!' )) && (rxBufPos > 1)) 
		{
			rxBuf[rxBufPos] = '\0';							// terminate string
			k_msgq_put(&uart_msgq, &rxBuf, K_NO_WAIT);		// if queue is full, message is silently dropped
			rxBufPos = 0;									// reset the buffer (it was copied to the msgq)
		} 
	}
}
#endif

static bool Sdi12DirectionInit()
{
	if (!device_is_ready(pinSdiDir.port)){
		printk("SDI-12: ERROR - GPIO device is not ready\r\n");
		return false;
	}

	if(!Power::InitializePin(&pinSdiDir, "SDI12 PowerOn")) return false;
	mSdi12_XmitDisable();
	return true;
}

void mRemoveParityBit(char* msg)
{
	while(*msg)
	{
		*msg &= 0x7F;
		msg++;
	}
}

static uint8_t mSetParity(uint8_t c)
{
    uint8_t n=c;
    n ^= (n >> 4);
    n &= 0x0F; 
    uint32_t p = (0x6996 >> n) & 0x01;
    return (p ? c | 0x80 : c & 0x7F);
}

static char mCommandString[20] = {0};
static char mAnswer[kBufferSize] = {0};

static void mBuildCommand(uint8_t address, const char* command)
{
	char* p = mCommandString;	// Initialize destination pointer
    if(address != '-')			// append address if exists
	{
		*p++ = address;
	}
	while(*command)				// append command
	{
		*p++ = *command++; 
	}
	*p++ = '!';					// append command terminator
	*p = '\0';					// terminate string
}

static void mSend(const char* msg)
{
	while(*msg)
	{
		uart_poll_out(mDevUart, mSetParity(*msg++));
	}
}

static void mSendCommand(uint8_t address, const char* command)
{
	mBuildCommand(address, command);
	printk("Sent: %s, ", mCommandString);
	mSdi12_XmitEnable();
	k_usleep(300);
	mSendBreak();
	mSend(mCommandString);
	k_msleep(10);
	mSdi12_XmitDisable();									// Disable transmitter 
}

static void mSendBreak()
{
	const int txPin = DT_PROP_BY_IDX(DT_CHILD(DT_NODELABEL(uart1_default), group1), psels, 0);

	int ret = pm_device_action_run(mDevUart, PM_DEVICE_ACTION_SUSPEND);	// Suspend UART1 device
	if (ret < 0) {
		printk("SDI-12: Failure to suspend UART!\n");
	} 
	else 
	{
		nrf_gpio_cfg_output(txPin);										// Re-configure TX pin as output
		nrf_gpio_pin_write(txPin, 1);
		nrf_gpio_pin_write(txPin, 0);									// Issue a break
		k_msleep(12);
		nrf_gpio_pin_write(txPin, 1);
		
		ret = pm_device_action_run(mDevUart, PM_DEVICE_ACTION_RESUME);	// Resume UART1 device
		if (ret < 0)
		{
			printk("SDI-12: Failure to resume device!\n");
		}
		k_msleep(8);													// Wait for 2ms
	}
}

static void printHex(char* msg)
{
	printk("Msg: [%s] = ", msg);
	for (int n=0 ; n<strlen(msg) ; n++)
		printk("%2X ", msg[n]);
	printk("\n");
}

static char* mGetAnswer()
{
	char rcvBuf[kBufferSize] = {0};					// queue msg receiver buffer
	mAnswer[0]='\0';								// invalidate answer
	int result = k_msgq_get(&uart_msgq, &rcvBuf, K_TIMEOUT_ABS_MS(kReadTimeout)); // Wait answer
	if(result == 0)
	{
		mRemoveParityBit(rcvBuf);
		printk("Queue msg1: {%s}\n", rcvBuf);
		if(strcmp(mCommandString, rcvBuf)==0)		// drop answer if same as command 
		{
			int result = k_msgq_get(&uart_msgq, &rcvBuf, K_TIMEOUT_ABS_MS(kReadTimeout)); // Wait answer
			if(result == 0)
			{
				mRemoveParityBit(rcvBuf);
				printk("Queue msg2: {%s}\n", rcvBuf);
			}
		}
		strcpy(mAnswer, rcvBuf);					// copy answer to field
	}
	return mAnswer;
}

static void mSdi12Worker()
{
	printk("SDI-12: Starting...\n");
	const int32_t sleep_ms = 5000;
	if (!SDI12::Initialize()) return;		    // Exit if initialization failed
    Power::Sdi12_Ch1_On();                      // For test only Power on SDI-12 interface on Channel 1
	uart_irq_rx_enable(mDevUart);				// Enable UART receiver
	printk("SDI-12: Successfully started!\n");
	k_msleep(1000);								// Wait device initialization

	while (true) {
        // TODO: SDI-12 state processing
        mSendCommand('-',"?");      				// Send Query command to bus
		char* p = mGetAnswer();
		printk("Received: %s\n", p);
		char address = *p;
		if(address >= '0' && address<= '9')
		{
		printk("Device address: %c\n", *p);



        mSendCommand(address,"I");      				// Send Query command to bus
		printk("Received: %s\n", mGetAnswer());
        mSendCommand('0',"M");      				// Send Query command to bus
		printk("Received: %s\n", mGetAnswer());
		k_msleep(3000);
        mSendCommand('0',"D0");      				// Send Query command to bus
		printk("Received: %s\n", mGetAnswer());
		}		
		k_msleep(sleep_ms);
	}
}

static void mSdi12_XmitEnable()  	{ gpio_pin_set_dt(&pinSdiDir, 0); }
static void mSdi12_XmitDisable()	{ gpio_pin_set_dt(&pinSdiDir, 1); }

static bool mUartConfig()
{
#if defined(CONFIG_UART_USE_RUNTIME_CONFIGURE)
	struct uart_config dev_config;

	dev_config.baudrate = 1200;
	dev_config.data_bits = UART_CFG_DATA_BITS_8;
	dev_config.parity = UART_CFG_PARITY_NONE;
	dev_config.stop_bits = UART_CFG_STOP_BITS_1;
	dev_config.flow_ctrl = UART_CFG_FLOW_CTRL_NONE;
	if(int err = uart_configure(mDevUart, &dev_config))
	{
		printk("SDI-12: ERROR - Failed to configure UART1! error %d\n", err);		
		return false;
	}
	return true;
#else 
	#error printk("Update prj.conf with 'CONFIG_UART_USE_RUNTIME_CONFIGURE=y'");
#endif
}

} // namespace Hardware::SDI12

using namespace Hardware;

bool SDI12::Initialize()
{
    printk("SDI-12: Initializing...\n");
    if (!device_is_ready(mDevUart)) {
        printk("SDI-12: ERROR - Device Initialization failed!\n");
        return false;
    }
	if(!Sdi12DirectionInit())	return false;
	mSdi12_XmitDisable();						// Disable transmitter 
	if (!mUartConfig())
	{
		return false;
	}

#if defined(CONFIG_UART_ASYNC_API)
	if (int err = uart_callback_set (mDevUart, uart_cb, NULL)) {
        printk("SDI-12: ERROR - Failed to set ASYNC API UART1 callback! error %d\n", err);
	 	return false;
	}
#elif defined(CONFIG_UART_INTERRUPT_DRIVEN)
	if (int err = uart_irq_callback_user_data_set(mDevUart, serial_cb, NULL)) {
        printk("SDI-12: ERROR - Failed to set INTERRUPT API UART1 callback! error %d\n", err);
	 	return false;
	}
#endif
    printk("SDI-12: Initializing done!\n");
    return true;
}

void SDI12::Start()
{
	k_thread_create(&sdi12_thread, sdi12_stack, kStackSize, 
					(k_thread_entry_t) mSdi12Worker, NULL, NULL, NULL, 
					K_PRIO_COOP(7), 0, K_NO_WAIT);
	k_thread_name_set(&sdi12_thread, "sdi-12");
}
