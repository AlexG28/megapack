import requests
import json
import signal
from megapack import Megapack
import time
import os

should_exit = False
instance_count = os.environ.get('INSTANCE_COUNT')

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
        # print(f"Data sent successfully for unit {data['unit_id']}. Status code: {response.status_code}")
    except requests.exceptions.RequestException as e:
        print(f"Error sending data for unit {data['unit_id']}: {e}")


if __name__ == "__main__":
    megapack_array = []
    for i in range(int(instance_count)):
        new = Megapack(f"new-simulated-megapack-{i}")
        megapack_array.append(new)

    while not should_exit: 
        for pack in megapack_array: 
            pack.loop()
            data = pack.get_data() 
            send_data_to_api_gateway(data)
        
        time.sleep(1)
