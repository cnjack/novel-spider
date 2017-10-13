# -*- coding:utf-8 -*-
# auth: jack
# date 03/30/2017


# 获取标签后面的value
def get_label_value(old_str):
    if len(old_str) == 0:
        return ''
    index = 0
    if index == 0:
        index = old_str.find("：")
    if index == 0:
        index = old_str.find(":")
    return old_str[index + 1:]
