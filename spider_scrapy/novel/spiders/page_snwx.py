# -*- coding:utf-8 -*-
# auth: jack
# date 03/30/2017

from scrapy.spiders import CrawlSpider, Rule
from scrapy.linkextractors import LinkExtractor
from novel.tool import filter
import re
from novel.items import NovelItem


class NovelSpider(CrawlSpider):
    name = "snwx"
    start_urls = [
        "http://www.snwx8.com/toplastupdate/1.html"
    ]

    allow_domains = ['www.snwx8.com', 'www.snwx.com']

    rules = (
        Rule(LinkExtractor(allow=(r'/toplastupdate/[0-9]+\.html'))),
        Rule(LinkExtractor(allow=(r'book/[0-9]+/[0-9]+'),deny="\.html"), callback='parse_item'),
    )

    title_xpath = './/div[@class="infotitle"]/h1/text()'
    auth_xpath = './/div[@class="infotitle"]/i[1]/text()'
    cover_xpath = '//div[@id="fmimg"]/img/@src'
    style_xpath = '//div[@class="infotitle"]/i[2]/text()'
    status_xpath = '//div[@class="infotitle"]/i[3]/text()'
    intro_xpath = '//div[@class="intro"]'
    chapters = '//div[@id="list"]/dl/dd/a'

    def parse_item(self, response):
        item = NovelItem()
        item["url"] = response.url
        title = response.xpath(self.title_xpath).extract()
        if len(title) > 0:
            item["title"] = title[0]
        else:
            item["title"] = ""
        auth_label = response.xpath(self.auth_xpath).extract()
        if len(auth_label) > 0:
            item['auth'] = filter.get_label_value(auth_label[0])
        else:
            item['auth'] = ""
        covers = response.xpath(self.cover_xpath).extract()
        if len(covers) > 0:
            item['cover'] = covers[0]
        else:
            item['cover'] = ""
        style_label = response.xpath(self.style_xpath).extract()
        if len(style_label) > 0:
            item['style'] = filter.get_label_value(style_label[0]).replace("小说", "")
        else:
            item['style'] = ""

        status_label = response.xpath(self.status_xpath).extract()
        if len(status_label) > 0:
            item['status'] = filter.get_label_value(status_label[0])
        else:
            item['status'] = ""

        intros = response.xpath(self.intro_xpath).extract()
        if len(intros) > 0:
            intro = re.compile(r'<[^>]*br[^>]*>', re.S).sub("\n", intros[0])
            item['intro'] = re.compile(r'<[^>]+>', re.S).sub("", intro).strip()
        else:
            item['intro'] = ""
        return item
