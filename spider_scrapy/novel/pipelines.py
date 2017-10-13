# -*- coding: utf-8 -*-

# Define your item pipelines here
#
# Don't forget to add your pipeline to the ITEM_PIPELINES setting
# See: http://doc.scrapy.org/en/latest/topics/item-pipeline.html

import json
import pymysql.cursors


class NovelPipeline(object):
    def __init__(self):
        self.file = open('items.jl', 'wb')

    def process_item(self, item, spider):
        line = json.dumps(dict(item), ensure_ascii=False) + "\n"
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
        sql = "INSERT INTO `novels`(title, url, style, auth, status, cover, introduction, created_at, updated_at) VALUES (%s, %s, %s, %s, %s, %s, %s, now(), now())"
        params = (item["title"], item["url"], item["style"], item["auth"], item["status"], item["cover"], item['intro'])
        try:
            with db_conn.cursor() as cursor:
                cursor.execute(sql, params)
            db_conn.commit()
        finally:
            pass
            # db_conn.close()

    def process_item(self, item, spider):
        self._conditional_insert(self.db_conn, item)
        return item
        pass
