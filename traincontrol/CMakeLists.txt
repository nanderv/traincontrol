cmake_minimum_required(VERSION 3.25)
project(traincommunications CXX)

set(CMAKE_CXX_STANDARD 14)

include_directories(.)

add_executable(traincommunications
        coms.h
        coms.cpp
        hexConv.h
        hexConv.cpp
        message.h
        traincommunications.cpp)
