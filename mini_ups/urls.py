from django.urls import path
from django.contrib.auth import views as auth_views

from . import views


urlpatterns = [
    path('login/', auth_views.LoginView.as_view(template_name='login.html'), name='login'),
    path('register/', views.RegisterView.as_view(), name='register'),
    path('logout/', auth_views.LogoutView.as_view(), name='logout'),
    path('home/', views.homepage, name='homepage'),
    path('track/', views.TrackView.as_view(), name='track'),
    path('ups/',views.ups, name='ups'),
    path('redirect/<int:package_id>',views.Redirectpage, name='Redirectpage'),
]
