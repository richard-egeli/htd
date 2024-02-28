
target_include_directories(app PRIVATE .)

target_sources(app PRIVATE ${CMAKE_CURRENT_SOURCE_DIR}/AliveLed.cpp)
target_sources(app PRIVATE ${CMAKE_CURRENT_SOURCE_DIR}/Application.cpp)
#target_sources(app PRIVATE ${CMAKE_CURRENT_SOURCE_DIR}/task.cpp)
target_sources(app PRIVATE ${CMAKE_CURRENT_SOURCE_DIR}/Sdi12.cpp)
target_sources(app PRIVATE ${CMAKE_CURRENT_SOURCE_DIR}/Modem.cpp)

