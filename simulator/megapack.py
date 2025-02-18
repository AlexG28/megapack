from enums import State
import datetime, random

class Megapack: 
    def __init__(self, id):
        self.state = State.STARTUP
        self.id = id
        self.running_time = 0
        self.charge = 1000
        self.power = 0 
        self.internal_temp = 22
        self.ambient_temp = 22
        self.cumulative_output = 0
        self.cycle = 1

    def loop(self): 
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
        elif random.random() < 0.005:
            self.state = State.FAULT
        else: 
            self.state = State.IDLE
    

    def charging(self): 
        if self.charge > 900: 
            self.state = State.IDLE
            self.cycle += 1
        elif random.random() < 0.005:
            self.state = State.FAULT
        else:
            self.internal_temp *= 1.05
            self.charge += 20
    
    def fault(self): 
        self.internal_temp *= 0.95
        if random.random() < 0.3: 
            self.state = State.MAINTENANCE
    
    def maintenance(self): 
        self.internal_temp *= 0.95
        if random.random() < 0.2: 
            self.state = State.STARTUP
    
    def idle(self): 
        self.internal_temp *= 0.95
        if self.charge < 200 and random.random() < 0.7:
            self.state = State.CHARGING
        elif random.random() < 0.5:
            self.state = State.DISCHARGING

    def discharging(self): 
        if self.charge < 100: 
            self.state = State.CHARGING
        elif random.random() < 0.005:
            self.state = State.FAULT
        else:
            self.charge -= 20
            self.internal_temp *= 1.05

    def get_data(self):
        return {
            "unit_id": self.id,
            "state": self.state.value,
            "timestamp": datetime.datetime.now().isoformat(),
            "temperature_celsius": self.internal_temp,
            "charge_level_percent": self.charge,
            "charge_cycle": self.cycle,
            "cumulative_power": self.cumulative_output,
        }
