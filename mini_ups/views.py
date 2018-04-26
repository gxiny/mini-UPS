from django.contrib.auth.decorators import login_required
from django.contrib.auth.forms import UserCreationForm
from django.contrib.auth.mixins import LoginRequiredMixin
from django.db import transaction
from django.shortcuts import render
from django.urls import reverse_lazy
from django.views.generic.base import TemplateView
from django.views.generic.edit import FormView

from . import ups_comm_pb2
from .forms import *
from .models import *
from .rpc import rpc

UPS_ADDRESS = ('vcm-3878.vm.duke.edu', 8080)

def rpc_ups(request):
    return rpc(UPS_ADDRESS, request, ups_comm_pb2.Response())

# Create your views here.

def index(request):
    return render(request, 'index.html')


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


@login_required
def packages(request):
    ups_id = UpsId.objects.get(user=request.user)
    req = ups_comm_pb2.Request()
    req.get_package_list.user_id = ups_id.value
    resp = rpc_ups(req)
    return render(request, 'packages.html', {'packages': resp.package_list.packages})


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
