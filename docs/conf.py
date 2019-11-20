from docs_conf.conf import *
linkcheck_ignore = [
  'http://localhost.*',
  'http://127.0.0.1.*',
  'https://gerrit.o-ran-sc.org.*'
]
branch = 'master'
intersphinx_mapping = {}

intersphinx_mapping['ric-plt-e2mgr'] = ('https://docs.o-ran-sc.org/projects/o-ran-sc-ric-plt-e2mgr/en/%s' % branch, None)
