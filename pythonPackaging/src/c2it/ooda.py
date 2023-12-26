import threading
import queue


class OperationPicture:
    def __init__(self):
        self.data = {}
        self.lock = threading.Lock()

    def update(self, key, value):
        with self.lock:
            self.data[key] = value


class OODALoop:
    def __init__(self, operationPicture):
        self.operationPicture = operationPicture
        self.observation_queue = queue.Queue()
        self.orientation_queue = queue.Queue()
        self.decision_queue = queue.Queue()
        self.action_queue = queue.Queue()
        self.is_running = True
        # Initialize the thread attributes here but do not start them yet
        self.observe_thread = threading.Thread(target=self.observe, daemon=True)
        self.orient_thread = threading.Thread(target=self.orient, daemon=True)
        self.decide_thread = threading.Thread(target=self.decide, daemon=True)
        self.act_thread = threading.Thread(target=self.act, daemon=True)

    def observe(self):
        # Observing the situation
        while self.is_running:
            newObservations = {}
            self.observation_queue.put(newObservations)

    def orient(self):
        # Wait for a new observation and then process it
        while self.is_running:
            observation = self.observation_queue.get()
            # Now results of observation are available and can be added to Operation Picture
            self.operationPicture.update('Observations', observation)
            orientation = {}
            self.operationPicture.update('Orientation', orientation)
            self.orientation_queue.put(True)

    def decide(self):
        # Wait for the signal of orient step finished
        while self.is_running:
            self.orientation_queue.get()
            decision = {}
            self.operationPicture.update('Decision', decision)
            self.decision_queue.put(True)

    def act(self):
        # Wait for the signal of decide step finished
        while self.is_running:
            self.decision_queue.get()
            actions = {}
            self.operationPicture.update('Action', actions)

    def start(self):
        # Start the threads
        self.observe_thread.start()
        self.orient_thread.start()
        self.decide_thread.start()
        self.act_thread.start()

    def stop(self):
        self.is_running = False
        self.observe_thread.join()
        self.orient_thread.join()
        self.decide_thread.join()
        self.act_thread.join()


# Here's your main function
if __name__ == "__main__":
    operation_picture = OperationPicture()
    ooda = OODALoop(operation_picture)
    try:
        ooda.start()
    except KeyboardInterrupt:
        # when you press Ctrl+C on your keyboard to break the program execution
        ooda.stop()
