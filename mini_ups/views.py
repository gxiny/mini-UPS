from django import forms
from django.contrib.auth import authenticate,login,logout
from django.contrib.auth.decorators import login_required
from django.contrib.auth.forms import UserCreationForm
from django.db import transaction
from django.shortcuts import render, redirect
from django.urls import reverse_lazy
from django.views.generic.edit import FormView

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


class RegisterView(FormView):
    template_name = 'register.html'
    form_class = UserCreationForm
    success_url = reverse_lazy('login')

    @transaction.atomic
    def form_valid(self, form):
        req = ups_comm_pb2.Request()
        req.new_user = form.cleaned_data['username']
        resp = rpc_ups(req)

        if resp.error:
            form.add_error(None, resp.error)
            return self.form_invalid(form)
        user = form.save()
        UpsId.objects.create(user=user, ups_id=resp.user_id)
        return super().form_valid(form)


@login_required
def homepage(request):
    ups_id = UpsId.objects.get(user=request.user)
    req = ups_comm_pb2.Request()
    req.get_packages = ups_id.value
    resp = rpc_ups(req)
    judge = True
    return render(request, 'homepage.html', {'test': resp.packages, 'judge': judge})


class TrackView(FormView):
    template_name = 'track.html'
    form_class = TrackForm

    def form_valid(self, form):
        req = ups_comm_pb2.Request()
        req.get_package_status.extend(form.cleaned_data['pkgids'])
        resp = rpc_ups(req)

        if resp.error:
            form.add_error(None, resp.error)
            return self.form_invalid(form)
        return self.render_to_response(self.get_context_data(results=resp.packages))

    def get(self, request, *args, **kwargs):
        q = request.GET.get('q')
        if not q:
            return super().get(request, *args, **kwargs)
        form = self.get_form_class()({'pkgids': q})
        if form.is_valid():
            return self.form_valid(form)
        else:
            return self.form_invalid(form)

@login_required    
def Redirectpage(request,package_id) :
    if request.method == "POST":    
        form = RedirectForm(request.POST)
        if form.is_valid():
            x = form.cleaned_data['x']
            y = form.cleaned_data['y']
            ups_id = UpsId.objects.get(user=request.user)
            req = ups_comm_pb2.Request()
            command = req.change_destination
            command.user_id = ups_id.value
            command.package_id = int(package_id)
            command.x = x
            command.y = y
            resp = rpc_ups(req)

            if resp.error:
                return render(request, 'redirect.html', {'form': form,'wrong_message': resp.error})
            return redirect('/home/')
    else :
        form = RedirectForm()
    return render(request, 'redirect.html', {'form': form})
                  
# vim: ts=4:sw=4:et
