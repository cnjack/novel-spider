# -*- coding: utf-8 -*-

# Define your item pipelines here
#
# Don't forget to add your pipeline to the ITEM_PIPELINES setting
# See: http://doc.scrapy.org/en/latest/topics/item-pipeline.html

import json
import MySQLdb
from twisted.enterprise import adbapi


class NovelPipeline(object):
    def __init__(self):
        self.file = open('items.jl', 'wb')

    def process_item(self, item, spider):
        line = json.dumps(dict(item), ensure_ascii=False) + "\n"
        self.file.write(line)
        return item


class MysqlPipeline(object):
    def __init__(self, db_pool):
        self.db_pool = db_pool
        pass

    @classmethod
    def from_settings(cls, settings):
        db_params = dict(
            host=settings['MYSQL_HOST'],
            db=settings['MYSQL_DB_NAME'],
            user=settings['MYSQL_USER'],
            passwd=settings['MYSQL_PASSWORD'],
            charset='utf8',
            cursorclass=MySQLdb.cursors.DictCursor,
            use_unicode=False,
        )
        db_pool = adbapi.ConnectionPool('MySQLdb', **db_params)
        return cls(db_pool)

    def _handle_error(self, failue, item, spider):
        print failue

    def _conditional_insert(self, tx, item):
        sql = "insert into novel(name,url) values(%s,%s)"
        params = (item["name"], item["url"])
        tx.execute(sql, params)

    def process_item(self, item, spider):
        query = self.db_pool.runInteraction(self._conditional_insert, item)  # 调用插入的方法
        query.addErrback(self._handle_error, item, spider)  # 调用异常处理方法
        return item
        pass
