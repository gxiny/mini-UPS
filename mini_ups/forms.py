from django import forms
from django.forms import ModelForm,modelformset_factory
from django.contrib.auth.models import User
from django.contrib.auth.forms import UserCreationForm
from .models import *


class SignUpForm(UserCreationForm): #UserCreationForm
    username   =  forms.CharField(max_length=50)
    first_name = forms.CharField()
    last_name  = forms.CharField()
    password1  = forms.CharField(widget=forms.PasswordInput())
    password2  = forms.CharField(widget=forms.PasswordInput())
    email = forms.EmailField()
    class Meta:
        model = User
        fields = ('username','first_name','last_name','password1', 'password2','email')
    def __init__(self,*args, **kwargs):
        super(SignUpForm, self).__init__(*args, **kwargs)

class SearchForm(forms.Form):
    tracking_number = forms.CharField(label="tracking_number(seperate name by ',')",widget = forms.Textarea)

class RedirectForm(forms.Form):
    #tracking_number = forms.CharField()
    x = forms.IntegerField()
    y = forms.IntegerField()