#!/usr/bin/env python
import pymongo
from pymongo import MongoClient
import sys


env_file = "../.env"
env = {}
with open(env_file) as f:
    for line in f:
        if line.startswith('#') or not line.strip():
            continue

        key, value = line.strip().split('=', 1)
        # os.environ[key] = value  # Load to local environ
        env[key] = value

client = MongoClient()
client = MongoClient('localhost', 27017)
db = client['kama']

collections = [collection for collection in db.list_collection_names()
               if not collection.startswith('system.')]

for c in collections:
    db[c].drop()
