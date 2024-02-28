#pragma once

namespace Hardware {
namespace LedAlive {

void Start();
bool Initialize();  // return false on error
void PauseResume();

} // namespace LedAlive
} // namespace Hardware

