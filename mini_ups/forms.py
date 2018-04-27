from django import forms
from django.core.exceptions import ValidationError
from .models import *

class UserUpdateForm(forms.ModelForm):
    upsid = forms.IntegerField(label="UPS ID", widget=forms.TextInput, disabled=True)
    class Meta:
        model = User
        fields = ('username', 'upsid', 'email', 'first_name', 'last_name')

    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)

class TrackForm(forms.Form):
    pkgids = forms.CharField(label="Tracking numbers (one per line)", widget=forms.Textarea)

    def clean_pkgids(self):
        data = self.cleaned_data['pkgids']
        cleaned = []
        errors = []
        for i in data.split('\n'):
            try:
                x = int(i)
                if not (-9223372036854775808 <= x <= 9223372036854775807):
                    raise ValueError # out of int64 range
                cleaned.append(x)
            except ValueError:
                errors.append("Invalid tracking number: " + i)
        if errors:
            raise ValidationError(errors)
        return cleaned


class RedirectForm(forms.Form):
    #tracking_number = forms.CharField()
    x = forms.IntegerField()
    y = forms.IntegerField()

# vim: ts=4:sw=4:et
