from django.db import models
from django.contrib.auth.models import User


# Create your models here.
class user_id_recv(models.Model):
    username = models.CharField(max_length = 20, null=True)
    user_id_recv = models.IntegerField(blank=True, null=True)
