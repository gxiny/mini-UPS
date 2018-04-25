from django.conf.urls import url
from . import views


urlpatterns = [
    url(r'^login/$', views.signin,name='login'),
    url(r'^regist/$', views.regist,name='regist'),
    url(r'^logout/$',views.signout,name='logout'),
    url(r'^wrong/$',views.wrong,name='wrong'),
    url(r'^home/$', views.homepage, name='homepage'),
    url(r'^search/$', views.searchpage, name='searchpage'),
    #url(r'^search_res/$',views.search_res, name='search_res'),
    #url(r'^redirect/$',views.Redirectpage, name='Redirectpage'),
    url(r'^ups/$',views.ups, name='ups'),
    url(r'^redirect/(?P<package_id>\d+)$',views.Redirectpage, name='Redirectpage'),
]
