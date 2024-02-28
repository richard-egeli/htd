#include "datetime.hpp"
#include "compile_time.h"
#include <unistd.h>

namespace Utilities {
namespace DateTime {

static tm curTm;


void printTime()
{
    timeval curTime;
    int res = gettimeofday(&curTime, NULL);
    time_t now = time(NULL);
    struct tm tm;
    localtime_r(&now, &tm);

    if (res < 0) {
        printk("\nError in gettimeofday(): %d\n", errno);
        return;
    }

    printk("\n%s\n", asctime(&tm));
}

uint32_t getSeconds()
{
    ReadRTC();
    return curTm.tm_sec;
}

uint32_t getMinutes()
{
    ReadRTC();
    return curTm.tm_min;
}

uint32_t getHours()
{
    ReadRTC();
    return curTm.tm_hour;
}

uint32_t getDay()
{
    ReadRTC();
    return curTm.tm_mday;
}

uint32_t getMonth()
{
    ReadRTC();
    return curTm.tm_mon + 1;            // Adjust Month
}

uint32_t getYear()
{
    ReadRTC();
    return curTm.tm_year + 1900;        // Adjust year
}

// initialize clock from the compile time GCC timestamp
void InitializeRTC()                  { InitializeRTC(UNIX_TIMESTAMP); }

void InitializeRTC(time_t unix_ts)
{
    timeval curTime;
	curTime.tv_sec =  unix_ts + kAdjustToNorway*3600; 
	//clock_settime(CLOCK_REALTIME, (const timespec*) &curTime);
}

void ReadRTC()
{
    time_t timeNow = time(nullptr);
    curTm = *localtime(&timeNow);
}


} // namespace DateTime
} // namespace Utilities
