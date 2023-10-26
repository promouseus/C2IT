"""Create and return HTTPSConnection to the push service."""
# open() is a basic function but this is more failsafe as everything is handled,
# for example File descriptors
import fileinput
from datetime import datetime, timezone
import time

import random # For testing errors

import http.client # https://docs.python.org/3/library/http.client.html
import os # https://docs.python.org/3/library/os.html
import ssl # https://docs.python.org/3/library/ssl.html
# Only needed if Pyhton installed in VENV (so paths are not correct)
import certifi # Seems required in venv, maybe not on normal OS installation? This needs to be installed with PIP (not part of Python lib).
# print(ssl.get_default_verify_paths()) shows set path, if in VENV then ENV vars from root OS can mess this up.
# Can verify if this file exists as work arround (so no need for pip certifi)

#print(sysconfig.get_config_vars())
# Fancy pritn with: python -m sysconfig

#print(certifi.where())

def add_one(number):
    return number + 1

def main():
    """Create and return HTTPSConnection to the push service."""
    #print("Promouseus")
    # https://www.idnt.net/en-US/kb/941772
    # ticks = os.sysconf(os.sysconf_names['SC_CLK_TCK']) # USER_HZ or Jiffies (normally 1/100 second)

    context = ssl.SSLContext(ssl.PROTOCOL_TLS)
    context.verify_mode = ssl.CERT_REQUIRED
    context.check_hostname = True
    context.minimum_version = ssl.TLSVersion.TLSv1_2
    ca_file = os.path.relpath(certifi.where())
    context.load_verify_locations(ca_file)

    headers = {"Content-type": "application/json"}
    
    #r1 = conn.getresponse()
    #print(r1.status, r1.reason)


    # https://developer.mozilla.org/en-US/docs/Web/HTTP/Methods
    #HEAD
    #DELETE
    #PATCH
    #OPTIONS
    # POST: each time is processed differntly
    # PUT: idempotent, same result for each resend of same data (for example use Unix Timestamp in message to check if data is already ingested at server side)
    # GET: use to get pointer of last Unix Timestamp and resends PUT if behind (in this case queue at server side will never result in duplicates)

    #print("The type of the secure socket created: " + str(context.sslsocket_class))
    #print("Maximum version of the TLS: " + str(context.maximum_version))
    #print("Minimum version of the TLS: " + str(context.minimum_version))
    #print("SSL options enabled in the context object: " + str(context.options));
    #print("Protocol set in the context: " + str(context.protocol))
    #print("Verify flags for certificates: " + str(context.verify_flags))
    #print("Verification mode(how to validate peer's certificate and handle failures if any):" + str(context.verify_mode))
    #print("SNI support: " + str(ssl.HAS_SNI))

    #print("OS name: " + str(os.uname()))

    #import socket
    #print(os.uname())
    #print(socket.gethostbyaddr(socket.gethostname())) # Not working on Mac it seems

    procStat = './proc-stat.txt' # Procstat dummy file for testing
    #procStat = '/proc/stat'

    starttime = time.time()
    record_counter = 0
    record_flush = 60 # Try flush after i new records/seconds
    record_backlog = 1440 * record_flush # 60 seconds for 24 hours backlog 

    record = ""
    while True:
        now = datetime.now(timezone.utc).isoformat(timespec="milliseconds")

        try:
            with fileinput.input(files=(procStat), mode='r') as f:
                for line in f:
                    if line.startswith("cpu "):
                        cpu_string = line.split(maxsplit=11)
                        record += now + "," + ",".join(cpu_string[1:11]) + "\n"
            record_counter += 1
                        
        except Exception as error:
            print(str(error) + " at " + str(time.asctime( time.localtime(time.time()))) )

        if (record_counter >= record_flush):
            record_length = record.count('\n')

            try:
                #conn = http.client.HTTPSConnection("api.promouseus.io", context=context)
                #conn.request("PUT", "/", record, headers=headers,)
                if(bool(random.getrandbits(1))):
                    raise NameError('Simulated fail')
                print(record)
                record = ""
                print("Flushed " + str(record_length) + " records at " + str(time.asctime( time.localtime(time.time()))) )
                record_length = 0
                # response = conn.getresponse()
                # print(response.status, response.reason)
            except Exception as error:
                print(str(error) + " at " + str(time.asctime( time.localtime(time.time()))))
            finally:
                True
                #conn.close()

            record_counter = 0
            # Clean-up oldest logs if more then seconds/lines of (back)log
            if(record_length >= record_backlog):
                record = record.split("\n", record_flush)[record_flush]
                print("Deleted " + str(record_flush) + " records from backlog at " + str(time.asctime( time.localtime(time.time()))) )
            
        time.sleep(1.0 - ((time.time() - starttime) % 1.0))

# https://docs.python.org/3/library/
