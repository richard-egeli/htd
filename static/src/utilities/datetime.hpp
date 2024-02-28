#pragma once

#include <zephyr/kernel.h>
#include <sys/time.h>

namespace Utilities::DateTime {

const time_t kAdjustToNorway = 0;   // Number of hours to Greenwich meridian

void InitializeRTC();
void InitializeRTC(time_t unix_ts);
void InitializeRTC(int year, int month, int day, int hour, int min, int sec);

uint32_t getYear();
uint32_t getMonth();
uint32_t getDay();
uint32_t getDayofWeek();
uint32_t getHours();
uint32_t getMinutes();
uint32_t getSeconds();

void     printTime();
void     ReadRTC();

} // namespace Utilities::DateTime
