from django.contrib.auth.decorators import login_required
from django.contrib.auth.forms import UserCreationForm
from django.contrib.auth.models import User
from django.contrib.auth.mixins import LoginRequiredMixin
from django.core.paginator import Paginator
from django.db import transaction
from django.shortcuts import render
from django.urls import reverse, reverse_lazy
from django.views.generic.base import TemplateView
from django.views.generic.edit import FormView, UpdateView

import os

from . import ups_comm_pb2
from .forms import *
from .models import *
from .rpc import rpc

UPS_ADDRESS = (
    os.getenv('UPS_HOST', 'localhost'),
    int(os.getenv('UPS_PORT', '8080')),
)

def rpc_ups(request):
    return rpc(UPS_ADDRESS, request, ups_comm_pb2.Response())

# Create your views here.

class PagedPkgList:
    def __init__(self, upsid=None):
        self.req = req = ups_comm_pb2.Request()
        self.cmd = cmd = req.get_package_list
        if upsid is not None:
            cmd.user_id = upsid

    def count(self):
        self.cmd.offset = 0
        self.cmd.limit = 0
        resp = rpc_ups(self.req)
        return resp.package_list.total

    def __getitem__(self, i):
        if not isinstance(i, slice):
            raise TypeError
        if i.step not in (None, 1):
            raise ValueError
        self.cmd.offset = i.start
        self.cmd.limit = i.stop - i.start
        resp = rpc_ups(self.req)
        return resp.package_list.packages

PAGE_SIZE = 10

def index(request):
    page = request.GET.get('page')
    paginator = Paginator(PagedPkgList(), PAGE_SIZE)
    return render(request, 'index.html', {'package_list': paginator.get_page(page)})


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
        UpsId.objects.create(user=user, value=resp.user_id)
        return super().form_valid(form)


class UserUpdateView(LoginRequiredMixin, UpdateView):
    model = User
    template_name = 'profile.html'
    form_class = UserUpdateForm
    success_url = reverse_lazy('profile')

    def get_initial(self):
        user = self.request.user
        return {'upsid': user.upsid.value}

    def get_object(self, queryset=None):
        return self.request.user


@login_required
def packages(request):
    ups_id = UpsId.objects.get(user=request.user)
    page = request.GET.get('page')
    paginator = Paginator(PagedPkgList(upsid=ups_id.value), PAGE_SIZE)
    return render(request, 'packages.html', {'packages': paginator.get_page(page)})


class TrackView(FormView):
    template_name = 'track.html'
    form_class = TrackForm

    def form_valid(self, form):
        req = ups_comm_pb2.Request()
        req.get_package_list.package_ids.extend(form.cleaned_data['pkgids'])
        resp = rpc_ups(req)

        if resp.error:
            form.add_error(None, resp.error)
            return self.form_invalid(form)
        return self.render_to_response(self.get_context_data(results=resp.package_list.packages))

    def get(self, request, *args, **kwargs):
        q = request.GET.get('q')
        if not q:
            return super().get(request, *args, **kwargs)
        form = self.get_form_class()({'pkgids': q})
        if form.is_valid():
            return self.form_valid(form)
        else:
            return self.form_invalid(form)

class PackageDetailView(TemplateView):
    template_name = 'package.html'

    def get(self, request, *args, **kwargs):
        package_id = self.kwargs['package_id']
        req = ups_comm_pb2.Request()
        req.get_package_detail = package_id
        resp = rpc_ups(req)
        return self.render_to_response(self.get_context_data(
            package_id=package_id,
            package=resp.package_detail))

class RedirectView(LoginRequiredMixin, FormView):
    template_name = 'redirect.html'
    form_class = RedirectForm
    success_url = reverse_lazy('packages')

    def get_success_url(self):
        return reverse('package_detail', kwargs={'package_id': self.kwargs['package_id']})

    def form_valid(self, form):
        x = form.cleaned_data['x']
        y = form.cleaned_data['y']
        ups_id = UpsId.objects.get(user=self.request.user)
        req = ups_comm_pb2.Request()
        command = req.change_destination
        command.user_id = ups_id.value
        command.package_id = int(self.kwargs['package_id'])
        command.x = x
        command.y = y
        resp = rpc_ups(req)

        if resp.error:
            form.add_error(None, resp.error)
            return self.form_invalid(form)
        return super().form_valid(form)

# vim: ts=4:sw=4:et
