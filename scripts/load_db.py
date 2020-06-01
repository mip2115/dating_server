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

# users
u1 = {
    "uuid": "eab85cb1-0a11-47d1-890d-93015dc1e6fa", "email": "testuser1@test.com", "password": "test123", "first_name": "michael",
    "last_name": "daniels", "mobile": "0001112222", "dob": "08/31/1998", "gender": "male", "age": 22,
    "drink": True, "smoke": False, "job": "janitor", "university": "Columbia University",
    "politics": "liberal", "religion": "judaism", "hometown": "West Hartford", "partner_gender": "female",
    "purpose": "casual", "city": "New York", "profile_image": "nil", "images": [],
    "last_update": "date", "past_dates": [], "future_dates": [], "matches": [], "recently_matched": [], "blocked_users": [], "userslikedme": []
}

u2 = {
    "uuid": "eab85cb1-0a11-47d1-890d-93015dc1e6fb", "email": "testuser2@test.com", "password": "test123", "first_name": "rob",
    "last_name": "smith", "mobile": "0001112222", "dob": "08/31/1998", "gender": "male", "age": 22,
    "drink": True, "smoke": False, "job": "janitor", "university": "Columbia University",
    "politics": "conservative", "religion": "judaism", "hometown": "West Hartford", "partner_gender": "female",
    "purpose": "casual", "city": "New York", "profile_image": "nil", "images": [],
    "last_update": "date", "past_dates": [], "future_dates": [], "matches": [], "recently_matched": [], "blocked_users": [], "userslikedme": []
}

u3 = {
    "uuid": "eab85cb1-0a11-47d1-890d-93015dc1e6fc", "email": "testuser3@test.com", "password": "test123", "first_name": "alex",
    "last_name": "jones", "mobile": "0001112222", "dob": "08/31/1998", "gender": "male", "age": 22,
    "drink": True, "smoke": False, "job": "janitor", "university": "Columbia University",
    "politics": "liberal", "religion": "hindu", "hometown": "West Hartford", "partner_gender": "female",
    "purpose": "casual", "city": "New York", "profile_image": "nil", "images": [],
    "last_update": "date", "past_dates": [], "future_dates": [], "matches": [], "recently_matched": [], "blocked_users": [], "userslikedme": []
}

u4 = {
    "uuid": "eab85cb1-0a11-47d1-890d-93015dc1e6fd", "email": "testuser4@test.com", "password": "test123", "first_name": "clara",
    "last_name": "mince", "mobile": "0001112222", "dob": "08/31/1998", "gender": "female", "age": 22,
    "drink": True, "smoke": False, "job": "janitor", "university": "Columbia University",
    "politics": "liberal", "religion": "hindu", "hometown": "West Hartford", "partner_gender": "male",
    "purpose": "casual", "city": "New York", "profile_image": "nil", "images": [],
    "last_update": "date", "past_dates": [], "future_dates": [], "matches": [], "recently_matched": [], "blocked_users": [], "userslikedme": []
}

u5 = {
    "uuid": "eab85cb1-0a11-47d1-890d-93015dc1e6fe", "email": "testuser5@test.com", "password": "test123", "first_name": "steph",
    "last_name": "jones", "mobile": "0001112222", "dob": "08/31/1998", "gender": "female", "age": 22,
    "drink": True, "smoke": False, "job": "janitor", "university": "Columbia University",
    "politics": "liberal", "religion": "hindu", "hometown": "West Hartford", "partner_gender": "male",
    "purpose": "casual", "city": "New York", "profile_image": "nil", "images": [],
    "last_update": "date", "past_dates": [], "future_dates": [], "matches": [], "recently_matched": [], "blocked_users": [], "userslikedme": []
}

u6 = {
    "uuid": "eab85cb1-0a11-47d1-890d-93015dc1e6ff", "email": "testuser6@test.com", "password": "test123", "first_name": "rick",
    "last_name": "sanchez", "mobile": "0001112222", "dob": "08/31/1998", "gender": "male", "age": 22,
    "drink": True, "smoke": False, "job": "janitor", "university": "New York University",
    "politics": "liberal", "religion": "hindu", "hometown": "West Hartford", "partner_gender": "female",
    "purpose": "casual", "city": "New York", "profile_image": "nil", "images": [],
    "last_update": "date", "past_dates": [], "future_dates": [], "matches": [], "recently_matched": [], "blocked_users": [], "userslikedme": ['eab85cb1-0a11-47d1-890d-93015dc1e6fe']
}


u7 = {
    "uuid": "eab85cb1-0a11-47d1-890d-93015dc1e6fg", "email": "testuser7@test.com", "password": "test123", "first_name": "sarah",
    "last_name": "redd", "mobile": "0001112222", "dob": "08/31/1998", "gender": "female", "age": 22,
    "drink": True, "smoke": False, "job": "janitor", "university": "New York University",
    "politics": "liberal", "religion": "hindu", "hometown": "West Hartford", "partner_gender": "male",
    "purpose": "casual", "city": "New York", "profile_image": "nil", "images": [],
    "last_update": "date", "past_dates": [], "future_dates": [], "matches": [], "recently_matched": [], "blocked_users": [], "userslikedme": ['eab85cb1-0a11-47d1-890d-93015dc1e6fa', 'eab85cb1-0a11-47d1-890d-93015dc1e6fb']
}

u8 = {
    "partnerGender": "male",
    "drink": "sometimes",
    "smoke": "sometimes",
    "city": "New York",
    "religion": "judaism",
    "university": "New York University",
    "politics": "conservative",
    "uuid": "eab85cb1-0a11-47d1-890d-93015dc1e6fz",
    "email": "testuser8@test.com",
    "password": "test123",
    "first_name": "raquel",
    "last_name": "snell",
    "mobile": "0009998888",
    "dob": "08/31/80",
    "gender": "female",
    "age": 20,
    "job": "cleaner",
    "hometown": "west hartford",
    "purpose": "casual",
    "profile_image": "nil",
    "images": [],
    "last_update": "date",
    "past_dates": [],
    "future_dates": [],
    "matches": [],
    "recently_matched": [],
    "blocked_users": [],
    "userslikedme": [],
}


# matches
match1 = {
    "uuid": "eab85cb1-0a11-47d1-890d-93015dc1e6fh",
    "user_a": "eab85cb1-0a11-47d1-890d-93015dc1e6fg",
    "user_b": "eab85cb1-0a11-47d1-890d-93015dc1e6ff",
    "date_created": "1598861408",
    "messagesID": ["eab85cb1-0a11-47d1-890d-93015dc1e6fi"]
}

# msgs
msg1 = {
    "uuid": "eab85cb1-0a11-47d1-890d-93015dc1e6fi",
    "userFrom": "eab85cb1-0a11-47d1-890d-93015dc1e6fg",
    "userTo": "eab85cb1-0a11-47d1-890d-93015dc1e6ff",
    "dateCreated": "1601541002",
    "content": "Hey whats up?"
}

msg2 = {
    "uuid": "eab85cb1-0a11-47d1-890d-93015dc1e6fj",
    "userFrom": "eab85cb1-0a11-47d1-890d-93015dc1e6ff",
    "userTo": "eab85cb1-0a11-47d1-890d-93015dc1e6fg",
    "dateCreated": "1603182600",
    "content": "All good"
}


users = db['users']
inserts = [
    u1, u2, u3, u4, u5, u6, u7, u8
]
result = users.insert_many(inserts)
