# CMAKE generated file: DO NOT EDIT!
# Generated by "Unix Makefiles" Generator, CMake Version 3.13

# Delete rule output on recipe failure.
.DELETE_ON_ERROR:


#=============================================================================
# Special targets provided by cmake.

# Disable implicit rules so canonical targets will work.
.SUFFIXES:


# Remove some rules from gmake that .SUFFIXES does not remove.
SUFFIXES =

.SUFFIXES: .hpux_make_needs_suffix_list


# Suppress display of executed commands.
$(VERBOSE).SILENT:


# A target that is always out of date.
cmake_force:

.PHONY : cmake_force

#=============================================================================
# Set environment variables for the build.

# The shell in which to execute make rules.
SHELL = /bin/sh

# The CMake executable.
CMAKE_COMMAND = /Applications/CLion.app/Contents/bin/cmake/mac/bin/cmake

# The command to remove a file.
RM = /Applications/CLion.app/Contents/bin/cmake/mac/bin/cmake -E remove -f

# Escaping for special characters.
EQUALS = =

# The top-level source directory on which CMake was run.
CMAKE_SOURCE_DIR = "/Users/wyn/Library/Mobile Documents/com~apple~CloudDocs/Daocloud/GOProject/src/github.com/IsolationWyn/paddle"

# The top-level build directory on which CMake was run.
CMAKE_BINARY_DIR = "/Users/wyn/Library/Mobile Documents/com~apple~CloudDocs/Daocloud/GOProject/src/github.com/IsolationWyn/paddle/cmake-build-debug"

# Include any dependencies generated for this target.
include nsenter/CMakeFiles/nsenter.dir/depend.make

# Include the progress variables for this target.
include nsenter/CMakeFiles/nsenter.dir/progress.make

# Include the compile flags for this target's objects.
include nsenter/CMakeFiles/nsenter.dir/flags.make

nsenter/CMakeFiles/nsenter.dir/nsenter.cc.o: nsenter/CMakeFiles/nsenter.dir/flags.make
nsenter/CMakeFiles/nsenter.dir/nsenter.cc.o: ../nsenter/nsenter.cc
	@$(CMAKE_COMMAND) -E cmake_echo_color --switch=$(COLOR) --green --progress-dir="/Users/wyn/Library/Mobile Documents/com~apple~CloudDocs/Daocloud/GOProject/src/github.com/IsolationWyn/paddle/cmake-build-debug/CMakeFiles" --progress-num=$(CMAKE_PROGRESS_1) "Building CXX object nsenter/CMakeFiles/nsenter.dir/nsenter.cc.o"
	cd "/Users/wyn/Library/Mobile Documents/com~apple~CloudDocs/Daocloud/GOProject/src/github.com/IsolationWyn/paddle/cmake-build-debug/nsenter" && /Library/Developer/CommandLineTools/usr/bin/c++  $(CXX_DEFINES) $(CXX_INCLUDES) $(CXX_FLAGS) -o CMakeFiles/nsenter.dir/nsenter.cc.o -c "/Users/wyn/Library/Mobile Documents/com~apple~CloudDocs/Daocloud/GOProject/src/github.com/IsolationWyn/paddle/nsenter/nsenter.cc"

nsenter/CMakeFiles/nsenter.dir/nsenter.cc.i: cmake_force
	@$(CMAKE_COMMAND) -E cmake_echo_color --switch=$(COLOR) --green "Preprocessing CXX source to CMakeFiles/nsenter.dir/nsenter.cc.i"
	cd "/Users/wyn/Library/Mobile Documents/com~apple~CloudDocs/Daocloud/GOProject/src/github.com/IsolationWyn/paddle/cmake-build-debug/nsenter" && /Library/Developer/CommandLineTools/usr/bin/c++ $(CXX_DEFINES) $(CXX_INCLUDES) $(CXX_FLAGS) -E "/Users/wyn/Library/Mobile Documents/com~apple~CloudDocs/Daocloud/GOProject/src/github.com/IsolationWyn/paddle/nsenter/nsenter.cc" > CMakeFiles/nsenter.dir/nsenter.cc.i

nsenter/CMakeFiles/nsenter.dir/nsenter.cc.s: cmake_force
	@$(CMAKE_COMMAND) -E cmake_echo_color --switch=$(COLOR) --green "Compiling CXX source to assembly CMakeFiles/nsenter.dir/nsenter.cc.s"
	cd "/Users/wyn/Library/Mobile Documents/com~apple~CloudDocs/Daocloud/GOProject/src/github.com/IsolationWyn/paddle/cmake-build-debug/nsenter" && /Library/Developer/CommandLineTools/usr/bin/c++ $(CXX_DEFINES) $(CXX_INCLUDES) $(CXX_FLAGS) -S "/Users/wyn/Library/Mobile Documents/com~apple~CloudDocs/Daocloud/GOProject/src/github.com/IsolationWyn/paddle/nsenter/nsenter.cc" -o CMakeFiles/nsenter.dir/nsenter.cc.s

# Object files for target nsenter
nsenter_OBJECTS = \
"CMakeFiles/nsenter.dir/nsenter.cc.o"

# External object files for target nsenter
nsenter_EXTERNAL_OBJECTS =

nsenter/libnsenter.a: nsenter/CMakeFiles/nsenter.dir/nsenter.cc.o
nsenter/libnsenter.a: nsenter/CMakeFiles/nsenter.dir/build.make
nsenter/libnsenter.a: nsenter/CMakeFiles/nsenter.dir/link.txt
	@$(CMAKE_COMMAND) -E cmake_echo_color --switch=$(COLOR) --green --bold --progress-dir="/Users/wyn/Library/Mobile Documents/com~apple~CloudDocs/Daocloud/GOProject/src/github.com/IsolationWyn/paddle/cmake-build-debug/CMakeFiles" --progress-num=$(CMAKE_PROGRESS_2) "Linking CXX static library libnsenter.a"
	cd "/Users/wyn/Library/Mobile Documents/com~apple~CloudDocs/Daocloud/GOProject/src/github.com/IsolationWyn/paddle/cmake-build-debug/nsenter" && $(CMAKE_COMMAND) -P CMakeFiles/nsenter.dir/cmake_clean_target.cmake
	cd "/Users/wyn/Library/Mobile Documents/com~apple~CloudDocs/Daocloud/GOProject/src/github.com/IsolationWyn/paddle/cmake-build-debug/nsenter" && $(CMAKE_COMMAND) -E cmake_link_script CMakeFiles/nsenter.dir/link.txt --verbose=$(VERBOSE)

# Rule to build all files generated by this target.
nsenter/CMakeFiles/nsenter.dir/build: nsenter/libnsenter.a

.PHONY : nsenter/CMakeFiles/nsenter.dir/build

nsenter/CMakeFiles/nsenter.dir/clean:
	cd "/Users/wyn/Library/Mobile Documents/com~apple~CloudDocs/Daocloud/GOProject/src/github.com/IsolationWyn/paddle/cmake-build-debug/nsenter" && $(CMAKE_COMMAND) -P CMakeFiles/nsenter.dir/cmake_clean.cmake
.PHONY : nsenter/CMakeFiles/nsenter.dir/clean

nsenter/CMakeFiles/nsenter.dir/depend:
	cd "/Users/wyn/Library/Mobile Documents/com~apple~CloudDocs/Daocloud/GOProject/src/github.com/IsolationWyn/paddle/cmake-build-debug" && $(CMAKE_COMMAND) -E cmake_depends "Unix Makefiles" "/Users/wyn/Library/Mobile Documents/com~apple~CloudDocs/Daocloud/GOProject/src/github.com/IsolationWyn/paddle" "/Users/wyn/Library/Mobile Documents/com~apple~CloudDocs/Daocloud/GOProject/src/github.com/IsolationWyn/paddle/nsenter" "/Users/wyn/Library/Mobile Documents/com~apple~CloudDocs/Daocloud/GOProject/src/github.com/IsolationWyn/paddle/cmake-build-debug" "/Users/wyn/Library/Mobile Documents/com~apple~CloudDocs/Daocloud/GOProject/src/github.com/IsolationWyn/paddle/cmake-build-debug/nsenter" "/Users/wyn/Library/Mobile Documents/com~apple~CloudDocs/Daocloud/GOProject/src/github.com/IsolationWyn/paddle/cmake-build-debug/nsenter/CMakeFiles/nsenter.dir/DependInfo.cmake" --color=$(COLOR)
.PHONY : nsenter/CMakeFiles/nsenter.dir/depend

