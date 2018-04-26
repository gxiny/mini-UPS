from django.urls import path
from django.contrib.auth import views as auth_views

from . import views


urlpatterns = [
    path('login/', auth_views.LoginView.as_view(template_name='login.html'), name='login'),
    path('regist/', views.regist,name='regist'),
    path('logout/',views.signout,name='logout'),
    path('home/', views.homepage, name='homepage'),
    path('search/', views.searchpage, name='searchpage'),
    path('ups/',views.ups, name='ups'),
    path('redirect/<int:package_id>',views.Redirectpage, name='Redirectpage'),
]
