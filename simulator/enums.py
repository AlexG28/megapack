from enum import Enum

class State(Enum): 
    STARTUP = "startup"
    CHARGING = "charging"
    DISCHARGING = "discharging"
    SHUTDOWN = "shutdown"
    IDLE = "idle"
    FAULT = "fault"
    MAINTENANCE = "maintenance"