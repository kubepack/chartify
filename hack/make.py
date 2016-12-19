#!/usr/bin/env python

# ref: https://github.com/ellisonbg/antipackage
import antipackage
from github.appscode.libbuild import libbuild

import datetime
import io
import json
import os
import os.path
import socket
import subprocess
import sys
from collections import OrderedDict
from os.path import expandvars


libbuild.REPO_ROOT = expandvars('$GOPATH') + '/src/github.com/appscode/chartify'
BUILD_METADATA = libbuild.metadata(libbuild.REPO_ROOT)
libbuild.BIN_MATRIX = {
    'chartify': {
        'type': 'go',
        'go_version': True,
        'distro': {
            'darwin': ['386', 'amd64'],
            'linux': ['arm', '386', 'amd64'],
            'windows': ['386', 'amd64']
        }
    }
}
libbuild.BUCKET_MATRIX = {
    'prod': 'gs://appscode-cdn',
    'dev': 'gs://appscode-dev'
}


def call(cmd, stdin=None, cwd=libbuild.REPO_ROOT):
    print(cmd)
    return subprocess.call([expandvars(cmd)], shell=True, stdin=stdin, cwd=cwd)


def die(status):
    if status:
        sys.exit(status)


def version():
    # json.dump(BUILD_METADATA, sys.stdout, sort_keys=True, indent=2)
    for k in sorted(BUILD_METADATA):
        print k + '=' + BUILD_METADATA[k]


def fmt():
    die(call('goimports -w pkg main.go'))
    call('gofmt -w main.go')
    call('go fmt ./pkg/...')


def lint():
    call('golint ./pkg/... min.go')


def vet():
    call('go vet ./pkg/...')


def build_cmd(name):
    cfg = libbuild.BIN_MATRIX[name]
    if cfg['type'] == 'go':
        if 'distro' in cfg.keys():
            for goos, archs in cfg['distro'].iteritems():
                for goarch in archs:
                    libbuild.go_build(name, goos, goarch, main='main.go')
        else:
            libbuild.go_build(name, libbuild.GOHOSTOS, libbuild.GOHOSTARCH, main='main.go')


def build_cmds():
    fmt()
    for name in libbuild.BIN_MATRIX.keys():
        build_cmd(name)


def build_all():
    build_cmds()


def push_all():
    dist = libbuild.REPO_ROOT + '/dist'
    for name in os.listdir(dist):
        d = dist + '/' + name
        if os.path.isdir(d):
            push(d)


def push_cmd(name):
    bindir = libbuild.REPO_ROOT + '/dist/' + name
    push(bindir)


def push(bindir):
    call('rm -f *.md5', cwd=bindir)
    call('rm -f *.sha1', cwd=bindir)
    for f in os.listdir(bindir):
        if os.path.isfile(bindir + '/' + f):
            libbuidl.upload_to_cloud(bindir, f, BUILD_METADATA['version'])


def update_registry():
    libbuiold.update_registry(BUILD_METADATA['version'])


def install():
    die(call('GO15VENDOREXPERIMENT=1 ' + libbuild.GOC + ' install .'))


def default():
    fmt()
    die(call('GO15VENDOREXPERIMENT=1 ' + libbuild.GOC + ' install .'))


if __name__ == "__main__":
    if len(sys.argv) > 1:
        # http://stackoverflow.com/a/834451
        # http://stackoverflow.com/a/817296
        globals()[sys.argv[1]](*sys.argv[2:])
    else:
        default()
