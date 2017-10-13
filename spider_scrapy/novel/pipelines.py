# -*- coding: utf-8 -*-

# Define your item pipelines here
#
# Don't forget to add your pipeline to the ITEM_PIPELINES setting
# See: http://doc.scrapy.org/en/latest/topics/item-pipeline.html

import json
import pymysql.cursors
from twisted.enterprise import adbapi


class NovelPipeline(object):
    def __init__(self):
        self.file = open('items.jl', 'wb')

    def process_item(self, item, spider):
        line = json.dumps(dict(item), ensure_ascii=False) + "\n"
        #print(line)
        self.file.write(line.encode())
        return item


class MysqlPipeline(object):
    def __init__(self, db_conn):
        self.db_conn = db_conn
        pass

    @classmethod
    def from_settings(cls, settings):
        db_conn = pymysql.connect(host=settings['MYSQL_HOST'],
                                     user=settings['MYSQL_USER'],
                                     password=settings['MYSQL_PASSWORD'],
                                     db=settings['MYSQL_DB_NAME'],
                                     charset='utf8mb4',
                                     cursorclass=pymysql.cursors.DictCursor)
        return cls(db_conn)

    def _conditional_insert(self, db_conn, item):
        sql = "insert into `novel`(`name`, `url`) values(%s, %s)"
        params = (item["name"], item["url"])
        try:
            with db_conn.cursor() as cursor:
                cursor.execute(sql, params)
            db_conn.commit()
        finally:
            db_conn.close()

    def process_item(self, item, spider):
        self._conditional_insert(self.db_conn, item)
        return item
        pass
