from django import forms
from .models import *


class SearchForm(forms.Form):
    tracking_number = forms.CharField(label="tracking_number(seperate name by ',')",widget = forms.Textarea)

class RedirectForm(forms.Form):
    #tracking_number = forms.CharField()
    x = forms.IntegerField()
    y = forms.IntegerField()
