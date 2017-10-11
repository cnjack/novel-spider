# -*- coding: utf-8 -*-
#auth jack
#date 03/31/2017

import json
import gnsq
import base64

def pub_message(msg):
    writer.publish('novel', bytes(msg))
    pass

def run():
    for line in file:
        novel = json.loads(line)
        task = dict()
        task["ID"] = 0
        task["CreatedAt"] = '0001-01-01T00:00:00Z'
        task["UpdatedAt"] = '0001-01-01T00:00:00Z'
        task["DeletedAt"] = None
        task["TType"] = 0
        task["Url"] = novel["url"]
        task["Status"] = 0
        task["Times"] = -1
        task["TargetID"] = 0
        msg = dict()
        msg["Header"] = None
        msg["Body"] = base64.b64encode(json.dumps(task, ensure_ascii=False))
        pub_message(json.dumps(msg, ensure_ascii=False))
    pass

file = open('../items.jl', 'r')
writer = gnsq.Nsqd(address='127.0.0.1', http_port=32776,tcp_port=32777)
writer.connect()
run()