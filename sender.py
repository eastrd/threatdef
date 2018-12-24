from datetime import datetime
from time import sleep
import os
import requests



        
# Cowrie log file
file_path = "sample2.txt"
url = "http://127.0.0.1:8080/signal"
http_auth = {
    "user": "plus",
    "pass": "midoriya"
}

get_today = lambda: datetime.now().strftime("%Y-%m-%d")


# def log(message):
#     logger_path = "/".join(os.path.realpath(__file__).split("/")[:-1] + ["logger.log"])
#     with open(logger_path, "a+") as f:
#         f.write(message)
#         f.write("\n" * 3)


def send_signal(line):
    # Sends the cowrie log message to backend url
    while True:
        # Keep attempting to send message
        try:
            requests.post(url, data=line, timeout=3, auth=(http_auth["user"], http_auth["pass"]))
            break
        except Exception as e:
            print(str(e), flush=True)
            sleep(5)


def main(today):
    with open(file_path, "r") as f:
        while True:
            if get_today() != today:
                # Exit condition: The moment of midnight. (Date changed)
                return
            line = f.readline().replace("\n", "")
            if len(line) > 1:
                # Send the log to threatdef's backend
                print("Sending line with size: %s" % len(line), flush=True)
                send_signal(line)
            else:
                sleep(5)

if __name__ == "__main__":
    while True:
        try:
            main(get_today())
        except Exception as e:
            print("Something went wrong: " + str(e), flush=True)
            sleep(5)
