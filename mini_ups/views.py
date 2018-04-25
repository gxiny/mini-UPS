from django.shortcuts import render, redirect
from django import forms
from django.contrib.auth import authenticate,login,logout
from django.contrib.auth.decorators import login_required
from django.db import transaction

import re

from . import ups_comm_pb2
from .forms import *
from .models import *
from .rpc import rpc

worng_login = "Your username or password is wrong"
worng_user = "The user is not alive"
wrong_format = "should be number"

UPS_ADDRESS = ('vcm-3878.vm.duke.edu', 8080)

def rpc_ups(request):
    return rpc(UPS_ADDRESS, request, ups_comm_pb2.Response())

# Create your views here.
class UserForm(forms.Form):
    username = forms.CharField(label = 'username',max_length=50)
    #email = forms.CharField(label = 'email',max_length=50)
    password = forms.CharField(label = 'password',max_length=50,widget=forms.PasswordInput())

def ups(request) :
    return render (request,'ups.html')

@transaction.atomic
def regist(request):
    if request.method == 'POST':
        uf = SignUpForm(request.POST)
        if uf.is_valid():
            uf.save()
            username = uf.cleaned_data.get('username')
            email    = uf.cleaned_data.get('email')
            command = ups_comm_pb2.Request()
            command.new_user = username
            resp = rpc_ups(command)
            if resp.error is not "":
                return render(request, 'regist.html', {'uf': uf,'wrong_message': resp.error})
            
            print(resp.user_id)
            user_id = user_id_recv (
                username = username,
                user_id_recv = resp.user_id,
            )
            user_id.save()
            #return render(request,'homepage.html',{'uf':uf})
            return redirect('/login/')        
    else:
        uf = SignUpForm()
    return render (request,'regist.html',{'uf':uf})

def signin(request):
    if request.user is not None:
        if request.user.is_active:
            login(request, request.user,backend='django.contrib.auth.backends.ModelBackend')
            return redirect('/home/')
            
    if request.method == 'POST':
        uf = UserForm(request.POST)
        if uf.is_valid():
            #get username and password
            username = uf.cleaned_data['username']
            password = uf.cleaned_data['password']
        username = request.POST['username']
        raw_password = request.POST['password']
        user = authenticate(username=username, password=raw_password)
        if user is not None:
            if user.is_active:
                login(request, user)
                return redirect('/home/')
            else:
                return render(request,'wrong.html',{'uf':uf,'wrong_message':worng_user})
        else:
            return render(request,'wrong.html',{'uf':uf,'wrong_message':worng_login})                    
    else:
        uf = UserForm()
    return render (request,'login.html',{'uf':uf})


def signout(request):
    logout(request)
    return render(request,'logout.html')
    
def wrong(request):
    return render(request,'wrong.html',{'wrong':wrong})

@login_required
def homepage(request):
    username = request.user.username
    user_id = user_id_recv.objects.get(username = username)
    print(user_id.user_id_recv)
    
    command = ups_comm_pb2.Request()
    command.get_packages = user_id.user_id_recv
    resp = rpc_ups(command)
    test = (resp.packages) 
    print(test)   
    return render (request,'homepage.html',{'username':username,'test':resp.packages,'user_id':user_id.user_id_recv})

def searchpage(request) :
    if request.method == "POST":    
        form = SearchForm(request.POST)
        if form.is_valid():
            tracking_num = form.cleaned_data['tracking_number']
            command = ups_comm_pb2.Request()
            track = tracking_num.split(',')
            for each in track :
                if re.match(r'^[-]?\d+$', each) == None :
                    return render(request, 'search.html', {'form': form,'wrong_message': wrong_format})
                else :
                    
                    com = command.get_package_status.append(int(each))
                    #command.get_package_status = int(tracking_num)
            resp = rpc_ups(command)
                
            if resp.error is not "":
                return render(request, 'search.html', {'wrong_message': resp.error, 'form':form})    
            test = (resp.packages)      
            return render(request, 'search_res.html',{'test':resp.packages})
    else :
        form = SearchForm()
        
    return render(request, 'search.html', {'form': form})   
    

def search_res(request) :
    return render(request,'search_res.html') 

@login_required    
def Redirectpage(request,package_id) :
    if request.method == "POST":    
        form = RedirectForm(request.POST)
        if form.is_valid():
            
            #package_id = form.cleaned_data['x']
            x = form.cleaned_data['x']
            y = form.cleaned_data['y']
            username = request.user.username
            user_id = user_id_recv.objects.get(username = username) 
            command = ups_comm_pb2.Request() 
            command.change_destination.user_id = user_id.user_id_recv
            command.change_destination.package_id = int(package_id)
            command.change_destination.x = x
            command.change_destination.y = y 
            resp = rpc_ups(command)
            
            if resp.error is not "":
                return render(request, 'redirect.html', {'form': form,'wrong_message': resp.error})
            return redirect('/home/')
    else :
        form = RedirectForm()
    return render(request, 'redirect.html', {'form': form})
                  
# vim: ts=4:sw=4:et
