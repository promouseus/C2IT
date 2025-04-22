import socket
import threading

def probe_sensor():
    UDP_IP = "127.0.0.1"
    UDP_PORT = 5005

    sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    sock.bind((UDP_IP, UDP_PORT))

    while True:
        # echo -n hello3 | nc -4u -w0 127.0.0.1 5005
        data, addr = sock.recvfrom(2) # buffer size is 1024 bytes for netcat
        print("received IPv4 message from " + str(addr) + ": %s" % data)


def probe_sensor_v6():
    UDP_IPv6 = "::1"
    UDP_PORT = 5005

    sock2 = socket.socket(socket.AF_INET6, socket.SOCK_DGRAM)
    sock2.bind((UDP_IPv6, UDP_PORT))

    while True:
        # echo -n "hello3" | nc -6u -w0 ::1 5005
        data, addr = sock2.recvfrom(2) # buffer size is 1024 bytes for netcat
        print("received IPv6 message from " + str(addr) + ": %s" % data)
        #print(addr)

#v4 = threading.Thread(target=probe_sensor, args=(0,))
v4 = threading.Thread(target=probe_sensor)
v6 = threading.Thread(target=probe_sensor_v6)

v4.start()
v6.start()
