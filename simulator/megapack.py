from enums import State
import datetime, random

class Megapack: 
    def __init__(self, id):
        self.state = State.STARTUP
        self.id = id
        self.running_time = 0
        self.charge = 1000  # Assuming 1000 = 100% charge
        self.power = 0 
        self.internal_temp = 22.0
        self.ambient_temp = 22.0
        self.cumulative_output = 0
        self.cycle = 1

    def loop(self): 
        prev_state = self.state
        
        if self.state == State.STARTUP: 
            self.startup()
        elif self.state == State.CHARGING: 
            self.charging()
        elif self.state == State.DISCHARGING: 
            self.discharging()
        elif self.state == State.IDLE: 
            self.idle()
        elif self.state == State.FAULT: 
            self.fault()
        elif self.state == State.MAINTENANCE: 
            self.maintenance()
        
        # Update running time if state didn't change to maintenance/reset
        if self.state not in {State.SHUTDOWN, State.MAINTENANCE}:
            self.running_time += 1

    def startup(self):
        if self.charge > 200:
            self.state = State.DISCHARGING
        elif random.random() < 0.005:
            self.state = State.FAULT
        else: 
            self.state = State.IDLE
        self._cool_down()

    def charging(self):
        if self.charge >= 900:  # Prevent overshooting
            self.state = State.IDLE
            self.cycle += 1
        elif random.random() < 0.005:
            self.state = State.FAULT
        else:
            self.charge = min(self.charge + 20, 1000)  # Cap at max charge
            self._heat_up()
            self.power = -20  # Power draw during charging

    def discharging(self):
        if self.charge <= 100:  # Prevent negative charge
            self.state = State.CHARGING
        elif random.random() < 0.005:
            self.state = State.FAULT
        else:
            self.charge = max(self.charge - 20, 0)
            self.cumulative_output += 20
            self._heat_up()
            self.power = 20  # Power output during discharging

    def idle(self):
        self._cool_down()
        self.power = 0
        
        if self.charge < 200:
            if random.random() < 0.7:
                self.state = State.CHARGING
        elif random.random() < 0.5:
            self.state = State.DISCHARGING

    def fault(self):
        self._cool_down()
        if random.random() < 0.3:
            self.state = State.MAINTENANCE

    def maintenance(self):
        self._cool_down()
        if random.random() < 0.2:
            self.state = State.STARTUP
            self.running_time = 0  # Reset counters after maintenance
            self.cycle = 1

    def _heat_up(self):
        self.internal_temp = min(self.internal_temp * 1.05, 100.0)  # Temp cap

    def _cool_down(self):
        self.internal_temp = max(self.internal_temp * 0.95, self.ambient_temp)

    def get_data(self):
        return {
            "unit_id": self.id,
            "state": self.state.value,
            "timestamp": datetime.datetime.now().isoformat(),
            "temperature": round(self.internal_temp, 1),
            "charge": self.charge,
            "cycle": self.cycle,
            "output": self.cumulative_output,
            "runtime": self.running_time,
            "power": self.power
        }
