from django.urls import path
from django.contrib.auth import views as auth_views

from . import views


urlpatterns = [
    path('', views.index, name='index'),

    path('register/', views.RegisterView.as_view(), name='register'),
    path('login/', auth_views.LoginView.as_view(template_name='login.html'), name='login'),
    path('logout/', auth_views.LogoutView.as_view(), name='logout'),

    path('packages/', views.packages, name='packages'),
    path('track/', views.TrackView.as_view(), name='track'),
    path('track/<int:package_id>/', views.PackageDetailView.as_view(), name='package_detail'),
    path('redirect/<int:package_id>/', views.RedirectView.as_view(), name='Redirectpage'),
]
