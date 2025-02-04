import requests
from faker import Faker
from faker.providers import internet
import time
import random
import json

fake = Faker()

def simulate_megapack_data(unit_id):
    timestamp = fake.iso8601()
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

def send_data_to_api_gateway(data, api_url="http://localhost:8080/telemetry"): # Adjust port if needed
    headers = {'Content-type': 'application/json'} # Important: Tell API Gateway data is JSON
    try:
        response = requests.post(api_url, data=json.dumps(data), headers=headers)
        response.raise_for_status() # Raise HTTPError for bad responses (4xx or 5xx)
        print(f"Data sent successfully for unit {data['unit_id']}. Status code: {response.status_code}")
    except requests.exceptions.RequestException as e:
        print(f"Error sending data for unit {data['unit_id']}: {e}")

if __name__ == "__main__":
    id = "simulated-megapack-001"
    
    while True: 
        data = simulate_megapack_data(id)
        print(data)
        send_data_to_api_gateway(data)
        time.sleep(1)

