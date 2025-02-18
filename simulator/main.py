import requests
from faker import Faker
from faker.providers import internet
import json
import signal
from megapack import Megapack

fake = Faker()

should_exit = False

def handle_exit_signal(signum, frame):
    global should_exit 
    print("Received exit signal. Shutting down...")
    should_exit = True

signal.signal(signal.SIGTERM, handle_exit_signal)
signal.signal(signal.SIGINT, handle_exit_signal)

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
    megapack = Megapack()

    while not should_exit: 
        megapack.loop()
        data = megapack.get_data()
        send_data_to_api_gateway(data)
