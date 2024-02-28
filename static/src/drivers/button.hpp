#pragma once

#include <zephyr/kernel.h>


namespace Hardware::Button {

const uint64_t kDebouncingMS = 20;
enum  ButtonEvent { keyPressed, keyReleased };

typedef void (*KeyPressEvent)(ButtonEvent evt);

bool Initialize(KeyPressEvent handler); // Initialize hardware. Return true if successeful

} // namespace Hardware::Button

