#include "Modem.hpp"
#include <modem/lte_lc.h>
#include <modem/nrf_modem_lib.h>
#include <modem/modem_info.h>
#include <nrf_modem.h>
#include <nrf_modem_at.h>

#include "../drivers/socket.h"
#include <stdio.h>
#include <zephyr/net/socket.h>


/* Define a callback */
LTE_LC_ON_CFUN(cfun_hook, on_cfun, NULL);

/* Callback implementation */
static void on_cfun(enum lte_lc_func_mode mode, void *context)
{
    printk("Functional mode changed to %d\n", mode);
}


namespace Hardware::Modem {

void modemWorker()
{

}

const int kStackSize = 500;
k_thread modemThread;
K_THREAD_STACK_DEFINE(modemStack, kStackSize);

int mSocketId=0;

} //namespace Hardware::Modem

using namespace Hardware;

bool Modem::Initialize() 
{
    printk("Initializing LTE network...\n");

    int error = nrf_modem_lib_init();

    if (error != 0) {
        printk("Failed to initialize modem library, error: %d\n", error);
        return false;
    }
    printk("Initialized modem library\n");

    error = lte_lc_init();
    if (error != 0) {
        printk("Failed to initialize LTE link controller, error: %d\n", error);
        return false;
    }
    printk("Initialized LTE link controller\n");

    return true;
}

bool Modem::ConnectNetwork()
{
   int error = lte_lc_connect();
    if (error != 0) {
        printk("Failed to connect to LTE network, error: %d\n", error);
        return false;
    }
    printk("Connected to LTE network\n");
    return true;
}


bool Modem::ConnectServer()
{
    int mSocketId = socket(AF_INET, SOCK_STREAM, IPPROTO_TCP);
    if (mSocketId < 0) {
        printk("Failed to create socket: %d\n", errno);
        return errno;
    }

    struct sockaddr_in addr = {
        .sin_family = AF_INET,
        .sin_port   = htons(80),
        .sin_addr =
            {
                .s_addr = htonl(0x0a00000a),
            },
    };

    int error = connect(mSocketId, (struct sockaddr *)&addr, sizeof(addr));

    if (error != 0) {
        printk("Failed to connect to server: %d\n", errno);
        return errno;
    }
    printk("Connected to server\n");
    return true;
}

void Modem::SendMessage(const char* msg)
{
    int bytes = send(mSocketId, msg, strlen(msg), 0);
    printk("Sent %d bytes: %s\n", bytes, msg);
}

bool Modem::Shutdown() {
    printk("Shutting down LTE network...\n");

    int error = lte_lc_offline();
    if (error != 0) {
        printk("Failed to disconnect from LTE network, error: %d\n", error);
        return false;
    }

    printk("Disconnected from LTE network\n");
    error = nrf_modem_lib_shutdown();
    if (error != 0) {
        printk("Failed to shutdown modem library, error: %d\n", error);
        return false;
    }

    printk("Shutdown modem library\n");
    return true;
}

void Modem::Start()
{
	k_thread_create(&modemThread, modemStack, kStackSize, 
					(k_thread_entry_t) modemWorker, NULL, NULL, NULL, 
					K_PRIO_COOP(7), 0, K_NO_WAIT);
}


void Modem::getInfo()
{
    modem_param_info modem_param;
	modem_info_init();
	modem_info_params_init(&modem_param);
    
    printk("====== Cell Network Info ======\n");
    char sbuf[128]={0};
    modem_info_string_get(MODEM_INFO_RSRP, sbuf, sizeof(sbuf));
    printk("Signal strength: %s\n", sbuf);
    modem_info_string_get(MODEM_INFO_CUR_BAND, sbuf, sizeof(sbuf));
    printk("Current LTE band: %s\n", sbuf);
    modem_info_string_get(MODEM_INFO_SUP_BAND, sbuf, sizeof(sbuf));
    printk("Supported LTE bands: %s\n", sbuf);
    modem_info_string_get(MODEM_INFO_AREA_CODE, sbuf, sizeof(sbuf));
    printk("Tracking area code: %s\n", sbuf);
    modem_info_string_get(MODEM_INFO_UE_MODE, sbuf, sizeof(sbuf));
    printk("Current mode: %s\n", sbuf);
    modem_info_string_get(MODEM_INFO_OPERATOR, sbuf, sizeof(sbuf));
    printk("Current operator name: %s\n", sbuf);
    modem_info_string_get(MODEM_INFO_CELLID, sbuf, sizeof(sbuf));
    printk("Cell ID of the device: %s\n", sbuf);
    modem_info_string_get(MODEM_INFO_IP_ADDRESS, sbuf, sizeof(sbuf));
    printk("IP address of the device: %s\n", sbuf);
    modem_info_string_get(MODEM_INFO_FW_VERSION, sbuf, sizeof(sbuf));
    printk("Modem firmware version: %s\n", sbuf);
    modem_info_string_get(MODEM_INFO_LTE_MODE, sbuf, sizeof(sbuf));
    printk("LTE-M support mode: %s\n", sbuf);
    modem_info_string_get(MODEM_INFO_NBIOT_MODE, sbuf, sizeof(sbuf));
    printk("NB-IoT support mode: %s\n", sbuf);
    modem_info_string_get(MODEM_INFO_GPS_MODE, sbuf, sizeof(sbuf));
    printk("GPS support mode: %s\n", sbuf);
    modem_info_string_get(MODEM_INFO_DATE_TIME, sbuf, sizeof(sbuf));
    printk("Mobile network time and date: %s\n", sbuf);
    printk("===============================\n");
}
void Modem::getSignalStrength(void)
{

    const char signal_strength_command[] = "AT+CESQ";
    char response[50] = {0};
    nrf_modem_at_cmd(response, 50, "%s", signal_strength_command);
    printk("%s: %s\n", signal_strength_command, response);
}

void Modem::Test()
{
    static int counter=123;
    if(Modem::ConnectNetwork()) 
    {
        //Modem::getInfo();
        //Modem::getSignalStrength();
        if(Modem::ConnectServer())
        {
            char buf[50] = "";
            sprintf(buf, "Hello: %d", counter++);
            Modem::SendMessage(buf);
        }
    }
#if 0
    char response[32] = {0};
    nrf_modem_at_cmd(response, sizeof(response), "%s", "AT");
    printk("nrf_modem_at_cmd: %s\n", response);

    nrf_modem_at_cmd(response, sizeof(response), "%s", "AT+CGMM");
    printk("nrf_modem_at_cmd: %s\n", response);
    memset(response, 0, sizeof(response));
    nrf_modem_at_cmd(response, sizeof(response), "%s", "AT+CGSN");
    printk("nrf_modem_at_cmd: %s\n", response);
    nrf_modem_at_cmd(response, sizeof(response), "%s", "AT+CFUN?");
    printk("nrf_modem_at_cmd: %s\n", response);

    nrf_modem_at_cmd(response, sizeof(response), "%s", "AT%HWVERSION");
    printk("nrf_modem_at_cmd: %s\n", response);

    nrf_modem_at_cmd(response, sizeof(response), "%s", "AT%SHORTSWVER");
    printk("nrf_modem_at_cmd: %s\n", response);
#endif
}