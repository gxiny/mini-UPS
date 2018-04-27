#!/usr/bin/env python

import bjoern

from mysite.wsgi import application

bjoern.run(application, '0.0.0.0', 8000)
