import requests
from faker import Faker
from faker.providers import internet
import time
import datetime
import random
import json
import signal

fake = Faker()

should_exit = False

def handle_exit_signal(signum, frame):
    global should_exit 
    print("Received exit signal. Shutting down...")
    should_exit = True

signal.signal(signal.SIGTERM, handle_exit_signal)
signal.signal(signal.SIGINT, handle_exit_signal)

def simulate_megapack_data(unit_id):
    timestamp = datetime.datetime.now().isoformat()
    
    temperature_celsius = 20 + random.uniform(-5, 5) # Example: 20 +/- 5 degrees C
    voltage_volts = 475 + random.uniform(-10, 10)   # Example: 475 +/- 10 Volts
    charge_level_percent = 30 + random.uniform(0, 70) # Example: Charge level between 30% and 100%

    data = {
        "unit_id": unit_id,
        "timestamp": timestamp,
        "temperature_celsius": temperature_celsius,
        "voltage_volts": voltage_volts,
        "charge_level_percent": charge_level_percent
    }
    return data

def send_data_to_api_gateway(data, api_url="http://api-gateway:8080/telemetry"): 
    headers = {
        "Content-Type": "application/json"
    }
    try:
        response = requests.post(api_url, data=json.dumps(data), headers=headers)
        response.raise_for_status() # Raise HTTPError for bad responses (4xx or 5xx)
        print(f"Data sent successfully for unit {data['unit_id']}. Status code: {response.status_code}")
    except requests.exceptions.RequestException as e:
        print(f"Error sending data for unit {data['unit_id']}: {e}")


if __name__ == "__main__":
    
    while not should_exit: 
        n = random.randint(1, 10)
        id = f"simulated-megapack-{n}"
        data = simulate_megapack_data(id)
        send_data_to_api_gateway(data)
        time.sleep(1)

    print("Simulator ending")

