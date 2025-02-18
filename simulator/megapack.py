from enums import State
import time 
import datetime

class Megapack: 
    def __init__(self):
        self.state = State.STARTUP
        self.id = "new-simulated-megapack-1"
        self.running_time = 0
        self.charge = 1000
        self.power = 0 
        self.internal_temp = 22
        self.ambient_temp = 22
        self.cumulative_output = 0
        self.cycle = 1

    def loop(self): 
        time.sleep(1)
        if self.state == State.STARTUP: 
            self.startup()
        if self.state == State.CHARGING: 
            self.charging()
        if self.state == State.DISCHARGING: 
            self.discharging()
        if self.state == State.IDLE: 
            self.state = State.STARTUP
        if self.state == State.SHUTDOWN: 
            self.state = State.STARTUP
        if self.state == State.FAULT: 
            return State.MAINTENANCE
        if self.state == State.MAINTENANCE: 
            return State.STARTUP
        
    def startup(self):
        if self.charge > 200: 
            self.state = State.DISCHARGING
        else: 
            self.state = State.CHARGING
    

    def charging(self): 
        if self.charge > 900: 
            self.state = State.DISCHARGING
            self.cycle += 1
        else:
            self.internal_temp *= 1.05
            self.charge += 20
    
    
    def discharging(self): 
        if self.charge < 100: 
            self.state = State.CHARGING
        else:
            self.charge -= 20
            self.internal_temp *= 1.05

    def get_data(self):
        return {
            "unit_id": self.id,
            "timestamp": datetime.datetime.now().isoformat(),
            "temperature_celsius": self.internal_temp,
            "charge_level_percent": self.charge,
            "charge_cycle": self.cycle,
            "cumulative_power": self.cumulative_output,
        }
