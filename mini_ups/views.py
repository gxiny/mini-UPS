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

try:
    from io import BytesIO
except ImportError:
    from StringIO import StringIO as BytesIO

# Create your views here.
class UserForm(forms.Form):
    username = forms.CharField(label = 'username',max_length=50)
    email = forms.CharField(label = 'email',max_length=50)
    password = forms.CharField(label = 'password',max_length=50,widget=forms.PasswordInput())

def encode(number):
    buf = []
    while True:
        towrite = number & 0x7f
        number >>= 7
        if number:
            buf.append(towrite | 0x80)
        else:
            buf.append(towrite)
            break
    return bytes(buf)
    
def decode_stream(stream):
    """Read a varint from `stream`"""
    shift = 0
    result = 0
    while True:
        i = _read_one(stream)
        result |= (i & 0x7f) << shift
        shift += 7
        if not (i & 0x80):
            break

    return result
    
def _read_one(stream):
    """Read a byte from the file (as an integer)
    raises EOFError if the stream ends while reading bytes.
    """
    c = stream.recv(1)
    if c == b'':
        raise EOFError("Unexpected EOF while reading bytes")
    return ord(c)
 
def conn() :
    clientsocket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    host = "vcm-3878.vm.duke.edu"
    print(host)
    port = 8080
    clientsocket.connect((host,port))
    return clientsocket

def index(request) :
    return render (request,'ups.html')

def regist(request):
    if request.method == 'POST':
        uf = SignUpForm(request.POST)
        if uf.is_valid():
            #get data
            uf.save()
            username = uf.cleaned_data.get('username')
            command = ups_comm_pb2.Request()
            command.new_user = username
            resp = comm_ups(command)
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

def comm_ups(command) :
    clientsocket = conn() 
    send_mess = command.SerializeToString()
    leng = encode(len(send_mess))
    print(leng)
    clientsocket.send(leng)
    clientsocket.send(send_mess)
    res = decode_stream(clientsocket)
    print(res)
    msg = clientsocket.recv(res)
    resp = ups_comm_pb2.Response()
    resp.ParseFromString(msg)
    
    #return resp.buser_id
    clientsocket.close()
    return resp
    


def signin(request):
    if request.user is not None:
        if request.user.is_active:
            login(request, request.user)
            return redirect('/home/')
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
    resp = comm_ups(command)
    test = (resp.packages)    
    return render (request,'homepage.html',{'username':username,'test':resp.packages})

def searchpage(request) :
    if request.method == "POST":    
        form = SearchForm(request.POST)
        if form.is_valid():
            tracking_num = form.cleaned_data['tracking_number']
            if re.match(r'^[-]?\d+$', tracking_num) == None :
                return render(request, 'search.html', {'form': form,'wrong_message': wrong_format})
            else :
                
                command = ups_comm_pb2.Request()
                command.get_package_status = int(tracking_num)
                resp = comm_ups(command)
                
                if resp.error is not "":
                    return render(request, 'search.html', {'wrong_message': resp.error, 'form':form})    
                test = (resp.packages)      
                return render(request, 'search_res.html',{'test':resp.packages})
    else :
        form = SearchForm()
        
    return render(request, 'search.html', {'form': form})   
    

def search_res(request) :
    return render(request,'search_res.html') 
    
def Redirectpage(request) :
    if request.method == "POST":    
        form = RedirectForm(request.POST)
        if form.is_valid():
            
            package_id = form.cleaned_data['x']
            x = form.cleaned_data['x']
            y = form.cleaned_data['y']
           
            command = ups_comm_pb2.Request() 
            command.change_destination.package_id = package_id
            command.change_destination.x = x
            command.change_destination.y = y 
            resp = comm_ups(command)
            
            if resp.error is not "":
                return render(request, 'redirect.html', {'form': form,'wrong_message': resp.error})
            return redirect('/home/')
    else :
        form = RedirectForm()
    return render(request, 'redirect.html', {'form': form})
                  
