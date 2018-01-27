import time

from DagCreator.dagFileEventHandler import DagFileEventHandler
from watchdog.observers import Observer

if __name__ == "__main__":
    event_handler = DagFileEventHandler()
    event_handler.add_current_dir("/home/yiftach/.dag_home")
    observer = Observer()
    observer.schedule(event_handler, "/home/yiftach/.dag_home", recursive=False)
    observer.start()
    try:
        while True:
            time.sleep(1)
    except KeyboardInterrupt:
        observer.stop()
    observer.join()
