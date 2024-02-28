#pragma once

#include <zephyr/kernel.h>

namespace Hardware {
namespace TaskWorker {

void Start();
void Enable();

k_tid_t getThreadId();

} // namespace LedAlive
} // namespace Hardware
