from enum import Enum, auto

class State(Enum): 
    STARTUP = auto()
    CHARGING = auto()
    DISCHARGING = auto()
    SHUTDOWN = auto()
    IDLE = auto()
    FAULT = auto()
    MAINTENANCE = auto()