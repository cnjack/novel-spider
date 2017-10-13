# -*- coding: utf-8 -*-

# Define here the models for your scraped items
#
# See documentation in:
# http://doc.scrapy.org/en/latest/topics/items.html

import scrapy


class NovelItem(scrapy.Item):
    title = scrapy.Field()
    url = scrapy.Field()
    cover = scrapy.Field()
    auth = scrapy.Field()
    style = scrapy.Field()
    intro = scrapy.Field()
    status = scrapy.Field()
    pass
