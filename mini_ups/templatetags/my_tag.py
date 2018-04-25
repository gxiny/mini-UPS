from django import template
import time

register = template.Library()

@register.filter
def get_time(t1):
    t2 = time.localtime(t1)
    t  = time.strftime('%Y-%m-%d %H:%M:%S',t2)
    return t
