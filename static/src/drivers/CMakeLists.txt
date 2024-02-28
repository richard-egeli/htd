#==================================================
# This folder has all Apparatus-v1 low level drivers
#==================================================

#add_subdirectory(xxxx)
target_include_directories(app PRIVATE .)

target_sources(app PRIVATE ${CMAKE_CURRENT_SOURCE_DIR}/button.cpp)
target_sources(app PRIVATE ${CMAKE_CURRENT_SOURCE_DIR}/eeprom.cpp)
target_sources(app PRIVATE ${CMAKE_CURRENT_SOURCE_DIR}/power.cpp)
target_sources(app PRIVATE ${CMAKE_CURRENT_SOURCE_DIR}/i420ma.cpp)