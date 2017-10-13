# -*- coding: utf-8 -*-

# Scrapy settings for novel project
#
# For simplicity, this file contains only settings considered important or
# commonly used. You can find more settings consulting the documentation:
#
#     http://doc.scrapy.org/en/latest/topics/settings.html
#     http://scrapy.readthedocs.org/en/latest/topics/downloader-middleware.html
#     http://scrapy.readthedocs.org/en/latest/topics/spider-middleware.html

BOT_NAME = 'novel'

SPIDER_MODULES = ['novel.spiders']
NEWSPIDER_MODULE = 'novel.spiders'

ITEM_PIPELINES = {
    # 'novel.pipelines.NovelPipeline': 300,
    'novel.pipelines.MysqlPipeline': 300,
}

LOG_LEVEL = 'ERROR'

# mysql setting
MYSQL_HOST = '127.0.0.1'
MYSQL_DB_NAME = 'novel'
MYSQL_USER = 'root'
MYSQL_PASSWORD = 'root'
MYSQL_PORT = 3306
