/*
 * Copyright (C) 2022 Hopsworks AB
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 2
 * of the License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, write to the Free Software
 * Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301,
 * USA.
 */
#ifndef STORAGE_NDB_REST_SERVER_DATA_ACCESS_RONDB_SRC_LOGGER_HPP_
#define STORAGE_NDB_REST_SERVER_DATA_ACCESS_RONDB_SRC_LOGGER_HPP_

#include <string.h>
#include <iostream>
#include <string>
#include "src/rdrs-dal.h"

#define PanicLevel 0
#define FatalLevel 1
#define ErrorLevel 2
#define WarnLevel  3
#define InfoLevel  4
#define DebugLevel 5
#define TraceLevel 6

inline void log(const int level, const char *msg) {
  std::cout << "Log Level: " + std::to_string(level) + "; Message: " + msg << std::endl;
}

inline void PANIC(const char *msg) {
  log(PanicLevel, msg);
}

inline void PANIC(const std::string msg) {
  log(PanicLevel, msg.c_str());
}

inline void FATAL(const char *msg) {
  log(FatalLevel, msg);
}

inline void FATAL(const std::string msg) {
  log(FatalLevel, msg.c_str());
}

inline void ERROR(const char *msg) {
  log(ErrorLevel, msg);
}

inline void ERROR(const std::string msg) {
  log(ErrorLevel, msg.c_str());
}

inline void WARN(const char *msg) {
  log(WarnLevel, msg);
}

inline void WARN(const std::string msg) {
  log(WarnLevel, msg.c_str());
}

inline void INFO(const char *msg) {
  log(InfoLevel, msg);
}

inline void INFO(const std::string msg) {
  log(InfoLevel, msg.c_str());
}

inline void DEBUG(const char *msg) {
  log(DebugLevel, msg);
}

inline void DEBUG(const std::string msg) {
  log(DebugLevel, msg.c_str());
}

inline void TRACE(char *msg) {
  log(TraceLevel, msg);
}

#endif  // STORAGE_NDB_REST_SERVER_DATA_ACCESS_RONDB_SRC_LOGGER_HPP_
