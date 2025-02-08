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

def send_data_to_api_gateway(data, api_url="http://api-gateway:8080/telemetry"): # Adjust port if needed
    headers = {
        "Content-Type": "application/json"
    }
    try:
        response = requests.post(api_url, data=json.dumps(data), headers=headers)
        response.raise_for_status() # Raise HTTPError for bad responses (4xx or 5xx)
        print(f"Data sent successfully for unit {data['unit_id']}. Status code: {response.status_code}")
    except requests.exceptions.RequestException as e:
        print(f"Error sending data for unit {data['unit_id']}: {e}")


# def temp(): 
#     url = "http://localhost:8080/telemetry"
    
#     headers = {
#         "Content-Type": "application/json"
#     }
    
#     payload = {
#         "unit_id": "manual-test-unit-001",
#         "timestamp": "2024-10-27T12:00:00Z",
#         "temperature_celsius": 22.5,
#         "voltage_volts": 485.0,
#         "charge_level_percent": 98.7
#     }

#     try:
#         response = requests.post(url, headers=headers, data=json.dumps(payload))
#         response.raise_for_status()  # Raises an HTTPError for bad responses (4xx or 5xx)
#         print(f"Success: {response.status_code}")
#         print(f"Response: {response.text}")
#     except requests.exceptions.RequestException as e:
#         print(f"Error: {e}")

if __name__ == "__main__":
    id = "simulated-megapack-001"
    
    while True: 
        data = simulate_megapack_data(id)
        print(data)
        send_data_to_api_gateway(data)
        time.sleep(1)

