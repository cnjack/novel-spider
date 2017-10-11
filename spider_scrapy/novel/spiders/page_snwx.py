# -*- coding:utf-8 -*-
#auth: jack
#date 03/30/2017

from scrapy.spiders import CrawlSpider,Rule
from scrapy.linkextractors import LinkExtractor
from novel.items import NovelItem

class NovelSpider(CrawlSpider):
    name = "snwxPage"
    start_urls = [
        "http://www.snwx.com/toplastupdate/1.html"
    ]

    rules = (
        Rule(LinkExtractor(allow=('/book/'),deny=('\.html')), callback='parse_item'),
        Rule(LinkExtractor(allow=('/toplastupdate/'))),
    )

    name_xpath = './/div[@class="infotitle"]/h1/text()'

    def parse_item(self,response):
        item = NovelItem()
        name = response.xpath(self.name_xpath).extract()
        if len(name) > 0:
            item["name"] = name[0].encode('utf-8')
        item["url"] = response.url
        return item
