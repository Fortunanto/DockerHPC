import logging
import threading
from importlib import __import__
import sys
import os
import dag

from watchdog.events import FileSystemEventHandler, FileSystemEvent


def _filter_relevant_events(event: FileSystemEvent) -> bool:
    return not event.is_directory and event.src_path.endswith("dag.py")


def get_all_dag_structures(file_name):
    module_name = os.path.splitext(file_name)[0]
    dag_dir = os.environ['HOME'] + "/.dag_home"
    if 'DAG_HOME' in os.environ:
        dag_dir = os.environ['DAG_HOME']
    sys.path.insert(0, dag_dir)
    dags = []
    imported_module = __import__(module_name)
    keys = imported_module.__dict__.keys()
    for key in keys:
        possible_dag = vars(imported_module)[key]
        if isinstance(possible_dag, dag.DAG):
            dags.append(possible_dag)
    sys.path.pop(0)
    return dags


class DagFileEventHandler(FileSystemEventHandler):
    def __init__(self):
        self.dag_container = {}
        self.filename_dag_id = {}

        self._dag_lock = threading.Lock()

    def add_current_dir(self, directory):
        paths = [path for path in os.listdir(directory) if path.endswith("dag.py")]
        for path in paths:
            self._add_dag(path)

    def on_any_event(self, event: FileSystemEvent):
        if _filter_relevant_events(event):
            if event.event_type == 'deleted':
                self._remove_dag(event.src_path)
                print(f"current container: {self.dag_container}, file: {event.src_path} removed")
            else:
                self._add_dag(event.src_path)
                print(f"current container: {self.dag_container}, file: {event.src_path} added")
        super().on_any_event(event)

    def _remove_dag(self, path):
        self._dag_lock.acquire()
        try:
            if path in self.filename_dag_id:
                dag_ids = self.filename_dag_id[path]
                for dag_id in dag_ids:
                    self.dag_container.pop(dag_id, None)
                self.filename_dag_id.pop(path, None)
        finally:
            self._dag_lock.release()

    def _add_dag(self, path):
        self._dag_lock.acquire()
        try:
            if path not in self.filename_dag_id:
                dags = get_all_dag_structures(path.split("/")[-1])
                dags = {dag.dag_id: dag for dag in dags}
                is_all_ids_unique = all([not (dag_id in self.dag_container) for dag_id in dags.keys()])
                if is_all_ids_unique:
                    self.filename_dag_id[path] = [dag for dag in dags.keys()]
                    for dag_id in dags.keys():
                        self.dag_container[dag_id] = dags[dag_id]
                else:
                    non_unique_dag_id = [dag_id for dag_id in dags.keys() if dag_id in self.dag_container]
                    logging.error(f"Dag id's {non_unique_dag_id} already exists, File: {path}")
        finally:
            self._dag_lock.release()

    def get_current_dags(self):
        self._dag_lock.acquire()
        try:
            return self.dag_container
        finally:
            self._dag_lock.release()
