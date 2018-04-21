from django.shortcuts import render,render_to_response,get_object_or_404
from django.utils import timezone
from django.shortcuts import redirect
from django import forms
from django.http import HttpResponse,HttpResponseRedirect
from django.template import RequestContext
from .forms import *
from  django.contrib.auth.models import User
from .models import *
import socket,re
from django.contrib.auth.models import AbstractUser
from django.contrib.auth import authenticate,login,logout
from django.contrib.auth.decorators import login_required
from django.core.exceptions import ObjectDoesNotExist
from . import ups_comm_pb2

worng_login = "Your username or email or password is wrong"
worng_user = "The user is not alive"
wrong_format = "should be number"

# Create your views here.
class UserForm(forms.Form):
    username = forms.CharField(label = 'username',max_length=50)
    email = forms.CharField(label = 'email',max_length=50)
    password = forms.CharField(label = 'password',max_length=50,widget=forms.PasswordInput())


def conn() :
    clientsocket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    host = socket.gethostname()
    port = 9005
    clientsocket.connect((host,port))
    return clientsocket
    
def regist(request):
    if request.method == 'POST':
        uf = SignUpForm(request.POST)
        if uf.is_valid():
            #get data
            uf.save()
            username = uf.cleaned_data.get('username')
            password = uf.cleaned_data.get('password1')
            first_name = uf.cleaned_data.get('first_name')
            last_name = uf.cleaned_data.get('last_name')
            email = uf.cleaned_data.get('email')
            comm_ups(username)
            #return render(request,'homepage.html',{'uf':uf})
            return redirect('/login/')        
    else:
        uf = SignUpForm()
    return render (request,'regist.html',{'uf':uf})

def comm_ups(username) :
    clientsocket = conn()
    
    command = ups_comm_pb2.FCommands()
    command.fuser.fusername = username
    send_mess = command.SerializeToString()
    clientsocket.send(send_mess)
    msg = clientsocket.recv(1024)
    resp = ups_comm_pb2.FResponse()
    resp.ParseFromString(msg)
    print(resp.buser.buser_id)
    user_id = user_id_recv (
        username = username,
        user_id_recv = resp.buser.buser_id,
        )
    user_id.save()
    #return resp.buser_id
    clientsocket.close()
    


def signin(request):
    if request.method == 'POST':
        uf = UserForm(request.POST)
       # uf = UserCreationForm(request.POST)
        if uf.is_valid():
            #get username and password
            username = uf.cleaned_data['username']
            password = uf.cleaned_data['password']
            #compare with database
            #user = User.objects.filter(username__exact = username,password__exact = password)
        username = request.POST['username']
        raw_password = request.POST['password']
        user = authenticate(username=username, password=raw_password)
        if user is not None:
            if user.is_active:
                login(request, user)
                return redirect('/home/')
            else:
                return render(request,'wrong.html',{'wrong_message':worng_user})
        else:
            return render(request,'wrong.html',{'wrong_message':worng_login})                    
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
    
    clientsocket = conn()
    
    command = ups_comm_pb2.FCommands()
    command.buser_id.buser_id = user_id.user_id_recv
    send_mess = command.SerializeToString()
    clientsocket.send(send_mess)
    
    msg = clientsocket.recv(1024)
    resp = ups_comm_pb2.FResponse()
    resp.ParseFromString(msg)
    test = (resp.pack_info)    
    return render (request,'homepage.html',{'username':username,'test':resp.pack_info})

def searchpage(request) :
    if request.method == "POST":    
        form = SearchForm(request.POST)
        if form.is_valid():
            tracking_num = form.cleaned_data['tracking_number']
            if re.match(r'^[-]?\d+$', tracking_num) == None :
                return render(request, 'search.html', {'wrong_message': wrong_format})
            else :
                clientsocket = conn() 
                command = ups_comm_pb2.FCommands()
                command.track.package_id = int(tracking_num)
                send_mess = command.SerializeToString()
                clientsocket.send(send_mess)
                msg = clientsocket.recv(1024)
                resp = ups_comm_pb2.FResponse()
                resp.ParseFromString(msg)
                test = (resp.pack_info)          
                return render(request, 'search_res.html',{'test':resp.pack_info})
    else :
        form = SearchForm()
        
    return render(request, 'search.html', {'form': form})   
    

def search_res(request) :
    return render(request,'search_res.html') 
    
def Redirectpage(request,package_id) :
    if request.method == "POST":    
        form = RedirectForm(request.POST)
        if form.is_valid():
           x = form.cleaned_data['x']
           y = form.cleaned_data['y']
           clientsocket = conn() 
           command = ups_comm_pb2.FCommands() 
           command.transfer.package_id = package_id
           command.transfer.x = x
           command.transfer.y = y 
           send_mess = command.SerializeToString()
           clientsocket.send(send_mess)
           msg = clientsocket.recv(1024)
           resp = ups_comm_pb2.FResponse()
           resp.ParseFromString(msg)
           redirect('/homepage/')
    else :
        form = RedirectForm()
    return render(request, 'redirect.html', {'form': form})
                  