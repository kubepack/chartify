package cmd

import "mime/multipart"

type ChartFile struct {
	Data *multipart.FileHeader `form:"Data"`
}

var serverIndex = `apiVersion: v1
entries:
  alp:
  - apiVersion: v1
    created: 2017-01-09T19:54:58.393502011+06:00
    description: A Helm chart for Kubernetes
    digest: 705fb5100b5245bc67e07de8a2bed47fac5375349d8c42e2cdcf0e0b57d5f328
    name: alp
    urls:
    - http://127.0.0.1:8879/alp-0.1.0.tgz
    version: 0.1.0
  alpine-pod:
  - apiVersion: v1
    created: 2017-01-09T19:54:58.393662203+06:00
    description: A Helm chart for Kubernetes
    digest: e5ce16d773562fa07d53a00a20f053126410a6d38d52c80fb4f7ec882e78ee87
    name: alpine-pod
    urls:
    - http://127.0.0.1:8879/alpine-pod-0.1.0.tgz
    version: 0.1.0
  base:
  yakuza:
  - apiVersion: v1
    created: 2017-01-09T19:54:58.396601502+06:00
    description: A Helm chart for Kubernetes
    digest: 4bf04b908910fe2bdbd4a9c1f982d78852a1ab5b839742cfc3bba8d5789c45ab
    name: yakuza
    urls:
    - http://127.0.0.1:8879/yakuza-0.1.0.tgz
    version: 0.1.0
generated: 2017-01-09T19:54:58.393236164+06:00`
