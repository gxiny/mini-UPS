from django.shortcuts import render, redirect
from django import forms
from django.contrib.auth import authenticate,login,logout
from django.contrib.auth.decorators import login_required
from django.db import transaction

from . import ups_comm_pb2
from .forms import *
from .models import *
from .rpc import rpc

wrong_login = "Your username or password is wrong"
wrong_user = "The user is not alive"
wrong_format = "should be number"

UPS_ADDRESS = ('vcm-3878.vm.duke.edu', 8080)

def rpc_ups(request):
    return rpc(UPS_ADDRESS, request, ups_comm_pb2.Response())

# Create your views here.

def ups(request) :
    if request.user.is_active:
        print(request.user)
        return redirect('/home/')  
    return render (request,'ups.html')

@transaction.atomic
def regist(request):
    if request.method == 'POST':
        form = SignUpForm(request.POST)
        if form.is_valid():
            form.save()
            req = ups_comm_pb2.Request()
            req.new_user = form.cleaned_data['username']
            resp = rpc_ups(req)

            if resp.error:
                return render(request, 'regist.html', {'uf': form, 'wrong_message': resp.error})

            user_id = user_id_recv(
                username = form.cleaned_data['username'],
                user_id_recv = resp.user_id,
            )
            user_id.save()
            return redirect('/login/')        
    else:
        form = SignUpForm()
    return render(request,'regist.html', {'uf':form})


def signout(request):
    logout(request)
    return render(request,'logout.html')

@login_required
def homepage(request):
    username = request.user.username
    if user_id_recv.objects.get(username = username): 
        user_id = user_id_recv.objects.get(username = username)
    else: 
        return redirect('/login/')
    #print(user_id.user_id_recv)
    
    command = ups_comm_pb2.Request()
    command.get_packages = user_id.user_id_recv
    resp = rpc_ups(command)
    test = (resp.packages) 
    #print(test)
    judge = True   
    return render (request,'homepage.html',{'username':username,'test':resp.packages,'user_id':user_id.user_id_recv,'judge':judge})

def searchpage(request):
    if request.method == "POST":    
        form = SearchForm(request.POST)
        if form.is_valid():
            pkg_ids = form.cleaned_data['tracking_number']
            req = ups_comm_pb2.Request()
            try:
                req.get_package_status.extend(int(pkg_id) for pkg_id in pkg_ids.split(','))
            except ValueError:
                return render(request, 'search.html', {'form': form,'wrong_message': wrong_format})
            resp = rpc_ups(req)
            if resp.error:
                return render(request, 'search.html', {'wrong_message': resp.error, 'form':form})    
            return render(request, 'homepage.html',{'test':resp.packages})
    else :
        form = SearchForm()
    return render(request, 'search.html', {'form': form})   


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
            
            if resp.error:
                return render(request, 'redirect.html', {'form': form,'wrong_message': resp.error})
            return redirect('/home/')
    else :
        form = RedirectForm()
    return render(request, 'redirect.html', {'form': form})
                  
# vim: ts=4:sw=4:et
