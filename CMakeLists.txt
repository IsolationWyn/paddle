cmake_minimum_required(VERSION 3.13)
project(paddlec)

set(CMAKE_C_STANDARD 99)

include_directories(
        "/usr/local/include"
        "/usr/include"
        "/Library/Developer/CommandLineTools/usr/lib/clang/10.0.0/include/"
        "/Applications/Xcode.app/Contents/Developer/Toolchains/XcodeDefault.xctoolchain/usr/include/c++/v1"
)

aux_source_directory(. DIR_SRCS)

add_subdirectory(nsenter)

add_executable(paddlec ${DIR_SRCS})

target_link_libraries(paddlec nsenter)